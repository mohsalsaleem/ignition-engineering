package workspace

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// gitattributes makes Git merge append-only logs by keeping both sides, so
// two branches that log on the same item merge without a conflict.
const gitattributes = "# Logs are append-only: keep entries from both sides when branches merge.\nlog.md merge=union\n"

// RepoRoot returns the Git work tree containing dir, or "" if there is none.
func RepoRoot(dir string) string {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// PrepareRepo readies an application repository to track its .ignition
// folder: records must not be ignored by Git, and must not be copied into
// Docker images. It returns notes on what it changed or couldn't.
func PrepareRepo(app string) ([]string, error) {
	var notes []string
	check := exec.Command("git", "check-ignore", "-q", ".ignition/PROJECT.md")
	check.Dir = app
	if check.Run() == nil {
		notes = append(notes, "Git ignores .ignition here (see .gitignore or .git/info/exclude), so the records won't be committed. Remove that rule to keep them in the repository.")
	}
	dockerfiles, _ := filepath.Glob(filepath.Join(app, "*Dockerfile*"))
	if len(dockerfiles) == 0 {
		return notes, nil
	}
	path := filepath.Join(app, ".dockerignore")
	b, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return notes, err
	}
	for _, l := range strings.Split(string(b), "\n") {
		if l = strings.TrimSpace(l); l == ".ignition" || l == ".ignition/" || l == "/.ignition" {
			return notes, nil
		}
	}
	add := ".ignition\n"
	if len(b) > 0 && !bytes.HasSuffix(b, []byte("\n")) {
		add = "\n" + add
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o644)
	if err != nil {
		return notes, err
	}
	if _, err := f.WriteString(add); err != nil {
		f.Close()
		return notes, err
	}
	if err := f.Close(); err != nil {
		return notes, err
	}
	return append(notes, "Added .ignition to .dockerignore, so Docker images don't include the records."), nil
}
