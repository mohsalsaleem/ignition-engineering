package workspace

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mohsalsaleem/ignition-engineering/internal/model"
)

var agent = Actor{Name: "claude-code", Agent: "claude-code", Model: "claude-opus-5-5", Session: "s-1"}

func newWorkspace(t *testing.T) *Workspace {
	t.Helper()
	root := t.TempDir()
	if err := Init(root, "example"); err != nil {
		t.Fatal(err)
	}
	w := Open(root)
	w.Now = func() time.Time { return time.Date(2026, 10, 27, 12, 0, 0, 0, time.UTC) }
	return w
}

// fill replaces a section's guidance with content.
func fill(t *testing.T, it *Item, file, section, content string) {
	t.Helper()
	path := filepath.Join(it.Dir, file)
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	heading := "## " + section + "\n"
	i := strings.Index(string(b), heading)
	if i < 0 {
		t.Fatalf("%s has no %s", file, heading)
	}
	at := i + len(heading)
	out := string(b[:at]) + "\n" + content + "\n" + string(b[at:])
	if err := os.WriteFile(path, []byte(out), 0o644); err != nil {
		t.Fatal(err)
	}
}

func advance(t *testing.T, w *Workspace, id string, cp *Checkpoint) model.Stage {
	t.Helper()
	s, err := w.Advance(id, cp, agent)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func mustFail(t *testing.T, w *Workspace, id string, cp *Checkpoint, want error, mention string) {
	t.Helper()
	_, err := w.Advance(id, cp, agent)
	if !errors.Is(err, want) || !strings.Contains(err.Error(), mention) {
		t.Fatalf("Advance = %v, want %v mentioning %q", err, want, mention)
	}
}

func TestBugLifecycle(t *testing.T) {
	w := newWorkspace(t)
	it, err := w.NewItem(model.Bug, "Install ships a local image", "normal", "small", "", agent)
	if err != nil {
		t.Fatal(err)
	}
	mustFail(t, w, it.ID, nil, ErrNotReady, `fill in "Request"`)
	fill(t, it, "brief.md", "Request", "The rc.1 install was refused.")
	if s := advance(t, w, it.ID, nil); s != model.Reproduce {
		t.Fatalf("after intake: %s", s)
	}

	fill(t, it, "brief.md", "Expected and actual", "Expected a pull; got an arch error.")
	fill(t, it, "brief.md", "Reproduction", "Install with a local arm64 copy.")
	mustFail(t, w, it.ID, nil, ErrNotReady, "no evidence with result fail")
	if _, err := w.AddEvidence(it.ID, EvidenceInput{Claim: "Install refuses a released image", Result: "fail",
		Revision: "abc123", Environment: "local", Method: "go test", By: agent}); err != nil {
		t.Fatal(err)
	}
	advance(t, w, it.ID, nil)

	fill(t, it, "brief.md", "Root cause", "shipImage checks the local copy first.")
	if s := advance(t, w, it.ID, nil); s != model.Build {
		t.Fatalf("normal-risk diagnose needed no owner, got %s", s)
	}
	mustFail(t, w, it.ID, nil, ErrNotReady, "no change entry")
	if err := w.Log(it.ID, "change", agent, "Commit c18a6ac skips the local check."); err != nil {
		t.Fatal(err)
	}
	advance(t, w, it.ID, nil)

	if _, err := w.AddEvidence(it.ID, EvidenceInput{Claim: "Released image is pulled", Result: "pass",
		Revision: "c18a6ac", Environment: "local", Method: "go test", By: agent}); err != nil {
		t.Fatal(err)
	}
	mustFail(t, w, it.ID, nil, ErrNeedsOwner, "release to users")
	mustFail(t, w, it.ID, &Checkpoint{Owner: "owner"}, ErrNeedsOwner, "release to users")
	if s := advance(t, w, it.ID, &Checkpoint{Owner: "owner", Decision: "Publish rc.2."}); s != model.Release {
		t.Fatalf("after verify: %s", s)
	}
	if _, err := os.Stat(filepath.Join(it.Dir, "release.md")); err != nil {
		t.Fatal("release.md wasn't created on entering release")
	}
	for _, s := range []string{"Revision", "Checks", "Rollback", "Production verification"} {
		fill(t, it, "release.md", s, "Done.")
	}
	if s := advance(t, w, it.ID, nil); s != model.Closed {
		t.Fatalf("after release: %s", s)
	}

	it, _ = w.Load(it.ID)
	var checkpoints int
	for _, e := range it.Log {
		if e.Type == "checkpoint" {
			checkpoints++
			if !strings.Contains(e.Text, "Decision: Publish rc.2.") {
				t.Errorf("checkpoint entry lacks the decision: %q", e.Text)
			}
		}
	}
	if checkpoints != 1 {
		t.Errorf("checkpoints = %d, want 1", checkpoints)
	}
	if probs := it.Check(); len(probs) != 0 {
		t.Errorf("closed item has problems: %v", probs)
	}
}

func TestEmptyingAPassedStageIsAProblem(t *testing.T) {
	w := newWorkspace(t)
	it, _ := w.NewItem(model.Task, "Rotate logs", "low", "small", "", agent)
	fill(t, it, "brief.md", "Request", "Logs fill the disk.")
	advance(t, w, it.ID, nil)

	path := filepath.Join(it.Dir, "brief.md")
	b, _ := os.ReadFile(path)
	os.WriteFile(path, bytes.Replace(b, []byte("Logs fill the disk."), nil, 1), 0o644)
	it, _ = w.Load(it.ID)
	if probs := it.Check(); len(probs) != 1 || !strings.Contains(probs[0], `passed intake, but brief.md: fill in "Request"`) {
		t.Fatalf("problems = %v", probs)
	}
	mustFail(t, w, it.ID, nil, ErrNotReady, "fix these first")
}

func TestHighRiskFeatureStopsForTheOwnerTwiceBeforeBuilding(t *testing.T) {
	w := newWorkspace(t)
	it, _ := w.NewItem(model.Feature, "Withdraw an application", "high", "medium", "", agent)
	fill(t, it, "brief.md", "Request", "Candidates want to withdraw.")
	advance(t, w, it.ID, nil)
	for _, s := range []string{"Problem", "Outcome", "Acceptance criteria", "Scope"} {
		fill(t, it, "brief.md", s, "Written.")
	}
	mustFail(t, w, it.ID, nil, ErrNeedsOwner, "right problem")
	advance(t, w, it.ID, &Checkpoint{Owner: "owner", Decision: "Yes."})
	for _, s := range []string{"Context", "Options", "Decision", "Impacts", "Test strategy"} {
		fill(t, it, "design.md", s, "Written.")
	}
	mustFail(t, w, it.ID, nil, ErrNeedsOwner, "right approach")
	if s := advance(t, w, it.ID, &Checkpoint{Owner: "owner", Decision: "Option B."}); s != model.Plan {
		t.Fatalf("after design: %s", s)
	}
	mustFail(t, w, it.ID, nil, ErrNotReady, "list at least one task")
}

func TestInitiativeClosesOnlyWhenChildrenAreClosed(t *testing.T) {
	w := newWorkspace(t)
	init, _ := w.NewItem(model.Initiative, "Build an ATS", "high", "large", "", agent)
	if _, err := w.NewItem(model.Task, "Child", "low", "small", "missing", agent); err == nil {
		t.Fatal("accepted a missing parent")
	}
	child, err := w.NewItem(model.Task, "Child", "low", "small", init.ID, agent)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.NewItem(model.Task, "Grandchild", "low", "small", child.ID, agent); err == nil {
		t.Fatal("a task accepted a child")
	}
	it, _ := w.Load(init.ID)
	if u := it.unmet(model.Rule{Kind: model.ChildrenClosed}); len(u) != 1 || !strings.Contains(u[0], child.ID) {
		t.Fatalf("unmet = %v", u)
	}
}

// repeatReader returns the same bytes for the first n reads, then counts up.
type repeatReader struct{ n, reads int }

func (r *repeatReader) Read(p []byte) (int, error) {
	r.reads++
	for i := range p {
		p[i] = 0
		if r.reads > r.n {
			p[i] = byte(r.reads)
		}
	}
	return len(p), nil
}

func TestANewItemNeverReusesAnExistingID(t *testing.T) {
	w := newWorkspace(t)
	w.Rand = &repeatReader{n: 2}
	a, err := w.NewItem(model.Task, "A", "low", "small", "", agent)
	if err != nil {
		t.Fatal(err)
	}
	b, err := w.NewItem(model.Task, "B", "low", "small", "", agent)
	if err != nil {
		t.Fatal(err)
	}
	if a.ID == b.ID || !strings.HasPrefix(b.ID, "W-20261027-") || len(b.ID) != len("W-20261027-")+6 {
		t.Fatalf("IDs %s and %s", a.ID, b.ID)
	}
	if a, _ := w.Load(a.ID); a.Title() != "A" {
		t.Fatal("the first item was overwritten")
	}
}

func TestLogRefusesStageEntries(t *testing.T) {
	w := newWorkspace(t)
	it, _ := w.NewItem(model.Task, "T", "low", "small", "", agent)
	if err := w.Log(it.ID, "checkpoint", agent, "Owner said yes"); err == nil {
		t.Fatal("an agent wrote a checkpoint entry directly")
	}
}

func TestFind(t *testing.T) {
	app := t.TempDir()
	if err := Init(filepath.Join(app, ".ignition"), "example"); err != nil {
		t.Fatal(err)
	}
	sub := filepath.Join(app, "src", "pkg")
	os.MkdirAll(sub, 0o755)
	got, err := Find(sub)
	if err != nil || got != filepath.Join(app, ".ignition") {
		t.Fatalf("Find = %q, %v", got, err)
	}
	if _, err := Find(t.TempDir()); err == nil {
		t.Fatal("found a workspace in an empty folder")
	}
}

func TestAttribution(t *testing.T) {
	w := newWorkspace(t)
	owner := Actor{Name: "Owner"}
	it, _ := w.NewItem(model.Bug, "Crash on empty export", "normal", "small", "", agent)
	fill(t, it, "brief.md", "Request", "Export crashes.")
	advance(t, w, it.ID, nil)
	fill(t, it, "brief.md", "Expected and actual", "A file; a crash.")
	fill(t, it, "brief.md", "Reproduction", "Export an empty report.")
	other := Actor{Name: "codex", Agent: "codex", Model: "gpt-x"}
	if _, err := w.AddEvidence(it.ID, EvidenceInput{Claim: "Crashes", Result: "fail", Revision: "a1", Environment: "local", Method: "test", By: other}); err != nil {
		t.Fatal(err)
	}
	if err := w.Log(it.ID, "change", agent, "b2c3d4e fix the crash"); err != nil {
		t.Fatal(err)
	}
	if err := w.Log(it.ID, "note", owner, "Looks right."); err != nil {
		t.Fatal(err)
	}
	if err := w.Log(it.ID, "evidence", agent, "I checked it"); err == nil {
		t.Fatal("an evidence entry was logged without evidence")
	}

	it, _ = w.Load(it.ID)
	contribs, _ := it.Attribution()
	if len(contribs) != 3 {
		t.Fatalf("contributors = %+v", contribs)
	}
	a, c, o := contribs[0], contribs[1], contribs[2]
	if a.Actor.Model != "claude-opus-5-5" || a.Entries["stage"] != 2 || a.Entries["change"] != 1 ||
		len(a.Commits) != 1 || a.Commits[0] != "b2c3d4e" || len(a.Sessions) != 1 {
		t.Errorf("agent contribution = %+v", a)
	}
	if c.Actor.Agent != "codex" || c.Entries["evidence"] != 1 {
		t.Errorf("codex contribution = %+v", c)
	}
	if o.Actor.Name != "Owner" || o.Actor.Agent != "" || o.Entries["note"] != 1 {
		t.Errorf("owner contribution = %+v", o)
	}
	if e := it.Evidence[0]; e.Get("agent") != "codex" || e.Get("model") != "gpt-x" || e.Get("by") != "codex" {
		t.Errorf("evidence attribution: agent %q model %q by %q", e.Get("agent"), e.Get("model"), e.Get("by"))
	}
}
