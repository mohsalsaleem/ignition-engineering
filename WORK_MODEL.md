# Work model

Work moves through stages. Each stage answers one question and leaves a
written artifact. An item can leave a stage only when that artifact is
complete, and at some stages only when the owner has decided. The model lives
in [internal/model](internal/model/model.go); this page explains it.

## Kinds

| Kind | Use it for | Route |
| --- | --- | --- |
| **Initiative** | An outcome too big for one change ("build an ATS") | intake → shape → design → plan → *deliver* → observe → closed |
| **Feature** | New observable behaviour | intake → shape → design → plan → build → verify → release → observe → closed |
| **Bug** | Behaviour that is wrong | intake → *reproduce* → *diagnose* → build → verify → release → closed |
| **Task** | Bounded work: operations, refactoring, docs | intake → plan → build → verify → release → closed |
| **Investigation** | A question to answer before deciding | intake → *method* → *findings* → *decision* → closed |
| **Incident** | An outage, data or security event, or near miss | *contain* → *timeline* → *recover* → *review* → closed |

How the route adapts:
- A **low-risk feature** skips design.
- An item whose brief says `release: none` skips release, because nothing
  ships to users.
- An **initiative** never goes straight to code. Its plan becomes child items
  (`ignition new ... --parent <ID>`). It can't leave *deliver* until every
  child is closed.
- An **investigation** ends in a decision. Any production change it leads to
  needs its own item.
- An **incident** ends in follow-up bugs and tasks.

## Classification

Every item is classified at intake by **size** (small, medium, large) and
**risk** (low, normal, high). The two are independent: a one-line
authorization fix is small and high risk. Size says whether to split the work;
it never removes a check. Risk decides which stages apply and where the owner
decides. An item with no risk set is treated as high.

| Risk | Typical change |
| --- | --- |
| Low | Copy, isolated presentation, documentation |
| Normal | A feature or bug in a user journey |
| High | Auth, personal data, permissions, payments, AI decisions, schema, infrastructure, release tooling |

## Stages and what completes them

| Stage | Question | Complete when |
| --- | --- | --- |
| intake | What is asked, and what kind of work is it? | `brief.md`: Request; title, size and risk set |
| shape | What outcome, for whom, how will we know? | `brief.md`: Problem, Outcome, Acceptance criteria, Scope |
| design | How, and why this way? | `design.md`: Context, Options, Decision, Impacts, Test strategy |
| plan | In what steps? | `plan.md`: at least one task (a task item also needs "Done when") |
| reproduce | Can we make it fail? | `brief.md`: Expected and actual, Reproduction; evidence with result **fail** |
| diagnose | Why does it fail? | `brief.md`: Root cause |
| build | Is the change made? | every plan task done or dropped; a `change` log entry naming the commit |
| verify | Does it work? | evidence with result **pass** |
| release | Is it running for users, and can we undo it? | `release.md`: Revision, Checks, Rollback, Production verification |
| observe | Did the outcome happen? | `brief.md`: Result |
| deliver | Are the child items done? | every child item closed |
| method / findings / decision | What's the question, what was found, what's decided? | `brief.md` sections; evidence for findings |
| contain / timeline / recover / review | Impact stopped, events known, service restored, lessons taken? | `brief.md` sections; evidence of recovery |

A section counts as filled when it has content beyond the template's
`<!-- guidance -->` comments. `ignition check` also re-checks every stage the
item has already passed. Emptying an earlier artifact is reported as a
problem, and the item can't move until it's fixed.

## Owner checkpoints

| Leaving | Low | Normal | High | Question |
| --- | --- | --- | --- | --- |
| shape | | owner | owner | Is this the right problem, outcome and scope? |
| design | | | owner | Is this the right approach? |
| diagnose | | | owner | Is the root cause right, and should the fix go ahead? |
| verify, into release | | owner | owner | Is the evidence enough to release to users? |
| decision | | owner | owner | Do you accept this decision and its follow-up? |
| review (incident) | owner | owner | owner | Do you accept the review and its follow-up? |

At a checkpoint, `ignition advance` stops and prints the question. The agent
asks the owner, then records their answer with `--owner` and `--decision`. A
change of scope at any stage is a `scope` log entry, and the agent stops to
ask. Today these decisions are *recorded*, not *verified*. See
[ROADMAP.md](ROADMAP.md).

## The agent loop

```text
ignition brief <ID>    → stage, what's missing, whether the owner decides
do that stage's work   → write the artifact; log findings, decisions, changes; record evidence
ignition advance <ID>  → moves on, or says what's missing, or asks for the owner
before stopping        → ignition log <ID> --type handoff "changed / verified / open"
```

`ignition guide` prints the full instructions an agent follows.
