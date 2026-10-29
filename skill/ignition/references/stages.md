# Stages

Each kind takes a route through these stages. `ignition brief <ID>` shows the
item's route.

| Kind | Route |
| --- | --- |
| initiative | intake → shape → design → plan → deliver → observe → closed |
| feature | intake → shape → design → plan → build → verify → release → observe → closed |
| bug | intake → reproduce → diagnose → build → verify → release → closed |
| task | intake → plan → build → verify → release → closed |
| investigation | intake → method → findings → decision → closed |
| incident | contain → timeline → recover → review → closed |

Two things change a route:
- A low-risk feature skips design.
- `release: none` in the brief's frontmatter skips release. Use it only when
  nothing ships to users, and log why.

## intake

**Record:** `brief.md`: Request; `title`, `size` and `risk` set.
**Do:** write the request in the requester's words, and say who asked. Set
the kind, risk and size (see SKILL.md). If the request mixes several changes,
make it an initiative, or create one item per change.

## shape

**Record:** `brief.md`: Problem, Outcome, Acceptance criteria, Scope.
**Owner decides** for normal and high risk.
**Do:**
- **Problem:** who is affected and what fails or is missing today. Use
  specifics, not "users want".
- **Outcome:** what will be observably true afterwards.
- **Acceptance criteria:** numbered, as "Given / when / then". Include at
  least one failure, permission or withdrawal path.
- **Scope:** the smallest useful slice, what is explicitly out, and
  dependencies.

**Good looks like:** someone else could test it from the criteria alone.

## design

**Record:** `design.md`. For a feature: Context, Options, Decision, Impacts,
Test strategy. Medium and large features, and high-risk ones, also need
Detailed design. A child of an initiative also needs Parent design. For an
initiative, see "Solution design" below.
**Owner decides** for high risk, and for an initiative at normal risk too.
**Do:**
- Read the code it touches first. Context is the real constraints: existing
  code, data, users, operations.
- Give at least two real options, each with costs and failure modes.
  "Do nothing" counts if it's real.
- Under Decision, say which option and why. A choice that affects the whole
  product also gets `ignition decision`, linked here.
- Under Impacts, cover data, security and privacy, migrations and
  compatibility, rollback, observability and cost. Write "none" where there
  is none.
- Under Test strategy, say which tests prove each criterion and each failure
  path.
- Under Detailed design, when required: API and interfaces (signatures,
  requests, responses, events), the data model (tables, fields, indexes,
  migrations, backfills), key flows including failure paths, and errors (what
  fails, what the user sees, retries). Write what a reviewer needs to judge
  the change before it's built, not a copy of the code.
- Under Parent design, for a child item: link the initiative's solution
  design, name the components and interfaces this item builds on, and say what
  it adds. Don't repeat the parent's design. If this item needs to change it,
  log `--type scope` and stop for the owner.

**Solution design (initiative).** The design of the whole initiative, which
every child item builds on. Design the whole before any slice.
- Context: the real constraints, and what the product already has. Read the
  code and flows first.
- Options and Decision: at least two real approaches to the whole solution,
  including reuse or buy where real.
- Architecture: the components, new and changed, what each owns, and how they
  fit the existing system.
- Interfaces: the contracts between components and between child items: APIs,
  events, shared data and its owner. Children rely on these.
- Cross-cutting concerns: security and permissions, personal data, data
  migration, observability, performance, cost, rollback.
- Slicing: the child items, in delivery order, and what each delivers on its
  own. The plan stage creates them.
- Risks: what could make the design wrong, how you'll find out early, and the
  fallback.
- Test strategy: the end-to-end journeys that prove the outcome across slices.

## plan

**Record:** `plan.md`: a task table with at least one row. A task item also
needs "Done when" in `brief.md`.
**Do:** split the work into tasks, each small enough to review on its own,
each with its own acceptance. Status is one of `todo`, `doing`, `done`,
`dropped`. For an initiative, the tasks are child items, following the solution design's
Slicing: create them with `ignition new <kind> "..." --parent <ID>`, and put
their IDs in the table. Child items can't be created before the initiative's
design is done.

## reproduce (bug)

**Record:** `brief.md`: Expected and actual, Reproduction; evidence with
result **fail**.
**Do:** write a test that fails because of the bug, run it, and record it:
`ignition evidence <ID> --result fail ...`. Quote exact errors. If you can't
reproduce it, log a finding and ask the owner. Don't guess a fix.

## diagnose (bug)

**Record:** `brief.md`: Root cause.
**Owner decides** for high risk.
**Do:** name the code and the condition that cause it, and say how you
confirmed it. "Probably" is not a root cause.

## build

**Record:** every plan task done or dropped; a `change` log entry naming each
commit.
**Do:** write the test first, then the change. Keep commits small and follow
the project's commit rules. After each commit:
`ignition log <ID> --type change "<sha> <what and why>"`. Update the task
statuses in `plan.md`. If the work outgrows the scope, log `--type scope` and
stop.

## verify

**Record:** evidence on the **build revision**, the commit the latest
`change` entry names: at least one pass, a pass for every numbered acceptance
criterion (`--criteria`), and no criterion or claim whose latest result is a
failure. At normal and high risk, only **QA evidence** counts: evidence from
someone other than the builder.
**Owner decides** before release, for normal and high risk.
**Do:** hand the item to a QA agent with `ignition qa <ID>` (see SKILL.md,
"Verify with a QA agent"). QA tests that exact revision, prefers automated
checks, writes the tests that are missing, and records every check, failures
included. For a bug, the test that failed in reproduce must now pass. If QA
finds a failure, fix it in this stage, log the `change`, and run QA again on
the new revision. Then ask the owner whether the evidence is enough to
release.

**Criteria** are the numbered lines of `brief.md`'s Acceptance criteria (a
feature) or Done when (a task). Number them, or verify can't link evidence to
them.

## release

**Record:** `release.md`: Revision, Checks, Rollback, Production
verification.
**Do:** release only through the project's route, and only with the owner's
approval when the project requires it. Name the exact revision, tag or
digest. Under Rollback, give the real route back and whether it's compatible
with data written since. Under Production verification, describe the real
user journey you checked after release, with evidence. A health check isn't
enough.

## observe

**Record:** `brief.md`: Result.
**Do:** say whether the outcome happened, using real signals: use, errors,
feedback. Link follow-up items.

## deliver (initiative)

**Record:** every child item closed.
**Do:** work through the child items. Run `ignition status` to see them, and
`ignition brief <child>` for each.

## method, findings, decision (investigation)

**Record:** `brief.md`: Question and Method; then Findings, with at least one
evidence record; then Decision and Follow-up.
**Owner decides** the decision for normal and high risk.
**Do:** time-box the question. Every claim in Findings needs evidence. A
production change the decision calls for gets its own work item, listed under
Follow-up.

## contain, timeline, recover, review (incident)

**Record:** `brief.md`: Impact and Containment; then Timeline; then Recovery,
with evidence of recovery; then Cause and Follow-up.
**Owner decides** the review, at every risk level.
**Do:** stop the impact first, and touch only what the project allows.
Timeline entries are "time with zone · what happened · evidence". Recovery
evidence is a real journey or data check. Cause separates facts from
hypotheses, without blame. Follow-up lists the bugs and tasks you created.
