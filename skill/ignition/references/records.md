# Records

## Where they live

Records live in the repository being worked on, at `<repo>/.ignition/`, and
are committed with the code they describe.

```text
.ignition/
├── PROJECT.md              # owner, repositories, environments, commands, behaviour checks
├── AGENTS.md               # product-specific rules for agents
├── architecture/           # the system as it is now (see diagrams.md)
├── decisions/D-….md        # product-wide decisions, referenced from many items
└── work/W-20261027-k3f9q2/ # one folder per item, of any kind
    ├── brief.md            # kind, risk, size, stage | request, problem, outcome, criteria, result
    ├── design.md           # created on entering design
    ├── plan.md             # created on entering plan
    ├── release.md          # created on entering release
    ├── evidence/E-….md     # one record per observation
    └── log.md              # append-only: stages, checkpoints, findings, changes, handoffs
```

The application repository stays the source of truth for code, tests,
migrations and configuration. Records link to exact commits, pull requests,
CI runs and deployments. They don't copy them.

`ignition init` makes `log.md` merge cleanly across branches (Git's union
driver) and keeps `.ignition` out of Docker images. For a public repository,
or a product that spans several, `ignition attach --app DIR --records-url URL`
keeps the records in a separate private repository checked out at
`.ignition/`, hidden from the application's Git. Use it only when the owner
asks.

## Format

Records are Markdown. Indexed records start with flat `key: value`
frontmatter between `---` lines: no nesting, lists or anchors. Sections are
`## ` headings. Template guidance sits in `<!-- -->` comments and doesn't
count as content: a section is filled when it has text beyond them.

IDs are `<prefix>-<YYYYMMDD>-<6 random characters>`: W for work, E for
evidence, D for decisions. Items link by ID.

## Evidence

An evidence record is an observation, not an intention. It names its claim,
its result (pass, fail or inconclusive), the exact revision (the commit
first, as `a1b2c3d` or `a1b2c3d (staging)`), the environment, the method, and
who recorded it. `criteria` lists the acceptance criteria it covers.
- Make the method a command someone else can run again on that revision.
- Fill `--limits` honestly: what this evidence doesn't show.
- Record failures as failures. Don't retry until green.
- "Merged", "deployed", "builds" or "returns 200" doesn't show that a user
  journey works. Check the journey.
- Evidence and log entries are appended, never edited. A correction is a new
  entry that says what it corrects.

## Attribution

Every log entry and evidence record names who acted: a name, and for an
agent its tool, exact model ID and session. `ignition whoami` shows what will
be recorded, and `ignition attribution` reports it. These are what the agent
reported. They aren't verified yet.

## Privacy

Never put secrets, credentials, private user data or raw production logs in
the records. Link to restricted evidence held in access-controlled storage.
