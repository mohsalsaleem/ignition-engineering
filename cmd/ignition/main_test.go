package main

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func ignition(t *testing.T, dir string, args ...string) (string, int) {
	t.Helper()
	var out, errw bytes.Buffer
	code := run(append(args, "--dir", dir, "--by", "agent"), &out, &errw)
	return out.String() + errw.String(), code
}

func TestTaskThroughTheCLI(t *testing.T) {
	dir := t.TempDir()
	if out, code := ignition(t, dir, "init", "--product", "example"); code != 0 {
		t.Fatal(out)
	}
	out, code := ignition(t, dir, "new", "task", "Rotate logs", "--risk", "low", "--size", "small")
	if code != 0 {
		t.Fatal(out)
	}
	id := regexp.MustCompile(`W-\d{8}-[0-9a-z]{6}`).FindString(out)
	if id == "" || !strings.Contains(out, `brief.md: fill in "Request"`) {
		t.Fatalf("new printed:\n%s", out)
	}
	brief := filepath.Join(dir, "work", id, "brief.md")
	b, _ := os.ReadFile(brief)
	b = bytes.Replace(b, []byte("## Request\n"), []byte("## Request\n\nLogs fill the disk.\n"), 1)
	b = bytes.Replace(b, []byte("## Done when\n"), []byte("## Done when\n\nLogs rotate daily.\n"), 1)
	os.WriteFile(brief, b, 0o644)

	if out, code := ignition(t, dir, "advance", id); code != 0 || !strings.Contains(out, "now in plan") {
		t.Fatalf("advance: %d\n%s", code, out)
	}
	if out, code := ignition(t, dir, "advance", id); code != 1 || !strings.Contains(out, "list at least one task") {
		t.Fatalf("advance without tasks: %d\n%s", code, out)
	}
	if out, code := ignition(t, dir, "log", id, "--type", "checkpoint", "Owner said yes"); code != 1 {
		t.Fatalf("log checkpoint: %d\n%s", code, out)
	}
	if out, code := ignition(t, dir, "evidence", id, "--claim", "Logs rotate", "--result", "maybe",
		"--revision", "abc", "--environment", "local", "--method", "ls"); code != 1 || !strings.Contains(out, "result must be") {
		t.Fatalf("bad evidence: %d\n%s", code, out)
	}
	if out, code := ignition(t, dir, "check"); code != 0 || !strings.Contains(out, "no problems") {
		t.Fatalf("check: %d\n%s", code, out)
	}
	if out, _ := ignition(t, dir, "status"); !strings.Contains(out, id) || !strings.Contains(out, "in progress") {
		t.Fatalf("status:\n%s", out)
	}
}

func TestGuideAndHelp(t *testing.T) {
	var out bytes.Buffer
	if run([]string{"guide"}, &out, &out) != 0 || !strings.Contains(out.String(), "ignition advance") {
		t.Fatalf("guide:\n%s", out.String())
	}
	out.Reset()
	if run([]string{"nope"}, &out, &out) != 2 {
		t.Fatal("unknown command didn't exit 2")
	}
}
