# Multi-Agent Codex MCP Runner

This is a standalone multi-agent harness that uses Codex CLI only (no API key).

## Setup
1. Ensure `codex` is on your PATH.
2. Run the script.

Example:
```bash
python3 run.py
```

Use a custom task list and project dir:
```bash
python3 run.py --project /path/to/project --tasks /path/to/tasks.txt
```

## What It Does
- Runs a Project Manager, Designer, Frontend, Backend, and Tester in sequence.
- Uses Codex CLI for each role.
- Produces output files in `/design`, `/frontend`, `/backend`, and `/tests`.
