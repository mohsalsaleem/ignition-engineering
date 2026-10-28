package templates

import (
	"strings"
	"testing"

	"github.com/mohsalsaleem/ignition-engineering/internal/model"
	"github.com/mohsalsaleem/ignition-engineering/internal/record"
)

// Every section an exit rule requires must exist in the template that creates
// the artifact, so a new item can always be completed.
func TestTemplatesHaveEveryRequiredSection(t *testing.T) {
	var facts []model.Facts
	for _, k := range model.Kinds {
		facts = append(facts, model.Facts{Kind: k}, model.Facts{Kind: k, Child: true})
	}
	for _, f := range facts {
		data := map[string]string{"ID": "W-1", "Kind": "k", "Title": "T", "Stage": "s"}
		if f.Child {
			data["Parent"] = "W-0"
		}
		for _, s := range model.Route(f) {
			for _, r := range model.Exit(f, s) {
				if r.Kind != model.Sections {
					continue
				}
				name := Name(r.Artifact, f.Kind)
				text, err := Render(name, data)
				if err != nil {
					t.Fatalf("%s: %v", name, err)
				}
				d, err := record.Parse(text)
				if err != nil {
					t.Fatalf("%s: %v", name, err)
				}
				for _, sec := range r.Names {
					if _, ok := d.Section(sec); !ok {
						t.Errorf("%s (%+v at %s) lacks section %q", name, f, s, sec)
					}
					if d.Filled(sec) {
						t.Errorf("%s section %q is filled in the template", name, sec)
					}
				}
			}
		}
	}
}

func TestParentDesignOnlyForChildItems(t *testing.T) {
	text, _ := Render("design.md", map[string]string{"ID": "W-1", "Title": "T"})
	if strings.Contains(text, "Parent design") {
		t.Error("a design without a parent has a Parent design section")
	}
	text, _ = Render("design.md", map[string]string{"ID": "W-1", "Title": "T", "Parent": "W-0"})
	if !strings.Contains(text, "## Parent design") || !strings.Contains(text, "W-0") {
		t.Error("a child's design doesn't name its parent")
	}
}
