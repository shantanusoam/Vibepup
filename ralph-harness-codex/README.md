# Ralph Harness for Codex CLI

This repo is a minimal, production-grade loop harness for Codex CLI.

## Files
- `prd.json`: Your scope and stop condition. Set `passes: false` until done.
- `progress.txt`: Short log appended each iteration.
- `ralph-prompt.md`: The reusable prompt for Codex.
- `ralph-once.sh`: Run one HITL iteration.
- `afk-ralph.sh`: Run a capped AFK loop.

## Setup
1. Edit `prd.json` with your real tasks.
2. Optionally edit `ralph-prompt.md` for repo-specific rules.
3. Run a single iteration:
   - `./ralph-once.sh`
4. Go AFK with a capped loop:
   - `./afk-ralph.sh 5`

## Notes
- Codex runs in sandboxed workspace-write mode (`--full-auto`).
- Scope: operate only within this repo unless the PRD explicitly expands it.
- If you want a different model, set `CODEX_MODEL`:
  - `CODEX_MODEL=o3 ./ralph-once.sh`
- If Codex is not on PATH, set `CODEX_BIN`:
  - `CODEX_BIN=/path/to/codex ./ralph-once.sh`
