package workspace

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/mohsalsaleem/ignition-engineering/internal/model"
	"github.com/mohsalsaleem/ignition-engineering/internal/record"
	"github.com/mohsalsaleem/ignition-engineering/templates"
)

// Unmet lists what keeps the item from meeting a stage's exit rules.
func (it *Item) Unmet(s model.Stage) []string {
	var out []string
	for _, r := range model.Exit(it.Facts(), s) {
		out = append(out, it.unmet(r)...)
	}
	return out
}

func (it *Item) unmet(r model.Rule) []string {
	switch r.Kind {
	case model.Sections:
		d := it.Docs[r.Artifact]
		if d == nil {
			return []string{fmt.Sprintf("%s is missing", r.Artifact)}
		}
		var out []string
		for _, n := range r.Names {
			if _, ok := d.Section(n); !ok {
				out = append(out, fmt.Sprintf("%s has no %q section", r.Artifact, n))
			} else if !d.Filled(n) {
				out = append(out, fmt.Sprintf("%s: fill in %q", r.Artifact, n))
			}
		}
		return out
	case model.Fields:
		var out []string
		for _, n := range r.Names {
			v := it.Brief.Get(n)
			switch {
			case v == "":
				out = append(out, fmt.Sprintf("brief.md: set %s", n))
			case n == "risk" && !slices.Contains(model.Risks, model.Risk(v)):
				out = append(out, fmt.Sprintf("brief.md: risk %q is not low, normal or high", v))
			case n == "size" && !slices.Contains(model.Sizes, v):
				out = append(out, fmt.Sprintf("brief.md: size %q is not small, medium or large", v))
			}
		}
		return out
	case model.EvidenceResult:
		for _, e := range it.Evidence {
			if e.Get("result") == r.Result {
				return nil
			}
		}
		return []string{fmt.Sprintf("no evidence with result %s (add it with `ignition evidence %s --result %s ...`)", r.Result, it.ID, r.Result)}
	case model.Verified:
		return it.unverified(r.Independent)
	case model.EvidenceAny:
		if len(it.Evidence) == 0 {
			return []string{fmt.Sprintf("no evidence recorded (add it with `ignition evidence %s ...`)", it.ID)}
		}
	case model.TasksListed, model.TasksDone:
		d := it.Docs[model.PlanDoc]
		if d == nil {
			return []string{"plan.md is missing"}
		}
		tasks, err := d.Tasks()
		if err != nil {
			return []string{"plan.md: " + err.Error()}
		}
		if len(tasks) == 0 {
			return []string{"plan.md: list at least one task"}
		}
		if r.Kind == model.TasksDone {
			var out []string
			for _, t := range tasks {
				if t.Status != "done" && t.Status != "dropped" {
					out = append(out, fmt.Sprintf("plan.md: task %s is %q, not done or dropped", t.ID, t.Status))
				}
			}
			return out
		}
	case model.LogEntry:
		for _, e := range it.Log {
			if e.Type == r.Entry {
				return nil
			}
		}
		return []string{fmt.Sprintf("log: no %s entry (record it with `ignition log %s --type %s \"...\"`)", r.Entry, it.ID, r.Entry)}
	case model.ChildrenClosed:
		if len(it.Children) == 0 {
			return []string{fmt.Sprintf("no child items (create them with `ignition new <kind> \"...\" --parent %s`)", it.ID)}
		}
		var out []string
		for _, c := range it.Children {
			if c.Stage() != model.Closed {
				out = append(out, fmt.Sprintf("child %s is in %s", c.ID, c.Stage()))
			}
		}
		return out
	}
	return nil
}

// Check lists everything wrong with an item's records: values the model
// doesn't know, malformed evidence or plans, and stages already passed whose
// exit rules no longer hold. What the current stage still needs is not a
// problem; see Unmet.
func (it *Item) Check() []string {
	out := slices.Clone(it.Problems)
	f := it.Facts()
	if !slices.Contains(model.Kinds, f.Kind) {
		return append(out, fmt.Sprintf("brief.md: unknown kind %q", f.Kind))
	}
	if it.Brief.Get("id") != it.ID {
		out = append(out, fmt.Sprintf("brief.md: id %q doesn't match the folder %s", it.Brief.Get("id"), it.ID))
	}
	if v := it.Brief.Get("risk"); v != "" && !slices.Contains(model.Risks, model.Risk(v)) {
		out = append(out, fmt.Sprintf("brief.md: risk %q is not low, normal or high", v))
	}
	if !model.OnRoute(f, it.Stage()) {
		return append(out, fmt.Sprintf("brief.md: stage %q isn't on the %s route (%s)", it.Stage(), f.Kind, joinStages(model.Route(f))))
	}
	for _, e := range it.Evidence {
		out = append(out, checkEvidence(it.ID, e)...)
		if c := e.Get("criteria"); c != "" {
			if _, err := it.normalizeCriteria(c); err != nil {
				out = append(out, fmt.Sprintf("evidence %s: %v", e.Get("id"), err))
			}
		}
	}
	if d := it.Docs[model.PlanDoc]; d != nil {
		if tasks, err := d.Tasks(); err != nil {
			out = append(out, "plan.md: "+err.Error())
		} else {
			for _, t := range tasks {
				if !slices.Contains(model.TaskStatuses, t.Status) {
					out = append(out, fmt.Sprintf("plan.md: task %s has status %q (%s)", t.ID, t.Status, strings.Join(model.TaskStatuses, ", ")))
				}
			}
		}
	}
	for _, s := range model.Passed(f, it.Stage()) {
		if s == model.Verify {
			// Evidence is append-only, so verify can't be emptied. A failure
			// found after release is new work: a bug, not a broken record.
			continue
		}
		for _, u := range it.Unmet(s) {
			out = append(out, fmt.Sprintf("passed %s, but %s", s, u))
		}
	}
	return out
}

func checkEvidence(work string, e *record.Doc) []string {
	id := e.Get("id")
	var out []string
	for _, k := range []string{"id", "claim", "revision", "environment", "method", "at"} {
		if e.Get(k) == "" {
			out = append(out, fmt.Sprintf("evidence %s: %s is empty", id, k))
		}
	}
	if !slices.Contains(model.Results, e.Get("result")) {
		out = append(out, fmt.Sprintf("evidence %s: result %q is not pass, fail or inconclusive", id, e.Get("result")))
	}
	if e.Get("work") != work {
		out = append(out, fmt.Sprintf("evidence %s: work %q doesn't match %s", id, e.Get("work"), work))
	}
	return out
}

func joinStages(ss []model.Stage) string {
	var parts []string
	for _, s := range ss {
		parts = append(parts, string(s))
	}
	return strings.Join(parts, " → ")
}

// Checkpoint is an owner decision recorded when leaving a stage.
type Checkpoint struct {
	Owner    string // who decided
	Decision string // what they decided, in their words
}

// ErrNotReady is returned when the current stage's exit rules aren't met.
var ErrNotReady = errors.New("stage not complete")

// ErrNeedsOwner is returned when leaving the stage needs an owner decision
// that wasn't given.
var ErrNeedsOwner = errors.New("owner decision needed")

// Advance moves an item to its next stage. The current stage's exit rules must
// hold and any earlier problem must be fixed; where the owner decides, their
// decision must be given and is recorded in the log.
func (w *Workspace) Advance(id string, cp *Checkpoint, by Actor) (model.Stage, error) {
	it, err := w.Load(id)
	if err != nil {
		return "", err
	}
	if probs := it.Check(); len(probs) > 0 {
		return "", fmt.Errorf("%w: fix these first:\n  - %s", ErrNotReady, strings.Join(probs, "\n  - "))
	}
	cur := it.Stage()
	if cur == model.Closed {
		return "", fmt.Errorf("%s is closed", id)
	}
	if unmet := it.Unmet(cur); len(unmet) > 0 {
		return "", fmt.Errorf("%w: %s needs:\n  - %s", ErrNotReady, cur, strings.Join(unmet, "\n  - "))
	}
	f := it.Facts()
	next := model.Next(f, cur)
	need, question := model.OwnerDecides(f, cur)
	if need && (cp == nil || strings.TrimSpace(cp.Owner) == "" || strings.TrimSpace(cp.Decision) == "") {
		return "", fmt.Errorf("%w: leaving %s needs the owner's answer to: %s\n  Ask them, then run `ignition advance %s --owner <name> --decision \"<their words>\"`", ErrNeedsOwner, cur, question, id)
	}
	if a := model.Opens(next); a != "" {
		if _, ok := it.Docs[a]; !ok {
			text, err := templates.Render(templates.Name(a, it.Kind()), map[string]string{"ID": id, "Title": it.Title(), "Parent": it.Brief.Get("parent")})
			if err != nil {
				return "", err
			}
			if err := writeNew(filepath.Join(it.Dir, string(a)), text); err != nil {
				return "", err
			}
		}
	}
	if need {
		text := fmt.Sprintf("Leaving %s: %s\nOwner: %s\nDecision: %s", cur, question, cp.Owner, cp.Decision)
		if err := w.appendEntry(it, "checkpoint", by, text); err != nil {
			return "", err
		}
	}
	it.Brief.Set("stage", string(next))
	if err := os.WriteFile(filepath.Join(it.Dir, string(model.BriefDoc)), []byte(it.Brief.String()), 0o644); err != nil {
		return "", err
	}
	return next, w.appendEntry(it, "stage", by, fmt.Sprintf("%s → %s", cur, next))
}
