# Product Engineering System

A small, product-independent way to plan, build, review, release, and operate software. The core is deliberately independent of programming language, cloud provider, issue tracker, and AI coding tool. Each product supplies its own commands, environments, risks, and owners in `engineering/PROJECT.md`.

## Start a product

1. Read [SYSTEM.md](SYSTEM.md) and copy [templates/PROJECT.md](templates/PROJECT.md) into the product repository as `engineering/PROJECT.md`. Fill every field that affects release safety. Link existing product documentation instead of duplicating it.
2. Put each change in a work item using [templates/WORK_ITEM.md](templates/WORK_ITEM.md). Keep it small enough to review and verify. Use [templates/DESIGN_NOTE.md](templates/DESIGN_NOTE.md) for a consequential choice, migration, new data flow, or external integration.
3. Implement on a short-lived branch. The product profile names its test commands and branch/release route. Use [templates/PULL_REQUEST.md](templates/PULL_REQUEST.md) to show evidence and review the diff.
4. Record staging and production evidence with [templates/RELEASE.md](templates/RELEASE.md). Treat a successful merge as a code event; a release is complete only after the production journey works and monitoring is normal.
5. Use [templates/INCIDENT.md](templates/INCIDENT.md) after an outage, data or security event, failed release, or a near miss worth learning from. Review the backlog and operational signals weekly.

The [GitHub adapter](adapters/github/README.md) is optional. A product can use another tracker or CI provider and still follow the core system. [AGENT_CONTRACT.md](AGENT_CONTRACT.md) gives AI coding agents the same evidence and review expectations.

[RECORDS.md](RECORDS.md) defines where requirements, decisions, agent instructions, review comments, test results, deployment facts, and incident evidence live. It distinguishes proposals from verified observations and keeps each claim tied to its revision and environment.

## What this kit does and does not do

This kit supplies shared decisions, templates, and gates. The product repository remains the source of truth for its own requirements, code, tests, configuration, release records, and incidents. Copy or link the relevant templates into a product; do not make products depend on this local folder at runtime.

The checklist is not a substitute for functioning CI, branch controls, backups, monitoring, and human judgment. The first adoption task for each product is to verify those controls actually run. If a hosting plan cannot enforce a gate, record the manual check and owner in its project profile rather than claiming automated protection.

## Versioning

Treat this as a versioned playbook. Change the core only when a lesson applies across products. Keep product-specific exceptions in that product's `engineering/PROJECT.md`. Current version: `0.2.0`; revise after the first real release and incident review.

## Basis

The emphasis on small changes, automation, and measurable delivery follows [DORA's small-batch guidance](https://dora.dev/capabilities/working-in-small-batches/) and [continuous-delivery capabilities](https://dora.dev/capabilities/continuous-delivery/). Use [DORA's current delivery metrics](https://dora.dev/guides/dora-metrics/) as learning signals, not targets for individuals.
