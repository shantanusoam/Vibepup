import argparse
import os
import subprocess
import sys
from pathlib import Path


DEFAULT_TASKS = """Goal: Build a tiny browser game to showcase a multi-agent workflow.

High-level requirements:
- Single-screen game called "Bug Busters".
- Player clicks a moving bug to earn points.
- Game ends after 20 seconds and shows final score.
- Optional: submit score to a simple backend and display a top-10 leaderboard.

Roles:
- Designer: create a one-page UI/UX spec and basic wireframe.
- Frontend Developer: implement the page and game logic.
- Backend Developer: implement a minimal API (GET /health, GET/POST /scores).
- Tester: write a quick test plan and a simple script to verify core routes.

Constraints:
- No external database—memory storage is fine.
- Keep everything readable for beginners; no frameworks required.
- All outputs should be small files saved in clearly named folders.
"""


def run_codex(prompt: str, project_dir: Path) -> int:
    cmd = [
        "codex",
        "exec",
        "--full-auto",
        "--skip-git-repo-check",
        "-C",
        str(project_dir),
        "-",
    ]
    proc = subprocess.run(
        cmd,
        input=prompt,
        text=True,
        capture_output=True,
    )
    if proc.stdout:
        print(proc.stdout)
    if proc.stderr:
        print(proc.stderr, file=sys.stderr)
    return proc.returncode


def require_files(project_dir: Path, rel_paths: list[str]) -> list[str]:
    missing = []
    for rel_path in rel_paths:
        if not (project_dir / rel_path).exists():
            missing.append(rel_path)
    return missing


def read_tasks(tasks_path: Path | None) -> str:
    if tasks_path is None:
        return DEFAULT_TASKS
    return tasks_path.read_text(encoding="utf-8")


def main() -> int:
    parser = argparse.ArgumentParser(description="Codex-only multi-agent runner.")
    parser.add_argument(
        "--project",
        type=Path,
        default=Path.cwd(),
        help="Target project directory (default: current directory).",
    )
    parser.add_argument(
        "--tasks",
        type=Path,
        default=None,
        help="Path to a task list file (plain text).",
    )
    parser.add_argument(
        "--max-attempts",
        type=int,
        default=3,
        help="Max attempts per role before failing.",
    )
    args = parser.parse_args()

    project_dir = args.project.resolve()
    tasks = read_tasks(args.tasks)

    pm_prompt = f"""You are the Project Manager.
Objective: Convert the input task list into three project-root files the team will execute against.

Deliverables (write in project root):
- REQUIREMENTS.md: concise summary of product goals, target users, key features, and constraints.
- TEST.md: tasks with [Owner] tags (Designer, Frontend, Backend, Tester) and clear acceptance criteria.
- AGENT_TASKS.md: one section per role containing:
  - Project name
  - Required deliverables (exact file names and purpose)
  - Key technical notes and constraints

Process:
- Resolve ambiguities with minimal, reasonable assumptions. Be specific so each role can act without guessing.
- Do not create folders. Only create REQUIREMENTS.md, TEST.md, AGENT_TASKS.md.
- Do not output status updates; just perform the work.

Task list:
{tasks}
"""

    designer_prompt = """You are the Designer.
Your only source of truth is AGENT_TASKS.md and REQUIREMENTS.md from the Project Manager.
Do not assume anything that is not written there.

Deliverables (write to /design):
- design_spec.md – a single page describing the UI/UX layout, main screens, and key visual notes as requested in AGENT_TASKS.md.
- wireframe.md – a simple text or ASCII wireframe if specified.

Keep the output short and implementation-friendly. Do not add features beyond the documents.
"""

    frontend_prompt = """You are the Frontend Developer.
Read AGENT_TASKS.md and design_spec.md. Implement exactly what is described there.

Deliverables (write to /frontend):
- index.html – main page structure
- styles.css or inline styles if specified
- main.js or game.js if specified

Follow the Designer’s DOM structure and any integration points given by the Project Manager.
Do not add features or branding beyond the provided documents.
"""

    backend_prompt = """You are the Backend Developer.
Read AGENT_TASKS.md and REQUIREMENTS.md. Implement the backend endpoints described there.

Deliverables (write to /backend):
- package.json – include a start script if requested
- server.js – implement the API endpoints and logic exactly as specified

Keep the code as simple and readable as possible. No external database.
"""

    tester_prompt = """You are the Tester.
Read AGENT_TASKS.md and TEST.md. Verify that the outputs of the other roles meet the acceptance criteria.

Deliverables (write to /tests):
- TEST_PLAN.md – bullet list of manual checks or automated steps as requested
- test.sh or a simple automated script if specified

Keep it minimal and easy to run.
"""

    steps = [
        ("Project Manager", pm_prompt, ["REQUIREMENTS.md", "TEST.md", "AGENT_TASKS.md"]),
        ("Designer", designer_prompt, ["design/design_spec.md"]),
        ("Frontend Developer", frontend_prompt, ["frontend/index.html"]),
        ("Backend Developer", backend_prompt, ["backend/server.js"]),
        ("Tester", tester_prompt, ["tests/TEST_PLAN.md"]),
    ]

    for role, prompt, required in steps:
        attempts = 0
        while attempts < args.max_attempts:
            print(f"== {role} (attempt {attempts + 1}/{args.max_attempts}) ==")
            exit_code = run_codex(prompt, project_dir)
            if exit_code != 0:
                attempts += 1
                continue
            missing = require_files(project_dir, required)
            if not missing:
                break
            attempts += 1
            print(f"Missing required files: {', '.join(missing)}")
        else:
            print(f"{role} failed to produce required files after {args.max_attempts} attempts.")
            return 1

    print("Multi-agent run complete.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
