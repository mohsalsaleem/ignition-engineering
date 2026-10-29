# Solution design: {{.Title}}

Work item: {{.ID}}

The design of the whole initiative. Every child item builds on it, and
designs only what differs from it. Read `.ignition/architecture/` first.

## Context

<!-- The constraints that drive the solution: users, data, scale, the existing system and what it already has, operations, team. -->

## Options

<!-- At least two real approaches to the whole solution (including buy, reuse or extend, where real), each with benefits, costs and failure modes. -->

## Decision

<!-- The chosen approach and why. Product-wide decisions also go in decisions/ (`ignition decision`) and are linked here. -->

## Architecture

<!-- The components, new and changed, what each owns, and how they fit the existing system in .ignition/architecture/overview.md. A Mermaid flowchart of the containers (apps, services, workers, stores, queues) and how they connect. Update architecture/overview.md to match. -->

## Interfaces

<!-- The contracts between components and between child items: APIs, events, shared data and who owns it. Child items rely on these, so changing one is a scope change. -->

## Data model

<!-- The entities across the whole initiative, as a Mermaid erDiagram, and which component owns each. For each: key fields, personal data, retention, invariants. Child items detail their own tables, within this model.
Update .ignition/architecture/data-model.md to the agreed model, naming this item, so the owner approves the model with the design. Write "none: <why>" if no data changes. -->

## Flows

<!-- The journeys that cross components or child items, each as a Mermaid sequenceDiagram including a failure path. Add lasting ones to architecture/flows/. -->

## Cross-cutting concerns

<!-- Security and permissions, personal data and privacy, data migration and compatibility, observability, performance, cost, rollback. Write "none" where there is none. -->

## Slicing

<!-- How the work splits into child items, why in that order, and what each delivers on its own. The plan stage creates them. -->

## Risks

<!-- What could make this design wrong, how you'll find out early (spikes, the first slice), and the fallback. -->

## Test strategy

<!-- The end-to-end journeys that prove the outcome across slices, and what each child item must prove on its own. -->
