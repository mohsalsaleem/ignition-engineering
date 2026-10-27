# Engineering profile: <product>

**Owner:** <person or team>
**Repository and tracker:** <links>
**System version adopted:** 0.3.0
**Application repository:** <URL>
**Ignite records repository:** <URL>
**Last reviewed:** <date>

## Record locations

- Work tracker and ID scheme: <link or product-repository path>
- Requirements/design/decision notes: <paths>
- Pull request reviews and CI runs: <systems>
- Release and incident records: <paths>
- Restricted evidence store, access owner, and retention: <reference; no credentials>

## Product and boundaries

- Users and critical journeys: <links or short list>
- Source of truth for requirements and decisions: <locations>
- Sensitive data, external providers, and spending limits: <summary and secure references, no secrets>
- Availability or recovery expectation: <explicit target or “not yet set”>

## Delivery route

| Stage | Branch or environment | Deployment trigger | Owner | Verification |
| --- | --- | --- | --- | --- |
| Local | <branch> | <command> | <owner> | <checks> |
| Staging | <branch/environment> | <trigger> | <owner> | <journey> |
| Production | <branch/environment> | <trigger> | <owner> | <journey and observation> |

- Direct push policy: <what is blocked, and what is manually checked if it cannot be enforced>
- Review policy: <normal/high-risk reviewers and self-review rule>
- Release authorization: <who decides, for which risk level>
- Hotfix route and merge-back: <steps>

## Checks and evidence

- Fast local check: `<command>`
- Full test command: `<command>`
- Build or package check: `<command>`
- Schema/config check: `<command or not applicable>`
- Critical journey smoke test: <steps or link>
- CI status and required checks: <actual configured state, not intended state>
- Test-data cleanup rule: <how synthetic records are removed or isolated>

## Operations

- Dashboard and alert location: <links>
- Backup schedule and last restore rehearsal: <links/dates>
- Rollback procedure: <link; include data and schema compatibility>
- Incident contact and communication route: <owner/channel>
- Secret store and rotation owner: <reference only, never secret values>
- AI/provider usage controls if applicable: <scope, caps, monitoring>

## Product-specific exceptions

| Core rule changed | Reason | Compensating control | Owner | Revisit date |
| --- | --- | --- | --- | --- |
| None | | | | |
