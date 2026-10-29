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
- An **initiative** never goes straight to code. Its design is a solution
  design for the whole initiative. Its plan becomes child items
  (`ignition new ... --parent <ID>`), which can only be created once that
  design is done. It can't leave *deliver* until every child is closed.
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
| design | How, and why this way? | `design.md`: see [Design depth](#design-depth) |
| plan | In what steps? | `plan.md`: at least one task (a task item also needs "Done when") |
| reproduce | Can we make it fail? | `brief.md`: Expected and actual, Reproduction; evidence with result **fail** |
| diagnose | Why does it fail? | `brief.md`: Root cause |
| build | Is the change made? | every plan task done or dropped; a `change` log entry naming the commit |
| verify | Does it work? | QA evidence on the built revision: every acceptance criterion passes, nothing failing. See [Verification](#verification) |
| release | Is it running for users, and can we undo it? | `release.md`: Revision, Checks, Rollback, Production verification |
| observe | Did the outcome happen? | `brief.md`: Result |
| deliver | Are the child items done? | every child item closed |
| method / findings / decision | What's the question, what was found, what's decided? | `brief.md` sections; evidence for findings |
| contain / timeline / recover / review | Impact stopped, events known, service restored, lessons taken? | `brief.md` sections; evidence of recovery |

A section counts as filled when it has content beyond the template's
`<!-- guidance -->` comments. `ignition check` also re-checks every stage the
item has already passed. Emptying an earlier artifact is reported as a
problem, and the item can't move until it's fixed.

## Design depth

There is one design stage. How much it must cover depends on the kind, size
and risk of the work, not on separate high-level and low-level design stages.

| Item | `design.md` must fill |
| --- | --- |
| Initiative: the **solution design** | Context, Options, Decision, Architecture, Interfaces, Data model, Flows, Cross-cutting concerns, Slicing, Risks, Test strategy |
| Feature | Context, Options, Decision, Impacts, Test strategy |
| Feature that is medium or large, or high risk | also **Interfaces**, **Data model** (with the migration), **Flows** and **Errors** |
| Feature that belongs to an initiative | also **Parent design**: what it builds on in the initiative's solution design, and what it adds |
| Low-risk feature | no design stage |

A section where nothing changes says `none: <why>`. A design written before
a section was required isn't re-checked for it once it has passed design.

### Architecture docs and the agreed data model

`.ignition/architecture/` describes the system as it is now: `overview.md`
(context, containers, trust boundaries), `data-model.md` and `flows/` (one
file per important journey). Designs read it first and describe their change
against it. The work that changes the system updates it in the same branch.

The data model is agreed in design, not discovered in a migration. A design
whose Data model section describes a change writes the target model into
`architecture/data-model.md`, naming its item, and the design stage can't
complete until it does. The owner approves the model with the design. The
migrations and ORM models then implement it, and stay the source of truth for
exact types. A build that must differ updates the model and the design, and
asks the owner if it changes what they approved.

Diagrams are Mermaid, which agents read and write as text and GitHub renders:
a `flowchart` for containers, an `erDiagram` for data, a `sequenceDiagram`
per flow with its failure paths. They're recommended, not checked. See
[skill/ignition/references/diagrams.md](skill/ignition/references/diagrams.md).

The initiative's solution design is the high-level design that all its child
items share, so the owner approves it at normal risk as well as high. A child
item designs only what differs from it. If a child needs to change the
parent's design, such as an interface between slices, that is a scope change:
the agent logs a `scope` entry and asks the owner.

## Verification

The builder doesn't verify its own work. Verify is done by a QA agent: a
fresh agent that didn't build the change, started with the handoff
`ignition qa <ID>` prints. It tests the exact revision, prefers automated
checks, writes the tests that are missing, and records every check as
evidence, failures included. It doesn't fix what it finds.

`verify` is complete when, on the **build revision** (the commit the latest
`change` entry names):
- every numbered acceptance criterion (Acceptance criteria for a feature,
  Done when for a task) has a passing record that names it with
  `--criteria`;
- no criterion or claim's latest result is a failure. A failure recorded
  after a pass reopens it;
- at least one record passes.

At normal and high risk, only **QA evidence** counts towards those passes:
evidence recorded by someone other than whoever logged the changes (a
different name, or the same name in a different session). Failures count
whoever records them. Low-risk work may verify itself.

A fix made during verify is a new `change`, so it moves the build revision:
QA runs again, and evidence on the old revision no longer counts. The item
stays in verify; it doesn't go back to build.

Who recorded the evidence is what the agent reported. Until agent runs are
bound to verified identities ([ROADMAP.md](ROADMAP.md), phase 5), keeping QA
independent is a rule agents follow, not one Ignition can prove.

## Owner checkpoints

| Leaving | Low | Normal | High | Question |
| --- | --- | --- | --- | --- |
| shape | | owner | owner | Is this the right problem, outcome and scope? |
| design (initiative) | | owner | owner | Is this the right solution design for the whole initiative? |
| design (feature) | | | owner | Is this the right approach? |
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
