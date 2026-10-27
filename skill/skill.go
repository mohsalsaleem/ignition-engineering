// Package skill holds the agent skill that teaches agents to work with
// Ignition, embedded so the installed skill always matches the binary.
package skill

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

//go:embed ignition
var files embed.FS

const root = "ignition"

// Files returns the skill's files by path relative to the skill folder.
func Files() (map[string][]byte, error) {
	out := map[string][]byte{}
	err := fs.WalkDir(files, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		b, err := files.ReadFile(path)
		if err != nil {
			return err
		}
		out[strings.TrimPrefix(path, root+"/")] = b
		return nil
	})
	return out, err
}

// Body returns SKILL.md without its frontmatter.
func Body() string {
	b, _ := files.ReadFile(root + "/SKILL.md")
	s := string(b)
	if rest, ok := strings.CutPrefix(s, "---\n"); ok {
		if _, body, ok := strings.Cut(rest, "\n---\n"); ok {
			return strings.TrimLeft(body, "\n")
		}
	}
	return s
}

// Name is the skill's name, which is also its folder name.
const Name = "ignition"

// Targets returns the folders to install the skill into for the given agents:
// "claude" (Claude Code reads .claude/skills), "agents" (the shared
// .agents/skills folder read by Codex, Cursor, Gemini CLI, GitHub Copilot and
// others), or "all". base is the home folder, or a project for a project-only
// install.
func Targets(agents, base string) ([]string, error) {
	claude := filepath.Join(base, ".claude", "skills", Name)
	shared := filepath.Join(base, ".agents", "skills", Name)
	switch agents {
	case "all", "":
		return []string{claude, shared}, nil
	case "claude":
		return []string{claude}, nil
	case "agents":
		return []string{shared}, nil
	}
	return nil, fmt.Errorf("unknown target %q (all, claude, agents)", agents)
}

// Install writes the skill into dir, replacing an earlier Ignition skill there.
// It refuses a folder that holds some other skill.
func Install(dir string) error {
	existing, err := os.ReadFile(filepath.Join(dir, "SKILL.md"))
	switch {
	case errors.Is(err, os.ErrNotExist):
	case err != nil:
		return err
	case !bytes.Contains(existing, []byte("\nname: ignition\n")):
		return fmt.Errorf("%s holds a different skill; choose another folder", dir)
	}
	fsFiles, err := Files()
	if err != nil {
		return err
	}
	for rel, b := range fsFiles {
		path := filepath.Join(dir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, b, 0o644); err != nil {
			return err
		}
	}
	return nil
}
