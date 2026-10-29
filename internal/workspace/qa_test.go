package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mohsalsaleem/ignition-engineering/internal/model"
)

// buildTask takes a task with two numbered criteria to verify, built by agent
// in commit a1b2c3d.
func buildTask(t *testing.T, w *Workspace, risk string) *Item {
	t.Helper()
	it, err := w.NewItem(model.Task, "Rotate logs", risk, "small", "", agent)
	if err != nil {
		t.Fatal(err)
	}
	fill(t, it, "brief.md", "Request", "Logs fill the disk.")
	fill(t, it, "brief.md", "Done when", "1. Logs rotate daily.\n2. Logs older than 7 days are deleted.")
	advance(t, w, it.ID, nil)
	plan := filepath.Join(it.Dir, "plan.md")
	b, _ := os.ReadFile(plan)
	b = []byte(strings.Replace(string(b), "|  |  |  |  |  |", "| T1 | Add logrotate config | Both criteria | done | - |", 1))
	if err := os.WriteFile(plan, b, 0o644); err != nil {
		t.Fatal(err)
	}
	advance(t, w, it.ID, nil)
	if err := w.Log(it.ID, "change", agent, "a1b2c3d adds a daily logrotate config"); err != nil {
		t.Fatal(err)
	}
	if s := advance(t, w, it.ID, nil); s != model.Verify {
		t.Fatalf("after build: %s", s)
	}
	return it
}

// tick makes each of w's timestamps a second after the last, so the order of
// evidence is known.
func tick(w *Workspace) {
	start, n := w.Now(), 0
	w.Now = func() time.Time {
		n++
		return start.Add(time.Duration(n) * time.Second)
	}
}

func evidence(t *testing.T, w *Workspace, id string, by Actor, rev, criteria, result string) {
	t.Helper()
	if _, err := w.AddEvidence(id, EvidenceInput{Claim: "Logs rotate", Criteria: criteria, Result: result,
		Revision: rev, Environment: "staging", Method: "qa/rotate_test.sh", By: by}); err != nil {
		t.Fatal(err)
	}
}

func TestQAVerifiesEveryCriterionOnTheBuiltRevision(t *testing.T) {
	w := newWorkspace(t)
	tick(w)
	it := buildTask(t, w, "normal")
	owner := &Checkpoint{Owner: "owner", Decision: "Ship it."}

	evidence(t, w, it.ID, agent, "a1b2c3d", "1,2", "pass")
	mustFail(t, w, it.ID, owner, ErrNotReady, "criterion 1 has no QA evidence on a1b2c3d")

	evidence(t, w, it.ID, qaAgent, "a1b2c3d", "1", "pass")
	evidence(t, w, it.ID, qaAgent, "9f9f9f9", "2", "pass") // another revision doesn't count
	mustFail(t, w, it.ID, owner, ErrNotReady, "criterion 2 has no QA evidence on a1b2c3d")

	evidence(t, w, it.ID, qaAgent, "a1b2c3d4e5f6 (staging deploy)", "2", "pass")
	mustFail(t, w, it.ID, nil, ErrNeedsOwner, "release to users")

	// A failure after the pass reopens the criterion, whoever recorded it.
	evidence(t, w, it.ID, agent, "a1b2c3d", "1", "fail")
	mustFail(t, w, it.ID, owner, ErrNotReady, "criterion 1 failed on a1b2c3d")

	// The fix is a new change: what was verified before no longer counts.
	if err := w.Log(it.ID, "change", agent, "e4f5a6b keeps the current log when rotating"); err != nil {
		t.Fatal(err)
	}
	mustFail(t, w, it.ID, owner, ErrNotReady, "criterion 2 has no QA evidence on e4f5a6b")
	evidence(t, w, it.ID, qaAgent, "e4f5a6b", "1,2", "pass")
	if s := advance(t, w, it.ID, owner); s != model.Release {
		t.Fatalf("after verify: %s", s)
	}
	it, _ = w.Load(it.ID)
	if probs := it.Check(); len(probs) != 0 {
		t.Errorf("verified item has problems: %v", probs)
	}

	if _, err := w.AddEvidence(it.ID, EvidenceInput{Claim: "x", Criteria: "3", Result: "pass", Revision: "e4f5a6b",
		Environment: "local", Method: "x", By: qaAgent}); err == nil || !strings.Contains(err.Error(), "isn't one of") {
		t.Errorf("evidence for an unknown criterion: %v", err)
	}
}

func TestLowRiskWorkMayVerifyItself(t *testing.T) {
	w := newWorkspace(t)
	it := buildTask(t, w, "low")
	evidence(t, w, it.ID, agent, "a1b2c3d", "1,2", "pass")
	if s := advance(t, w, it.ID, nil); s != model.Release {
		t.Fatalf("after verify: %s", s)
	}
}

func TestQAIsSomeoneOtherThanTheBuilder(t *testing.T) {
	w := newWorkspace(t)
	it := buildTask(t, w, "high")
	for _, a := range []Actor{
		{Name: agent.Name, Session: agent.Session}, // the builder, under another model
		{Name: agent.Name, Agent: "codex", Session: agent.Session},
	} {
		evidence(t, w, it.ID, a, "a1b2c3d", "1,2", "pass")
	}
	mustFail(t, w, it.ID, nil, ErrNotReady, "no passing QA evidence")
	evidence(t, w, it.ID, Actor{Name: "qa", Agent: "claude-code", Session: agent.Session}, "a1b2c3d", "1,2", "pass")
	mustFail(t, w, it.ID, nil, ErrNeedsOwner, "release to users")
}

func TestBuildRevisionIsTheFirstCommitInTheChangeEntry(t *testing.T) {
	for text, want := range map[string]string{
		"a1b2c3d adds rotation":                    "a1b2c3d",
		"Commit c18a6ac skips the local check.":    "c18a6ac",
		"defaced feed, then 0123abcd fixes it":     "0123abcd",
		"rewrites the docs; no commit in this one": "",
	} {
		if got := commitIn(text); got != want {
			t.Errorf("commitIn(%q) = %q, want %q", text, got, want)
		}
	}
}
