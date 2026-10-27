# Ignition

Ignition gives agents, and the people they work for, one structured way to do
engineering work. A request becomes a work item of a known kind. It moves
through stages: shape, design, plan, build, verify, release, observe. Each
stage leaves a written artifact, and the owner decides at fixed checkpoints.
Any new session can pick the work up from the records alone. Chat is working
memory; the records are what last.

```text
Intake → Shape → Design → Plan → Build → Verify → Release → Observe → Closed
```

- [WORK_MODEL.md](WORK_MODEL.md): kinds, stages, routes, checkpoints and the
  agent loop.
- [RECORDS.md](RECORDS.md): where records live, their format, and what counts
  as evidence.
- [ROADMAP.md](ROADMAP.md): what comes next, including enforced governance.

## Install

```sh
go install github.com/mohsalsaleem/ignition-engineering/cmd/ignition@latest
```

## Use

A product's records live in their own private Git repository, checked out at
`<app>/.ignition/` so record changes never become application commits or
deployments:

```sh
ignition attach --app ~/src/product --records-url git@github.com:me/product-records.git
cd ~/src/product/.ignition && ignition init --product product   # first time only
```

Then, from anywhere in the application:

```sh
ignition new bug "Export fails for an empty report" --risk normal --size small
ignition brief W-20261027-k3f9q2      # stage, route, what's missing, owner checkpoint
ignition evidence W-20261027-k3f9q2 --claim "Empty export returns 500" --result fail \
    --revision 4f2a9c1 --environment local --method "go test ./export -run Empty"
ignition log W-20261027-k3f9q2 --type change "4f2a9c1 writes headers for empty reports"
ignition advance W-20261027-k3f9q2
ignition status
ignition check
```

## For agents

Ignition ships an agent skill, [skill/ignition](skill/ignition/SKILL.md),
written to the [Agent Skills](https://agentskills.io/specification) format.
It teaches an agent the work model, the loop, owner checkpoints, evidence
rules and attribution. Install it with the CLI, so the skill always matches
the binary:

```sh
ignition skill install                    # ~/.claude/skills and ~/.agents/skills
ignition skill install --project .        # this project only
ignition skill install --for claude       # or --for agents; --to DIR for any other folder
```

`.claude/skills` is read by Claude Code. `.agents/skills` is the shared folder
read by Codex, Cursor, Gemini CLI, GitHub Copilot and others. Agents without
skill support can read the same instructions with `ignition guide`.

## Attribution

Every log entry and evidence record stores who acted: a name, the agent tool,
the exact model ID, and the session.
- The agent and session are detected where the tool exposes them; Claude Code
  sets `AI_AGENT` and a session ID.
- An agent can also give them with `--agent` and `--session` (or
  `IGNITION_AGENT`, `IGNITION_SESSION`).
- No tool exposes its model, so agents pass `--model` (or `IGNITION_MODEL`).
  A missing model is recorded as unknown, with a warning, and never guessed.

`ignition attribution [ID]` reports who built what: each person, agent and
model, their entries, commits and sessions, and who made each owner decision.
`ignition whoami` shows what will be recorded.

## What it does and doesn't do

Ignition checks that each stage's artifact is complete before an item moves
on. It refuses to skip a stage. It also records the owner's decision at each
checkpoint. It doesn't yet *enforce* governance: an agent with write access
to the records could still type an owner decision itself, or report a
different agent or model. Signed approvals, verifying evidence at its
source, and gates in CI are the next phases, in [ROADMAP.md](ROADMAP.md).

## Develop

```sh
make check    # vet and test
make build    # bin/ignition
```

Version: see [VERSION](VERSION).
