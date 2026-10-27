# Ignition — Product Engineering System

A small, product-independent way to plan, build, review, release, and operate software. The core is deliberately independent of programming language, cloud provider, issue tracker, and AI coding tool. Each product supplies its own commands, environments, risks, and owners in `.ignition/PROJECT.md`. The product's Ignition space is a separate Git checkout; see [LAYOUT.md](LAYOUT.md).

## Start a product

1. Read [SYSTEM.md](SYSTEM.md), create a private Git repository for the product's records, and attach it with `python3 scripts/attach.py --app /path/to/product --records-url <records-git-url>`. This checks it out at `.ignition/` and excludes it from application Git without a tracked app change. Copy [templates/PROJECT.md](templates/PROJECT.md) into it as `.ignition/PROJECT.md`. Fill every field that affects release safety and link existing product documentation instead of duplicating it.
2. Put each change in a work item using [templates/WORK_ITEM.md](templates/WORK_ITEM.md). Keep it small enough to review and verify. Use [templates/DESIGN_NOTE.md](templates/DESIGN_NOTE.md) for a consequential choice, migration, new data flow, or external integration.
3. Implement on a short-lived branch. The product profile names its test commands and branch/release route. Use [templates/PULL_REQUEST.md](templates/PULL_REQUEST.md) to show evidence and review the diff.
4. Record staging and production evidence with [templates/RELEASE.md](templates/RELEASE.md). Treat a successful merge as a code event; a release is complete only after the production journey works and monitoring is normal.
5. Use [templates/INCIDENT.md](templates/INCIDENT.md) after an outage, data or security event, failed release, or a near miss worth learning from. Review the backlog and operational signals weekly.

The [GitHub adapter](adapters/github/README.md) is optional. A product can use another tracker or CI provider and still follow the core system. [AGENT_CONTRACT.md](AGENT_CONTRACT.md) gives AI coding agents the same evidence and review expectations.

[RECORDS.md](RECORDS.md) defines where requirements, decisions, agent instructions, review comments, test results, deployment facts, and incident evidence live. It distinguishes proposals from verified observations and keeps each claim tied to its revision and environment.

[DESIGN.md](DESIGN.md) is the proposed product architecture and first workbench design. Its status is **Proposed** until the owner accepts it.

[FORMAT.md](FORMAT.md) proposes Markdown with small YAML frontmatter for narrative records, standalone YAML for machine-only records, and plain Markdown for general instructions and notes. The current v0.4 templates do not yet include those schemas or validation.

[AGENT_GOVERNANCE.md](AGENT_GOVERNANCE.md) proposes how Ignition will attribute agent runs, audit consequential actions, and account for development usage and cost. Automatic capture is planned, not present in v0.4.

[ARCHITECTURE_MODEL.md](ARCHITECTURE_MODEL.md) proposes a product-wide architecture graph and C4 views so one work item can account for web, workers, data stores, deployments, and dependencies across repositories.

## What this kit does and does not do

This kit supplies shared rules and templates. Each product's `.ignition/` checkout holds its engineering records without becoming part of the application repository. The application repository remains the source of truth for code, tests, migrations, and runtime configuration. Neither application runtime nor deployment depends on the local Ignition checkout.

The checklist is not a substitute for functioning CI, branch controls, backups, monitoring, and human judgment. The first adoption task for each product is to verify those controls actually run. If a hosting plan cannot enforce a gate, record the manual check and owner in its project profile rather than claiming automated protection.

## Versioning

Treat this as a versioned playbook. Change the core only when a lesson applies across products. Keep product-specific exceptions in that product's `.ignition/PROJECT.md`. Current version: `0.4.0`; revise after the first real release and incident review.

## Basis

The emphasis on small changes, automation, and measurable delivery follows [DORA's small-batch guidance](https://dora.dev/capabilities/working-in-small-batches/) and [continuous-delivery capabilities](https://dora.dev/capabilities/continuous-delivery/). Use [DORA's current delivery metrics](https://dora.dev/guides/dora-metrics/) as learning signals, not targets for individuals.
