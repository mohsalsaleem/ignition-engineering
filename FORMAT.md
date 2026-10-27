# Proposed record formats

**Status:** Design proposal, not yet implemented in the v0.4 templates.

Choose the file format by what the record needs to express. Narrative records use Markdown, with a small YAML frontmatter block only when Ignition needs to index fields. Machine-only records use standalone YAML. Standing instructions and general notes can be plain Markdown. All live in the product's separate `.ignition/` Git repository.

| Content | Proposed format | Reason |
| --- | --- | --- |
| Work item, initiative, design, decision, run handoff, evidence, release, incident | `.md` with minimal YAML frontmatter | Keep explanation and searchable identity/status/links together. |
| Project configuration, architecture graph and environment maps, consequential audit event, usage/cost observation | `.yaml` | Structured facts need validation but little narrative. |
| Agent instructions and general notes | Plain `.md` | No schema when the content needs no machine fields. |
| Full agent traces, raw logs, invoices, sensitive or large artifacts | Source system or controlled storage, linked from Ignition | Avoid copying high-volume or restricted data into Git. |

```text
.ignition/
├── project.yaml                   # IDs, tool mappings, environments, policies
├── PROJECT.md                     # product context and operating explanation
├── AGENTS.md                      # standing agent instructions
├── architecture/                  # systems, runtimes, flows, deployments
│   ├── model.yaml
│   ├── environments/production.yaml
│   └── flows/advance-applicant.md
├── initiatives/<INIT-ID>/
│   ├── <INIT-ID>.md               # initiative outcome and child work
│   └── work/<WORK-ID>/
│       ├── <WORK-ID>.md           # work item and acceptance criteria
│       ├── design/<DES-ID>.md
│       ├── runs/<RUN-ID>/
│       │   ├── <RUN-ID>.md        # agent handoff; parent run linked by ID
│       │   └── usage/<USE-ID>.yaml
│       ├── evidence/<EV-ID>.md
│       └── audit/<EVT-ID>.yaml
├── work/<WORK-ID>/                # work without a parent initiative
│   └── <WORK-ID>.md
├── decisions/<DEC-ID>.md          # product-wide decision
├── releases/<REL-ID>/
│   ├── <REL-ID>.md
│   ├── evidence/<EV-ID>.md        # release-specific verification
│   └── audit/<EVT-ID>.yaml
├── incidents/<INC-ID>/            # incident owns its own evidence/actions
│   └── <INC-ID>.md
└── costs/<POOL-ID>.yaml           # shared bills and unallocated pool
```

The folder follows the **primary owner of the record**. A work item owns its design, sessions, local evidence, and local audit events. A run owns its metering observations. An initiative groups its child work and any initiative-level research. A release or incident owns evidence specific to that event. Product-wide architecture, decisions, and shared costs stay at the top level.

IDs, not paths, define relationships. Every record keeps a collision-resistant ID and references its work, run, release, or incident by ID. When one item supports several work items, keep one canonical copy under a primary owner and link to it from the others. Moving a work folder does not change its ID or create a second cost or evidence record.

## ID minting

Use `<PRODUCT-PREFIX>-<TYPE>-<ULID>`, for example `EXATS-WORK-01M5YKPJG08ANDD3TRV7QG69ZW`. `project.yaml` defines the stable uppercase product prefix. A ULID is a 26-character time-and-random identifier; agents can mint one independently without reserving a counter in a shared branch. The validator rejects malformed IDs and duplicates across the product records repository. Display labels may be shorter, but links and filenames use the full ID. Existing numbered records are retained as legacy IDs until a documented migration. See the [ULID specification](https://github.com/ulid/spec).

```yaml
# project.yaml excerpt
product_id: example-ats
id_prefix: EXATS
```

## Narrative example

`initiatives/<INIT-ID>/work/<WORK-ID>/<WORK-ID>.md` has a short frontmatter block and a normal Markdown body:

```markdown
---
schema_version: 1
id: "EXATS-WORK-01M5YKPJG08ANDD3TRV7QG69ZW"
record_type: work_item
product: example-ats
initiative: "EXATS-INIT-01M5YKPJG061GGETCYTTY0ZHHB"
title: Export fails for an empty report
status: ready
owner: "person:owner"
created_at: "2026-10-27T12:00:00Z"
kind: bug
size: small
risk: normal
urgency: normal
links:
  - relation: decision
    target: "EXATS-DEC-01M5YKPJG0XMHR91Q742TQR9HA"
---

# Export fails for an empty report

## Expected and observed

An empty report should download with headers. It currently returns an error.

## Acceptance criteria

- The download succeeds for an empty report.
- A regression test covers the empty case.
```

The frontmatter owns status and links; prose explains them without keeping a second status value. General notes that need no index can omit frontmatter entirely.

## Machine-only example

A delegated run records its parent in frontmatter, even when the parent lives under another work folder:

```markdown
---
schema_version: 1
id: "EXATS-RUN-01M5YKPJG028T5YQWMWRG0RKEN"
record_type: agent_run
product: example-ats
primary_work: "EXATS-WORK-01M5YKPJG08ANDD3TRV7QG69ZW"
parent_run: "EXATS-RUN-01M5YKPJG06QHBP0BH1HGSHN2D"
---

# Delegated implementation run

The parent assigned the empty-report export regression fix to this run.
```

`initiatives/<INIT-ID>/work/<WORK-ID>/runs/<RUN-ID>/usage/<USE-ID>.yaml` can hold a provider observation without an empty Markdown body:

```yaml
schema_version: 1
id: "EXATS-USE-01M5YKPJG002BH0HYTW34BNVJM"
record_type: usage
product: example-ats
run: "EXATS-RUN-01M5YKPJG028T5YQWMWRG0RKEN"
source_id: example-provider-usage-18
metered_at: "2026-10-27T13:10:00Z"
unit: tokens
quantity: 25000
usage_state: measured
cost_state: estimated
amount: "0.42"               # illustrative, not a provider price
currency: USD
rate_ref: example-pricing-snapshot
allocations:
  - work: "EXATS-WORK-01M5YKPJG08ANDD3TRV7QG69ZW"
    basis_points: 10000
allocation_basis: Session dedicated to this work item
```

The provider source ID prevents duplicate imports. Each usage or shared-cost record owns its own allocation list; `basis_points` must total 10,000 (100%), with an explicit `unallocated` target for any remainder. The run's work links describe contribution but do not duplicate the accounting split. If no reliable amount exists, omit `amount` and set `cost_state: unknown`; unknown is never zero. An audit event is also a small YAML record with action, actor, time, affected work/run/revision, and a link to the original PR, CI, deployment, or approval event.

## Validation and ownership rules

- Use a constrained YAML 1.2 subset in both frontmatter and `.yaml` files: no custom tags, anchors, aliases, or duplicate keys. Quote timestamps and decimal money strings so parsers do not reinterpret them. Use a safe parser and versioned schemas.
- Frontmatter starts on the first line with `---` and ends at the next line containing only `---`. Parse only that section as YAML; treat the rest as Markdown. Keep metadata small; long rationale belongs in the body.
- For ID-bearing records, IDs must match filenames; `project.yaml` and the architecture model use their own schema keys. Validate required fields, enums, relationships, parent-run references and cycles, allocation sums, and duplicate provider metering IDs. Derive indexes from records rather than maintaining a hand-edited master index. Schema changes require a documented migration.
- If an external tracker owns the backlog, a work folder may still hold its engineering dossier and child records. The tracker owns priority and workflow status; Ignition links to it and does not maintain a duplicate backlog. Without a tracker, the work record owns those fields.
- Git history shows edits and authorship. Explicit timestamps represent the actual event, not file modification time. Git alone is not tamper-proof; formal audit needs protected access and an immutable or signed export.
- Keep secrets, raw user files, and sensitive production data out of Ignition. Link restricted evidence with access and retention metadata. Native PR, CI, deployment, provider, and billing systems remain authoritative for their own events.

The current v0.4 kit still has Markdown-only templates. Implementing this proposal requires schemas, a validator, template conversion, and migration for any existing product records. The attach script does not yet create or validate these formats.

[ARCHITECTURE_MODEL.md](ARCHITECTURE_MODEL.md) proposes the separate architecture graph. Its YAML describes structured system relationships; its Markdown flow notes explain important journeys.
