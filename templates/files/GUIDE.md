How an agent works with Ignition

Ignition holds a product's engineering records in its workspace (usually
<app>/.ignition, a separate Git repository). Chat is working memory; the
records are what last. A new session must be able to continue from them alone.

Start
  1. Read the workspace's PROJECT.md and AGENTS.md.
  2. Run `ignition brief <ID>`. If you were given a request with no ID, create
     the item: `ignition new <kind> "<title>" --risk <r> --size <s>`.
     Kinds: initiative (an outcome made of child items), feature (new
     behaviour), bug (wrong behaviour), task (bounded work), investigation (a
     question to answer), incident (an outage or near miss).

Loop
  1. `ignition brief <ID>` says the stage and what it still needs.
  2. Do the work of that stage, and only that stage:
       shape      problem, outcome, acceptance criteria, scope in brief.md
       design     context, at least two options, decision, impacts, test
                  strategy in design.md
       plan       tasks with acceptance in plan.md
       reproduce  a failing test, recorded as evidence with result fail
       diagnose   the verified root cause
       build      the change, test first; `ignition log <ID> --type change`
                  with the commit
       verify     evidence with result pass, on a named revision and
                  environment
       release    revision, checks, rollback, production verification in
                  release.md
       observe    the result: did the outcome happen?
  3. Record as you go:
       `ignition log <ID> --type finding "..."`   something learned
       `ignition log <ID> --type decision "..."`  a choice made within scope
       `ignition log <ID> --type scope "..."`     scope must change: stop
                                                  and ask the owner
       `ignition evidence <ID> --claim ... --result ... --revision ...
          --environment ... --method ...`
  4. `ignition advance <ID>`. If it says the owner decides, ask them the
     question it gives, in plain words, with your recommendation. Advance with
     `--owner <name> --decision "<their words>"` only after they answer. Never
     answer for them.
  5. Before your session ends: `ignition log <ID> --type handoff "..."` with
     what changed, what was verified, and what is open.

Rules
  - Evidence is what was observed on a named revision and environment, not
    what should be true. "Merged" or "returns 200" is not "works for users".
  - Report failures as failures. A failing check is evidence too.
  - Don't edit earlier log entries or evidence; add new ones.
  - Never put secrets, credentials or private user data in the records.
  - Product-wide choices go in a decision record: `ignition decision "..."`.
