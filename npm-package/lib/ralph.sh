#!/bin/bash
set -e
set -o pipefail


# --- Hardened Environment (Anti-Interactive) ---
export CI=1
export GIT_TERMINAL_PROMPT=0
export npm_config_yes=true
export npm_config_audit=false
export npm_config_fund=false
export FORCE_COLOR=0
export TERM=dumb

# --- Configuration ---
ITERATIONS=5
WATCH_MODE="false"

# Watchdog Tunables
RALPH_MAX_TURN_SECONDS=${RALPH_MAX_TURN_SECONDS:-900}     # 15 min hard cap
RALPH_NO_OUTPUT_SECONDS=${RALPH_NO_OUTPUT_SECONDS:-180}   # 3m without output => stuck

# Model Priority: Build Phase (Implementation)
# (Defaults - can be overridden by config)
BUILD_MODELS_PREF=(
    "github-copilot/gpt-5.2-codex"
    "github-copilot/claude-sonnet-4.5"
    "github-copilot/gemini-3-pro-preview"
    "github-copilot-enterprise/gpt-5.2-codex"
    "github-copilot-enterprise/claude-sonnet-4.5"
    "github-copilot-enterprise/gemini-3-pro-preview"
    "openai/gpt-5.2-codex"
    "openai/gpt-5.1-codex-max"
    "google/gemini-3-pro-preview"
    "opencode/grok-code"
)

# Model Priority: Plan Phase (Reasoning/Architecture)
PLAN_MODELS_PREF=(
    "github-copilot/claude-opus-4.5"
    "github-copilot/gemini-3-pro-preview"
    "github-copilot-enterprise/claude-opus-4.5"
    "github-copilot-enterprise/gemini-3-pro-preview"
    "openai/gpt-5.2"
    "google/antigravity-claude-opus-4-5-thinking"
    "google/gemini-3-pro-preview"
    "opencode/glm-4.7-free"
)

# --- Setup Directories ---
# Resolve the directory where this script lives (lib/)
SOURCE=${BASH_SOURCE[0]}
while [ -L "$SOURCE" ]; do
  DIR=$( cd -P "$( dirname "$SOURCE" )" >/dev/null 2>&1 && pwd )
  SOURCE=$(readlink "$SOURCE")
  [[ $SOURCE != /* ]] && SOURCE=$DIR/$SOURCE
done
ENGINE_DIR=$( cd -P "$( dirname "$SOURCE" )" >/dev/null 2>&1 && pwd )

# PROJECT_DIR is where the user is running the command from
PROJECT_DIR="$(pwd)"
RUNS_DIR="$PROJECT_DIR/.ralph/runs"

mkdir -p "$RUNS_DIR"

# Cleanup trap
trap "pkill -P $$; exit" SIGINT SIGTERM

# --- Parse Args ---
MODE="default"
PROJECT_IDEA=""
FREE_MODE="false"

while [[ "$#" -gt 0 ]]; do
    case $1 in
        free)
            FREE_MODE="true"
            shift
            ;;
        new)
            MODE="new"
            PROJECT_IDEA="$2"
            shift 2
            ;;
        --watch) WATCH_MODE="true"; shift ;;
        *) ITERATIONS="$1"; shift ;;
    esac
done

echo "🐾 Vibepup v1.0.3 (CLI Mode)"
echo "   Engine:  $ENGINE_DIR"
echo "   Context: $PROJECT_DIR"

echo "   Tips:"
echo "   - Run 'vibepup free' for free-tier setup"
echo "   - Run 'vibepup new ""My idea""' to bootstrap a project"
echo "   - Run 'vibepup --tui' for a guided interface"

type -p opencode >/dev/null 2>&1 || {
    if command -v curl >/dev/null 2>&1; then
        UNAME=$(uname -s)
        if [ "$UNAME" = "Linux" ] || [ "$UNAME" = "Darwin" ]; then
            echo "⚠️  opencode not found. Installing..."
            curl -fsSL https://opencode.ai/install | bash || true
        fi
    fi
}

if ! command -v opencode >/dev/null 2>&1; then
    if [[ "$FREE_MODE" == "true" ]]; then
        echo "🔧 Free setup: installing opencode..."
        if command -v npm >/dev/null 2>&1; then
            npm install -g opencode-ai opencode-antigravity-auth || true
        else
            echo "❌ npm not found. Install Node.js or use:"
            echo "   curl -fsSL https://opencode.ai/install | bash"
            exit 127
        fi
    else
        echo "❌ opencode not found. Vibepup requires opencode to run."
        echo "   Install with one of:"
        echo "   - curl -fsSL https://opencode.ai/install | bash"
        echo "   - npm install -g opencode-ai"
        echo "   - brew install anomalyco/tap/opencode"
        echo "   Free-tier option:"
        echo "   - vibepup free"
        exit 127
    fi
fi

if [[ "$FREE_MODE" == "true" ]]; then
    echo "✨ Vibepup Free Setup"
    echo "   1) Installing auth plugin"
    if command -v npm >/dev/null 2>&1; then
        npm install -g opencode-antigravity-auth || true
    fi
    echo "   2) Starting Google auth"
    opencode auth login antigravity || true
    echo "   3) Refreshing models"
    opencode models --refresh || true
    echo "✅ Free setup complete. Run 'vibepup --watch' next."
    exit 0
fi

exec node "$(dirname "$0")/runner/index.js" "$@"
