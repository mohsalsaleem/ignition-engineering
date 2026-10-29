package model

import (
	"slices"
	"testing"
)

func TestRoutes(t *testing.T) {
	cases := []struct {
		f    Facts
		want []Stage
	}{
		{Facts{Kind: Feature, Risk: Normal}, []Stage{Intake, Shape, Design, Plan, Build, Verify, Release, Observe, Closed}},
		{Facts{Kind: Feature, Risk: Low}, []Stage{Intake, Shape, Plan, Build, Verify, Release, Observe, Closed}},
		{Facts{Kind: Task, Risk: Low, NoRelease: true}, []Stage{Intake, Plan, Build, Verify, Closed}},
		{Facts{Kind: Bug, Risk: High}, []Stage{Intake, Reproduce, Diagnose, Build, Verify, Release, Closed}},
		{Facts{Kind: Incident}, []Stage{Contain, Timeline, Recover, Review, Closed}},
	}
	for _, c := range cases {
		if got := Route(c.f); !slices.Equal(got, c.want) {
			t.Errorf("Route(%+v) = %v, want %v", c.f, got, c.want)
		}
	}
}

func TestUnclassifiedRiskNeverSkipsTheOwner(t *testing.T) {
	f := Facts{Kind: Feature}
	if ok, _ := OwnerDecides(f, Design); !ok {
		t.Fatal("an item without a risk skipped the design decision")
	}
	if !OnRoute(f, Design) {
		t.Fatal("an item without a risk skipped design")
	}
}

func TestOwnerDecisions(t *testing.T) {
	cases := []struct {
		f     Facts
		stage Stage
		want  bool
	}{
		{Facts{Kind: Feature, Risk: Low}, Shape, false},
		{Facts{Kind: Feature, Risk: Normal}, Shape, true},
		{Facts{Kind: Feature, Risk: Normal}, Design, false},
		{Facts{Kind: Feature, Risk: High}, Design, true},
		{Facts{Kind: Initiative, Risk: Low}, Design, false},
		{Facts{Kind: Initiative, Risk: Normal}, Design, true},
		{Facts{Kind: Bug, Risk: Normal}, Diagnose, false},
		{Facts{Kind: Bug, Risk: High}, Diagnose, true},
		{Facts{Kind: Bug, Risk: Normal}, Verify, true},
		{Facts{Kind: Task, Risk: Normal, NoRelease: true}, Verify, false},
		{Facts{Kind: Task, Risk: Low}, Verify, false},
		{Facts{Kind: Incident, Risk: Low}, Review, true},
	}
	for _, c := range cases {
		if got, _ := OwnerDecides(c.f, c.stage); got != c.want {
			t.Errorf("OwnerDecides(%+v, %s) = %v, want %v", c.f, c.stage, got, c.want)
		}
	}
}

func TestEveryStageHasExitRules(t *testing.T) {
	for _, k := range Kinds {
		for _, s := range Route(Facts{Kind: k}) {
			if s != Closed && len(Exit(Facts{Kind: k}, s)) == 0 {
				t.Errorf("%s has no exit rules for %s", k, s)
			}
		}
	}
}

func designSections(f Facts) []string {
	for _, r := range Exit(f, Design) {
		if r.Kind == Sections && r.Artifact == DesignDoc {
			return r.Names
		}
	}
	return nil
}

func TestDesignDepthFollowsSizeAndRisk(t *testing.T) {
	cases := []struct {
		f       Facts
		section string
		want    bool
	}{
		{Facts{Kind: Feature, Risk: Normal, Size: "small"}, "Data model", false},
		{Facts{Kind: Feature, Risk: Normal, Size: "medium"}, "Data model", true},
		{Facts{Kind: Feature, Risk: Normal, Size: "large"}, "Data model", true},
		{Facts{Kind: Feature, Risk: High, Size: "small"}, "Errors", true},
		{Facts{Kind: Feature, Risk: Normal}, "Interfaces", true},
		{Facts{Kind: Feature, Risk: Normal, Size: "small"}, "Parent design", false},
		{Facts{Kind: Feature, Risk: Normal, Size: "small", Child: true}, "Parent design", true},
		{Facts{Kind: Initiative, Risk: Normal, Size: "large"}, "Architecture", true},
		{Facts{Kind: Initiative, Risk: Normal, Size: "large"}, "Data model", true},
		{Facts{Kind: Initiative, Risk: Normal, Size: "large"}, "Flows", true},
		{Facts{Kind: Initiative, Risk: Normal, Size: "large"}, "Errors", false},
	}
	for _, c := range cases {
		if got := slices.Contains(designSections(c.f), c.section); got != c.want {
			t.Errorf("%+v requires %q = %v, want %v", c.f, c.section, got, c.want)
		}
	}
}
