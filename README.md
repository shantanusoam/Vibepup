# 🤖 Ralph v3.1 - The "Split-Brain" Autonomous Agent

Ralph is a robust, global CLI harness that turns any directory into an autonomous coding environment. It is designed for **Developer Experience (DX)**, safety, and resilience.

Unlike standard agents that get stuck in loops or overwrite their own memory, Ralph uses a **Split-Brain Architecture** to separate *Human Intent* from *Agent State*.

## ✨ Key Features

### 🧠 Split-Brain Architecture
*   **`prd.md` (The Boss)**: A simple Markdown checklist for you. Edit this anytime.
    *   `- [ ] Add login page`
*   **`prd.state.json` (The Worker)**: A machine-managed state file.
    *   Tracks attempts, verification status, and errors without polluting your checklist.
    *   *Result:* You can edit requirements mid-run without breaking the agent's brain.

### 🌓 Phase Detection
Ralph automatically switches modes based on project maturity:
1.  **PLAN Mode 🏗️**: Triggered when `repo-map.md` is empty. The agent explores `ls -R`, reads key files, and maps the architecture. **No coding allowed.**
2.  **BUILD Mode 🔨**: Triggered when the map exists. The agent executes tasks from `prd.md` one by one.

### ⚡ Hybrid Model Priority
Ralph uses a "Gold Standard" fallback chain. It automatically discovers which models you have access to (via `opencode`) and uses the best available:
1.  **Github Copilot** (Standard & Enterprise) - *Speed & Context*
2.  **OpenAI** (GPT-5.2/4o) - *Raw Power*
3.  **Google** (Gemini 3 Pro) - *Huge Context Window*
4.  **Zen** (Grok/OpenCode) - *Resilient Fallback*

### 👁️ Real-Time Visibility
*   **Streaming Output**: See exactly what the agent is thinking and running in real-time. No more staring at a blank screen.
*   **Anti-Wizard Protocol**: Ralph is strictly forbidden from running interactive CLIs (like `npm init` without `-y`) to prevent hanging.

---

## 🚀 Installation

### 1. Clone & Setup
Clone this repository to your preferred location (e.g., `~/Projects/personal/ralph-project`):

```bash
git clone https://github.com/shantanusoam/ralph-project.git ~/Projects/personal/ralph-project
```

### 2. Global Symlink
Make Ralph accessible from anywhere. **Important:** Use the absolute path.

```bash
# Fix permissions first
chmod +x ~/Projects/personal/ralph-project/global/ralph

# Create the link
sudo ln -sf ~/Projects/personal/ralph-project/global/ralph /usr/local/bin/ralph
```

### 3. Verify
```bash
ralph --help
# Should output: 🤖 Ralph v3.1 (Split-Brain Architecture) ...
```

---

## 🎮 Usage

### Initialize a New Project
Navigate to any folder (empty or existing) and run Ralph:

```bash
cd ~/my-new-app
ralph 1
```

Ralph will detect missing files and initialize:
- `prd.md` (Your task list)
- `repo-map.md` (Architecture memory)
- `prd.state.json` (Internal state)

### The Workflow

1.  **Edit `prd.md`**: Add your tasks.
    ```markdown
    - [ ] Create Next.js app structure
    - [ ] Add Tailwind CSS
    ```
2.  **Run Ralph**:
    ```bash
    ralph 5   # Run for 5 iterations
    ```
3.  **Watch Mode (Recommended)**:
    ```bash
    ralph --watch
    ```
    In this mode, Ralph runs tasks until done. If you edit `prd.md` (e.g., add a new feature), Ralph **automatically detects the change**, resets the loop, and starts working on the new task immediately.

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

**"Command not found: ralph"**
- Check your PATH: `echo $PATH`
- Verify the symlink: `ls -l /usr/local/bin/ralph`
- Ensure the target file is executable: `chmod +x ...`

**"Agent gets stuck on `npm init`"**
- Ralph v3.1 has "Anti-Wizard" rules, but if it happens, kill the process (`Ctrl+C`) and add the config file manually (e.g., `package.json`) so Ralph can skip the interactive step.

**"ModelNotFoundError"**
- Run `opencode models --refresh` to update your local model cache. Ralph auto-discovers available models at startup.
