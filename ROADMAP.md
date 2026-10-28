# Roadmap

Ignition 0.8 is the groundwork: the work model (including solution designs
for initiatives and design depth by size and risk), records kept in the work
repository, the CLI and the agent skill. Today it structures the work and
records decisions. The next phases make governance enforced rather than
advisory.

The principle is to govern through capabilities and chokepoints, not
instructions. An agent can ignore a document, but it can't use a credential
it doesn't hold, and it can't forge a signature.

| Phase | Adds | Done when |
| --- | --- | --- |
| 1. Agent identity | Agents work under their own restricted Git identity: they can push branches and open PRs, but can't merge protected branches, publish, or deploy to production | An agent token is shown unable to merge or publish |
| 2. Signed decisions | A signed `policy.yaml` (risk floors computed from the diff, owner-only actions); `ignition approve`, which signs a checkpoint decision with the owner's SSH key, bound to the item, revision and action | An agent-written checkpoint without a signature is rejected |
| 3. Gate | `ignition gate merge` as a required CI check: work item linked, stages complete, signed approvals for the tier, diff within scope | A PR missing any of these can't merge |
| 4. Evidence at source | `ignition verify` confirms CI runs and deployments against their APIs. Only verified evidence counts for gates | Evidence citing a run for another commit fails |
| 5. Agent binding | Run IDs, commit trailers, harness hooks, an MCP server; asks that notify the owner. Agent, model and session come from the harness, not from what the agent reports | A fresh agent session goes from brief to merged PR with every step recorded |
| 6. Release gate | `ignition gate release` in deploy and publish pipelines; signed break-glass exceptions with expiry | A production release without its gates is blocked |

Later, if real use calls for them:
- development cost attributed to work items;
- an architecture graph for computing blast radius;
- a local workbench UI over the same files.

Each phase is tested with an attack suite: the ways an agent could get around
the gate (writing its own approval, reusing one on another revision, lowering
its risk, citing evidence for the wrong commit). Each attack must be blocked.
