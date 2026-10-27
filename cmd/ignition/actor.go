package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mohsalsaleem/ignition-engineering/internal/workspace"
)

// getenv is replaced in tests, which run inside agents too.
var getenv = os.Getenv

// detectAgent names the agent tool running Ignition, from variables the tools
// set themselves. It returns "" when no agent is recognised.
func detectAgent() string {
	if a := getenv("AI_AGENT"); a != "" {
		name, _, _ := strings.Cut(a, "_") // e.g. claude-code_2-1-304_agent
		return name
	}
	if getenv("CLAUDECODE") == "1" {
		return "claude-code"
	}
	return ""
}

// actingAs resolves who is acting: flags first, then IGNITION_* variables,
// then what the environment shows. A model is never guessed.
func (c *cli) actingAs() workspace.Actor {
	a := workspace.Actor{Agent: first(c.agent, getenv("IGNITION_AGENT"), detectAgent())}
	a.Model = first(c.model, getenv("IGNITION_MODEL"))
	a.Session = first(c.session, getenv("IGNITION_SESSION"))
	if a.Session == "" && a.Agent == "claude-code" {
		a.Session = getenv("CLAUDE_CODE_SESSION_ID")
	}
	a.Name = first(c.by, getenv("IGNITION_ACTOR"), a.Agent)
	if a.Name == "" {
		if out, err := exec.Command("git", "config", "user.name").Output(); err == nil {
			a.Name = strings.TrimSpace(string(out))
		}
	}
	if a.Agent != "" && a.Model == "" {
		fmt.Fprintf(c.errw, "Warning: recording %s with an unknown model. Pass --model <exact model ID> or set IGNITION_MODEL.\n", a.Agent)
	}
	return a
}

func first(vals ...string) string {
	for _, v := range vals {
		if v = strings.TrimSpace(v); v != "" {
			return v
		}
	}
	return ""
}

func (c *cli) whoami(fs *flag.FlagSet, args []string) error {
	pos, err := parse(fs, args)
	if err != nil {
		return err
	}
	if err := want(pos, 0, "no arguments"); err != nil {
		return err
	}
	a := c.actingAs()
	fmt.Fprintf(c.out, "name:    %s\nagent:   %s\nmodel:   %s\nsession: %s\n",
		orUnknown(a.Name), orNone(a.Agent), orUnknown(a.Model), orNone(a.Session))
	return nil
}

func orUnknown(s string) string {
	if s == "" {
		return "unknown"
	}
	return s
}

func orNone(s string) string {
	if s == "" {
		return "none"
	}
	return s
}

func (c *cli) attribution(fs *flag.FlagSet, args []string) error {
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
	type key struct{ name, agent, model string }
	totals := map[key]int{}
	var order []key
	for _, id := range ids {
		it, err := w.Load(id)
		if err != nil {
			return err
		}
		contribs, decisions := it.Attribution()
		fmt.Fprintf(c.out, "%s  %s\n", it.ID, it.Title())
		for _, ct := range contribs {
			who := orUnknown(ct.Actor.Name)
			if ct.Actor.Agent != "" || ct.Actor.Model != "" {
				who += fmt.Sprintf(" · %s · %s", orNone(ct.Actor.Agent), orUnknown(ct.Actor.Model))
			}
			fmt.Fprintf(c.out, "  %s\n    %s\n", who, ct.EntrySummary())
			if len(ct.Commits) > 0 {
				fmt.Fprintf(c.out, "    commits: %s\n", strings.Join(ct.Commits, ", "))
			}
			if len(ct.Sessions) > 0 {
				fmt.Fprintf(c.out, "    sessions: %s\n", strings.Join(ct.Sessions, ", "))
			}
			k := key{ct.Actor.Name, ct.Actor.Agent, ct.Actor.Model}
			if _, ok := totals[k]; !ok {
				order = append(order, k)
			}
			for _, n := range ct.Entries {
				totals[k] += n
			}
		}
		for _, d := range decisions {
			fmt.Fprintf(c.out, "  owner decision leaving %s: %s (recorded by %s)\n", d.Leaving, orUnknown(d.Owner), d.RecordedBy)
		}
	}
	if len(ids) > 1 {
		fmt.Fprintln(c.out, "\nAcross these items:")
		for _, k := range order {
			fmt.Fprintf(c.out, "  %-24s %-12s %-20s %d entries\n", orUnknown(k.name), orNone(k.agent), orUnknown(k.model), totals[k])
		}
	}
	fmt.Fprintln(c.out, "\nAgent, model and session are as each actor reported them.")
	return nil
}
