// Package templates holds the record templates Ignition creates records from.
package templates

import (
	"embed"
	"strings"
	"text/template"
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
