// Package record reads and writes Ignition's Markdown records: a flat
// frontmatter block, then a body of `## ` sections. It has no file I/O.
package record

import (
	"fmt"
	"regexp"
	"strings"
)

// Doc is one Markdown record.
type Doc struct {
	keys   []string
	values map[string]string
	Body   string
}

var keyPattern = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

// Parse reads a record. Frontmatter is optional; when present it starts on the
// first line with `---`, ends at the next `---` line, and holds only
// `key: value` lines.
func Parse(text string) (*Doc, error) {
	d := &Doc{values: map[string]string{}}
	text = strings.ReplaceAll(text, "\r\n", "\n")
	if !strings.HasPrefix(text, "---\n") {
		d.Body = text
		return d, nil
	}
	rest := text[len("---\n"):]
	end := strings.Index(rest, "\n---\n")
	var block string
	switch {
	case strings.HasPrefix(rest, "---\n"):
		block, d.Body = "", rest[len("---\n"):]
	case end >= 0:
		block, d.Body = rest[:end], rest[end+len("\n---\n"):]
	case strings.HasSuffix(rest, "\n---"):
		block, d.Body = strings.TrimSuffix(rest, "\n---"), ""
	default:
		return nil, fmt.Errorf("frontmatter is not closed with ---")
	}
	for i, line := range strings.Split(block, "\n") {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}
		k, v, ok := strings.Cut(line, ":")
		k = strings.TrimSpace(k)
		if !ok || !keyPattern.MatchString(k) {
			return nil, fmt.Errorf("frontmatter line %d: want `key: value`, got %q", i+2, line)
		}
		if _, dup := d.values[k]; dup {
			return nil, fmt.Errorf("frontmatter line %d: duplicate key %q", i+2, k)
		}
		d.keys = append(d.keys, k)
		d.values[k] = unquote(strings.TrimSpace(v))
	}
	return d, nil
}

func unquote(v string) string {
	if len(v) >= 2 && (v[0] == '"' && v[len(v)-1] == '"' || v[0] == '\'' && v[len(v)-1] == '\'') {
		return v[1 : len(v)-1]
	}
	return v
}

// Get returns a frontmatter value, or "" if the key is absent.
func (d *Doc) Get(key string) string { return d.values[key] }

// Set sets a frontmatter value, keeping the position of an existing key.
func (d *Doc) Set(key, value string) {
	if _, ok := d.values[key]; !ok {
		d.keys = append(d.keys, key)
	}
	d.values[key] = value
}

// String renders the record.
func (d *Doc) String() string {
	if len(d.keys) == 0 {
		return d.Body
	}
	var b strings.Builder
	b.WriteString("---\n")
	for _, k := range d.keys {
		v := d.values[k]
		if strings.ContainsAny(v, ":#") || strings.TrimSpace(v) != v {
			v = fmt.Sprintf("%q", v)
		}
		if v == "" {
			fmt.Fprintf(&b, "%s:\n", k)
			continue
		}
		fmt.Fprintf(&b, "%s: %s\n", k, v)
	}
	b.WriteString("---\n")
	b.WriteString(d.Body)
	return b.String()
}

var comments = regexp.MustCompile(`(?s)<!--.*?-->`)

// Section returns the content of the `## name` section with guidance comments
// removed and whitespace trimmed, and whether the section exists.
func (d *Doc) Section(name string) (string, bool) {
	body := comments.ReplaceAllString(d.Body, "")
	var (
		in    bool
		found bool
		out   []string
	)
	for _, line := range strings.Split(body, "\n") {
		if h, ok := strings.CutPrefix(line, "## "); ok {
			if in {
				break
			}
			if strings.EqualFold(strings.TrimSpace(h), name) {
				in, found = true, true
			}
			continue
		}
		if in {
			out = append(out, line)
		}
	}
	return strings.TrimSpace(strings.Join(out, "\n")), found
}

// Filled reports whether a section exists and has content beyond guidance.
func (d *Doc) Filled(name string) bool {
	s, ok := d.Section(name)
	return ok && s != ""
}

// Task is one row of a plan's task table.
type Task struct {
	ID, Task, Acceptance, Status, Evidence string
}

// Tasks reads the table in the `## Tasks` section. The header row and the
// separator row are skipped; rows whose cells are all empty are ignored.
func (d *Doc) Tasks() ([]Task, error) {
	s, ok := d.Section("Tasks")
	if !ok {
		return nil, fmt.Errorf("plan has no `## Tasks` section")
	}
	var out []Task
	header := true
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "|") {
			continue
		}
		cells := strings.Split(strings.Trim(line, "|"), "|")
		for i := range cells {
			cells[i] = strings.TrimSpace(cells[i])
		}
		if header {
			header = false
			continue
		}
		if strings.Trim(strings.Join(cells, ""), "-: ") == "" {
			continue
		}
		if len(cells) != 5 {
			return nil, fmt.Errorf("task row %q: want 5 cells (ID, Task, Acceptance, Status, Evidence)", line)
		}
		out = append(out, Task{cells[0], cells[1], cells[2], strings.ToLower(cells[3]), cells[4]})
	}
	return out, nil
}

// Entry is one log entry.
type Entry struct {
	At, Type, By, Text string
}

const entrySep = " · "

// ParseLog reads a log: entries are `## <time> · <type> · <actor>` headings
// followed by text.
func ParseLog(text string) ([]Entry, error) {
	var (
		out []Entry
		cur *Entry
	)
	for i, line := range strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n") {
		if h, ok := strings.CutPrefix(line, "## "); ok {
			parts := strings.Split(h, entrySep)
			if len(parts) != 3 {
				return nil, fmt.Errorf("log line %d: want `## time · type · actor`, got %q", i+1, line)
			}
			out = append(out, Entry{At: parts[0], Type: parts[1], By: parts[2]})
			cur = &out[len(out)-1]
			continue
		}
		if cur != nil {
			cur.Text += line + "\n"
		}
	}
	for i := range out {
		out[i].Text = strings.TrimSpace(out[i].Text)
	}
	return out, nil
}

// FormatEntry renders a log entry for appending.
func FormatEntry(e Entry) string {
	return fmt.Sprintf("\n## %s%s%s%s%s\n\n%s\n", e.At, entrySep, e.Type, entrySep, e.By, strings.TrimSpace(e.Text))
}
