# Vibepup: The Loyal AI Harness for Vibe Coders

> "Fetch Code. Sit. Stay. Good Pup."

Vibepup is a Split-Brain autonomous agent harness designed for the vibe-coding era. It turns any directory into a safe, managed coding environment where Human Intent and Agent State never conflict.

![Pummy the Cyberpunk Corgi](https://raw.githubusercontent.com/shantanusoam/ralph-project/refs/heads/gh-images/assets/corgi_Loop.png)

## Why Vibepup

Most AI agents are black boxes that overwrite your prompts, get stuck in infinite loops, or hallucinate dependencies. Vibepup is different: loyal, safe, and DX-first.

### Split-Brain Architecture
Vibepup separates the Boss from the Worker.
- The Boss (`prd.md`): A clean Markdown checklist of what you want. Edit this anytime.
- The Worker (`prd.state.json`): A hidden JSON file where the agent tracks its retries and failures.
- Result: You can change requirements mid-flight without breaking the agent's brain.

### Anti-Wizard Protocol
Vibepup refuses to run interactive commands that hang your terminal. If a process hangs for 3 minutes, the Watchdog stops it.

### DX-First
No complex config. No Docker containers required. Just a CLI that respects your workflow.

## Quick Start

No install required. Run:

```bash
bunx vibepup new "A retro-style pomodoro timer using React and Tailwind"
```

Or install globally:

```bash
npm install -g vibepup
```

## How It Works

1. Plan Mode: If the project is new, Vibepup explores the folder and builds a `repo-map.md` of your architecture.
2. Build Mode: Once mapped, it reads your `prd.md` and starts checking off boxes.
3. Watch Mode: Run `vibepup --watch`. It sits quietly until you edit the PRD, then wakes up to code.

![Pummy Checklist](https://raw.githubusercontent.com/shantanusoam/ralph-project/refs/heads/gh-images/assets/corgi_checklist.png)
