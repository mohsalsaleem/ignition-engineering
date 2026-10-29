---
name: ignition
description: Do engineering work through Ignition's staged work model. Work items (initiative, feature, bug, task, investigation, incident) move through stages such as shape, design, plan, build, verify and release. Each stage leaves written records and evidence, and the owner decides at checkpoints. Use when a repository has an `.ignition/` folder, when the user mentions Ignition or a work item ID like W-20261027-k3f9q2, or when the user asks to start, continue or hand off a task, feature, bug, initiative, investigation or incident under Ignition.
compatibility: Requires the ignition CLI and git. Works with any agent that reads Agent Skills, including Claude Code, Codex, Cursor, Gemini CLI and GitHub Copilot.
---

# Ignition

Ignition keeps a product's engineering records: work items, designs, plans,
evidence, releases and decisions. They live in the repository being worked
on, at `<repo>/.ignition/`, and are committed with the code they describe.
The repository itself then knows what is happening, why, and what was
proven, for every future agent. The `ignition` CLI moves each item through
its stages. It refuses to advance while a stage's record is incomplete, and it
stops when the owner must decide.

Chat is working memory and the records are what last. Write as you go, so
that a new session can continue from the records alone.

## Before you start

1. Run `ignition version`. If the command is missing, tell the user and suggest
   `go install github.com/mohsalsaleem/ignition-engineering/cmd/ignition@latest`.
2. Run `ignition status` from the repository. If it reports no workspace and
   the user asked you to work under Ignition, run `ignition init` at the
   repository root and tell the user. It creates a tracked `.ignition/`,
   keeps logs mergeable across branches, and keeps the records out of Docker
   images. Commit it with your first change.

   Ask first instead if the repository is public, or if the user might want
   the records kept separately (`ignition attach`). Never create another
   repository on your own.
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

**For an initiative, shape the whole of it before any slice.**
1. Shape: the problem, the outcome and how it will be measured, the
   boundaries, and what the product already has. Read the existing code and
   flows first.
2. Solution design: the approach for the whole initiative: architecture,
   the interfaces between slices, cross-cutting concerns, how it's sliced,
   and the risks. The owner approves it.
3. Plan: create the child items (`ignition new feature "..." --parent <ID>`),
   in delivery order, and list them in `plan.md`.

Only then start the first child. The owner approves the shape and the
solution design before anything is built. Each child's design builds on the
solution design and covers only what it adds. A child that needs to change
the solution design is a scope change: log it and ask the owner.

There are no separate high-level and low-level design stages. The design
stage asks for more as size and risk grow: medium, large and high-risk
features also fill a Detailed design (interfaces, data, flows, errors).

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

## Verify with a QA agent

The builder doesn't verify its own work. At normal and high risk, `verify`
counts only evidence from someone other than whoever logged the changes: in
practice, a QA agent you start for it. Every numbered acceptance criterion
needs a pass on the revision the latest `change` entry names, and a newer
failure reopens it.

1. When build is done, run `ignition qa <ID>`. It prints the QA agent's whole
   task: the revision, the criteria, the test strategy, what's still missing,
   and how to record evidence.
2. Start a fresh agent with that output as its prompt: a subagent, or a new
   session. Give it no other context, so it tests what was asked rather than
   what you meant to build. It records as `--by qa`. A subagent may share
   your session ID, so the name is what tells it apart from you.
3. Wait for its handoff. Don't record evidence for it, and don't argue its
   failures away.
4. If it found failures, fix them while the item stays in verify. Commit, log
   the `change`, and run QA again: evidence on the old revision no longer
   counts.
5. When `ignition brief <ID>` shows verify complete, ask the owner (see
   below).

At low risk you may verify your own work, with the same criteria and revision
rules.

## Committing the records

The records are part of the change. Commit `.ignition/` changes in the same
commits and pull request as the code they describe, following the project's
commit rules. A reviewer then sees the change and its evidence together.
- Work on one item per branch where you can. Logs merge cleanly across
  branches, but two branches changing the same item's stage will conflict.
- Records written after a merge (release and observe) go in a small
  records-only commit, or in the next pull request if the project prefers.
- Before your session ends, commit what you recorded. Uncommitted records are
  lost with the working copy.

## Recording as you go

| When | Command |
| --- | --- |
| You learned something that matters later | `ignition log <ID> --type finding "..."` |
| You made a choice inside the agreed scope | `ignition log <ID> --type decision "..."` |
| You committed a change | `ignition log <ID> --type change "<sha> <what and why>"` |
| The scope must change | `ignition log <ID> --type scope "..."`, then **stop and ask the owner** |
| You checked something | `ignition evidence <ID> --claim ... --criteria <N,N> --result pass\|fail\|inconclusive --revision ... --environment ... --method ... --observed ... --limits ...` |
| Build is done and the item needs QA | `ignition qa <ID>`, given to a fresh QA agent |
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
- Name the acceptance criteria a check covers with `--criteria`, and give
  `--revision` the commit you actually tested.
- Never edit earlier evidence or log entries. Add a new entry that corrects
  them.

## Never

- Put secrets, tokens, passwords, private keys, private user data or raw
  production logs in records or commands.
- Answer an owner checkpoint, or record approval the owner didn't give.
- Guess your model, or record yourself as another agent or person. Never
  record as `qa` for work you built.
- Edit frontmatter, earlier log entries or evidence to get past a check.
- Run `ignition advance` past a stage you haven't actually done.
- Push, merge or release unless the project's rules or the user say to.
  Committing the records with your work is expected (see above).

## Reference

- [references/stages.md](references/stages.md): each stage, its record, and
  what good looks like.
- [references/commands.md](references/commands.md): every command, its flags,
  and what its errors mean.
