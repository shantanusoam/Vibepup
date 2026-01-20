# Repository Architecture Map

## Overview
Ralph is a global CLI tool for autonomous coding. It works by orchestrating AI models to complete tasks defined in a PRD (Product Requirements Document), with built-in fallback chains, phase detection, and code review.

## Core Architecture: "Split-Brain" Design (v3.0)
- **Human Layer**: `prd.md` - read-only checklist that humans edit
- **Agent Layer**: `prd.state.json` - agent's private state tracking (attempts, verification status)
- **Phase Detection**: Automatic PLAN vs BUILD mode based on `repo-map.md` state

## Directory Structure

```
ralph-project/
├── global/                          # Engine (global CLI logic)
│   ├── ralph                        # Main CLI executable (bash, 215 lines)
│   ├── prompt.md                    # System prompt for agent instructions
│   └── agents/
│       └── reviewer.md              # Code reviewer agent prompt
├── ralph-harness/                   # Legacy v1 harness (opencode-based)
│   ├── ralph.sh                     # Old loop script with fallback chain
│   ├── prompt.md                    # Agent prompt
│   ├── prd.json                     # Task definitions
│   ├── helloworld.py                # Demo target
│   └── test_helloworld.py           # Demo test
├── ralph-harness-codex/             # Legacy v2 harness (codex CLI-based)
│   ├── ralph-once.sh                # Single iteration runner
│   ├── afk-ralph.sh                 # AFK loop wrapper
│   └── ralph-prompt.md              # Codex-specific prompt
├── .ralph/                          # Runtime data (per-project)
│   └── runs/
│       ├── iter-XXXX/               # Per-iteration logs
│       │   ├── agent_response.txt   # Raw agent output
│       │   └── progress.tail.log    # Context window (last 200 lines)
│       └── latest -> iter-XXXX      # Symlink for debugging
├── prd.md                           # Human task checklist (current project)
├── prd.state.json                   # Agent state tracking
├── repo-map.md                      # This file (architecture cache)
├── progress.log                     # Full history log
└── README.md                        # Project documentation
```

## Key Files

### `global/ralph` (Main CLI)
- Entry point for `ralph` command
- Phase detection: PLAN (empty repo-map) vs BUILD (repo-map exists)
- Model priority chains:
  - **PLAN_MODELS**: opus-4.5, gpt-5.2, gemini-2.5-pro (deep reasoning)
  - **BUILD_MODELS**: gpt-5.2-codex, claude-sonnet-4.5, gemini-3-pro, grok-code, gpt-4o (fast coding)
- Watch mode: `ralph --watch` - monitors prd.md for changes
- Iteration limit: `ralph 10` - runs N iterations

### `global/prompt.md`
System prompt that instructs the agent on:
- Phase awareness (PLAN vs BUILD)
- PRD contract (read prd.md, update prd.state.json)
- Verification requirement before marking tasks complete
- Completion signal: `<promise>COMPLETE</promise>`

### `global/agents/reviewer.md`
Code review agent that:
- Analyzes diffs before commit
- Checks against PRD requirements
- Outputs `<review>PASS</review>` or `<review>FAIL</review>`

## Workflow

1. **Initialization**: `ralph` creates prd.md, repo-map.md, prd.state.json if missing
2. **Phase Detection**: Empty repo-map.md triggers PLAN mode
3. **PLAN Mode**: Agent explores codebase, populates repo-map.md
4. **BUILD Mode**: Agent picks first unchecked task from prd.md, implements it
5. **Verification**: Agent runs tests/build, updates prd.state.json
6. **Completion**: When all tasks checked, agent outputs `<promise>COMPLETE</promise>`

## Model Fallback Strategy
Models tried in priority order until one succeeds:
1. Primary model for current phase
2. Fallback models (phase-specific)
3. External codex CLI (legacy support)

## Legacy Harnesses
- `ralph-harness/`: Demo using opencode with multi-model fallback
- `ralph-harness-codex/`: Demo using codex CLI with `--full-auto`
