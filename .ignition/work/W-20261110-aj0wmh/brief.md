---
id: W-20261110-aj0wmh
kind: initiative
title: One platform where an owner's agents work end to end
stage: shape
risk: high
size: large
parent:
created: "2026-11-10T08:56:06Z"
---
# One platform where an owner's agents work end to end

## Request

From the owner, Mohamed Saleem, on 2026-11-10, about Ignition and OpenAppPlatform (OAP):

> "Should we merge the 2? How one builds and how one operates."

> "I wanna build both into a platform. Where the agents operate e2e"

His decisions the same day:
1. It is for one owner on their own server.
2. Agents run locally or anywhere, but not in Claude cloud sessions.
3. Two repositories: Ignition stays its own repository and OAP uses it as a library.

## Problem

The platform is for one person who owns a few products and a server, with agents doing most of the building and the operating.

Today that person has two tools that don't meet:

- **Each built the same ideas separately, and only one enforces them.** Ignition records who the agent is, what the owner decided and what was proven, but an agent could type any of it. OAP verifies the agent by device login, takes the owner's answers through an owner key, and enforces its policy envelope on the server.
- **Nothing connects building to operating.** Ignition's release stage is text an agent writes. OAP will deploy any commit an approved agent proposes, whether or not that commit was verified against a work item.
- **Nothing starts the agents.** Every stage of every item waits for the owner to open a session and say "continue". The owner is the scheduler.
- **Operating doesn't feed back into building.** A failure OAP sees in production becomes work only if the owner notices and files it.

So the owner still carries the work between tools and between stages by hand, and the guarantees stop at the boundary between the two.

## Outcome

The owner states a request once. Agents shape, design, build, verify, release, operate and observe it. The owner is asked only where a decision is theirs, in one place, and answers with one key.

It will be measured on one real product running on OAP:

- **Owner actions per item** equal its checkpoints. No stage is started by hand.
- **No owner decision exists that the owner's key didn't make**, and no log entry or evidence record names an agent the server didn't verify.
- **Every production deploy** traces to a work item whose verify passed on that exact commit, or to an approval the owner gave for deploying without one.

## Acceptance criteria

1. **One identity.** On the platform, every log entry and evidence record names an agent that OAP's device login verified. `ignition check` reports an entry that doesn't.
2. **One inbox.** Every owner checkpoint reaches the owner as an OAP `ask`. The answer, given with the owner key, is written into the item's log. `ignition check` reports a checkpoint that has no such answer.
3. **Release goes through OAP.** The release stage deploys exactly the commit verify passed on. For a commit with no verified work item, OAP holds the deploy for the owner or refuses it, as the envelope says.
4. **Release evidence comes from the source.** The deployment, image digest and health in the item's release evidence are written from OAP's journal, not typed by an agent.
5. **Agents start without the owner.** Once the owner approves an item's shape, a runner starts the agents for each following stage, including a separate QA agent, and stops at each checkpoint. It runs on the owner's machine or their server, never in Claude cloud sessions, on the owner's own API credentials, under a spending limit the owner sets.
6. **Operating feeds building.** A production failure OAP detects opens a bug that names the application and the evidence.
7. **Ignition still works alone.** Without OAP, the `ignition` CLI behaves as it does in 0.10: checkpoints are recorded, and nothing requires a server.
8. **Proven end to end.** One real feature on a real OAP-hosted application goes from request to observed in production with the owner acting only at checkpoints, and the records show it.

## Scope

**In**
- Ignition's work model as a Go library that OAP imports (its packages are private today).
- The seam in OAP: reading a product's work items from its repository, checkpoints as asks, a release gate fed by work-item facts, and release evidence written back.
- A runner that starts and supervises agents.
- The feedback from OAP's health and events into new work items.

**Reused, not rebuilt**
- From OAP: device login and scoped agent sessions, the owner key, `ask`, the policy envelope, the extension point for change facts, the operation journal, the Git repository adapter, notifications.
- From Ignition: the work model, records and evidence rules, the QA handoff, the skill.
- Ignition's roadmap phases 1, 2, 5 and 6 (agent identity, signed decisions, agent binding, release gate) are met by OAP's enforcement and won't be built a second time inside Ignition.

**Out**
- Teams, several owners, a hosted service, more than one server.
- Running agents in Claude cloud sessions.
- Merging the two repositories.
- A web workbench, cost attribution, an architecture graph.

**Constraints and unknowns**
- OAP changes go through OAP's own plan, which the owner approves there. The child items for OAP's side are proposals to that plan.
- The runner drives an installed, unmodified `claude` (or another agent tool) on API keys the owner supplies. It never handles a Claude subscription login.
- OAP's change-facts contract has fixed fields today. Whether work-item facts extend it or need a second source is a design question.
- How the server reaches a product's records (it knows the application's source repository and commit) is a design question.
- How a runner limits what a builder agent can touch is the least understood part. It is delivered last.

**Order.** The seam first, with the owner still starting sessions by hand: identity, inbox, release gate, evidence. Then the feedback loop. Then the runner.

## Result

<!-- Filled at observe: was the outcome measured? What follows? -->
