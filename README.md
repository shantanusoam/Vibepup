# Ralph - Autonomous Coding Agent

Ralph is a global CLI tool that turns any directory into an autonomous coding environment. It works on both new (greenfield) and existing (brownfield) projects.

## Features

- **Global CLI**: Run `ralph` anywhere.
- **Brownfield Optimized**: Uses `repo-map.md` to navigate large codebases efficiently.
- **Safety First**: Includes Git Sync and Peer Review steps.
- **Watch Mode**: `ralph --watch` runs indefinitely, waking up when you edit your `prd.json`.
- **Robust Fallback**: Falls back to Codex/Gemini if the primary model fails.

## Installation

1.  Clone this repo:
    ```bash
    git clone https://github.com/shantanusoam/ralph-project.git ~/ralph-project
    ```
2.  Symlink the binary:
    ```bash
    sudo ln -s ~/ralph-project/global/ralph /usr/local/bin/ralph
    ```
    *(Or add `~/ralph-project/global` to your `$PATH`)*

## Usage

### 1. Initialize a Project
Go to any folder and run:
```bash
cd ~/my-project
ralph
```
This creates:
- `prd.json`: Your task list.
- `repo-map.md`: Architecture cache.
- `progress.txt`: Log file.

### 2. Run the Loop
Run for 10 iterations:
```bash
ralph 10
```

Run in **Watch Mode** (Ralph waits for you to add tasks):
```bash
ralph --watch
```

## Configuration

Ralph uses a local configuration strategy:
- **Global Logic**: `~/ralph-project/global/ralph` & `prompt.md`.
- **Local Context**: `prd.json` in your project folder.

## Architecture

- **Primary Brain**: Antigravity (Claude Sonnet 4.5 Thinking)
- **Reviewer**: Gemini 2.5 Pro
- **Fallback**: GPT-5.1 Codex
