# Ignition spaces

Ignition has two versioned layers:

1. **Shared Ignition kit:** this repository holds the product-independent workflow, record model, and templates.
2. **Product Ignition space:** a separate private Git repository for one product's work items, decisions, review outcomes, evidence, releases, incidents, and agent operating context. It is checked out at `<application-repo>/.ignition/` for convenience.

The application repository remains the source of truth for application code, tests, migrations, configuration, and runtime documentation. The Ignition space links to exact application commit SHAs, pull requests, CI runs, and deployment IDs. It does not copy those artifacts or change application behavior.

```text
product/                        # application Git repository
├── .git/
├── .ignition/                    # separate Git checkout; excluded from app Git
│   ├── .git/
│   ├── PROJECT.md              # product profile and record locations
│   ├── AGENTS.md               # operating instructions for agents on this product
│   ├── work/                   # only if no external issue tracker is canonical
│   ├── design/
│   ├── decisions/
│   ├── evidence/
│   ├── releases/
│   └── incidents/
├── src/ ...                    # existing application layout is unchanged
└── tests/ ...
```

The application repository must exclude `.ignition/` from its Git status. [The attach script](scripts/attach.py) adds a local `.git/info/exclude` entry and checks the result without a tracked application change. Local build tools that scan the whole working directory may also need an exclusion. Hosted builds from the application Git repository do not contain the separate Ignition checkout.

This isolation is deliberate: a new release note or design decision in `.ignition/` must not create an application commit or trigger an application deployment. An app PR can still link to its Ignition records, and an Ignition record must name the app revision it describes. Use the product profile to name both repositories and the owner of each.

## Integration points

- **Agent:** give the agent `.ignition/PROJECT.md`, `.ignition/AGENTS.md`, and the work item. A global agent skill or a one-line root pointer can automate discovery, but neither is required for the records to exist.
- **Review and CI:** pull request comments and machine checks remain with the application repository. The Ignition evidence record links to them and retains consequential conclusions.
- **Deployments:** the deployment platform remains authoritative for runtime state. Ignition releases record deployment IDs, commit/image revisions, configuration scope, user-journey results, and rollback decisions.
- **Backlog:** an external issue tracker may be canonical. When it is, `.ignition/work/` holds only supplementary notes or exports, never a second backlog.

Do not use a Git submodule by default. Updating a submodule pointer in the application repo would create application commits for Ignition-only changes. If a product later chooses a tracked `.ignition/` folder instead, it must configure build/deploy path filters and record that exception in its profile.
