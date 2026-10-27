package record

import (
	"strings"
	"testing"
)

const brief = `---
id: W-20261027-k3f9q2
kind: bug
title: "Install ships a local image: not the release"
stage: intake
---
# Install ships a local image

## Request

<!-- What was asked, in the requester's words. -->
The rc.1 install was refused.

## Root cause

<!-- Why it happens. -->
`

func TestParseAndRender(t *testing.T) {
	d, err := Parse(brief)
	if err != nil {
		t.Fatal(err)
	}
	if got := d.Get("title"); got != "Install ships a local image: not the release" {
		t.Fatalf("title = %q", got)
	}
	d.Set("stage", "reproduce")
	d.Set("risk", "normal")
	again, err := Parse(d.String())
	if err != nil {
		t.Fatal(err)
	}
	if again.Get("stage") != "reproduce" || again.Get("risk") != "normal" || again.Get("title") != d.Get("title") {
		t.Fatalf("round trip lost values:\n%s", d.String())
	}
	d.Set("parent", "")
	if !strings.Contains(d.String(), "\nparent:\n") {
		t.Fatalf("an empty value isn't rendered as `key:`:\n%s", d.String())
	}
	if !strings.HasPrefix(again.Body, "# Install ships") {
		t.Fatalf("body changed: %q", again.Body)
	}
}

func TestSectionsIgnoreGuidance(t *testing.T) {
	d, _ := Parse(brief)
	if !d.Filled("Request") {
		t.Error("Request has content but reads as empty")
	}
	if d.Filled("Root cause") {
		t.Error("Root cause has only guidance but reads as filled")
	}
	if _, ok := d.Section("Missing"); ok {
		t.Error("a missing section was found")
	}
}

func TestRejectsBadFrontmatter(t *testing.T) {
	for _, text := range []string{
		"---\nid: a\nid: b\n---\n",
		"---\nnot a pair\n---\n",
		"---\nid: a\n",
	} {
		if _, err := Parse(text); err == nil {
			t.Errorf("accepted %q", text)
		}
	}
}

func TestTasks(t *testing.T) {
	d, _ := Parse("## Tasks\n\n| ID | Task | Acceptance | Status | Evidence |\n| --- | --- | --- | --- | --- |\n| T1 | Fix | Test passes | Done | E-1 |\n|  |  |  |  |  |\n")
	tasks, err := d.Tasks()
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 1 || tasks[0].Status != "done" || tasks[0].Evidence != "E-1" {
		t.Fatalf("tasks = %+v", tasks)
	}
}

func TestLogRoundTrip(t *testing.T) {
	text := "# Log\n" + FormatEntry(Entry{At: "2026-10-27T10:00:00Z", Type: "finding", By: "agent", Text: "Port 80 is open.\nSecond line."}) +
		FormatEntry(Entry{At: "2026-10-27T11:00:00Z", Type: "stage", By: "agent", Text: "intake → reproduce"})
	entries, err := ParseLog(text)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 || entries[0].Text != "Port 80 is open.\nSecond line." || entries[1].Type != "stage" {
		t.Fatalf("entries = %+v", entries)
	}
}

func TestLogEntriesCarryAgentAndModel(t *testing.T) {
	text := FormatEntry(Entry{At: "t1", Type: "change", By: "codex", Agent: "codex", Model: "gpt-x", Text: "abc fix"}) +
		FormatEntry(Entry{At: "t2", Type: "note", By: "Owner", Text: "ok"})
	if !strings.Contains(text, "## t1 · change · codex · codex · gpt-x · -\n") || !strings.Contains(text, "## t2 · note · Owner\n") {
		t.Fatalf("formatted:\n%s", text)
	}
	entries, err := ParseLog(text)
	if err != nil {
		t.Fatal(err)
	}
	if e := entries[0]; e.Agent != "codex" || e.Model != "gpt-x" || e.Session != "" {
		t.Fatalf("entry = %+v", e)
	}
	if e := entries[1]; e.Agent != "" || e.By != "Owner" {
		t.Fatalf("entry = %+v", e)
	}
	if _, err := ParseLog("## t · note · a · b\n"); err == nil {
		t.Fatal("accepted a four-part heading")
	}
}
