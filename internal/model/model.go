// Package model defines Ignition's work model: kinds of work, the stages each
// kind passes through, what must be true to leave a stage, and where the owner
// decides. It has no I/O.
package model

import "fmt"

type Kind string

const (
	Initiative    Kind = "initiative"
	Feature       Kind = "feature"
	Bug           Kind = "bug"
	Task          Kind = "task"
	Investigation Kind = "investigation"
	Incident      Kind = "incident"
)

var Kinds = []Kind{Initiative, Feature, Bug, Task, Investigation, Incident}

type Stage string

const (
	Intake    Stage = "intake"
	Shape     Stage = "shape"
	Design    Stage = "design"
	Plan      Stage = "plan"
	Build     Stage = "build"
	Verify    Stage = "verify"
	Release   Stage = "release"
	Observe   Stage = "observe"
	Deliver   Stage = "deliver"
	Reproduce Stage = "reproduce"
	Diagnose  Stage = "diagnose"
	Method    Stage = "method"
	Findings  Stage = "findings"
	Decide    Stage = "decision"
	Contain   Stage = "contain"
	Timeline  Stage = "timeline"
	Recover   Stage = "recover"
	Review    Stage = "review"
	Closed    Stage = "closed"
)

type Risk string

const (
	Low    Risk = "low"
	Normal Risk = "normal"
	High   Risk = "high"
)

var Risks = []Risk{Low, Normal, High}

var Sizes = []string{"small", "medium", "large"}

func (r Risk) rank() int {
	switch r {
	case Low:
		return 1
	case Normal:
		return 2
	case High:
		return 3
	}
	return 0
}

// AtLeast reports whether r is at least min. An unset risk counts as high, so
// an unclassified item never skips an owner decision.
func (r Risk) AtLeast(min Risk) bool {
	if r.rank() == 0 {
		return true
	}
	return r.rank() >= min.rank()
}

var routes = map[Kind][]Stage{
	Initiative:    {Intake, Shape, Design, Plan, Deliver, Observe, Closed},
	Feature:       {Intake, Shape, Design, Plan, Build, Verify, Release, Observe, Closed},
	Bug:           {Intake, Reproduce, Diagnose, Build, Verify, Release, Closed},
	Task:          {Intake, Plan, Build, Verify, Release, Closed},
	Investigation: {Intake, Method, Findings, Decide, Closed},
	Incident:      {Contain, Timeline, Recover, Review, Closed},
}

// Facts are the item attributes that change its route.
type Facts struct {
	Kind      Kind
	Risk      Risk
	Size      string
	NoRelease bool // the brief says `release: none`: nothing ships to users
	Child     bool // the item belongs to an initiative
}

func ParseKind(s string) (Kind, error) {
	for _, k := range Kinds {
		if string(k) == s {
			return k, nil
		}
	}
	return "", fmt.Errorf("unknown kind %q (one of initiative, feature, bug, task, investigation, incident)", s)
}

// Route is the ordered list of stages an item with these facts passes through.
func Route(f Facts) []Stage {
	var out []Stage
	for _, s := range routes[f.Kind] {
		if s == Design && f.Kind == Feature && f.Risk == Low {
			continue
		}
		if s == Release && f.NoRelease {
			continue
		}
		out = append(out, s)
	}
	return out
}

// First is the stage a new item starts in.
func First(k Kind) Stage { return routes[k][0] }

// Next is the stage after s on the item's route, or "" if s is the last stage
// or isn't on the route.
func Next(f Facts, s Stage) Stage {
	r := Route(f)
	for i, st := range r {
		if st == s && i+1 < len(r) {
			return r[i+1]
		}
	}
	return ""
}

// Passed reports the stages before s on the item's route.
func Passed(f Facts, s Stage) []Stage {
	var out []Stage
	for _, st := range Route(f) {
		if st == s {
			return out
		}
		out = append(out, st)
	}
	return nil
}

func OnRoute(f Facts, s Stage) bool {
	for _, st := range Route(f) {
		if st == s {
			return true
		}
	}
	return false
}

// Artifact is a document an item keeps in its folder.
type Artifact string

const (
	BriefDoc   Artifact = "brief.md"
	DesignDoc  Artifact = "design.md"
	PlanDoc    Artifact = "plan.md"
	ReleaseDoc Artifact = "release.md"
)

// Opens is the artifact a stage writes, created when the item enters it.
func Opens(s Stage) Artifact {
	switch s {
	case Design:
		return DesignDoc
	case Plan:
		return PlanDoc
	case Release:
		return ReleaseDoc
	}
	return ""
}

// RuleKind is the type of check a rule makes.
type RuleKind int

const (
	// Sections: every named section of Artifact has content.
	Sections RuleKind = iota
	// Fields: every named frontmatter field of the brief has a valid value.
	Fields
	// EvidenceResult: at least one evidence record has Result.
	EvidenceResult
	// EvidenceAny: at least one evidence record exists.
	EvidenceAny
	// TasksDone: the plan has tasks and every one is done or dropped.
	TasksDone
	// TasksListed: the plan has at least one task.
	TasksListed
	// LogEntry: the log has an entry of type Entry.
	LogEntry
	// ChildrenClosed: the item has child items and all are closed.
	ChildrenClosed
	// Verified: evidence on the built revision passes, every acceptance
	// criterion passes, and nothing on that revision has failed since. With
	// Independent, the passes must come from someone other than the builder.
	Verified
)

type Rule struct {
	Kind        RuleKind
	Artifact    Artifact
	Names       []string
	Result      string
	Entry       string
	Independent bool
}

// CriteriaSection is the brief section whose numbered lines are an item's
// acceptance criteria, which verify evidence must cover. Other kinds have none.
func CriteriaSection(k Kind) string {
	switch k {
	case Feature:
		return "Acceptance criteria"
	case Task:
		return "Done when"
	}
	return ""
}

// IndependentQA reports whether verifying an item needs evidence from someone
// other than its builder: a QA agent or person. Low-risk work may verify
// itself.
func IndependentQA(f Facts) bool { return f.Risk.AtLeast(Normal) }

func sections(a Artifact, names ...string) Rule {
	return Rule{Kind: Sections, Artifact: a, Names: names}
}

// DetailedDesign reports whether a feature's design must include a detailed
// design: interfaces, data, flows and errors. Small features that aren't high
// risk don't need one. An unset size or risk counts as needing it.
func DetailedDesign(f Facts) bool {
	return f.Kind == Feature && (f.Size != "small" || f.Risk.AtLeast(High))
}

// Exit is what must be true for an item with facts f to leave stage s.
func Exit(f Facts, s Stage) []Rule {
	k := f.Kind
	switch s {
	case Intake:
		return []Rule{sections(BriefDoc, "Request"), {Kind: Fields, Names: []string{"title", "size", "risk"}}}
	case Shape:
		return []Rule{sections(BriefDoc, "Problem", "Outcome", "Acceptance criteria", "Scope")}
	case Design:
		if k == Initiative {
			// The solution design: the whole of the initiative, which its
			// child items build on.
			return []Rule{sections(DesignDoc, "Context", "Options", "Decision", "Architecture", "Interfaces",
				"Cross-cutting concerns", "Slicing", "Risks", "Test strategy")}
		}
		names := []string{"Context", "Options", "Decision", "Impacts", "Test strategy"}
		if f.Child {
			names = append(names, "Parent design")
		}
		if DetailedDesign(f) {
			names = append(names, "Detailed design")
		}
		return []Rule{sections(DesignDoc, names...)}
	case Plan:
		if k == Task {
			return []Rule{sections(BriefDoc, "Done when"), {Kind: TasksListed}}
		}
		return []Rule{{Kind: TasksListed}}
	case Build:
		if k == Bug {
			return []Rule{{Kind: LogEntry, Entry: "change"}}
		}
		return []Rule{{Kind: TasksDone}, {Kind: LogEntry, Entry: "change"}}
	case Verify:
		return []Rule{{Kind: Verified, Independent: IndependentQA(f)}}
	case Release:
		return []Rule{sections(ReleaseDoc, "Revision", "Checks", "Rollback", "Production verification")}
	case Observe:
		return []Rule{sections(BriefDoc, "Result")}
	case Deliver:
		return []Rule{{Kind: ChildrenClosed}}
	case Reproduce:
		return []Rule{sections(BriefDoc, "Expected and actual", "Reproduction"), {Kind: EvidenceResult, Result: "fail"}}
	case Diagnose:
		return []Rule{sections(BriefDoc, "Root cause")}
	case Method:
		return []Rule{sections(BriefDoc, "Question", "Method")}
	case Findings:
		return []Rule{sections(BriefDoc, "Findings"), {Kind: EvidenceAny}}
	case Decide:
		return []Rule{sections(BriefDoc, "Decision", "Follow-up")}
	case Contain:
		return []Rule{sections(BriefDoc, "Impact", "Containment"), {Kind: Fields, Names: []string{"title", "risk"}}}
	case Timeline:
		return []Rule{sections(BriefDoc, "Timeline")}
	case Recover:
		return []Rule{sections(BriefDoc, "Recovery"), {Kind: EvidenceResult, Result: "pass"}}
	case Review:
		return []Rule{sections(BriefDoc, "Cause", "Follow-up")}
	}
	return nil
}

// OwnerDecides reports whether leaving stage s needs a decision from the owner,
// and the question the owner answers.
func OwnerDecides(f Facts, s Stage) (bool, string) {
	next := Next(f, s)
	switch {
	case s == Shape && f.Risk.AtLeast(Normal):
		return true, "Is this the right problem, outcome and scope?"
	case s == Design && f.Kind == Initiative && f.Risk.AtLeast(Normal):
		return true, "Is this the right solution design for the whole initiative?"
	case s == Design && f.Risk.AtLeast(High):
		return true, "Is this the right approach?"
	case s == Diagnose && f.Risk.AtLeast(High):
		return true, "Is the root cause right, and should the fix go ahead?"
	case s == Verify && next == Release && f.Risk.AtLeast(Normal):
		return true, "Is the evidence enough to release to users?"
	case s == Decide && f.Risk.AtLeast(Normal):
		return true, "Do you accept this decision and its follow-up?"
	case s == Review:
		return true, "Do you accept the incident review and its follow-up?"
	}
	return false, ""
}

// Evidence results.
var Results = []string{"pass", "fail", "inconclusive"}

// Plan task statuses.
var TaskStatuses = []string{"todo", "doing", "done", "dropped"}

// Log entry types.
var EntryTypes = []string{"note", "finding", "decision", "change", "handoff", "scope", "stage", "checkpoint", "evidence"}
