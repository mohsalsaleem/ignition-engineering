package workspace

import (
	"fmt"
	"regexp"
	"slices"
	"sort"
	"strings"

	"github.com/mohsalsaleem/ignition-engineering/internal/model"
	"github.com/mohsalsaleem/ignition-engineering/internal/record"
)

// Criterion is one numbered acceptance criterion from the brief.
type Criterion struct {
	Num, Text string
}

var criterionLine = regexp.MustCompile(`^(\d+)[.)]\s+(.*)$`)

// Criteria lists the item's acceptance criteria: the top-level numbered lines
// of its criteria section (see model.CriteriaSection). Unnumbered criteria
// can't be linked to evidence, so they aren't listed.
func (it *Item) Criteria() []Criterion {
	name := model.CriteriaSection(it.Kind())
	if name == "" {
		return nil
	}
	s, _ := it.Brief.Section(name)
	var out []Criterion
	for _, line := range strings.Split(s, "\n") {
		if m := criterionLine.FindStringSubmatch(strings.TrimRight(line, " \t")); m != nil {
			out = append(out, Criterion{m[1], m[2]})
		}
	}
	return out
}

var commitPattern = regexp.MustCompile(`\b[0-9a-f]{7,40}\b`)

// commitIn returns the first commit hash in text: 7 to 40 hex characters, with
// at least one digit so that words like "defaced" aren't taken for one.
func commitIn(text string) string {
	for _, m := range commitPattern.FindAllString(text, -1) {
		if strings.ContainsAny(m, "0123456789") {
			return m
		}
	}
	return ""
}

func sameCommit(a, b string) bool {
	return a != "" && b != "" && (strings.HasPrefix(a, b) || strings.HasPrefix(b, a))
}

// BuildRevision is the commit the latest change entry names: what verify must
// test. A fix made during verify is a new change entry, so earlier evidence no
// longer counts.
func (it *Item) BuildRevision() string {
	for i := len(it.Log) - 1; i >= 0; i-- {
		if it.Log[i].Type == "change" {
			return commitIn(it.Log[i].Text)
		}
	}
	return ""
}

type actorKey struct{ name, session string }

// Builders are who logged changes on the item.
func (it *Item) Builders() []Actor {
	var out []Actor
	seen := map[actorKey]bool{}
	for _, e := range it.Log {
		k := actorKey{e.By, e.Session}
		if e.Type == "change" && !seen[k] {
			seen[k] = true
			out = append(out, Actor{Name: e.By, Agent: e.Agent, Model: e.Model, Session: e.Session})
		}
	}
	return out
}

// independent reports whether evidence came from someone other than the
// builders: a different name, or the same name in a different session. Both
// are what the actors reported; see ROADMAP.md, phase 5.
func (it *Item) independent(e *record.Doc) bool {
	for _, b := range it.Builders() {
		if b.Name == e.Get("by") && b.Session == e.Get("session") {
			return false
		}
	}
	return true
}

// evidenceCriteria returns the criteria numbers an evidence record names.
func evidenceCriteria(e *record.Doc) []string {
	var out []string
	for _, c := range strings.Split(e.Get("criteria"), ",") {
		if c = strings.TrimSpace(c); c != "" {
			out = append(out, c)
		}
	}
	return out
}

// OnBuild returns the evidence recorded on the build revision, oldest first.
func (it *Item) OnBuild() []*record.Doc {
	rev := it.BuildRevision()
	var out []*record.Doc
	for _, e := range it.Evidence {
		if sameCommit(commitIn(e.Get("revision")), rev) {
			out = append(out, e)
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Get("at") < out[j].Get("at") })
	return out
}

// unverified lists what keeps the item's verification from holding. The
// latest evidence for a criterion, or for a claim, is what counts: a failure
// recorded after a pass reopens it.
func (it *Item) unverified(independent bool) []string {
	rev := it.BuildRevision()
	if rev == "" {
		return []string{fmt.Sprintf("log: the latest change entry names no commit (log it as `ignition log %s --type change \"<sha> <what>\"`)", it.ID)}
	}
	qa, hint := "evidence", ""
	if independent {
		qa = "QA evidence"
		hint = fmt.Sprintf("; QA is someone other than the builder, such as a QA agent given `ignition qa %s`", it.ID)
	}
	var (
		out       []string
		order     []string
		latest    = map[string]*record.Doc{} // by criterion or claim
		latestQA  = map[string]*record.Doc{} // the same, independent evidence only
		anyPassQA bool
	)
	for _, e := range it.OnBuild() {
		keys := evidenceCriteria(e)
		for i := range keys {
			keys[i] = "criterion " + keys[i]
		}
		if len(keys) == 0 {
			keys = []string{fmt.Sprintf("%q", e.Get("claim"))}
		}
		indep := !independent || it.independent(e)
		for _, k := range keys {
			if latest[k] == nil {
				order = append(order, k)
			}
			latest[k] = e
			if indep {
				latestQA[k] = e
			}
		}
		if indep && e.Get("result") == "pass" {
			anyPassQA = true
		}
	}
	for _, k := range order {
		if e := latest[k]; e.Get("result") == "fail" {
			out = append(out, fmt.Sprintf("%s failed on %s (%s); fix it, log the change, and verify again", k, rev, e.Get("id")))
		}
	}
	for _, c := range it.Criteria() {
		k := "criterion " + c.Num
		switch e := latestQA[k]; {
		case latest[k] != nil && latest[k].Get("result") == "fail":
			// reported above
		case e == nil:
			out = append(out, fmt.Sprintf("%s has no %s on %s (record it with `--criteria %s`)", k, qa, rev, c.Num))
		case e.Get("result") != "pass":
			out = append(out, fmt.Sprintf("%s is %s on %s (%s); it needs a pass", k, e.Get("result"), rev, e.Get("id")))
		}
	}
	if !anyPassQA {
		out = append(out, fmt.Sprintf("no passing %s on %s, the revision the latest change entry names%s", qa, rev, hint))
	}
	return out
}

// normalizeCriteria validates and normalizes an evidence record's criteria: a comma
// separated list of the item's criterion numbers.
func (it *Item) normalizeCriteria(s string) (string, error) {
	if strings.TrimSpace(s) == "" {
		return "", nil
	}
	have := it.Criteria()
	var nums []string
	for _, c := range strings.Split(s, ",") {
		c = strings.TrimSpace(c)
		if !slices.ContainsFunc(have, func(h Criterion) bool { return h.Num == c }) {
			if len(have) == 0 {
				return "", fmt.Errorf("criteria: %s has no numbered acceptance criteria", it.ID)
			}
			return "", fmt.Errorf("criteria: %q isn't one of %s's acceptance criteria (%s)", c, it.ID, criteriaNums(have))
		}
		if !slices.Contains(nums, c) {
			nums = append(nums, c)
		}
	}
	return strings.Join(nums, ","), nil
}

func criteriaNums(cs []Criterion) string {
	var out []string
	for _, c := range cs {
		out = append(out, c.Num)
	}
	return strings.Join(out, ", ")
}
