# Design: {{.Title}}

Work item: {{.ID}}

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

## Detailed design

<!-- Required for medium and large features, and for high risk. Otherwise optional.
API and interfaces: signatures, requests and responses, events.
Data model: tables, fields, indexes, migrations and backfills.
Key flows: the main path and the failure paths, step by step.
Errors: what fails, what the user sees, retries. -->

## Impacts

<!-- Data, security and privacy, migrations and compatibility, rollback, observability, cost. Write "none" where there is none. -->

## Test strategy

<!-- Which tests prove the acceptance criteria and failure paths, at which level. -->
