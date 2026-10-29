package workspace

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mohsalsaleem/ignition-engineering/internal/model"
	"github.com/mohsalsaleem/ignition-engineering/templates"
)

// The product's architecture docs: the system as it is now, which designs
// read first and update with the change they make.
var architectureDocs = map[string]string{
	"architecture/overview.md":     "architecture-overview.md",
	"architecture/data-model.md":   "architecture-data-model.md",
	"architecture/flows/README.md": "architecture-flows.md",
}

// DataModelDoc is the agreed data model, relative to the workspace.
const DataModelDoc = "architecture/data-model.md"

// EnsureArchitecture creates whichever architecture docs a workspace lacks.
// Workspaces made before they existed get them when an item enters design.
func EnsureArchitecture(root string) error {
	for path, tmpl := range architectureDocs {
		full := filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			return err
		}
		text, err := templates.Render(tmpl, nil)
		if err != nil {
			return err
		}
		if err := writeNew(full, text); err != nil && !errors.Is(err, os.ErrExist) {
			return err
		}
	}
	return nil
}

// dataModelUnagreed reports a design that changes the data without updating
// the agreed data model. A Data model section that is empty or starts with
// "none" changes nothing; Sections reports an empty one.
func (it *Item) dataModelUnagreed() []string {
	d := it.Docs[model.DesignDoc]
	if d == nil {
		return nil
	}
	s, _ := d.Section("Data model")
	if s == "" || strings.HasPrefix(strings.ToLower(s), "none") {
		return nil
	}
	b, err := os.ReadFile(filepath.Join(it.Root, DataModelDoc))
	if err != nil || !strings.Contains(string(b), it.ID) {
		return []string{fmt.Sprintf("%s: the design changes the data model; update it there to the agreed model, naming %s", DataModelDoc, it.ID)}
	}
	return nil
}
