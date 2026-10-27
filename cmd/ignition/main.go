// Command ignition runs Ignition's work model over a product's records.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mohsalsaleem/ignition-engineering/internal/attach"
	"github.com/mohsalsaleem/ignition-engineering/internal/model"
	"github.com/mohsalsaleem/ignition-engineering/internal/workspace"
	"github.com/mohsalsaleem/ignition-engineering/skill"
)

var version = "0.7.0"

const usage = `ignition: structured engineering work for people and agents.

Usage: ignition <command> [arguments] [flags]

Workspace
  init [--product NAME]             create the workspace at <repo>/.ignition, tracked with the code
                                    (outside a Git repository: in the current folder)
  attach --app DIR --records-url URL
                                    instead keep records in a separate repository, checked out at
                                    DIR/.ignition (for public or multi-repository products)
Work
  new KIND "TITLE" [--risk R] [--size S] [--parent ID]
                                    KIND: initiative feature bug task investigation incident
  brief ID                          where an item stands and what it needs next
  status                            every item, its stage and what it waits on
  advance ID [--owner NAME --decision "TEXT"]
                                    move to the next stage once this one is complete
  log ID --type TYPE "TEXT"         TYPE: note finding decision change handoff scope
  evidence ID --claim C --result pass|fail|inconclusive --revision R --environment E --method M [--observed O] [--limits L]
  decision "TITLE"                  create a product-wide decision record
  check [ID]                        report problems in the records; exits 1 if any
  attribution [ID]                  who built what: people, agents and models, from the log
  whoami                            who Ignition will record as acting
Help
  guide                             how an agent works with Ignition (the skill's instructions)
  skill install [--for all|claude|agents] [--project DIR] [--to DIR]
                                    install the agent skill for Claude Code (.claude/skills)
                                    and other agents (.agents/skills); default: both, in ~
  skill print                       print the skill's SKILL.md
  version

Flags for every command:
  --dir DIR       the workspace (default: found from the current folder, or $IGNITION_DIR)
  --by NAME       who is acting (default: $IGNITION_ACTOR, then the agent, then git user.name)
  --agent TOOL    the agent tool, e.g. claude-code, codex, cursor, gemini-cli
                  (default: $IGNITION_AGENT, then detected from the environment)
  --model ID      the exact model ID, e.g. claude-opus-5-5 (default: $IGNITION_MODEL)
  --session ID    the agent's session ID (default: $IGNITION_SESSION, then detected)
`

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

type cli struct {
	out, errw io.Writer
	dir, by   string
	agent     string
	model     string
	session   string
}

// parse reads flags wherever they appear among the positional arguments.
func parse(fs *flag.FlagSet, args []string) ([]string, error) {
	var pos []string
	for {
		if err := fs.Parse(args); err != nil {
			return nil, err
		}
		args = fs.Args()
		if len(args) == 0 {
			return pos, nil
		}
		pos = append(pos, args[0])
		args = args[1:]
	}
}

func run(args []string, out, errw io.Writer) int {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		fmt.Fprint(out, usage)
		return 0
	}
	c := &cli{out: out, errw: errw}
	cmd, rest := args[0], args[1:]
	fs := flag.NewFlagSet(cmd, flag.ContinueOnError)
	fs.SetOutput(errw)
	fs.StringVar(&c.dir, "dir", os.Getenv("IGNITION_DIR"), "workspace folder")
	fs.StringVar(&c.by, "by", "", "who is acting")
	fs.StringVar(&c.agent, "agent", "", "the agent tool")
	fs.StringVar(&c.model, "model", "", "the exact model ID")
	fs.StringVar(&c.session, "session", "", "the agent's session ID")
	var err error
	switch cmd {
	case "init":
		err = c.init(fs, rest)
	case "attach":
		err = c.attach(fs, rest)
	case "new":
		err = c.newItem(fs, rest)
	case "brief":
		err = c.brief(fs, rest)
	case "status":
		err = c.status(fs, rest)
	case "advance":
		err = c.advance(fs, rest)
	case "log":
		err = c.log(fs, rest)
	case "evidence":
		err = c.evidence(fs, rest)
	case "decision":
		err = c.decision(fs, rest)
	case "check":
		err = c.check(fs, rest)
	case "attribution":
		err = c.attribution(fs, rest)
	case "whoami":
		err = c.whoami(fs, rest)
	case "guide":
		fmt.Fprint(out, skill.Body())
	case "skill":
		err = c.skill(fs, rest)
	case "version":
		fmt.Fprintln(out, "ignition", version)
	default:
		fmt.Fprintf(errw, "ignition: unknown command %q\n\n%s", cmd, usage)
		return 2
	}
	if errors.Is(err, flag.ErrHelp) {
		return 0
	}
	if errors.Is(err, errProblems) {
		return 1
	}
	if err != nil {
		fmt.Fprintln(errw, "ignition:", err)
		return 1
	}
	return 0
}

var errProblems = errors.New("problems found")

func (c *cli) workspace() (*workspace.Workspace, error) {
	if c.dir != "" {
		return workspace.Open(c.dir), nil
	}
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	root, err := workspace.Find(wd)
	if err != nil {
		return nil, err
	}
	return workspace.Open(root), nil
}

func want(pos []string, n int, what string) error {
	if len(pos) != n {
		return fmt.Errorf("want %s", what)
	}
	return nil
}

func (c *cli) init(fs *flag.FlagSet, args []string) error {
	product := fs.String("product", "", "product name")
	pos, err := parse(fs, args)
	if err != nil {
		return err
	}
	if err := want(pos, 0, "no arguments"); err != nil {
		return err
	}
	// By default the records live in the repository being worked on, at
	// <repo>/.ignition, tracked with the code. --dir puts them anywhere else.
	root, app := c.dir, ""
	if root == "" {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		root = wd
		if top := workspace.RepoRoot(wd); top != "" {
			if filepath.Base(top) == ".ignition" {
				root = top // a separate records repository, checked out by attach
			} else {
				root, app = filepath.Join(top, ".ignition"), top
			}
		}
	}
	if *product == "" {
		*product = filepath.Base(root)
		if app != "" {
			*product = filepath.Base(app)
		}
	}
	if err := workspace.Init(root, *product); err != nil {
		return err
	}
	fmt.Fprintf(c.out, "Created an Ignition workspace in %s.\n", root)
	if app != "" {
		notes, err := workspace.PrepareRepo(app)
		for _, n := range notes {
			fmt.Fprintln(c.out, n)
		}
		if err != nil {
			return err
		}
		fmt.Fprintln(c.out, "The records are part of this repository: commit .ignition/ with the changes it describes.")
	}
	fmt.Fprintln(c.out, "Fill in PROJECT.md, then `ignition new`.")
	return nil
}

func (c *cli) attach(fs *flag.FlagSet, args []string) error {
	app := fs.String("app", "", "application repository root")
	url := fs.String("records-url", "", "the product's records repository")
	if _, err := parse(fs, args); err != nil {
		return err
	}
	if *app == "" || *url == "" {
		return errors.New("attach needs --app and --records-url")
	}
	res, err := attach.Attach(*app, *url)
	if err != nil {
		return err
	}
	fmt.Fprintf(c.out, "Records %s at %s. The application's Git status is unchanged.\n", res.Outcome, res.Path)
	for _, w := range res.Warnings {
		fmt.Fprintln(c.out, "Warning:", w)
	}
	return nil
}

func (c *cli) newItem(fs *flag.FlagSet, args []string) error {
	risk := fs.String("risk", "", "low, normal or high")
	size := fs.String("size", "", "small, medium or large")
	parent := fs.String("parent", "", "the initiative this belongs to")
	pos, err := parse(fs, args)
	if err != nil {
		return err
	}
	if err := want(pos, 2, `KIND "TITLE"`); err != nil {
		return err
	}
	kind, err := model.ParseKind(pos[0])
	if err != nil {
		return err
	}
	w, err := c.workspace()
	if err != nil {
		return err
	}
	it, err := w.NewItem(kind, pos[1], *risk, *size, *parent, c.actingAs())
	if err != nil {
		return err
	}
	fmt.Fprintf(c.out, "Created %s (%s) in %s.\n", it.ID, kind, it.Dir)
	return c.printBrief(w, it.ID)
}

func (c *cli) oneID(fs *flag.FlagSet, args []string) (*workspace.Workspace, string, error) {
	pos, err := parse(fs, args)
	if err != nil {
		return nil, "", err
	}
	if err := want(pos, 1, "a work item ID"); err != nil {
		return nil, "", err
	}
	w, err := c.workspace()
	return w, pos[0], err
}

func (c *cli) brief(fs *flag.FlagSet, args []string) error {
	w, id, err := c.oneID(fs, args)
	if err != nil {
		return err
	}
	return c.printBrief(w, id)
}

func (c *cli) printBrief(w *workspace.Workspace, id string) error {
	it, err := w.Load(id)
	if err != nil {
		return err
	}
	f := it.Facts()
	risk, size := it.Brief.Get("risk"), it.Brief.Get("size")
	if risk == "" {
		risk = "unset"
	}
	if size == "" {
		size = "unset"
	}
	fmt.Fprintf(c.out, "\n%s · %s · %s risk · %s\n%s\n", it.ID, f.Kind, risk, size, it.Title())
	if p := it.Brief.Get("parent"); p != "" {
		fmt.Fprintf(c.out, "Part of %s\n", p)
	}
	var route []string
	for _, s := range model.Route(f) {
		if s == it.Stage() {
			route = append(route, "["+string(s)+"]")
		} else {
			route = append(route, string(s))
		}
	}
	fmt.Fprintf(c.out, "\nStage: %s\n", it.Stage())
	fmt.Fprintf(c.out, "Route: %s\n", strings.Join(route, " → "))
	if probs := it.Check(); len(probs) > 0 {
		fmt.Fprintln(c.out, "\nProblems to fix first:")
		for _, p := range probs {
			fmt.Fprintln(c.out, "  -", p)
		}
	}
	if it.Stage() != model.Closed {
		if unmet := it.Unmet(it.Stage()); len(unmet) > 0 {
			fmt.Fprintf(c.out, "\nTo finish %s:\n", it.Stage())
			for _, u := range unmet {
				fmt.Fprintln(c.out, "  -", u)
			}
		} else {
			fmt.Fprintf(c.out, "\n%s is complete.\n", it.Stage())
		}
		if need, q := model.OwnerDecides(f, it.Stage()); need {
			fmt.Fprintf(c.out, "Leaving %s needs the owner's decision: %s\n", it.Stage(), q)
		} else {
			fmt.Fprintf(c.out, "Leaving %s needs no owner decision; run `ignition advance %s` when complete.\n", it.Stage(), it.ID)
		}
	}
	var files []string
	for _, a := range []model.Artifact{model.BriefDoc, model.DesignDoc, model.PlanDoc, model.ReleaseDoc} {
		if _, ok := it.Docs[a]; ok {
			files = append(files, string(a))
		}
	}
	files = append(files, "log.md")
	if n := len(it.Evidence); n > 0 {
		files = append(files, fmt.Sprintf("evidence/ (%d)", n))
	}
	fmt.Fprintf(c.out, "\nFiles in %s: %s\n", it.Dir, strings.Join(files, ", "))
	if len(it.Children) > 0 {
		fmt.Fprintln(c.out, "\nChild items:")
		for _, ch := range it.Children {
			fmt.Fprintf(c.out, "  %s  %-8s %-10s %s\n", ch.ID, ch.Kind(), ch.Stage(), ch.Title())
		}
	}
	if n := len(it.Log); n > 0 {
		fmt.Fprintln(c.out, "\nRecent log:")
		for _, e := range it.Log[max(0, n-5):] {
			first, _, _ := strings.Cut(e.Text, "\n")
			fmt.Fprintf(c.out, "  %s  %-10s %s\n", e.At, e.Type, first)
		}
	}
	return nil
}

func (c *cli) status(fs *flag.FlagSet, args []string) error {
	pos, err := parse(fs, args)
	if err != nil {
		return err
	}
	if err := want(pos, 0, "no arguments"); err != nil {
		return err
	}
	w, err := c.workspace()
	if err != nil {
		return err
	}
	items, err := w.List()
	if err != nil {
		return err
	}
	if len(items) == 0 {
		fmt.Fprintln(c.out, "No work items. Create one with `ignition new`.")
		return nil
	}
	for _, it := range items {
		waits := ""
		full, err := w.Load(it.ID)
		switch {
		case err != nil:
			waits = "unreadable: " + err.Error()
		case it.Stage() == model.Closed:
			waits = ""
		case len(full.Check()) > 0:
			waits = "has problems (ignition check " + it.ID + ")"
		case len(full.Unmet(full.Stage())) > 0:
			waits = "in progress"
		default:
			if need, _ := model.OwnerDecides(full.Facts(), full.Stage()); need {
				waits = "waits for the owner"
			} else {
				waits = "ready to advance"
			}
		}
		fmt.Fprintf(c.out, "%s  %-13s %-6s %-10s %-24s %s\n", it.ID, it.Kind(), it.Brief.Get("risk"), it.Stage(), waits, it.Title())
	}
	return nil
}

func (c *cli) advance(fs *flag.FlagSet, args []string) error {
	owner := fs.String("owner", "", "who made the checkpoint decision")
	decision := fs.String("decision", "", "the owner's decision, in their words")
	w, id, err := c.oneID(fs, args)
	if err != nil {
		return err
	}
	var cp *workspace.Checkpoint
	if *owner != "" || *decision != "" {
		cp = &workspace.Checkpoint{Owner: *owner, Decision: *decision}
	}
	next, err := w.Advance(id, cp, c.actingAs())
	if err != nil {
		return err
	}
	fmt.Fprintf(c.out, "%s is now in %s.\n", id, next)
	return c.printBrief(w, id)
}

func (c *cli) log(fs *flag.FlagSet, args []string) error {
	typ := fs.String("type", "note", "note, finding, decision, change, handoff or scope")
	pos, err := parse(fs, args)
	if err != nil {
		return err
	}
	if err := want(pos, 2, `ID "TEXT"`); err != nil {
		return err
	}
	w, err := c.workspace()
	if err != nil {
		return err
	}
	if err := w.Log(pos[0], *typ, c.actingAs(), pos[1]); err != nil {
		return err
	}
	fmt.Fprintf(c.out, "Logged a %s entry on %s.\n", *typ, pos[0])
	return nil
}

func (c *cli) evidence(fs *flag.FlagSet, args []string) error {
	var in workspace.EvidenceInput
	fs.StringVar(&in.Claim, "claim", "", "what was checked")
	fs.StringVar(&in.Result, "result", "", "pass, fail or inconclusive")
	fs.StringVar(&in.Revision, "revision", "", "commit, tag, digest or deployment ID")
	fs.StringVar(&in.Environment, "environment", "", "where it was checked")
	fs.StringVar(&in.Method, "method", "", "command, test or steps")
	fs.StringVar(&in.Observed, "observed", "", "what was seen")
	fs.StringVar(&in.Limits, "limits", "", "what this does not show")
	w, id, err := c.oneID(fs, args)
	if err != nil {
		return err
	}
	in.By = c.actingAs()
	eid, err := w.AddEvidence(id, in)
	if err != nil {
		return err
	}
	fmt.Fprintf(c.out, "Recorded evidence %s on %s.\n", eid, id)
	return nil
}

func (c *cli) decision(fs *flag.FlagSet, args []string) error {
	pos, err := parse(fs, args)
	if err != nil {
		return err
	}
	if err := want(pos, 1, `"TITLE"`); err != nil {
		return err
	}
	w, err := c.workspace()
	if err != nil {
		return err
	}
	path, err := w.NewDecision(pos[0])
	if err != nil {
		return err
	}
	fmt.Fprintf(c.out, "Created %s.\n", path)
	return nil
}

func (c *cli) check(fs *flag.FlagSet, args []string) error {
	pos, err := parse(fs, args)
	if err != nil {
		return err
	}
	w, err := c.workspace()
	if err != nil {
		return err
	}
	ids := pos
	if len(ids) == 0 {
		items, err := w.List()
		if err != nil {
			return err
		}
		for _, it := range items {
			ids = append(ids, it.ID)
		}
	}
	found := 0
	for _, id := range ids {
		it, err := w.Load(id)
		if err != nil {
			fmt.Fprintf(c.out, "%s: %v\n", id, err)
			found++
			continue
		}
		for _, p := range it.Check() {
			fmt.Fprintf(c.out, "%s: %s\n", id, p)
			found++
		}
	}
	if found > 0 {
		fmt.Fprintf(c.out, "%d problem(s).\n", found)
		return errProblems
	}
	fmt.Fprintf(c.out, "%d item(s), no problems.\n", len(ids))
	return nil
}

func (c *cli) skill(fs *flag.FlagSet, args []string) error {
	target := fs.String("for", "all", "all, claude or agents")
	project := fs.String("project", "", "install into this project instead of the home folder")
	to := fs.String("to", "", "install into this folder")
	pos, err := parse(fs, args)
	if err != nil {
		return err
	}
	if len(pos) != 1 || (pos[0] != "install" && pos[0] != "print") {
		return errors.New("want `skill install [--for all|claude|agents] [--project DIR] [--to DIR]` or `skill print`")
	}
	if pos[0] == "print" {
		files, err := skill.Files()
		if err != nil {
			return err
		}
		fmt.Fprint(c.out, string(files["SKILL.md"]))
		return nil
	}
	var dirs []string
	switch {
	case *to != "":
		dir := *to
		if filepath.Base(dir) != skill.Name {
			dir = filepath.Join(dir, skill.Name) // a skill's folder must be named after it
		}
		dirs = []string{dir}
	default:
		base := *project
		if base == "" {
			if base, err = os.UserHomeDir(); err != nil {
				return err
			}
		}
		if dirs, err = skill.Targets(*target, base); err != nil {
			return err
		}
	}
	for _, dir := range dirs {
		if err := skill.Install(dir); err != nil {
			return err
		}
		fmt.Fprintf(c.out, "Installed the Ignition skill (ignition %s) in %s.\n", version, dir)
	}
	return nil
}
