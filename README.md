# 🐾 Vibepup (formerly Ralph)

> "Fetch Code. Sit. Stay. Good Pup."

![Corgi Pup Illustration](https://raw.githubusercontent.com/shantanusoam/Vibepup/refs/heads/gh-images/assets/corgi_pup_ilustration.png)

[![npm version](https://badge.fury.io/js/vibepup.svg)](https://badge.fury.io/js/vibepup)
![License](https://img.shields.io/npm/l/vibepup)

npm: https://www.npmjs.com/package/vibepup

Vibepup is a robust, global CLI harness that turns any directory into an autonomous coding environment. It is designed for **Developer Experience (DX)**, safety, and vibe-coding resilience.

**Mascot:** Pummy the cyberpunk corgi.

![Pummy Loop](https://raw.githubusercontent.com/shantanusoam/ralph-project/refs/heads/gh-images/assets/corgi_Loop.png)

Unlike standard agents that get stuck in loops or overwrite their own memory, Vibepup uses a **Split-Brain Architecture** to separate *Human Intent* from *Agent State*.

## ✨ The Vibe

**Selling Points:**
- DX-first onboarding
- Vibe-coding friendly
- Safe, loop-resistant agent harness
- Minimal setup, works everywhere
- Loyal helper with a cyberpunk-cute mascot

![Pummy Checklist](https://raw.githubusercontent.com/shantanusoam/ralph-project/refs/heads/gh-images/assets/corgi_checklist.png)

## ✨ Key Features

### 🧠 Split-Brain Architecture
*   **`prd.md` (The Boss)**: A simple Markdown checklist for you. Edit this anytime.
    *   `- [ ] Add login page`
*   **`prd.state.json` (The Worker)**: A machine-managed state file.
    *   Tracks attempts, verification status, and errors without polluting your checklist.
    *   *Result:* You can edit requirements mid-run without breaking the agent's brain.

### 🛡️ Safety & Hardening
*   **Watchdog Protocol**: A built-in safety monitor that automatically kills any process taking longer than 15 minutes or hanging without output for 3 minutes.
*   **Infinite Loop Protection**: In Watch Mode, if Ralph completes all tasks, he enters a low-power "Wait State", pausing execution until you modify `prd.md`. No more burning tokens on empty cycles.

### 🌓 Phase Detection
Ralph automatically switches modes based on project maturity:
1.  **PLAN Mode 🏗️**: Triggered when `repo-map.md` is empty. The agent explores `ls -R`, reads key files, and maps the architecture. **No coding allowed.**
2.  **BUILD Mode 🔨**: Triggered when the map exists. The agent executes tasks from `prd.md` one by one.

### 🧠 Specialized Model Roles
Ralph assigns the "Right Brain" to the right task (configurable in `~/.config/opencode/oh-my-opencode.json`):
*   **The Architect (Sisyphus)**: `gpt-5.2-codex` for core logic and orchestration.
*   **The Designer (Frontend)**: `gemini-3-pro-preview` for massive context window and UI tasks.
*   **The Explorer (Explore)**: `grok-code-fast-1` for rapid codebase search.
*   **The Sage (Oracle)**: `claude-opus-4.5` for deep reasoning and architecture validation.

### 👁️ Real-Time Visibility
*   **Streaming Output**: See exactly what the agent is thinking and running in real-time. No more staring at a blank screen.
*   **Anti-Wizard Protocol**: Ralph is strictly forbidden from running interactive CLIs (like `npm init` without `-y`) to prevent hanging.

---

## 🚀 Installation

### 1. npm (recommended)
```bash
npm install -g vibepup
```

### 2. bunx (no global install)
```bash
bunx vibepup --watch
```

### 2b. TUI mode (optional)
```bash
vibepup --tui
```

### 3. Clone & Setup (engine-only)
Clone this repository to your preferred location (e.g., `~/Projects/personal/ralph-project`):

```bash
git clone https://github.com/shantanusoam/ralph-project.git ~/Projects/personal/ralph-project
```

### 4. Global Symlink (engine-only)
Make Ralph accessible from anywhere. **Important:** Use the absolute path.

```bash
# Fix permissions first
chmod +x ~/Projects/personal/ralph-project/global/ralph

# Create the link
sudo ln -sf ~/Projects/personal/ralph-project/global/ralph /usr/local/bin/ralph
```

### 5. Verify (engine-only)
```bash
ralph --help
# Should output: 🤖 Ralph v3.2 (Split-Brain Architecture) ...
```

---

## 🎮 Usage

### Initialize a New Project
Navigate to any folder (empty or existing) and run Vibepup:

```bash
cd ~/my-new-app
vibepup 1
```

Vibepup will detect missing files and initialize:
- `prd.md` (Your task list)
- `repo-map.md` (Architecture memory)
- `prd.state.json` (Internal state)

### The Workflow

1.  **Edit `prd.md`**: Add your tasks.
    ```markdown
    - [ ] Create Next.js app structure
    - [ ] Add Tailwind CSS
    ```
2.  **Run Vibepup**:
    ```bash
    vibepup 5   # Run for 5 iterations
    ```
3.  **Watch Mode (Recommended)**:
    ```bash
    vibepup --watch
    ```
    In this mode, Vibepup runs tasks until done. If you edit `prd.md` (e.g., add a new feature), Vibepup **automatically detects the change**, resets the loop, and starts working on the new task immediately.

---

## ⚙️ Configuration

Ralph's "Brain" is located in `global/ralph`. You can tweak the model priority lists directly in the script:

```bash
# global/ralph

# Prioritize your preferred models here
BUILD_MODELS_PREF=(
    "github-copilot/gpt-5.2-codex"
    "openai/gpt-5.2-codex"
    ...
)
```

## 🛠️ Troubleshooting

**"Command not found: vibepup"**
- Check your PATH: `echo $PATH`
- Reinstall: `npm install -g vibepup`
- Ensure `vibepup` is on your PATH

**"Agent gets stuck on `npm init`"**
- Vibepup has "Anti-Wizard" rules, but if it happens, kill the process (`Ctrl+C`) and add the config file manually (e.g., `package.json`) so Vibepup can skip the interactive step.

**"ModelNotFoundError"**
- Run `opencode models --refresh` to update your local model cache. Vibepup auto-discovers available models at startup.
