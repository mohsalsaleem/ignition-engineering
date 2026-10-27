// Package workspace reads and changes a product's Ignition records: its work
// items, their artifacts, evidence and log, and its decisions.
package workspace

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/mohsalsaleem/ignition-engineering/internal/model"
	"github.com/mohsalsaleem/ignition-engineering/internal/record"
	"github.com/mohsalsaleem/ignition-engineering/templates"
)

// Workspace is a product's records folder: PROJECT.md, AGENTS.md, work/ and
// decisions/.
type Workspace struct {
	Root string
	Now  func() time.Time
	Rand io.Reader
}

func Open(root string) *Workspace {
	return &Workspace{Root: root, Now: time.Now, Rand: rand.Reader}
}

// Find locates the workspace for a working directory: a folder holding
// PROJECT.md and work/, or such a folder at .ignition/ in it or a parent.
func Find(dir string) (string, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	for d := dir; ; d = filepath.Dir(d) {
		for _, c := range []string{d, filepath.Join(d, ".ignition")} {
			if isWorkspace(c) {
				return c, nil
			}
		}
		if filepath.Dir(d) == d {
			return "", fmt.Errorf("no Ignition workspace in %s or its parents (run `ignition init`, or pass --dir)", dir)
		}
	}
}

func isWorkspace(dir string) bool {
	if _, err := os.Stat(filepath.Join(dir, "PROJECT.md")); err != nil {
		return false
	}
	st, err := os.Stat(filepath.Join(dir, "work"))
	return err == nil && st.IsDir()
}

// Init creates a workspace's folders and profile. It refuses to overwrite.
func Init(root, product string) error {
	if isWorkspace(root) {
		return fmt.Errorf("%s is already an Ignition workspace", root)
	}
	for _, d := range []string{"work", "decisions"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			return err
		}
	}
	for _, name := range []string{"PROJECT.md", "AGENTS.md"} {
		text, err := templates.Render(name, map[string]string{"Product": product})
		if err != nil {
			return err
		}
		if err := writeNew(filepath.Join(root, name), text); err != nil {
			return err
		}
	}
	return writeNew(filepath.Join(root, ".gitattributes"), gitattributes)
}

func writeNew(path, text string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	if _, err := f.WriteString(text); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

const idAlphabet = "0123456789abcdefghjkmnpqrstvwxyz"

// NewID returns prefix-YYYYMMDD-xxxx: sortable by day, and random enough that
// agents working in parallel branches don't collide.
func (w *Workspace) NewID(prefix string) (string, error) {
	b := make([]byte, 6)
	if _, err := io.ReadFull(w.Rand, b); err != nil {
		return "", err
	}
	for i := range b {
		b[i] = idAlphabet[int(b[i])%len(idAlphabet)]
	}
	return fmt.Sprintf("%s-%s-%s", prefix, w.Now().UTC().Format("20060102"), b), nil
}

func (w *Workspace) stamp() string { return w.Now().UTC().Format(time.RFC3339) }

// Actor is who acted: a person, or an agent working for one. Agent, Model and
// Session are set when an agent acted; an empty value means unknown.
type Actor struct {
	Name    string // the person or agent's name
	Agent   string // the agent tool, e.g. claude-code, codex, cursor, gemini-cli
	Model   string // the exact model ID, e.g. claude-opus-5-5
	Session string // the agent's session ID, to find its transcript
}

func (a Actor) name() string {
	if strings.TrimSpace(a.Name) == "" {
		return "unknown"
	}
	return a.Name
}

// Item is one loaded work item.
type Item struct {
	ID       string
	Dir      string
	Brief    *record.Doc
	Docs     map[model.Artifact]*record.Doc
	Evidence []*record.Doc
	Log      []record.Entry
	Children []*Item
	Problems []string // records that couldn't be read
}

func (it *Item) Kind() model.Kind   { return model.Kind(it.Brief.Get("kind")) }
func (it *Item) Stage() model.Stage { return model.Stage(it.Brief.Get("stage")) }
func (it *Item) Title() string      { return it.Brief.Get("title") }

func (it *Item) Facts() model.Facts {
	return model.Facts{
		Kind:      it.Kind(),
		Risk:      model.Risk(it.Brief.Get("risk")),
		NoRelease: it.Brief.Get("release") == "none",
	}
}

func (w *Workspace) itemDir(id string) string { return filepath.Join(w.Root, "work", id) }

// NewItem creates a work item in its kind's first stage.
func (w *Workspace) NewItem(kind model.Kind, title, risk, size, parent string, by Actor) (*Item, error) {
	if strings.TrimSpace(title) == "" || strings.ContainsAny(title, "\r\n") {
		return nil, errors.New("a work item needs a one-line title")
	}
	if risk != "" && !slices.Contains(model.Risks, model.Risk(risk)) {
		return nil, fmt.Errorf("unknown risk %q (low, normal, high)", risk)
	}
	if size != "" && !slices.Contains(model.Sizes, size) {
		return nil, fmt.Errorf("unknown size %q (small, medium, large)", size)
	}
	if parent != "" {
		p, err := w.Load(parent)
		if err != nil {
			return nil, fmt.Errorf("parent: %w", err)
		}
		if p.Kind() != model.Initiative {
			return nil, fmt.Errorf("parent %s is a %s; only an initiative has child items", parent, p.Kind())
		}
	}
	if err := os.MkdirAll(filepath.Join(w.Root, "work"), 0o755); err != nil {
		return nil, err
	}
	var id, dir string
	for attempt := 0; ; attempt++ {
		var err error
		if id, err = w.NewID("W"); err != nil {
			return nil, err
		}
		dir = w.itemDir(id)
		err = os.Mkdir(dir, 0o755)
		if err == nil {
			break
		}
		if !errors.Is(err, os.ErrExist) || attempt == 4 {
			return nil, err
		}
	}
	stage := model.First(kind)
	data := map[string]string{"ID": id, "Kind": string(kind), "Title": title, "Stage": string(stage),
		"Risk": risk, "Size": size, "Parent": parent, "Created": w.stamp()}
	text, err := templates.Render("brief-"+string(kind)+".md", data)
	if err != nil {
		return nil, err
	}
	d, err := record.Parse(text)
	if err != nil {
		return nil, err
	}
	d.Set("title", title) // re-set so titles with colons are quoted
	if err := writeNew(filepath.Join(dir, string(model.BriefDoc)), d.String()); err != nil {
		return nil, err
	}
	logText, err := templates.Render("log.md", data)
	if err != nil {
		return nil, err
	}
	if err := writeNew(filepath.Join(dir, "log.md"), logText); err != nil {
		return nil, err
	}
	it, err := w.Load(id)
	if err != nil {
		return nil, err
	}
	return it, w.appendEntry(it, "stage", by, fmt.Sprintf("Created as a %s in %s.", kind, stage))
}

// Load reads one work item.
func (w *Workspace) Load(id string) (*Item, error) {
	if strings.ContainsAny(id, `/\`) || id == "" || id == "." || id == ".." {
		return nil, fmt.Errorf("bad work item ID %q", id)
	}
	dir := w.itemDir(id)
	text, err := os.ReadFile(filepath.Join(dir, string(model.BriefDoc)))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("no work item %s", id)
		}
		return nil, err
	}
	brief, err := record.Parse(string(text))
	if err != nil {
		return nil, fmt.Errorf("%s/brief.md: %w", id, err)
	}
	it := &Item{ID: id, Dir: dir, Brief: brief, Docs: map[model.Artifact]*record.Doc{model.BriefDoc: brief}}
	for _, a := range []model.Artifact{model.DesignDoc, model.PlanDoc, model.ReleaseDoc} {
		text, err := os.ReadFile(filepath.Join(dir, string(a)))
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, err
		}
		d, err := record.Parse(string(text))
		if err != nil {
			it.Problems = append(it.Problems, fmt.Sprintf("%s: %v", a, err))
			continue
		}
		it.Docs[a] = d
	}
	paths, _ := filepath.Glob(filepath.Join(dir, "evidence", "*.md"))
	sort.Strings(paths)
	for _, p := range paths {
		text, err := os.ReadFile(p)
		if err != nil {
			return nil, err
		}
		d, err := record.Parse(string(text))
		if err != nil {
			it.Problems = append(it.Problems, fmt.Sprintf("evidence/%s: %v", filepath.Base(p), err))
			continue
		}
		it.Evidence = append(it.Evidence, d)
	}
	if text, err := os.ReadFile(filepath.Join(dir, "log.md")); err == nil {
		entries, err := record.ParseLog(string(text))
		if err != nil {
			it.Problems = append(it.Problems, fmt.Sprintf("log.md: %v", err))
		}
		it.Log = entries
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	if it.Kind() == model.Initiative {
		all, err := w.List()
		if err != nil {
			return nil, err
		}
		for _, c := range all {
			if c.Brief.Get("parent") == id {
				it.Children = append(it.Children, c)
			}
		}
	}
	return it, nil
}

// List reads every work item, oldest first. Items whose brief can't be read are
// returned with a problem instead of failing the whole list.
func (w *Workspace) List() ([]*Item, error) {
	entries, err := os.ReadDir(filepath.Join(w.Root, "work"))
	if err != nil {
		return nil, err
	}
	var out []*Item
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		text, err := os.ReadFile(filepath.Join(w.itemDir(e.Name()), string(model.BriefDoc)))
		if err != nil {
			continue
		}
		d, err := record.Parse(string(text))
		it := &Item{ID: e.Name(), Dir: w.itemDir(e.Name()), Brief: d}
		if err != nil {
			it.Brief, _ = record.Parse("")
			it.Problems = []string{fmt.Sprintf("brief.md: %v", err)}
		}
		out = append(out, it)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

func (w *Workspace) appendEntry(it *Item, typ string, by Actor, text string) error {
	if !slices.Contains(model.EntryTypes, typ) {
		return fmt.Errorf("unknown log entry type %q (%s)", typ, strings.Join(model.EntryTypes, ", "))
	}
	if strings.TrimSpace(text) == "" {
		return errors.New("a log entry needs text")
	}
	f, err := os.OpenFile(filepath.Join(it.Dir, "log.md"), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	e := record.Entry{At: w.stamp(), Type: typ, By: by.name(), Agent: by.Agent, Model: by.Model, Session: by.Session, Text: text}
	if _, err := f.WriteString(record.FormatEntry(e)); err != nil {
		f.Close()
		return err
	}
	it.Log = append(it.Log, e)
	return f.Close()
}

// Log appends an entry to a work item's log. Stage and checkpoint entries are
// written only by Advance.
func (w *Workspace) Log(id, typ string, by Actor, text string) error {
	switch typ {
	case "stage", "checkpoint":
		return fmt.Errorf("%s entries are written by `ignition advance`", typ)
	case "evidence":
		return fmt.Errorf("evidence entries are written by `ignition evidence`")
	}
	it, err := w.Load(id)
	if err != nil {
		return err
	}
	return w.appendEntry(it, typ, by, text)
}

// EvidenceInput is what an evidence record states.
type EvidenceInput struct {
	Claim, Result, Revision, Environment, Method, Observed, Limits string
	By                                                             Actor
}

// AddEvidence records evidence for a work item and returns its ID.
func (w *Workspace) AddEvidence(id string, in EvidenceInput) (string, error) {
	it, err := w.Load(id)
	if err != nil {
		return "", err
	}
	if !slices.Contains(model.Results, in.Result) {
		return "", fmt.Errorf("result must be one of %s", strings.Join(model.Results, ", "))
	}
	for name, v := range map[string]string{"claim": in.Claim, "revision": in.Revision, "environment": in.Environment, "method": in.Method} {
		if strings.TrimSpace(v) == "" {
			return "", fmt.Errorf("evidence needs a %s", name)
		}
	}
	eid, err := w.NewID("E")
	if err != nil {
		return "", err
	}
	text, err := templates.Render("evidence.md", map[string]string{"ID": eid, "Work": id, "Claim": "x", "Result": in.Result,
		"Revision": "x", "Environment": "x", "Method": "x", "By": "x", "At": w.stamp(), "Observed": in.Observed, "Limits": in.Limits})
	if err != nil {
		return "", err
	}
	d, err := record.Parse(text)
	if err != nil {
		return "", err
	}
	for k, v := range map[string]string{"claim": in.Claim, "revision": in.Revision, "environment": in.Environment, "method": in.Method,
		"by": in.By.name(), "agent": in.By.Agent, "model": in.By.Model, "session": in.By.Session} {
		d.Set(k, v)
	}
	d.Body = strings.Replace(d.Body, "# Evidence: x", "# Evidence: "+in.Claim, 1)
	if err := os.MkdirAll(filepath.Join(it.Dir, "evidence"), 0o755); err != nil {
		return "", err
	}
	if err := writeNew(filepath.Join(it.Dir, "evidence", eid+".md"), d.String()); err != nil {
		return "", err
	}
	return eid, w.appendEntry(it, "evidence", in.By, fmt.Sprintf("%s (%s): %s", eid, in.Result, in.Claim))
}

// NewDecision creates a product-wide decision record and returns its path.
func (w *Workspace) NewDecision(title string) (string, error) {
	if strings.TrimSpace(title) == "" {
		return "", errors.New("a decision needs a title")
	}
	id, err := w.NewID("D")
	if err != nil {
		return "", err
	}
	text, err := templates.Render("decision.md", map[string]string{"ID": id, "Title": "x", "Created": w.stamp()})
	if err != nil {
		return "", err
	}
	d, err := record.Parse(text)
	if err != nil {
		return "", err
	}
	d.Set("title", title)
	d.Body = strings.Replace(d.Body, "# Decision: x", "# Decision: "+title, 1)
	if err := os.MkdirAll(filepath.Join(w.Root, "decisions"), 0o755); err != nil {
		return "", err
	}
	path := filepath.Join(w.Root, "decisions", id+".md")
	return path, writeNew(path, d.String())
}
