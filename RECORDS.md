# Records, evidence, and source of truth

The system is a linked set of records, not one giant document or a transcript of every conversation. Each product maps these record types to its own tools in `.ignite/PROJECT.md`. Its `.ignite/` folder is a separate Git checkout, so process records do not trigger application commits or deployments.

## What counts as a fact?

Label the nature of a statement before relying on it:

- **Intent:** a requested outcome or acceptance criterion. It is not evidence that behavior exists.
- **Proposal:** an option or analysis awaiting a decision.
- **Decision:** an accepted choice, owner, date, reason, and supersession path.
- **Observation:** a test or production result with revision, environment, time, method, and evidence link.
- **Current rule:** an instruction or configuration known to be in force at a stated revision.

For example, “matching works” is too broad. “The employer comparison returned the expected result for test case X on production deployment Y at time Z” is an observation. Record its limits separately.

## Where each record lives

| Record | Canonical home | Link from | Rule |
| --- | --- | --- | --- |
| Work item, acceptance criteria, priority | One issue tracker; `.ignite/work/` if there is no tracker | PR, design, release | One ID per change; update status, do not maintain duplicate backlogs |
| Product requirements and design | `.ignite/design/` for new engineering records; link existing app documentation | Work item and PR | Review changes in Git; show status and owner |
| Architecture/product decision | `.ignite/decisions/`, as a numbered decision note | Work item, PR, later decisions | Mark Accepted or Superseded; keep historical rationale |
| Standing agent instructions and test/release commands | `.ignite/AGENTS.md`, `.ignite/PROJECT.md`; executable scripts stay in the app repo | Work item and agent task | Version both; task-specific instructions belong in the work item |
| Code, tests, migrations, runtime config, prompts | Application repository | PR and release | Identify the exact commit and built artifact |
| Review discussion and comments | Pull request or equivalent code-review system | Work item | Resolve threads; move consequential accepted decisions into a durable design or work item record |
| Automated test feedback | CI run and committed regression tests | PR and release | Save run URL, commit, check result; archive essential expiring artifacts |
| Manual test and analysis evidence | `.ignite/evidence/`; large/restricted files in controlled artifact storage | PR and release | Record method, result, environment, revision, time, and limits |
| Deployment and runtime state | Deployment platform and monitoring system | Release record | Record deployment ID, image/commit, config version, time, health and real journey result |
| Release decision and production verification | `.ignite/releases/` | Work item and incident | Append outcome and follow-up; do not silently rewrite an earlier result |
| Incident timeline and corrective actions | `.ignite/incidents/` plus tracked follow-up items | Release and work items | Preserve facts, hypotheses, customer impact, and recovery evidence |

Use native IDs and permanent links where possible. The work item is the root index for a change, but it does not duplicate every test log. A release may include several work items. The product profile lists the actual paths and tools used by that product.

## Minimal evidence envelope

Every consequential test, release, or production claim needs:

1. **Claim and expected result:** what was checked and what success means.
2. **Scope:** product, environment, user role, data fixture, and affected journey.
3. **Revision:** commit SHA and, for deployed tests, image digest or deployment ID plus relevant config/flag version.
4. **Method and actor:** command, browser steps, monitor, or reviewer; who/what ran it.
5. **Time and result:** timestamp with timezone, pass/fail, observed output, and any uncertainty.
6. **Evidence location and retention:** CI URL, redacted screenshot/report, or restricted artifact reference; when it expires or must be preserved.

Use [templates/EVIDENCE.md](templates/EVIDENCE.md) for a manual or unusual test. Routine CI need not generate one file per test: its run URL and a short PR/release summary are enough. A green health check proves health, not the affected user journey.

## Analysis and blast radius

The work item records the initial impact hypothesis. A design note expands it for high-risk work: users and roles; data and permissions; APIs and jobs; schema and migrations; integrations, prompts, models, and cost; deployment and rollback. The PR confirms what actually changed and links the tests that cover each risk. Mark unverified assumptions explicitly.

Review comments stay with the reviewed diff. A comment that changes scope, trust, architecture, or a release decision becomes a short entry in the work item or decision note. That prevents an important decision from surviving only in a hidden thread.

## Retention and privacy

Version-control small, redacted, durable evidence. Keep secrets, private user files, raw production logs, and sensitive screenshots out of Git and issue comments. Store restricted evidence in an access-controlled location with an owner and retention rule; link to it from the record. CI artifacts can expire, so a release that needs long-term proof must preserve a redacted summary or archive before expiry. GitHub, for example, documents configurable workflow-artifact retention and a 90-day default for build logs and artifacts [here](https://docs.github.com/en/actions/how-tos/manage-workflow-runs/download-workflow-artifacts).

An agent chat is working memory, not the system of record. Before closing work, transfer accepted decisions, evidence, limitations, and follow-ups to the records above. Do not preserve every intermediate thought or copy sensitive material into a new location merely for completeness.

The preference for versioning requirements, tests, configuration, release procedures, and AI instructions follows [DORA's version-control guidance](https://dora.dev/capabilities/version-control/). [NIST SSDF](https://csrc.nist.gov/projects/ssdf) also calls for tracking security requirements, risks, design decisions, and release provenance.
