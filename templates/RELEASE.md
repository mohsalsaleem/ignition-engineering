# Release: <product> — <date and revision>

**Release owner:** <name>
**Work items and PRs:** <links>
**Risk:** Low / Normal / High

## Before deployment

- [ ] Accepted commit and built artifact identified
- [ ] Required tests and review passed on that revision
- [ ] Config, secrets, migrations, data movement, provider costs, and flags reviewed
- [ ] Backup and restore position known; rollback procedure compatible with current data
- [ ] Staging real journey passed; synthetic data cleanup planned
- [ ] Release authorization recorded if required

## Deployment and verification

| Stage | Revision/deployment ID | Time | Health and worker result | Real journey result |
| --- | --- | --- | --- | --- |
| Staging | | | | |
| Production | | | | |

## Observation

Errors, latency, jobs, external-provider usage, and customer feedback checked at <times>.
Rollback trigger: <specific condition>.
Known limitations and follow-up: <links or none>.

**Final decision:** Accepted / Rolled back / Monitoring
**Decision owner and time:** <name and time>
