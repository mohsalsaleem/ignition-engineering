package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/mohsalsaleem/ignition-engineering/internal/model"
	"github.com/mohsalsaleem/ignition-engineering/internal/workspace"
)

// qa prints the handoff for an independent QA agent: what to test, on which
// revision, against which criteria, and how to record what it finds. The
// builder gives it to a fresh agent as that agent's whole task.
func (c *cli) qa(fs *flag.FlagSet, args []string) error {
	w, id, err := c.oneID(fs, args)
	if err != nil {
		return err
	}
	it, err := w.Load(id)
	if err != nil {
		return err
	}
	f := it.Facts()
	if !model.OnRoute(f, model.Verify) {
		return fmt.Errorf("%s is a %s, which has no verify stage", id, f.Kind)
	}
	rev := it.BuildRevision()
	if rev == "" {
		return fmt.Errorf("%s has no change entry naming a commit yet: there is nothing to test", id)
	}
	p := func(format string, a ...any) { fmt.Fprintf(c.out, format, a...) }

	p("# QA: %s\n\n", it.Title())
	p("You are the QA agent for work item %s (%s, %s risk). You didn't build it.\n", id, f.Kind, orUnknown(it.Brief.Get("risk")))
	p("Your job is to find out whether it works, not to make it pass.\n\n")
	if it.Stage() != model.Verify {
		p("Note: the item is in %s, not verify. Its records may still change.\n\n", it.Stage())
	}
	p("Revision to test: %s (the latest change entry)\n", rev)
	var builders []string
	for _, b := range it.Builders() {
		s := b.Name
		if b.Session != "" {
			s += " (session " + b.Session + ")"
		}
		builders = append(builders, s)
	}
	p("Built by: %s\n", strings.Join(builders, ", "))
	p("Records: %s\n", it.Dir)

	if crit := it.Criteria(); len(crit) > 0 {
		p("\n## Acceptance criteria\n\n")
		for _, cr := range crit {
			p("%s. %s\n", cr.Num, cr.Text)
		}
	} else if name := model.CriteriaSection(f.Kind); name != "" {
		p("\n## %s\n\n%s\n", name, section(it, model.BriefDoc, name))
	}
	for _, s := range []struct {
		doc  model.Artifact
		name string
	}{
		{model.BriefDoc, "Outcome"},
		{model.BriefDoc, "Scope"},
		{model.BriefDoc, "Expected and actual"},
		{model.BriefDoc, "Reproduction"},
		{model.BriefDoc, "Root cause"},
		{model.DesignDoc, "Test strategy"},
	} {
		if text := section(it, s.doc, s.name); text != "" {
			p("\n## %s (%s)\n\n%s\n", s.name, s.doc, text)
		}
	}

	if on := it.OnBuild(); len(on) > 0 {
		p("\n## Evidence already on %s\n\n", rev)
		for _, e := range on {
			crit := e.Get("criteria")
			if crit == "" {
				crit = "-"
			}
			p("- %s %s, criteria %s, by %s: %s\n", e.Get("id"), e.Get("result"), crit, e.Get("by"), e.Get("claim"))
		}
	}
	if unmet := it.Unmet(model.Verify); len(unmet) > 0 {
		p("\n## Verify still needs\n\n")
		for _, u := range unmet {
			p("- %s\n", u)
		}
	}

	p(`
## How to work

Record everything as QA: pass --by qa, --model <your exact model ID>, and
--session <your session ID> if your tool has one, on every ignition command.
`)
	if model.IndependentQA(f) {
		p("Evidence the builder records doesn't count at this risk; yours does.\n")
	}
	p(`
1. Test exactly %[1]s. Check it out in a separate worktree or clone; don't test
   uncommitted changes.
2. Read the criteria, then the code and the existing tests. Decide how you'd
   show each criterion holds, including its failure paths.
3. Prefer automated checks: run the project's test suite, and write the tests
   that are missing (unit, integration, end-to-end, or driving the running
   app). Commit tests you write on their own branch, and say so in --method.
4. Exercise the change the way a user or caller would, beyond the happy path:
   bad input, permissions, empty and large data, concurrency, rollback.
5. Record every check, passing or failing, as soon as you've seen it:

   ignition evidence %[2]s --criteria <N[,N]> --claim "<what you checked>" \
     --result pass|fail|inconclusive --revision %[1]s \
     --environment <where> --method "<command, test or steps>" \
     --observed "<what you saw>" --limits "<what this doesn't show>" --by qa

   Name the criteria each check covers. A check that covers none (a
   regression, a crash) leaves out --criteria.
6. Don't fix the product code, and don't change the brief, design or plan. A
   problem you find is a fail record plus a finding:
   ignition log %[2]s --type finding "<what fails, how to reproduce>" --by qa
   The builder fixes it and logs a new change; QA then runs again on the new
   revision.
7. Finish with a handoff naming what passed, what failed, and what you
   couldn't test:
   ignition log %[2]s --type handoff "verified: ... failed: ... untested: ..." --by qa
`, rev, id)
	return nil
}

func section(it *workspace.Item, doc model.Artifact, name string) string {
	d := it.Docs[doc]
	if d == nil {
		return ""
	}
	s, _ := d.Section(name)
	return s
}
