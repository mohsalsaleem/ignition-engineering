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

Agents start with `ignition guide`, then `ignition brief <ID>`.

[examples/oap](examples/oap) is a real bug, the OAP v0.2.0-rc.1 install
failure, recorded through every stage.

## What it does and doesn't do

Ignition checks that each stage's artifact is complete before an item moves
on. It refuses to skip a stage. It also records the owner's decision at each
checkpoint. It doesn't yet *enforce* governance: an agent with write access
to the records could still type an owner decision itself. Signed approvals,
verifying evidence at its source, and gates in CI are the next phases, in
[ROADMAP.md](ROADMAP.md).

## Develop

```sh
make check    # vet and test
make build    # bin/ignition
```

Version: see [VERSION](VERSION).
