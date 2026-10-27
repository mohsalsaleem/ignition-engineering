// Package attach checks a product's separate records repository out at
// <app>/.ignition/ without changing the application repository's Git status.
package attach

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func git(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	var out, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s: %s", args[0], strings.TrimSpace(errb.String()))
	}
	return strings.TrimSpace(out.String()), nil
}

// Redact removes credentials from a URL so it can be shown.
func Redact(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.User == nil {
		return raw
	}
	u.User = url.User("***")
	return u.String()
}

func sameRemote(a, b string) bool {
	norm := func(s string) string {
		s = strings.TrimSuffix(strings.TrimSuffix(s, "/"), ".git")
		if u, err := url.Parse(s); err == nil && u.User != nil {
			u.User = nil
			s = u.String()
		}
		if p, err := filepath.EvalSymlinks(s); err == nil {
			s = p
		}
		return s
	}
	return norm(a) == norm(b)
}

// Result says what Attach did.
type Result struct {
	Path     string
	Outcome  string // "attached" or "already attached"
	Warnings []string
}

// Attach clones recordsURL to <app>/.ignition and excludes it from app's Git.
// An existing .ignition must be a separate checkout of the same repository.
// If anything fails after cloning, the clone is removed.
func Attach(app, recordsURL string) (res Result, err error) {
	app, err = filepath.Abs(app)
	if err != nil {
		return res, err
	}
	if app, err = filepath.EvalSymlinks(app); err != nil {
		return res, fmt.Errorf("application directory: %w", err)
	}
	top, err := git(app, "rev-parse", "--show-toplevel")
	if err != nil {
		return res, err
	}
	if top, _ = filepath.EvalSymlinks(top); top != app {
		return res, fmt.Errorf("--app must be the application repository root: %s", top)
	}
	target := filepath.Join(app, ".ignition")
	res.Path = target
	if tracked, err := git(app, "ls-files", "--", ".ignition"); err != nil {
		return res, err
	} else if tracked != "" {
		return res, errors.New("the application repository already tracks .ignition; stop and review it")
	}

	if _, statErr := os.Stat(target); statErr == nil {
		nested, err := git(target, "rev-parse", "--show-toplevel")
		if nested, _ = filepath.EvalSymlinks(nested); err != nil || nested != target {
			return res, errors.New(".ignition exists but is not a separate Git checkout")
		}
		origin, err := git(target, "remote", "get-url", "origin")
		if err != nil || !sameRemote(origin, recordsURL) {
			return res, fmt.Errorf(".ignition is a checkout of %s, not %s", Redact(origin), Redact(recordsURL))
		}
		res.Outcome = "already attached"
	} else {
		if _, err := git(app, "clone", "--quiet", "--", recordsURL, target); err != nil {
			return res, fmt.Errorf("cloning %s failed: %s", Redact(recordsURL), strings.ReplaceAll(err.Error(), recordsURL, Redact(recordsURL)))
		}
		res.Outcome = "attached"
		defer func() {
			if err != nil {
				os.RemoveAll(target)
			}
		}()
	}

	// --git-path resolves info/exclude to the file Git actually reads, which in
	// a linked worktree is in the main repository's Git directory.
	exclude, err := git(app, "rev-parse", "--git-path", "info/exclude")
	if err != nil {
		return res, err
	}
	if !filepath.IsAbs(exclude) {
		exclude = filepath.Join(app, exclude)
	}
	if err := addLine(exclude, ".ignition/"); err != nil {
		return res, err
	}
	if status, err := git(app, "status", "--porcelain", "--untracked-files=all", "--", ".ignition"); err != nil {
		return res, err
	} else if status != "" {
		return res, errors.New(".ignition still appears in the application's Git status")
	}
	res.Warnings = buildContextWarnings(app)
	return res, nil
}

func addLine(path, line string) error {
	b, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	for _, l := range strings.Split(string(b), "\n") {
		if strings.TrimSpace(l) == line {
			return nil
		}
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	if len(b) > 0 && !bytes.HasSuffix(b, []byte("\n")) {
		line = "\n" + line
	}
	if _, err := f.WriteString(line + "\n"); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

// buildContextWarnings flags Docker builds whose context would include the
// private records: Git's exclude file doesn't apply to `docker build`.
func buildContextWarnings(app string) []string {
	matches, _ := filepath.Glob(filepath.Join(app, "*Dockerfile*"))
	if len(matches) == 0 {
		return nil
	}
	b, _ := os.ReadFile(filepath.Join(app, ".dockerignore"))
	for _, l := range strings.Split(string(b), "\n") {
		if l = strings.TrimSpace(l); l == ".ignition" || l == ".ignition/" || l == "/.ignition" {
			return nil
		}
	}
	return []string{"This repository has a Dockerfile, and a local `docker build .` would copy the private records into the image. Add `.ignition` to .dockerignore."}
}
