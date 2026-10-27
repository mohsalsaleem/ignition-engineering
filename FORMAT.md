# Proposed metadata format

**Status:** Design proposal, not yet implemented in the v0.4 Markdown templates.

Use one Markdown file per record in the product's separate `.ignition/` Git repository. Put TOML metadata between `+++` fences at the start of the file; put the explanation below it. This keeps a record easy to open, review, link, and move while giving Ignition structured fields it can validate. `PROJECT.md` uses the same format for product identity and tool mappings. `AGENTS.md` holds standing agent instructions and is versioned by Git.

```text
.ignition/
├── PROJECT.md                    # product metadata + operating context
├── AGENTS.md                     # standing agent instructions
├── initiatives/INIT-0001.md
├── work/WORK-0042.md             # only when no tracker owns the backlog
├── design/DES-0003.md
├── decisions/DEC-0007.md
├── runs/RUN-0018.md              # agent session and handoff
├── audit/EVT-0031.md             # consequential action or decision
├── usage/USE-0024.md             # metering/cost observation
├── evidence/EV-0021.md           # test or production observation
├── releases/REL-0004.md
└── incidents/INC-0002.md
```

For example, a hypothetical `work/WORK-0042.md` contains both metadata and narrative:

```markdown
+++
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
+++

# Export fails for an empty report

## Expected and observed

An empty report should download with headers. It currently returns an error.

## Acceptance criteria

- The download succeeds for an empty report.
- A regression test covers the empty case.
```

The frontmatter gives the CLI and future workbench reliable fields for lists, links, and gates. The Markdown body gives people enough context to judge the work. A field such as status has one canonical location; prose may explain a status change but does not maintain a second status value.

An agent run links to the work item and to usage observations through IDs. A usage record stores the provider's native usage ID so the same charge cannot be imported twice. Its amount is a decimal string such as `"0.42"` with a currency such as `"USD"`; if no reliable amount is available, the amount field is absent and the cost state is `"unknown"`. Its Markdown body explains the source and any allocation or estimation method.

## Parsing and validation rules

- The first line must be `+++`; a matching line containing only `+++` ends the frontmatter. Parse only that section as TOML 1.0; treat everything after it as Markdown. Keep the frontmatter small and avoid multiline TOML strings so the fence remains unambiguous.
- Use UTF-8, UTC timestamps with an explicit `Z` offset, stable IDs, and a numeric `schema_version`. The file basename must match the record ID. A validator must reject malformed frontmatter, missing required fields, invalid enum values, duplicate IDs, and broken local references.
- Discover records from their file paths and derive indexes rather than maintaining a hand-edited master index. A record may have a sibling assets directory named after its ID for small, redacted supporting files.
- The shared schema defines common fields; record-type schemas add their own fields. Evidence, for example, adds claim, method, environment, application revision, observed result, time, and evidence/retention references. A deployed result cannot be marked **Verified** without the deployment or image identity and relevant configuration scope.
- If an external issue tracker is the canonical backlog, keep its work item ID and status there. Ignition links to it and stores only supplementary design or evidence records, avoiding a duplicate work queue.
- Git history provides change history and authorship. Explicit timestamps represent the real event being recorded, such as a test run or release decision; they are not inferred from file modification time. Schema changes require a documented migration; old records are not silently reinterpreted.
- Store links and concise conclusions for PR reviews, CI runs, deployment events, and monitoring signals. Those systems retain their own authoritative records. Keep secrets, raw user files, and sensitive production data out of Ignition; reference restricted evidence with access and retention metadata.
- Agent runs, audit events, and usage observations follow the same format. The validator checks parent/child relationships, allocation shares, and duplicate provider metering IDs. Monetary values use decimal strings and an explicit currency, never binary floating-point numbers. Unknown quantities or costs are omitted and marked unknown, never written as zero.

The current v0.4 kit still has Markdown-only templates. Implementing this format requires versioned schemas, a frontmatter validator, template conversion, and a migration for any existing product records. Until then, it is a proposal in [DESIGN.md](DESIGN.md), not a capability of the attach script.
