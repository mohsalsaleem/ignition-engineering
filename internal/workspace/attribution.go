package workspace

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

// Contribution is what one actor, as one agent and model, did on an item.
type Contribution struct {
	Actor    Actor          // Session is left empty; see Sessions
	Entries  map[string]int // log entries by type
	Commits  []string       // the revisions their change entries name
	Sessions []string
	First    string
	Last     string
}

// Decision is an owner checkpoint decision and who recorded it.
type Decision struct {
	Owner, Leaving, RecordedBy string
}

// Attribution lists who contributed to an item, from its log, and the owner
// decisions recorded on it. Agent, model and session are what the actor
// reported; until runs are bound to verified identities they are claims.
func (it *Item) Attribution() ([]Contribution, []Decision) {
	byKey := map[Actor]*Contribution{}
	var order []Actor
	var decisions []Decision
	for _, e := range it.Log {
		key := Actor{Name: e.By, Agent: e.Agent, Model: e.Model}
		c, ok := byKey[key]
		if !ok {
			c = &Contribution{Actor: key, Entries: map[string]int{}, First: e.At}
			byKey[key] = c
			order = append(order, key)
		}
		c.Entries[e.Type]++
		c.Last = e.At
		if e.Session != "" && !slices.Contains(c.Sessions, e.Session) {
			c.Sessions = append(c.Sessions, e.Session)
		}
		switch e.Type {
		case "change":
			if f := strings.Fields(e.Text); len(f) > 0 {
				c.Commits = append(c.Commits, f[0])
			}
		case "checkpoint":
			d := Decision{RecordedBy: e.By}
			for _, line := range strings.Split(e.Text, "\n") {
				if v, ok := strings.CutPrefix(line, "Owner: "); ok {
					d.Owner = v
				}
				if v, ok := strings.CutPrefix(line, "Leaving "); ok {
					d.Leaving, _, _ = strings.Cut(v, ":")
				}
			}
			decisions = append(decisions, d)
		}
	}
	out := make([]Contribution, 0, len(order))
	for _, k := range order {
		out = append(out, *byKey[k])
	}
	return out, decisions
}

// EntrySummary renders entry counts like "2 change, 4 evidence".
func (c Contribution) EntrySummary() string {
	var types []string
	for t := range c.Entries {
		types = append(types, t)
	}
	sort.Strings(types)
	var parts []string
	for _, t := range types {
		parts = append(parts, fmt.Sprintf("%d %s", c.Entries[t], t))
	}
	return strings.Join(parts, ", ")
}
