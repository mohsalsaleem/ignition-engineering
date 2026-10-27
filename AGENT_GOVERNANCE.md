# Agent work, audit, and cost model

**Status:** Design proposal. The v0.4 kit does not yet capture runs or costs automatically.

Ignition must answer four questions for each product change: which sessions contributed, what they actually produced, which decisions and checks accepted the work, and what resources it consumed. A session log alone cannot answer those questions. The durable record links agent activity to work items, Git revisions, PRs, tests, approvals, releases, and provider metering.

## Records and relationships

| Record | Purpose | Minimum provenance |
| --- | --- | --- |
| **Agent run** | One bounded session or delegated subtask. Stores task, agent/provider, model/configuration if exposed, instruction and policy revision, allowed scope, start/end, outcome, parent run, and artifact links. | Ignition run ID, native session ID when available, actor, source, timestamps. |
| **Audit event** | A consequential action or decision: delegation, scope change, PR, review, failed or passed gate, exception, approval, deployment, rollback. | Event ID, actor, timestamp, source link, affected work/run/revision, result. |
| **Usage/cost observation** | Metered resource use and monetary amount, if available, tied to a run, work item, release, or declared shared pool. | Source ID and time, unit/quantity, provider/model or service, currency/amount when known, pricing basis, confidence/status. |

An agent run belongs to one product. It can contribute to several work items and may have child runs. A work item can have many runs. The run ID is passed through task prompts, handoffs, PR descriptions, and evidence records so attribution survives a tool or agent change. Child-run usage is counted once from its own metering source; a parent summary does not duplicate it.

File an agent run under its primary work item or incident, and file per-run usage beneath that run. Other work items reference the same run ID and declare allocation shares. Shared invoices and subscriptions belong in a product-level cost pool until their allocation can be justified.

The run's Markdown body gives a short handoff: goal, changed files and revisions, tests attempted and observed results, decisions proposed, open questions, and next owner. It does not retain a raw transcript, secrets, private reasoning, or copied user data. Exact prompts or tool traces are retained only when a product has a justified, access-controlled evidence need and a retention period.

## Attribution and accounting

At run start, name a primary work item or mark the run **Unallocated** with a reason. If a run legitimately spans work items, record allocation shares and a basis; the shares must total 100%. Initiative totals roll up from child work items without adding the same run again. Attribution can be corrected later through a reviewed Git change, preserving the earlier version in history.

Usage observations record native metering units such as tokens, minutes, or compute time. Monetary amounts are separate facts. Use these states:

- **Measured:** a provider or local meter supplied the quantity.
- **Estimated:** a documented rate was applied to measured quantity; store the rate/version and calculation.
- **Reconciled:** an invoice or billing export supports the amount and allocation rule.
- **Unknown:** the source exposes no reliable quantity or per-run cost. Unknown is never recorded as zero.

The same ledger can include CI minutes, staging infrastructure, and other development services. Human review time may be recorded, but only as self-reported time; monetary labor cost needs an explicit rate and consented policy. Subscription fees, included allowances, shared infrastructure, and taxes are kept in a separate cost pool until an explicit allocation rule can apportion them. Do not count both an estimated run charge and its later reconciled charge in totals. A cost report shows currency, period, source coverage, estimated versus reconciled amounts, and the unallocated portion. Development cost is distinct from runtime AI cost incurred by product users.

## Governance and audit

The product profile sets agent permissions and budget policy: allowed repositories, environments, data classes, tools, spending caps or alerts, and which actions require an owner or independent reviewer. A run records the policy revision in force. A high-risk change still follows its design, review, staged verification, and release gates. An agent cannot approve its own exception or turn a self-reported result into **Verified** evidence.

The audit view should answer who authorized a scope or production change, which revision was inspected, what evidence was available, and what was deferred. Link to source PR, CI, deployment, and provider records. Git records preserve the proposed and corrected versions; protected branches and access controls reduce accidental or unauthorized rewrites. Git alone is not tamper-proof: products needing formal audit guarantees must also export signed or immutable snapshots to controlled storage. Provider data may arrive late, so cost reconciliation can update the cost status without changing the original run outcome.

Useful management measures include time from Ready to production, review and rework cycles, run count per work item, missing-evidence rate, gate exceptions, failed changes, and cost per accepted work item or initiative. Show the distribution and quality context. Token volume, number of sessions, or low cost alone are not performance scores.

## First implementation slice

Start with explicit run IDs, manual start/close, Git/PR/evidence links, an event record for consequential actions, and a usage record that can say **Unknown**. Then add read-only metering adapters and invoice reconciliation for providers that expose usable data. This gives an auditable path before any dashboard claims comprehensive cost coverage.
