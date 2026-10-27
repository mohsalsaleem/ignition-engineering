package skill

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/mohsalsaleem/ignition-engineering/internal/model"
	"github.com/mohsalsaleem/ignition-engineering/internal/record"
)

func TestFrontmatter(t *testing.T) {
	files, err := Files()
	if err != nil {
		t.Fatal(err)
	}
	d, err := record.Parse(string(files["SKILL.md"]))
	if err != nil {
		t.Fatal(err)
	}
	if d.Get("name") != "ignition" {
		t.Errorf("name = %q", d.Get("name"))
	}
	if desc := d.Get("description"); desc == "" || len(desc) > 1024 {
		t.Errorf("description is %d characters, want 1 to 1024", len(desc))
	}
	// The Agent Skills specification (agentskills.io/specification).
	if !regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`).MatchString(d.Get("name")) || len(d.Get("name")) > 64 || d.Get("name") != Name {
		t.Errorf("name %q breaks the spec or doesn't match the folder", d.Get("name"))
	}
	if c := d.Get("compatibility"); len(c) > 500 {
		t.Errorf("compatibility is %d characters, want at most 500", len(c))
	}
	// Agents parse frontmatter as YAML; unquoted ": " or " #" would break it.
	for _, k := range []string{"name", "description", "compatibility"} {
		if v := d.Get(k); strings.Contains(v, ": ") || strings.Contains(v, " #") {
			t.Errorf("%s has YAML-breaking text: %q", k, v)
		}
	}
	if strings.HasPrefix(Body(), "---") || !strings.HasPrefix(Body(), "# Ignition") {
		t.Errorf("Body starts %q", Body()[:20])
	}
}

func TestLinksResolve(t *testing.T) {
	files, _ := Files()
	link := regexp.MustCompile(`\]\(([^)#]+)\)`)
	for name, b := range files {
		for _, m := range link.FindAllStringSubmatch(string(b), -1) {
			target := filepath.ToSlash(filepath.Join(filepath.Dir(name), m[1]))
			if strings.HasPrefix(m[1], "http") {
				continue
			}
			if _, ok := files[target]; !ok {
				t.Errorf("%s links to %s, which isn't in the skill", name, m[1])
			}
		}
	}
}

// Every stage an item can be in must be explained in references/stages.md.
func TestEveryStageIsExplained(t *testing.T) {
	files, _ := Files()
	var headings []string
	for _, line := range strings.Split(string(files["references/stages.md"]), "\n") {
		if h, ok := strings.CutPrefix(line, "## "); ok {
			headings = append(headings, h)
		}
	}
	for _, k := range model.Kinds {
		for _, s := range model.Route(model.Facts{Kind: k}) {
			if s == model.Closed {
				continue
			}
			found := false
			for _, h := range headings {
				if regexp.MustCompile(`\b` + string(s) + `\b`).MatchString(h) {
					found = true
				}
			}
			if !found {
				t.Errorf("stages.md has no heading for %s", s)
			}
		}
	}
}

func TestInstall(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "ignition")
	if err := Install(dir); err != nil {
		t.Fatal(err)
	}
	if err := Install(dir); err != nil {
		t.Fatalf("reinstalling over the Ignition skill: %v", err)
	}
	for _, f := range []string{"SKILL.md", "references/stages.md", "references/commands.md"} {
		if _, err := os.Stat(filepath.Join(dir, f)); err != nil {
			t.Error(err)
		}
	}
	other := t.TempDir()
	os.WriteFile(filepath.Join(other, "SKILL.md"), []byte("---\nname: other\n---\n"), 0o644)
	if err := Install(other); err == nil {
		t.Fatal("installed over another skill")
	}
}

func TestTargets(t *testing.T) {
	dirs, err := Targets("all", "/home/u")
	if err != nil || len(dirs) != 2 || dirs[0] != "/home/u/.claude/skills/ignition" || dirs[1] != "/home/u/.agents/skills/ignition" {
		t.Fatalf("Targets = %v, %v", dirs, err)
	}
	if _, err := Targets("nope", "/home/u"); err == nil {
		t.Fatal("accepted an unknown target")
	}
}
