package templates

import (
	"testing"

	"github.com/mohsalsaleem/ignition-engineering/internal/model"
	"github.com/mohsalsaleem/ignition-engineering/internal/record"
)

// Every section an exit rule requires must exist in the template that creates
// the artifact, so a new item can always be completed.
func TestTemplatesHaveEveryRequiredSection(t *testing.T) {
	data := map[string]string{"ID": "W-1", "Kind": "k", "Title": "T", "Stage": "s"}
	for _, k := range model.Kinds {
		for _, s := range model.Route(model.Facts{Kind: k}) {
			for _, r := range model.Exit(k, s) {
				if r.Kind != model.Sections {
					continue
				}
				name := string(r.Artifact)
				if r.Artifact == model.BriefDoc {
					name = "brief-" + string(k) + ".md"
				}
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
						t.Errorf("%s (%s at %s) lacks section %q", name, k, s, sec)
					}
					if d.Filled(sec) {
						t.Errorf("%s section %q is filled in the template", name, sec)
					}
				}
			}
		}
	}
}
