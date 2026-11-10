# Project: ignition

## Owner

Mohamed Saleem (GitHub: mohsalsaleem). He decides at every checkpoint.

## Repositories

- `mohsalsaleem/ignition-engineering` (this one, public): the Ignition CLI, the agent skill, and these records.
- `mohsalsaleem/OpenAppPlatform` (public): OAP, the runtime the platform initiative builds on. It keeps its own plan in `docs/architecture/`.

## Environments

Local only. Ignition is a CLI installed with `go install`; there is no server. A release is a `vX.Y.Z` tag on `main`, pushed by the owner.

## Commands

- Fast check and full test: `make check` (`go vet` and `go test ./...`).
- Build: `make build`, which writes `bin/ignition`.
- Critical journey: `go run ./cmd/ignition check --dir examples/oap` must report no problems.

## Behaviour checks

`go test ./...` runs in a few seconds with no network. `internal/workspace` tests drive whole work items through their stages against a temporary workspace; `cmd/ignition` tests drive the CLI itself through `run(args, out, errw)`. To try a flow by hand, build the binary and point it at a scratch folder with `--dir`.

## Risks and boundaries

- The repository is public: records here are public too. No secrets, credentials or private data.
- Agents push only to their own branch. The owner merges, tags and releases.
- A removed kind, stage or command breaks existing workspaces: say so in the pull request.
