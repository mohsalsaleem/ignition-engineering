# Operating model

## Routes from idea to learning

Use the release path below for features, bugs, and tasks that change a running product. A step is recorded as not applicable when the work does not affect that environment; it is never claimed as passed without evidence.

| State | Required evidence | Exit decision |
| --- | --- | --- |
| Proposed | Problem, intended user, desired outcome | Worth doing now? |
| Ready | Acceptance criteria, scope, risk, verification plan | Small and clear enough to build? |
| In progress | Branch, tests, implementation notes | Works locally? |
| In review | Pull request or equivalent diff and test evidence | Correct, maintainable, safe? |
| On staging | Deployed revision and real journey evidence | Safe to release? |
| Released | Production revision, health, rollback path | Works for the intended user? |
| Observed | Monitoring and customer feedback | Done, follow-up, or rollback? |
| Done | Accepted outcome and links to evidence or follow-up | Closed with a reason and owner |

Other kinds use shorter routes with distinct completion evidence:

| Kind | Route | Completion evidence |
| --- | --- | --- |
| Initiative | Proposed → Ready → In progress → Observed → Done | Child work and measured outcome or documented decision to stop |
| Investigation | Proposed → Ready → In progress → Answered → Done | Method, findings, decision, and linked follow-up work if needed |
| Task without runtime change | Proposed → Ready → In progress → In review → Done | Deliverable and relevant review/check |
| Incident | Detected → Containing → Recovering → Monitoring → Reviewed → Done | Impact, recovery evidence, review, and corrective work |

Use one prioritized backlog. Limit active implementation work to what can be reviewed promptly; for a solo operator, start with at most two active work items. A work item is not done because code was merged.

## Size the process to the risk

| Level | Typical change | Minimum gates |
| --- | --- | --- |
| Low | Copy, isolated presentation, documentation | Work item, focused check, diff review, deployment check if live |
| Normal | Feature or bug affecting a user journey | Acceptance criteria, relevant automated tests, reviewer pass, staging journey, production verification |
| High | Auth, personal data, permissions, payments, AI decisioning, schema or infrastructure changes | Design note, threat/data review, failure and rollback plan, migration or restore rehearsal when relevant, explicit release owner, staged verification and close monitoring |

The product profile may raise these gates. A hotfix may compress the sequence, but record the omitted checks, restore them afterward, and hold a short incident review.

## Definition of ready

- The work item names a user or operational problem and a measurable outcome.
- Acceptance criteria describe observable behavior, including a failure or withdrawal path when relevant.
- Dependencies, data access, external services, and cost are identified.
- The smallest useful slice is selected, and a verification method is known.
- The risk level and release owner are named.

## Definition of done

- The change meets the acceptance criteria and relevant tests pass on the actual revision.
- The diff is understandable and reviewed. Config, migrations, prompts, model choices, and data changes receive the same attention as application code.
- Documentation, observability, and operational steps are updated where behavior changed.
- Staging verification is recorded when the change affects a deployed user or operational journey.
- Production is verified through the affected journey after deployment; health alone is insufficient.
- A release record names the revision, result, known limitations, and rollback route for deployed changes. Investigations and non-runtime tasks instead link their findings or deliverable.

## Roles and decisions

One person can fill several roles, but keep the decisions distinct:

- **Product owner:** chooses priority and accepts the user outcome.
- **Implementer:** builds the change and supplies evidence.
- **Reviewer:** challenges scope, design, tests, security, and readability. A solo operator performs a separate review pass after implementation; an independent reviewer is preferred for high-risk work.
- **Release owner:** confirms gates, deploys, watches production, and decides rollback.

Record who made a consequential decision. Do not treat a successful AI response or a green test as approval of a product or security decision.

## The release rule

Prepare a release from an identified commit. Check migration compatibility, secrets/configuration, backups, health checks, and rollback before promotion. Deploy to the product's staging environment and run the affected real journey. Promote the same accepted change to production, verify the real journey again, and observe errors, latency, workers, and external-provider usage for a period suited to risk. Stop or roll back when an agreed failure condition appears. Preserve user data and audit history during rollback.

High-risk changes need explicit release-owner authorization. A product may add stricter approvals. A release record is the evidence; a deployment dashboard's green label is one input.

## Operating cadence

- **At intake:** write the work item before implementation. Keep product ideas separate from committed work.
- **Daily when active:** review the active work, CI, staging, and production alerts; clear blockers before starting more work.
- **Weekly:** reorder the backlog, close obsolete items, inspect incidents and customer feedback, and choose one improvement to the process.
- **After every release:** record the revision, verification, unexpected behavior, and follow-up.
- **Monthly:** review restore ability, dependency/security updates, access, costs, and the trend in delivery and reliability signals.

Track lead time, deployment frequency, change failure, recovery time, and rework as system-level signals. Also track product outcomes and one or two reliability indicators that matter to customers. Use trends to find friction; do not reward more deployments or tests for their own sake.
