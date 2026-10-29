# Design: {{.Title}}

Work item: {{.ID}}

Read `.ignition/architecture/` first. This design describes a change to the
system described there.

## Context

<!-- The constraints that drive the design: users, data, scale, existing system, operations. -->
{{- if .Parent}}

## Parent design

<!-- This item belongs to {{.Parent}}. Link its solution design, name the parts this item builds on (components, interfaces), and say what this design adds. If this item needs to change the parent design, log a `scope` entry and ask the owner. -->
{{- end}}

## Options

<!-- At least two real options, each with benefits, costs and failure modes. -->

## Decision

<!-- The chosen option and why. Product-wide decisions go in decisions/ and are linked here. -->

## Interfaces

<!-- Required for medium and large features, and for high risk. Otherwise optional; write "none: <why>" where nothing changes.
Signatures, requests and responses, events, error codes, and compatibility with existing callers. -->

## Data model

<!-- Required for medium and large features, and for high risk; write "none: <why>" where the data doesn't change.
The tables or collections this adds or changes, as a Mermaid erDiagram of just those, then for each: fields and types, keys, indexes, constraints, personal data, retention.
Then the migration: how old and new code run side by side, backfill, and rollback.
Also update .ignition/architecture/data-model.md to the agreed model, naming this item, so the owner approves the model with the design.
-->

## Flows

<!-- Required for medium and large features, and for high risk. Otherwise optional.
Each key flow as a Mermaid sequenceDiagram, one per flow, including at least one failure path. Under each: retries, timeouts, idempotency, and where the change commits.
-->

## Errors

<!-- Required for medium and large features, and for high risk. Otherwise optional.
What can fail, what the user or caller sees, what is retried, and what is logged or alerted. -->

## Impacts

<!-- Data, security and privacy, migrations and compatibility, rollback, observability, cost. Write "none" where there is none. If this changes the system's components or a journey, update architecture/overview.md or architecture/flows/ too. -->

## Test strategy

<!-- Which tests prove the acceptance criteria and failure paths, at which level. Each failure path in Flows needs a test. -->
