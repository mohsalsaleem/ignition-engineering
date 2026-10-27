# Product architecture and system topology

**Status:** Design proposal. The v0.4 kit does not yet validate or render an architecture model.

Ignition models a **product as a connected system**, not as a repository. One repository may build a browser app, API, worker, scheduler, and migrations; one system may also span several repositories. Source paths, runtime boundaries, deployment units, data stores, queues, external providers, and user journeys are different concepts and need different IDs.

Use a small, versioned architecture graph in `.ignition/architecture/` as the declared model. Generate or review diagrams from that graph. The graph is the stable vocabulary used by work items, agent runs, risk analysis, tests, releases, incidents, and cost attribution. Git/CI/deployment/monitoring integrations supply **observed state** with source and time; they do not silently overwrite the declared model.

## Core elements

| Element | What it identifies | Example |
| --- | --- | --- |
| Product and software system | Business capability and its boundary | Hiring product, ATS subsystem |
| Runtime element | Independently running application or data store; roughly a C4 container | Web app, API, matching worker, queue, database |
| Repository and source path | Where code lives | `apps/web`, `services/worker` in one repo |
| Build and deployment unit | What is built, promoted, rolled back, and configured | Web image, worker image, migration job |
| Environment | Where a unit runs, with its own dependencies and policy | Staging, production |
| Relationship | Communication or data movement between elements | Web publishes job to queue; worker reads it |
| Journey or flow | Important behavior across elements | Candidate applies; employer advances stage |

Give each element a stable ID, owner, purpose, and links to its source of truth. Architecture element IDs are scoped to their product graph; engineering record IDs use the product-prefixed ULID format in [FORMAT.md](FORMAT.md). Relationships name direction, interaction, data classification, trust boundary, and failure effect where relevant. Deployment mappings name the intended environment and release unit; the platform's live deployment ID and image digest belong in observed evidence.

```text
.ignition/architecture/
├── model.yaml                     # systems, runtimes, stores, repositories, links
├── environments/
│   ├── staging.yaml               # intended deployment units and dependencies
│   └── production.yaml
├── flows/
│   ├── candidate-applies.md       # important cross-element journey
│   └── advance-applicant.md
└── OVERVIEW.md                    # human explanation and generated/view links
```

For example, one repository can own two separately deployed runtime elements:

```yaml
schema_version: 1
product: example-ats
systems:
  - id: ats
    name: Applicant tracking
repositories:
  - id: app-repo
    url: "https://example.invalid/example-ats.git"
runtimes:
  - id: web
    system: ats
    kind: application
    repository: app-repo
    source_paths: ["apps/web"]
    owner: "team:hiring"
  - id: worker
    system: ats
    kind: application
    repository: app-repo
    source_paths: ["services/worker"]
    owner: "team:hiring"
  - id: jobs
    system: ats
    kind: queue
    owner: "team:hiring"
  - id: applicant-db
    system: ats
    kind: datastore
    data_classification: personal
    owner: "team:hiring"
relationships:
  - from: web
    to: jobs
    interaction: publishes-application-job
    data_classification: personal
  - from: worker
    to: jobs
    interaction: consumes-application-job
  - from: worker
    to: applicant-db
    interaction: writes-application
    data_classification: personal
```

This declares architecture, not what is currently running. `environments/production.yaml` maps `web` and `worker` to separate release units:

```yaml
schema_version: 1
environment: production
deployments:
  - id: web-production
    runtime: web
    release_unit: web-image
  - id: worker-production
    runtime: worker
    release_unit: worker-image
```

A production evidence record links the actual deployment IDs, revisions, configuration scope, and observed journey result. The same repository commit can produce both images while only one unit is deployed; Ignition must not mark the other as released.

## Views, based on C4

Use [C4](https://c4model.com/) as the **viewing language**, with the architecture graph as Ignition's source for identities and relationships:

1. **System context:** people, owned systems, and external systems. Required once per product and updated when boundaries change.
2. **Container view:** runtime applications and data stores inside a system, including web, workers, queues, and databases. Required for a product with multiple runtime parts. In C4, “container” means a runtime boundary, not necessarily a Docker container.
3. **Deployment view:** map those runtime elements to staging and production nodes or services, and show which units can be released independently.
4. **Dynamic view:** show one important journey or failure path across elements when the interaction is not obvious. Keep this selective.
5. **Component view:** zoom inside a runtime only when its internal boundaries materially affect design, ownership, or risk. Do not create a diagram for every class or file.

The diagram is a projection, not a second inventory to maintain. Ignition should be able to render a context, container, and deployment view from the graph; a hand-drawn view must reference element IDs and be reviewed against the graph. Architecture decisions remain in decision notes. API contracts, queue schemas, migrations, and runnable configuration remain in the application repositories and are linked from graph elements.

## How the graph governs work

- **Intake and blast radius:** a work item identifies affected runtimes, relationships, data stores, and journeys. Ignition can list upstream/downstream dependencies, owners, trust boundaries, and release units for review.
- **Agent scope:** a run names the elements and source paths it is allowed to change. A changed element outside that scope is a visible exception, not a silently broadened task.
- **Testing:** tests and evidence identify the journey, element, revision, and environment they cover. A worker change prompts queue, retry, failure, and data-store checks when those relationships exist.
- **Release:** a release names the build/deployment units it changes, their order, migrations, compatibility constraints, verification journeys, and rollback route. A monorepo commit does not imply every runtime was deployed.
- **Operations:** incidents and monitoring signals attach to elements and flows. The architecture view can show affected users and dependencies before a fix is chosen.
- **Cost:** development effort attaches to work and runs; runtime infrastructure and provider usage attach to elements and environments. A feature's cost view keeps these categories distinct.

Begin with a small manual inventory and one critical journey. Validate IDs and references first. Read-only discovery from build manifests, deployment platforms, and monitoring can later flag differences between declared and observed state, with a human confirming architectural changes. This avoids claiming the model is current merely because a diagram exists.
