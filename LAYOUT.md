# Ignition spaces

Ignition has two versioned layers:

1. **Shared Ignition kit:** this repository holds the product-independent workflow, record model, and templates.
2. **Product Ignition space:** a separate private Git repository for one product's work items, decisions, review outcomes, evidence, releases, incidents, and agent operating context. It is checked out at `<application-repo>/.ignition/` for convenience. The work-centered structure below is proposed in [FORMAT.md](FORMAT.md); the current attach script does not scaffold it.

The application repository remains the source of truth for application code, tests, migrations, configuration, and runtime documentation. The Ignition space links to exact application commit SHAs, pull requests, CI runs, and deployment IDs. It does not copy those artifacts or change application behavior.

```text
product/                        # application Git repository
├── .git/
├── .ignition/                    # separate Git checkout; excluded from app Git
│   ├── .git/
│   ├── project.yaml            # proposed machine-readable product configuration
│   ├── PROJECT.md              # product context and operating notes
│   ├── AGENTS.md               # operating instructions for agents on this product
│   ├── architecture/           # product-wide system graph, deployments, flows
│   ├── initiatives/<INIT-ID>/
│   │   └── work/<WORK-ID>/     # design/, runs/, evidence/, audit/ under work
│   ├── work/<WORK-ID>/         # standalone work or tracker-linked dossier
│   ├── decisions/              # product-wide decisions
│   ├── releases/<REL-ID>/      # release-specific evidence and audit
│   ├── incidents/<INC-ID>/    # incident-specific evidence and actions
│   └── costs/                  # shared or unallocated bills
├── src/ ...                    # existing application layout is unchanged
└── tests/ ...
```

The application repository must exclude `.ignition/` from its Git status. [The attach script](scripts/attach.py) adds a local Git exclude entry and checks the result without a tracked application change. Worktree-safe exclude resolution and failed-attach cleanup are in [PR #2](https://github.com/mohsalsaleem/ignition-engineering/pull/2).

## Build-context safety decision

An attached `.ignition/` checkout sits inside the application's local filesystem. A local Docker build using the application directory as its context can send private Ignition records to the builder, and `COPY . .` can put them in an image. Therefore, a Docker-based product **must commit an effective `.dockerignore` rule for `.ignition` in the application repository before building from a local checkout**. Check every effective ignore file, including Dockerfile-specific ignore files, and verify the built context/image excludes the records. [Docker documents how build contexts and ignore files work](https://docs.docker.com/build/concepts/context/#dockerignore-files). This tracked safety change is an explicit exception to the usual goal of leaving application Git untouched. The attach script does not make Docker builds safe by itself.

Other packaging, archive, sync, and deployment tools that read the working directory need the same exclusion or a source-only build context. A hosted build that checks out only the application Git repository will not contain this separate records checkout, but local build paths still require review.

This isolation is deliberate: a new release note or design decision in `.ignition/` must not create an application commit or trigger an application deployment. An app PR can still link to its Ignition records, and an Ignition record must name the app revision it describes. Use the product profile to name both repositories and the owner of each.

## Integration points

- **Agent:** give the agent `.ignition/PROJECT.md`, `.ignition/AGENTS.md`, and the work item. A global agent skill or a one-line root pointer can automate discovery, but neither is required for the records to exist.
- **Review and CI:** pull request comments and machine checks remain with the application repository. The Ignition evidence record links to them and retains consequential conclusions.
- **Deployments:** the deployment platform remains authoritative for runtime state. Ignition releases record deployment IDs, commit/image revisions, configuration scope, user-journey results, and rollback decisions.
- **Backlog:** an external issue tracker may be canonical. When it is, Ignition work folders hold only an engineering dossier and links; priority and workflow status remain in the tracker.

Do not use a Git submodule by default. Updating a submodule pointer in the application repo would create application commits for Ignition-only changes. If a product later chooses a tracked `.ignition/` folder instead, it must configure build/deploy path filters and record that exception in its profile.
