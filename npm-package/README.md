# 🐾 Vibepup

> **"Fetch Code. Sit. Stay. Good Pup."**

Vibepup is a **Split-Brain Autonomous Agent Harness** that lives in your terminal. It turns any directory into an autonomous coding environment with **DX-first, vibe-coding energy**.

**Mascot:** Pummy the cyberpunk corgi.

![Vibepup Mascot](https://raw.githubusercontent.com/shantanusoam/ralph-project/refs/heads/gh-images/assets/corgi_Loop.png)

[![npm version](https://badge.fury.io/js/vibepup.svg)](https://badge.fury.io/js/vibepup)
![License](https://img.shields.io/npm/l/vibepup)

## ✨ The Vibe
Most AI agents are black boxes that overwrite your files and get stuck in loops. Vibepup is different. Loyal, friendly, a little meme-y — and built for real dev workflows.

**Selling Points:**
- DX-first onboarding
- Vibe-coding friendly
- Safe, loop-resistant agent harness
- Minimal setup, works everywhere

*   **🧠 Split-Brain**: Keeps your instructions (`prd.md`) separate from internal state (`prd.state.json`). Edit tasks mid-run without breaking the agent.
*   **🛡️ Anti-Wizard**: Refuses to run interactive commands that hang shells. Vibepup forces clarity.
*   **⚡ DX-First**: Optimized for fast iteration, visibility, and zero-friction adoption.
*   **🌈 Cyberpunk Corgi**: Cute, loyal, and ready to ship.
*   **🧩 Friendly + Meme-y**: The tool feels fun without being unserious.

![Checklist Mode](https://raw.githubusercontent.com/shantanusoam/ralph-project/refs/heads/gh-images/assets/corgi_checklist.png)

## 🚀 Get Started

### ✅ Works Everywhere
Linux, macOS, and Windows (via Git Bash or WSL). Vibepup is portable and requires only Bash + `opencode`.

### 1. Install
```bash
npm install -g vibepup
```

### 1b. bunx (no global install)
```bash
bunx vibepup --watch
```

### 1c. bun global install (optional)
```bash
bun add -g vibepup
```

### 2. Fetch!
Go to any empty folder and tell Vibepup what to build.

```bash
mkdir my-app
cd my-app
vibepup new "A react app for tracking my plant watering schedule"
```

Using bunx:
```bash
bunx vibepup new "A react app for tracking my plant watering schedule"
```

Using bun global install:
```bash
vibepup new "A react app for tracking my plant watering schedule"
```

Vibepup will:
1.  🏗️ **Plan**: Map out the architecture in `repo-map.md`.
2.  📝 **Draft**: Create a `prd.md` checklist.
3.  🔨 **Build**: Start checking off items one by one.

### 3. Watch Him Work
```bash
vibepup --watch
```
In watch mode, Vibepup keeps working until the PRD is done. If you edit `prd.md` (e.g., add "- [ ] Add dark mode"), he smells the change and gets back to work immediately.

## ⚙️ Configuration
Vibepup works out of the box, but you can tune his brain in `~/.config/ralph/config.json`:

```json
{
  "build_models": [
    "github-copilot/gpt-5.2-codex",
    "openai/gpt-4o"
  ],
  "plan_models": [
    "github-copilot/claude-opus-4.5"
  ]
}
```

## 🏗️ Architecture
*   **Plan Mode**: When `repo-map.md` is missing, Vibepup explores and plans.
*   **Build Mode**: When `repo-map.md` exists, Vibepup executes tasks from `prd.md`.

## 🎨 Mascot Prompt (NanoBanan)
Use this prompt to generate more assets for **Pummy**:

> "Cute loyal corgi mascot wearing a tiny tool belt, holding a glowing code scroll, subtle cyberpunk neon accents (teal + magenta), warm orange fur, friendly smile, modern flat vector style, clean bold outlines, no text."

## License
MIT
