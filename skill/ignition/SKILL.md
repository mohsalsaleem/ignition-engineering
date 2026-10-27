---
name: ignition
description: Do engineering work through Ignition's staged work model. Work items (initiative, feature, bug, task, investigation, incident) move through stages such as shape, design, plan, build, verify and release. Each stage leaves written records and evidence, and the owner decides at checkpoints. Use when a repository has an `.ignition/` folder, when the user mentions Ignition or a work item ID like W-20261027-k3f9q2, or when the user asks to start, continue or hand off a task, feature, bug, initiative, investigation or incident under Ignition.
compatibility: Requires the ignition CLI and git. Works with any agent that reads Agent Skills, including Claude Code, Codex, Cursor, Gemini CLI and GitHub Copilot.
---

# Ignition

Ignition keeps a product's engineering records: work items, designs, plans,
evidence, releases and decisions. They live in a separate Git repository
checked out at `<app>/.ignition/`. The `ignition` CLI moves each item through
its stages. It refuses to advance while a stage's record is incomplete, and it
stops when the owner must decide.

Chat is working memory and the records are what last. Write as you go, so
that a new session can continue from the records alone.

## Before you start

1. Run `ignition version`. If the command is missing, tell the user and suggest
   `go install github.com/mohsalsaleem/ignition-engineering/cmd/ignition@latest`.
2. Run `ignition status` from the application repository. If it reports no
   workspace, ask the user whether to set one up. Setup means creating a
   private records repository and running `ignition attach`. Don't create
   repositories or records anywhere on your own.
3. Read the workspace's `PROJECT.md` and `AGENTS.md`. Their rules, commands
   and boundaries apply to everything below.
4. Identify yourself, so the records show which agent and model did the work.
   Run `ignition whoami --model <your exact model ID>`. If it doesn't show
   your agent tool, add `--agent <tool>` (claude-code, codex, cursor,
   gemini-cli, copilot, ...).

   Then pass the same `--model` (and `--agent` if needed) on **every**
   ignition command that writes. Shell variables don't last between commands
   in many agents, so use the flags, or prefix each command with
   `IGNITION_MODEL=<id>`.

   If you don't know your exact model ID, leave it out. It's then recorded as
   unknown and the CLI warns. Never guess it, and never record yourself as
   another agent or person.

## Starting from a request

If the user gives a work item ID, run `ignition brief <ID>` and continue
from the stage it reports.

Otherwise, create the item:

```sh
ignition new <kind> "<one-line title>" --risk <low|normal|high> --size <small|medium|large>
```

| Kind | When |
| --- | --- |
| initiative | An outcome that needs several changes ("build an ATS"). It never goes straight to code. |
| feature | New observable behaviour |
| bug | Existing behaviour is wrong |
| task | Bounded work with a clear done condition: operations, refactoring, docs, upgrades |
| investigation | A question must be answered before anyone can decide |
| incident | Something broke or nearly broke in a live system |

**Risk** decides which stages apply and where the owner decides. Use **high**
for auth, permissions, personal data, payments, AI decisions, schema or data
migrations, infrastructure, and release or install tooling. Use **normal** for
user-facing behaviour, and **low** for copy, docs and isolated presentation.
When unsure, choose the higher level. Never lower risk to skip a stage.

**Size** says whether to split the work. Large work usually means an
initiative with child items (`--parent <ID>`).

Then write the user's request, in their words, under `## Request` in
`brief.md`.

## The loop

```text
ignition brief <ID>    → stage, what it still needs, whether the owner decides
do that stage's work   → edit the record files; log; add evidence
ignition advance <ID>  → next stage, or the list of what's missing, or an owner checkpoint
```

- Do one stage at a time, and only that stage's work. Don't write code during
  shape or design. The exception is a throwaway spike, and you must log it as
  a finding.
- Fill a record section by writing under its `##` heading. You may delete the
  `<!-- guidance -->` comment, but don't rename headings.
- When `advance` lists what's missing, fix exactly those things and run it
  again. Never edit `stage:` or other frontmatter by hand to get past a check.
- [references/stages.md](references/stages.md) explains what each stage needs
  and what good looks like. Read the section for your current stage before
  doing it.

## Owner checkpoints

When `ignition advance` says the owner decides, stop and ask the user. Use
the question it prints, then add:
- a short summary of what they are approving, with the record files to look at;
- your recommendation and the realistic alternatives;
- anything unresolved or uncertain.

Wait for their answer. Then record it in their own words:

```sh
ignition advance <ID> --owner "<their name>" --decision "<their words>"
```

- Never answer a checkpoint yourself, and never take silence or an unrelated
  reply as a yes.
- If their reply is ambiguous ("cool", "ok"), ask once to confirm what it
  approves. Then record both their words and your reading of them.
- If they say no or ask for changes, stay in the stage. Make the changes, then
  log the decision with `ignition log <ID> --type decision`.

## Recording as you go

| When | Command |
| --- | --- |
| You learned something that matters later | `ignition log <ID> --type finding "..."` |
| You made a choice inside the agreed scope | `ignition log <ID> --type decision "..."` |
| You committed a change | `ignition log <ID> --type change "<sha> <what and why>"` |
| The scope must change | `ignition log <ID> --type scope "..."`, then **stop and ask the owner** |
| You checked something | `ignition evidence <ID> --claim ... --result pass\|fail\|inconclusive --revision ... --environment ... --method ... --observed ... --limits ...` |
| A choice affects the whole product | `ignition decision "<title>"`, then fill in the file it creates |
| Your session is ending | `ignition log <ID> --type handoff "changed: ... verified: ... open: ..."` |
| Someone asks who built something | `ignition attribution <ID>` (or no ID, for every item) |

## Evidence rules

- Evidence is what you **observed** on a named **revision** (commit, tag,
  digest, deployment ID) in a named **environment**, by a stated **method**.
  It's not what should be true.
- Record failures as failures. A failing test before a fix is required
  evidence for a bug.
- "Merged", "deployed", "builds" or "returns 200" doesn't show that a user
  journey works. Check the journey.
- Fill `--limits` honestly: what this evidence doesn't show.
- Never edit earlier evidence or log entries. Add a new entry that corrects
  them.

## Never

- Put secrets, tokens, passwords, private keys, private user data or raw
  production logs in records or commands.
- Answer an owner checkpoint, or record approval the owner didn't give.
- Guess your model, or record yourself as another agent or person.
- Edit frontmatter, earlier log entries or evidence to get past a check.
- Run `ignition advance` past a stage you haven't actually done.
- Commit or push the records repository, or the application, unless the
  project's rules or the user say to.

## Reference

- [references/stages.md](references/stages.md): each stage, its record, and
  what good looks like.
- [references/commands.md](references/commands.md): every command, its flags,
  and what its errors mean.
