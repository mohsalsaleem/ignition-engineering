# Ignite spaces

Ignite has two versioned layers:

1. **Shared Ignite kit:** this repository holds the product-independent workflow, record model, and templates.
2. **Product Ignite space:** a separate private Git repository for one product's work items, decisions, review outcomes, evidence, releases, incidents, and agent operating context. It is checked out at `<application-repo>/.ignite/` for convenience.

The application repository remains the source of truth for application code, tests, migrations, configuration, and runtime documentation. The Ignite space links to exact application commit SHAs, pull requests, CI runs, and deployment IDs. It does not copy those artifacts or change application behavior.

```text
product/                        # application Git repository
├── .git/
├── .ignite/                    # separate Git checkout; excluded from app Git
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

The application repository must exclude `.ignite/` from its Git status. [The attach script](scripts/attach.py) adds a local `.git/info/exclude` entry and checks the result without a tracked application change. Local build tools that scan the whole working directory may also need an exclusion. Hosted builds from the application Git repository do not contain the separate Ignite checkout.

This isolation is deliberate: a new release note or design decision in `.ignite/` must not create an application commit or trigger an application deployment. An app PR can still link to its Ignite records, and an Ignite record must name the app revision it describes. Use the product profile to name both repositories and the owner of each.

## Integration points

- **Agent:** give the agent `.ignite/PROJECT.md`, `.ignite/AGENTS.md`, and the work item. A global agent skill or a one-line root pointer can automate discovery, but neither is required for the records to exist.
- **Review and CI:** pull request comments and machine checks remain with the application repository. The Ignite evidence record links to them and retains consequential conclusions.
- **Deployments:** the deployment platform remains authoritative for runtime state. Ignite releases record deployment IDs, commit/image revisions, configuration scope, user-journey results, and rollback decisions.
- **Backlog:** an external issue tracker may be canonical. When it is, `.ignite/work/` holds only supplementary notes or exports, never a second backlog.

Do not use a Git submodule by default. Updating a submodule pointer in the application repo would create application commits for Ignite-only changes. If a product later chooses a tracked `.ignite/` folder instead, it must configure build/deploy path filters and record that exception in its profile.
