# Proposed metadata format

**Status:** Design proposal, not yet implemented in the v0.4 Markdown templates.

Ignition records use two adjacent files: `record.toml` for structured facts and `README.md` for the narrative. Both are ordinary text files in the product's separate `.ignition/` Git repository. A product has one root `project.toml` for its identity and tool mappings, plus `PROJECT.md` for human context and `AGENTS.md` for standing agent instructions.

```text
.ignition/
├── project.toml
├── PROJECT.md
├── AGENTS.md
├── work/                         # only when no external tracker owns the backlog
│   └── WORK-0042/
│       ├── record.toml
│       └── README.md
├── decisions/
│   └── DEC-0007/
│       ├── record.toml
│       └── README.md
├── evidence/
│   └── EV-0021/
│       ├── record.toml
│       └── README.md
├── runs/RUN-0018/                # agent session and handoff
│   ├── record.toml
│   └── README.md
├── audit/EVT-0031/               # consequential action or decision
│   ├── record.toml
│   └── README.md
├── usage/USE-0024/               # metering/cost observation
│   ├── record.toml
│   └── README.md
└── releases/ ...
```

For example, a hypothetical work item could have this `record.toml`:

```toml
schema_version = 1
id = "WORK-0042"
record_type = "work_item"
product = "example-app"
title = "Export fails for an empty report"
status = "ready"
owner = "person:owner"
created_at = 2026-10-27T12:00:00Z
kind = "bug"
size = "small"
risk = "normal"
urgency = "normal"

[[links]]
relation = "decision"
target = "DEC-0007"
```

Its `README.md` explains the expected and observed behavior, reproduction, impact, acceptance criteria, and test plan. The metadata gives the CLI and future workbench reliable fields for lists, links, and gates. The Markdown gives people enough context to judge the work. A field such as status has one canonical location; prose may explain a status change but does not maintain a second status value.

An agent run links to the work item and to usage observations through IDs. A usage record stores the provider's native usage ID so the same charge cannot be imported twice. Its amount is a decimal string such as `"0.42"` with a currency such as `"USD"`; if no reliable amount is available, the amount field is absent and the cost state is `"unknown"`. Each `README.md` explains the source and any allocation or estimation method.

## Rules

- Use TOML 1.0 syntax, UTF-8, UTC timestamps with an explicit `Z` offset, stable IDs, and a numeric `schema_version`. A validator must reject malformed files, missing required fields, invalid enum values, duplicate IDs, and broken local references. Schema changes require a documented migration; old records are not silently reinterpreted.
- Keep one record per directory. Discover records from their `record.toml` files and derive indexes rather than maintaining a hand-edited master index.
- The shared schema defines common fields; record-type schemas add their own fields. Evidence, for example, adds claim, method, environment, application revision, observed result, time, and evidence/retention references. A deployed result cannot be marked **Verified** without the deployment or image identity and relevant configuration scope.
- If an external issue tracker is the canonical backlog, keep its work item ID and status there. Ignition links to it and stores only supplementary design or evidence records, avoiding a duplicate work queue.
- Git history provides change history and authorship. Explicit timestamps represent the real event being recorded, such as a test run or release decision; they are not inferred from file modification time.
- Store links and concise conclusions for PR reviews, CI runs, deployment events, and monitoring signals. Those systems retain their own authoritative records. Keep secrets, raw user files, and sensitive production data out of Ignition; reference restricted evidence with access and retention metadata.
- Agent runs, audit events, and usage observations follow the same versioned TOML rules. Run and cost links use stable IDs; the validator checks parent/child relationships, allocation shares, and duplicate provider metering IDs. Monetary values use decimal strings and an explicit currency, never binary floating-point numbers. Unknown quantities or costs are omitted and marked unknown, never written as zero.

The current v0.4 kit still has single-file Markdown templates. Implementing this format requires versioned schemas, a validator, template conversion, and a migration for any existing product records. Until then, the format is a design decision proposed in [DESIGN.md](DESIGN.md), not a capability of the attach script.
