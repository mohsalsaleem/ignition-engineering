package attach

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func run(t *testing.T, dir string, args ...string) string {
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

func repo(t *testing.T, dir string) string {
	t.Helper()
	os.MkdirAll(dir, 0o755)
	run(t, dir, "init", "-q", "-b", "main")
	run(t, dir, "commit", "-q", "--allow-empty", "-m", "init")
	return dir
}

func setup(t *testing.T) (root, app, records string) {
	root, _ = filepath.EvalSymlinks(t.TempDir())
	return root, repo(t, filepath.Join(root, "app")), repo(t, filepath.Join(root, "records"))
}

func TestRecordsAreSeparateAndInvisibleToAppGit(t *testing.T) {
	_, app, records := setup(t)
	res, err := Attach(app, records)
	if err != nil || res.Outcome != "attached" {
		t.Fatalf("Attach = %+v, %v", res, err)
	}
	res, err = Attach(app, records)
	if err != nil || res.Outcome != "already attached" {
		t.Fatalf("second Attach = %+v, %v", res, err)
	}
	os.WriteFile(filepath.Join(app, ".ignition", "PROJECT.md"), []byte("# Changed\n"), 0o644)
	if s := run(t, app, "status", "--porcelain"); s != "" {
		t.Fatalf("app status = %q", s)
	}
	b, _ := os.ReadFile(filepath.Join(app, ".git", "info", "exclude"))
	if strings.Count(string(b), ".ignition/") != 1 {
		t.Fatalf("exclude has %d entries", strings.Count(string(b), ".ignition/"))
	}
}

// A linked worktree's own Git directory has no exclude file Git reads; the
// entry must go to the main repository's.
func TestWorksInALinkedWorktree(t *testing.T) {
	root, app, records := setup(t)
	wt := filepath.Join(root, "wt")
	run(t, app, "worktree", "add", "-q", wt)
	if _, err := Attach(wt, records); err != nil {
		t.Fatal(err)
	}
	if s := run(t, wt, "status", "--porcelain"); s != "" {
		t.Fatalf("worktree status = %q", s)
	}
}

func TestRefusesACheckoutOfAnotherRepository(t *testing.T) {
	root, app, records := setup(t)
	other := repo(t, filepath.Join(root, "other"))
	if _, err := Attach(app, other); err != nil {
		t.Fatal(err)
	}
	if _, err := Attach(app, records); err == nil || !strings.Contains(err.Error(), "is a checkout of") {
		t.Fatalf("Attach = %v", err)
	}
}

func TestRefusesAnExistingNonGitFolder(t *testing.T) {
	_, app, records := setup(t)
	os.Mkdir(filepath.Join(app, ".ignition"), 0o755)
	if _, err := Attach(app, records); err == nil || !strings.Contains(err.Error(), "separate Git checkout") {
		t.Fatalf("Attach = %v", err)
	}
}

func TestRemovesTheCloneWhenAttachFails(t *testing.T) {
	_, app, records := setup(t)
	// A tracked .gitignore negation keeps .ignition visible despite the exclude.
	os.WriteFile(filepath.Join(app, ".gitignore"), []byte("!.ignition/\n!.ignition/**\n"), 0o644)
	os.MkdirAll(filepath.Join(app, ".git", "info"), 0o755)
	// Make the exclude unwritable so the attach fails after cloning.
	exclude := filepath.Join(app, ".git", "info", "exclude")
	os.WriteFile(exclude, nil, 0o444)
	os.Chmod(filepath.Join(app, ".git", "info"), 0o555)
	t.Cleanup(func() { os.Chmod(filepath.Join(app, ".git", "info"), 0o755) })
	if _, err := Attach(app, records); err == nil {
		t.Fatal("Attach succeeded with an unwritable exclude file")
	}
	if _, err := os.Stat(filepath.Join(app, ".ignition")); !os.IsNotExist(err) {
		t.Fatalf("the clone was left behind: %v", err)
	}
}

func TestCloneFailureDoesNotPrintCredentials(t *testing.T) {
	_, app, _ := setup(t)
	_, err := Attach(app, "https://user:s3cret@127.0.0.1:1/records.git")
	if err == nil || strings.Contains(err.Error(), "s3cret") {
		t.Fatalf("Attach = %v", err)
	}
}

func TestWarnsWhenDockerWouldCopyTheRecords(t *testing.T) {
	_, app, records := setup(t)
	os.WriteFile(filepath.Join(app, "Dockerfile"), []byte("FROM scratch\nCOPY . .\n"), 0o644)
	res, err := Attach(app, records)
	if err != nil || len(res.Warnings) != 1 {
		t.Fatalf("Attach = %+v, %v", res, err)
	}
	os.WriteFile(filepath.Join(app, ".dockerignore"), []byte(".ignition\n"), 0o644)
	if res, _ := Attach(app, records); len(res.Warnings) != 0 {
		t.Fatalf("warned with .dockerignore set: %v", res.Warnings)
	}
}
