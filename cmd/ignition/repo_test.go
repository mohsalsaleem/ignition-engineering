package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/mohsalsaleem/ignition-engineering/internal/record"
)

func git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=T", "GIT_AUTHOR_EMAIL=t@example.invalid",
		"GIT_COMMITTER_NAME=T", "GIT_COMMITTER_EMAIL=t@example.invalid")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func ign(t *testing.T, args ...string) string {
	t.Helper()
	var out bytes.Buffer
	if code := run(append(args, "--by", "agent"), &out, &out); code != 0 {
		t.Fatalf("ignition %v: %d\n%s", args, code, out.String())
	}
	return out.String()
}

func newRepo(t *testing.T) string {
	t.Helper()
	repo, _ := filepath.EvalSymlinks(t.TempDir())
	git(t, repo, "init", "-q", "-b", "main")
	return repo
}

func TestInitTracksTheRecordsInTheWorkRepository(t *testing.T) {
	fakeEnv(t, nil)
	repo := newRepo(t)
	os.WriteFile(filepath.Join(repo, "Dockerfile"), []byte("FROM scratch\nCOPY . .\n"), 0o644)
	os.WriteFile(filepath.Join(repo, ".dockerignore"), []byte(".env"), 0o644)
	sub := filepath.Join(repo, "src", "app")
	os.MkdirAll(sub, 0o755)
	t.Chdir(sub)

	out := ign(t, "init")
	if !strings.Contains(out, filepath.Join(repo, ".ignition")) || !strings.Contains(out, "commit .ignition/") {
		t.Fatalf("init printed:\n%s", out)
	}
	if b, _ := os.ReadFile(filepath.Join(repo, ".dockerignore")); string(b) != ".env\n.ignition\n" {
		t.Fatalf(".dockerignore = %q", b)
	}
	if p, _ := os.ReadFile(filepath.Join(repo, ".ignition", "PROJECT.md")); !strings.Contains(string(p), "# Project: "+filepath.Base(repo)) {
		t.Fatalf("PROJECT.md names the wrong product:\n%s", p)
	}
	if s := git(t, repo, "status", "--porcelain", "--untracked-files=all"); !strings.Contains(s, "?? .ignition/PROJECT.md") || !strings.Contains(s, "?? .ignition/.gitattributes") {
		t.Fatalf("records aren't visible to Git:\n%s", s)
	}
	if out := ign(t, "init", "--dir", filepath.Join(t.TempDir(), "x")); strings.Contains(out, ".dockerignore") {
		t.Fatalf("init --dir touched the repository:\n%s", out)
	}
}

func TestInitWarnsWhenGitIgnoresTheRecords(t *testing.T) {
	fakeEnv(t, nil)
	repo := newRepo(t)
	os.WriteFile(filepath.Join(repo, ".gitignore"), []byte(".ignition/\n"), 0o644)
	t.Chdir(repo)
	if out := ign(t, "init"); !strings.Contains(out, "Git ignores .ignition here") {
		t.Fatalf("init printed:\n%s", out)
	}
}

func TestInitInsideASeparateRecordsCheckout(t *testing.T) {
	fakeEnv(t, nil)
	records := filepath.Join(newRepo(t), ".ignition")
	os.Mkdir(records, 0o755)
	git(t, records, "init", "-q", "-b", "main")
	t.Chdir(records)
	ign(t, "init")
	if _, err := os.Stat(filepath.Join(records, "PROJECT.md")); err != nil {
		t.Fatalf("init in a records checkout didn't use it: %v", err)
	}
}

// Two branches that both log on one item must merge without a conflict and
// keep every entry.
func TestBranchesThatLogTheSameItemMergeCleanly(t *testing.T) {
	fakeEnv(t, nil)
	repo := newRepo(t)
	t.Chdir(repo)
	ign(t, "init")
	id := regexp.MustCompile(`W-\d{8}-[0-9a-z]{6}`).FindString(ign(t, "new", "feature", "Pipeline view", "--risk", "normal", "--size", "medium"))
	git(t, repo, "add", "-A")
	git(t, repo, "commit", "-q", "-m", "Start the pipeline view")

	git(t, repo, "switch", "-q", "-c", "a")
	ign(t, "log", id, "--type", "finding", "Stages are stored per job.")
	git(t, repo, "commit", "-q", "-am", "a")
	git(t, repo, "switch", "-q", "main")
	git(t, repo, "switch", "-q", "-c", "b")
	ign(t, "log", id, "--type", "decision", "Use a board layout.")
	git(t, repo, "commit", "-q", "-am", "b")
	git(t, repo, "switch", "-q", "main")
	git(t, repo, "merge", "-q", "a")
	git(t, repo, "merge", "-q", "--no-edit", "b")

	b, _ := os.ReadFile(filepath.Join(repo, ".ignition", "work", id, "log.md"))
	entries, err := record.ParseLog(string(b))
	if err != nil || strings.Contains(string(b), "<<<<<<<") {
		t.Fatalf("log after merge (%v):\n%s", err, b)
	}
	var texts []string
	for _, e := range entries {
		texts = append(texts, e.Text)
	}
	all := strings.Join(texts, "\n")
	if !strings.Contains(all, "Stages are stored per job.") || !strings.Contains(all, "Use a board layout.") {
		t.Fatalf("an entry was lost:\n%s", b)
	}
	ign(t, "check")
}
