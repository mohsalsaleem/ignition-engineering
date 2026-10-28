# Records

## Where they live

Records live in the repository being worked on, at `<repo>/.ignition/`, and
are committed in the same commits and pull requests as the code they
describe. Anyone who clones the repository, person or agent, gets the full
history of what was asked, decided, built and proven.

```text
.ignition/
├── PROJECT.md              # owner, repositories, environments, commands, boundaries
├── AGENTS.md               # product-specific rules for agents
├── decisions/D-….md        # product-wide decisions, referenced from many items
└── work/W-20261027-k3f9q2/ # one folder per item, of any kind
    ├── brief.md            # kind, risk, size, stage | request, problem, outcome, criteria, result
    ├── design.md           # created on entering design; a solution design for an initiative
    ├── plan.md             # created on entering plan
    ├── release.md          # created on entering release
    ├── evidence/E-….md     # one record per observation
    └── log.md              # append-only: stages, checkpoints, findings, changes, handoffs
```

The application repository stays the source of truth for code, tests,
migrations and configuration. Records link to exact commits, pull requests,
CI runs and deployments. They don't copy them.

Two settings keep tracked records out of the way, and `ignition init` sets
both:
- `.ignition/.gitattributes` merges `log.md` with Git's union driver. Two
  branches that log on the same item keep both sides' entries instead of
  conflicting.
- `.ignition` goes in `.dockerignore` when there's a Dockerfile, so images
  don't include the records.

Records-only commits (for example a release record written after merge) are
ordinary commits. If they would trigger CI or a deploy you don't want, filter
`.ignition/**` out of those triggers.

**Separate records repository.** For a public repository, or a product that
spans several repositories, `ignition attach` checks a private records
repository out at `<app>/.ignition/` instead. It excludes the folder from the
application's Git, and warns when Docker would copy it.

## Format

Records are Markdown. Records that Ignition indexes start with a flat
frontmatter block of `key: value` lines between `---` lines. Nested values,
lists and anchors aren't allowed. Sections are `## ` headings. Template
guidance sits in `<!-- -->` comments and doesn't count as content.

IDs are `<prefix>-<YYYYMMDD>-<6 random characters>`: W for work, E for
evidence, D for decisions. They sort by day. Parallel agents are very unlikely
to collide, and a new item never reuses an existing folder. Folders are named
by ID. Items link by ID: an initiative's children name it in `parent`.

## Claims and evidence

Label a statement before relying on it:
- **Intent:** what is wanted (acceptance criteria). Not evidence that it
  exists.
- **Proposal:** an option awaiting a decision.
- **Decision:** what was chosen, by whom, and why. A checkpoint entry in the
  log, or a decision record.
- **Evidence:** an observation, with its claim, result, revision, environment,
  method, actor and time.

Every evidence record names:
- the exact revision: commit, tag, image digest or deployment ID;
- the environment;
- the method: command, test or steps.

Its result is pass, fail or inconclusive. A failing check is evidence too;
record it rather than retrying until green. "Merged", "deployed" or "returns
200" is not evidence that a user journey works. Evidence and log entries are
appended, never edited. A correction is a new entry that says what it
corrects.

## Attribution

Every log entry heading names who acted. A person's entry is
`## <time> · <type> · <name>`. An agent's entry is
`## <time> · <type> · <name> · <agent> · <model> · <session>`, with `-` for an
unknown value. Evidence frontmatter has the same `by`, `agent`, `model` and
`session` fields. Owner decisions name the owner and who recorded them.
`ignition attribution` reports all of it. Until agent runs are bound to
verified identities ([ROADMAP.md](ROADMAP.md)), an agent's name, model and
session are what it reported.

## Privacy

Never put secrets, credentials, private user data or raw production logs in
the records. Link to restricted evidence held in access-controlled storage,
and state who owns it and how long it's kept.
