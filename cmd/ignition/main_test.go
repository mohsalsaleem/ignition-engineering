package main

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/mohsalsaleem/ignition-engineering/skill"
)

// fakeEnv makes the CLI see only these variables. Tests run inside agents,
// whose own variables would otherwise change what is recorded.
func fakeEnv(t *testing.T, env map[string]string) {
	t.Helper()
	old := getenv
	getenv = func(k string) string { return env[k] }
	t.Cleanup(func() { getenv = old })
}

func ignition(t *testing.T, dir string, args ...string) (string, int) {
	t.Helper()
	var out, errw bytes.Buffer
	code := run(append(args, "--dir", dir, "--by", "agent"), &out, &errw)
	return out.String() + errw.String(), code
}

func TestTaskThroughTheCLI(t *testing.T) {
	fakeEnv(t, nil)
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

// The skill must only teach commands and errors that exist.
func TestSkillMatchesTheCLI(t *testing.T) {
	files, err := skill.Files()
	if err != nil {
		t.Fatal(err)
	}
	cmd := regexp.MustCompile("`ignition ([a-z]+)")
	seen := map[string]bool{}
	for name, b := range files {
		for _, m := range cmd.FindAllStringSubmatch(string(b), -1) {
			if seen[m[1]] {
				continue
			}
			seen[m[1]] = true
			var out bytes.Buffer
			if code := run([]string{m[1], "--help"}, &out, &out); code == 2 {
				t.Errorf("%s teaches `ignition %s`, which isn't a command", name, m[1])
			}
		}
	}
	if len(seen) < 10 {
		t.Fatalf("found only %d commands in the skill", len(seen))
	}

	var source strings.Builder
	for _, dir := range []string{".", "../../internal/workspace"} {
		paths, _ := filepath.Glob(filepath.Join(dir, "*.go"))
		for _, p := range paths {
			if !strings.HasSuffix(p, "_test.go") {
				b, _ := os.ReadFile(p)
				source.Write(b)
			}
		}
	}
	for _, msg := range []string{"stage not complete", "fix these first", "owner decision needed", "needs the owner's answer to",
		"fill in", "no evidence with result", "no %s entry", "entries are written by", "no Ignition workspace in",
		"no passing %s on", "has no %s on", "failed on %s", "names no commit", "isn't one of"} {
		if !strings.Contains(source.String(), msg) {
			t.Errorf("commands.md explains %q, but no code produces it", msg)
		}
	}
}

func TestSkillInstallAndGuide(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "skills", "ignition")
	var out bytes.Buffer
	if code := run([]string{"skill", "install", "--to", dir}, &out, &out); code != 0 {
		t.Fatal(out.String())
	}
	if _, err := os.Stat(filepath.Join(dir, "references", "stages.md")); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	if run([]string{"skill", "print"}, &out, &out) != 0 || !strings.HasPrefix(out.String(), "---\nname: ignition\n") {
		t.Fatalf("skill print:\n%.80s", out.String())
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

func TestAttributionFromTheEnvironment(t *testing.T) {
	dir := t.TempDir()
	fakeEnv(t, map[string]string{"AI_AGENT": "claude-code_2-1-304_agent", "CLAUDE_CODE_SESSION_ID": "sess-1"})
	var out, errw bytes.Buffer
	if run([]string{"init", "--dir", dir}, &out, &errw) != 0 {
		t.Fatal(errw.String())
	}
	out.Reset()
	if run([]string{"new", "task", "Rotate logs", "--risk", "low", "--size", "small", "--dir", dir}, &out, &errw) != 0 {
		t.Fatal(errw.String())
	}
	if !strings.Contains(errw.String(), "Warning: recording claude-code with an unknown model") {
		t.Errorf("no unknown-model warning: %q", errw.String())
	}
	id := regexp.MustCompile(`W-\d{8}-[0-9a-z]{6}`).FindString(out.String())
	log, _ := os.ReadFile(filepath.Join(dir, "work", id, "log.md"))
	if !strings.Contains(string(log), " · stage · claude-code · claude-code · - · sess-1\n") {
		t.Fatalf("log:\n%s", log)
	}

	out.Reset()
	errw.Reset()
	run([]string{"log", id, "--type", "finding", "Disk is 90% full", "--model", "claude-opus-5-5", "--dir", dir}, &out, &errw)
	if errw.Len() != 0 {
		t.Errorf("warned with a model given: %q", errw.String())
	}
	out.Reset()
	if run([]string{"attribution", id, "--dir", dir}, &out, &errw) != 0 {
		t.Fatal(errw.String())
	}
	for _, want := range []string{"claude-code · claude-code · unknown", "claude-code · claude-code · claude-opus-5-5", "sessions: sess-1"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("attribution lacks %q:\n%s", want, out.String())
		}
	}
}

func TestWhoamiPrefersExplicitValues(t *testing.T) {
	fakeEnv(t, map[string]string{"AI_AGENT": "claude-code_2-1-304_agent", "IGNITION_MODEL": "m-env", "IGNITION_ACTOR": "Reviewer"})
	var out bytes.Buffer
	run([]string{"whoami", "--agent", "codex", "--model", "gpt-x"}, &out, &out)
	for _, want := range []string{"name:    Reviewer", "agent:   codex", "model:   gpt-x", "session: none"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("whoami lacks %q:\n%s", want, out.String())
		}
	}
	fakeEnv(t, map[string]string{"IGNITION_AGENT": "cursor"})
	out.Reset()
	run([]string{"whoami"}, &out, &out)
	if !strings.Contains(out.String(), "agent:   cursor") || !strings.Contains(out.String(), "model:   unknown") {
		t.Errorf("whoami:\n%s", out.String())
	}
}

func TestSkillInstallsForEveryAgentFamily(t *testing.T) {
	project := t.TempDir()
	var out bytes.Buffer
	if code := run([]string{"skill", "install", "--project", project}, &out, &out); code != 0 {
		t.Fatal(out.String())
	}
	for _, dir := range []string{".claude/skills/ignition", ".agents/skills/ignition"} {
		if _, err := os.Stat(filepath.Join(project, dir, "SKILL.md")); err != nil {
			t.Error(err)
		}
	}
	to := t.TempDir()
	out.Reset()
	if code := run([]string{"skill", "install", "--to", to}, &out, &out); code != 0 {
		t.Fatal(out.String())
	}
	if _, err := os.Stat(filepath.Join(to, "ignition", "SKILL.md")); err != nil {
		t.Errorf("--to a skills folder didn't create an ignition folder: %v", err)
	}
	out.Reset()
	if code := run([]string{"skill", "install", "--for", "nope", "--project", project}, &out, &out); code != 1 {
		t.Errorf("unknown target exit %d", code)
	}
}

func TestQAHandoffThroughTheCLI(t *testing.T) {
	fakeEnv(t, nil)
	dir := t.TempDir()
	ignition(t, dir, "init", "--product", "example")
	out, _ := ignition(t, dir, "new", "task", "Rotate logs", "--risk", "normal", "--size", "small")
	id := regexp.MustCompile(`W-\d{8}-[0-9a-z]{6}`).FindString(out)
	if out, code := ignition(t, dir, "qa", id); code != 1 || !strings.Contains(out, "nothing to test") {
		t.Fatalf("qa before a change: %d\n%s", code, out)
	}
	brief := filepath.Join(dir, "work", id, "brief.md")
	b, _ := os.ReadFile(brief)
	b = bytes.Replace(b, []byte("## Request\n"), []byte("## Request\n\nLogs fill the disk.\n"), 1)
	b = bytes.Replace(b, []byte("## Done when\n"), []byte("## Done when\n\n1. Logs rotate daily.\n2. Old logs are deleted.\n"), 1)
	os.WriteFile(brief, b, 0o644)
	ignition(t, dir, "advance", id)
	plan := filepath.Join(dir, "work", id, "plan.md")
	b, _ = os.ReadFile(plan)
	os.WriteFile(plan, bytes.Replace(b, []byte("|  |  |  |  |  |"), []byte("| T1 | Config | Both | done | - |"), 1), 0o644)
	ignition(t, dir, "advance", id)
	ignition(t, dir, "log", id, "--type", "change", "a1b2c3d adds logrotate")
	ignition(t, dir, "advance", id)

	out, code := ignition(t, dir, "qa", id)
	for _, want := range []string{"Revision to test: a1b2c3d", "1. Logs rotate daily.", "criterion 2 has no QA evidence", "--by qa"} {
		if code != 0 || !strings.Contains(out, want) {
			t.Fatalf("qa: %d, missing %q\n%s", code, want, out)
		}
	}
	var buf bytes.Buffer
	if code := run([]string{"evidence", id, "--criteria", "1,2", "--claim", "Rotation and retention", "--result", "pass",
		"--revision", "a1b2c3d", "--environment", "local", "--method", "qa/rotate_test.sh", "--dir", dir, "--by", "qa"}, &buf, &buf); code != 0 {
		t.Fatalf("qa evidence: %d\n%s", code, buf.String())
	}
	if out, _ := ignition(t, dir, "status"); !strings.Contains(out, "waits for the owner") {
		t.Fatalf("status after QA:\n%s", out)
	}
}
