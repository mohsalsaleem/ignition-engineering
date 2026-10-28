// Package templates holds the record templates Ignition creates records from.
package templates

import (
	"embed"
	"strings"
	"text/template"

	"github.com/mohsalsaleem/ignition-engineering/internal/model"
)

//go:embed files/*.md
var files embed.FS

// Render fills the named template.
func Render(name string, data any) (string, error) {
	t, err := template.ParseFS(files, "files/"+name)
	if err != nil {
		return "", err
	}
	var b strings.Builder
	if err := t.Execute(&b, data); err != nil {
		return "", err
	}
	return b.String(), nil
}

// Name is the template an artifact of an item of kind k is created from. Each
// kind has its own brief, and an initiative's design is a solution design.
func Name(a model.Artifact, k model.Kind) string {
	switch {
	case a == model.BriefDoc:
		return "brief-" + string(k) + ".md"
	case a == model.DesignDoc && k == model.Initiative:
		return "design-initiative.md"
	}
	return string(a)
}
