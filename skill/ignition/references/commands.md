# Commands

Every command accepts these flags:
- `--dir DIR`: the workspace. By default it's found from the current folder,
  or taken from `$IGNITION_DIR`.
- `--agent TOOL`: the agent tool, such as claude-code, codex, cursor,
  gemini-cli or copilot. Defaults to `$IGNITION_AGENT`, then detection from the
  environment (`AI_AGENT`, `CLAUDECODE`).
- `--model ID`: your exact model ID. Defaults to `$IGNITION_MODEL`. It's never
  detected or guessed.
- `--session ID`: your session ID. Defaults to `$IGNITION_SESSION`, then
  detection (Claude Code's session ID).
- `--by NAME`: the name recorded as acting. Defaults to `$IGNITION_ACTOR`,
  then the agent tool, then `git config user.name`.

Every log entry and evidence record stores the name, agent, model and
session. `ignition whoami` shows what will be recorded.

| Command | Does |
| --- | --- |
| `ignition status` | Lists every item: ID, kind, risk, stage, what it waits on, title |
| `ignition brief ID` | Shows the stage and route, problems, what the stage still needs, the owner checkpoint, files and recent log |
| `ignition new KIND "TITLE" --risk R --size S [--parent ID]` | Creates an item in its first stage. A child (`--parent`) can only be created once the initiative's solution design is done |
| `ignition advance ID [--owner NAME --decision "TEXT"]` | Moves to the next stage, creating its record file. The owner flags record a checkpoint decision. |
| `ignition log ID --type TYPE "TEXT"` | Appends to the log. TYPE is one of note, finding, decision, change, handoff, scope. |
| `ignition evidence ID --claim C --result R --revision V --environment E --method M [--observed O] [--limits L]` | Records evidence; R is one of pass, fail, inconclusive |
| `ignition decision "TITLE"` | Creates a product-wide decision record in `decisions/` |
| `ignition check [ID]` | Reports problems in the records; exits 1 if there are any |
| `ignition attribution [ID]` | Who built what: each person, agent and model, their entries, commits and sessions; owner decisions |
| `ignition whoami` | Shows the name, agent, model and session that will be recorded |
| `ignition init [--product NAME]` | Creates the workspace at `<repo>/.ignition`, tracked with the code. It adds `.ignition` to `.dockerignore` when there's a Dockerfile, and warns if Git ignores the folder. Run it once per repository. |
| `ignition attach --app DIR --records-url URL` | Alternative to `init` for public or multi-repository products: keeps records in a separate repository, checked out at `DIR/.ignition`. Only when the owner asks. |
| `ignition skill install [--for all\|claude\|agents] [--project DIR] [--to DIR]` | Installs this skill into `~/.claude/skills` (Claude Code) and `~/.agents/skills` (Codex, Cursor, Gemini CLI, Copilot and others), or into a project's folders with `--project` |
| `ignition guide` | Prints this skill's instructions |
| `ignition version` | Prints the version |

## What errors mean

| Message | Meaning and what to do |
| --- | --- |
| `stage not complete: <stage> needs: …` | Fix each listed item, then advance again. |
| `fix these first: passed <stage>, but …` | A record from an earlier stage is now incomplete or invalid. Restore it. Don't work around it. |
| `owner decision needed: leaving <stage> needs the owner's answer to: …` | Ask the user that question. Advance with `--owner` and `--decision` only after they answer. |
| `fill in "<Section>"` | The section has only guidance. Write content under the heading. |
| `no evidence with result fail/pass` | Run the check and record it with `ignition evidence`. |
| `log: no change entry` | Record your commit with `ignition log ID --type change`. |
| ``stage entries are written by `ignition advance` `` (or checkpoint) | You can't log these directly. Use `advance`. |
| ``evidence entries are written by `ignition evidence` `` | Use `ignition evidence`, not `ignition log`. |
| `Warning: recording <agent> with an unknown model` | Pass `--model <exact model ID>` if you know it. If you don't, leave it unknown. |
| `no Ignition workspace in …` | Run from the application repository, or pass `--dir`. If none exists, ask the user. |
