#!/bin/bash
set -e


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

while [[ "$#" -gt 0 ]]; do
    case $1 in
        new)
            MODE="new"
            PROJECT_IDEA="$2"
            shift 2
            ;;
        --watch) WATCH_MODE="true"; shift ;;
        *) ITERATIONS="$1"; shift ;;
    esac
done

echo "🐾 Vibepup v1.0 (CLI Mode)"
echo "   Engine:  $ENGINE_DIR"
echo "   Context: $PROJECT_DIR"

# --- Smart Model Discovery ---
get_available_models() {
    local PREF_MODELS=("$@")
    local AVAILABLE_MODELS=()
    
    echo "🔍 Verifying available models..." >&2
    
    local ALL_MODELS=$(opencode models --refresh | grep -E "^[a-z0-9-]+\/[a-z0-9.-]+" || true)
    
    for PREF in "${PREF_MODELS[@]}"; do
        if echo "$ALL_MODELS" | grep -q "^$PREF$"; then
            AVAILABLE_MODELS+=("$PREF")
        fi
    done
    
    if [ ${#AVAILABLE_MODELS[@]} -eq 0 ]; then
        echo "⚠️  No preferred models found. Falling back to generic discovery." >&2
        AVAILABLE_MODELS+=($(echo "$ALL_MODELS" | grep "gpt-4o" | head -n 1))
        AVAILABLE_MODELS+=($(echo "$ALL_MODELS" | grep "claude-sonnet" | head -n 1))
    fi
    
    echo "${AVAILABLE_MODELS[@]}"
}

echo "   Configuring Build Models..."
read -r -a BUILD_MODELS <<< "$(get_available_models "${BUILD_MODELS_PREF[@]}")"
echo "   Build Chain: ${BUILD_MODELS[*]}"

echo "   Configuring Plan Models..."
read -r -a PLAN_MODELS <<< "$(get_available_models "${PLAN_MODELS_PREF[@]}")"
echo "   Plan Chain:  ${PLAN_MODELS[*]}"


# --- Phase 0: The Architect (Genesis) ---
if [ "$MODE" == "new" ]; then
    echo ""
    echo "🏗️  Phase 0: The Architect"
    echo "   Idea: $PROJECT_IDEA"
    
    ARCHITECT_MODEL="${PLAN_MODELS[0]}"
    echo "   Using: $ARCHITECT_MODEL"
    
    # NOTE: We assume agents/architect.md is in lib/agents/
    opencode run "PROJECT IDEA: $PROJECT_IDEA" \
        --file "$ENGINE_DIR/agents/architect.md" \
        --agent general \
        --model "$ARCHITECT_MODEL"
        
    echo "✅ Architect initialization complete."
fi


# --- Initialization & Migration ---
if [ ! -f "$PROJECT_DIR/prd.md" ]; then
    if [ -f "$PROJECT_DIR/prd.json" ]; then
        echo "🔄 Migrating legacy prd.json to prd.md..."
        jq -r '.[] | "- [ ] " + .description' "$PROJECT_DIR/prd.json" > "$PROJECT_DIR/prd.md"
        mv "$PROJECT_DIR/prd.json" "$PROJECT_DIR/prd.json.bak"
    else
        echo "⚠️  No prd.md found. Initializing..."
        cat > "$PROJECT_DIR/prd.md" <<INNEREOF
# Product Requirements Document (PRD)

- [ ] Initialize repo-map.md with project architecture
- [ ] Setup initial project structure
INNEREOF
    fi
fi

if [ ! -f "$PROJECT_DIR/repo-map.md" ]; then
    touch "$PROJECT_DIR/repo-map.md"
fi

if [ ! -f "$PROJECT_DIR/prd.state.json" ]; then
    echo "{}" > "$PROJECT_DIR/prd.state.json"
fi

touch "$PROJECT_DIR/progress.log"

# --- Helper Functions ---

detect_phase() {
    if [ ! -s "$PROJECT_DIR/repo-map.md" ]; then
        echo "PLAN"
        return
    fi
    echo "BUILD"
}

get_current_prd_hash() {
    md5sum "$PROJECT_DIR/prd.md" | awk '{print $1}'
}

prepare_iteration_context() {
    local ITER_ID="$1"
    local ITER_DIR="$RUNS_DIR/$ITER_ID"
    mkdir -p "$ITER_DIR"
    
    tail -n 200 "$PROJECT_DIR/progress.log" > "$ITER_DIR/progress.tail.log"
    
    rm -f "$RUNS_DIR/latest"
    ln -s "$ITER_DIR" "$RUNS_DIR/latest"
    
    echo "$ITER_DIR"
}

run_with_watchdog () {
  local log="$1"; shift
  : > "$log"

  # run in background and capture output
  ( "$@" 2>&1 | tee -a "$log" ) &
  local pid=$!
  local start=$(date +%s)
  local last=$start
  local last_size=0

  while kill -0 "$pid" 2>/dev/null; do
    sleep 5
    local now=$(date +%s)
    local size=$(wc -c < "$log" 2>/dev/null || echo 0)

    if (( size != last_size )); then
      last="$now"
      last_size="$size"
    fi

    if (( now - start > RALPH_MAX_TURN_SECONDS )); then
      echo "[RALPH] TIMEOUT: killing opencode turn" >> "$log"
      kill -INT "$pid" 2>/dev/null || true
      sleep 3
      kill -TERM "$pid" 2>/dev/null || true
      sleep 1
      kill -KILL "$pid" 2>/dev/null || true
      return 124
    fi

    if (( now - last > RALPH_NO_OUTPUT_SECONDS )); then
      echo "[RALPH] NO OUTPUT: likely waiting for input / hung tool" >> "$log"
      kill -INT "$pid" 2>/dev/null || true
      sleep 3
      kill -TERM "$pid" 2>/dev/null || true
      sleep 1
      kill -KILL "$pid" 2>/dev/null || true
      return 125
    fi
  done

  wait "$pid"
}

run_agent() {
    local MODEL="$1"
    local PHASE="$2"
    local ITER_DIR="$3"
    
    local SYSTEM_PROMPT="$ENGINE_DIR/prompt.md"
    local PROMPT_SUFFIX=""
    local EXTRA_ARGS=()

    if [ "$DESIGN_MODE" == "true" ]; then
        echo "   🎨 Design Mode Active: Injecting frontend-design skill..."
        EXTRA_ARGS+=( "--file" "$HOME/.config/opencode/skills/frontend-design.md" )
        PROMPT_SUFFIX="MODE: DESIGN + BUILD. Apply the frontend-design skill guidelines to all work."
    elif [ "$PHASE" == "PLAN" ]; then
        PROMPT_SUFFIX="MODE: PLAN. Focus on exploring and mapping. Do NOT write implementation code yet."
    else
        PROMPT_SUFFIX="MODE: BUILD. Focus on completing tasks in prd.md."
    fi

    echo "   Thinking..."
    
    run_with_watchdog "$ITER_DIR/agent_response.txt" \
        opencode run "Proceed with task. $PROMPT_SUFFIX" \
        --file "$SYSTEM_PROMPT" \
        --file "$PROJECT_DIR/prd.md" \
        --file "$PROJECT_DIR/prd.state.json" \
        --file "$PROJECT_DIR/repo-map.md" \
        --file "$ITER_DIR/progress.tail.log" \
        "${EXTRA_ARGS[@]}" \
        --agent general \
        --model "$MODEL"
}

# --- Main Loop ---
LAST_HASH=$(get_current_prd_hash)
i=1

while true; do
    # Watch Mode Check
    # Only restart if hash changed EXTERNALLY (since last accepted state)
    CURRENT_HASH=$(get_current_prd_hash)
    if [[ "$CURRENT_HASH" != "$LAST_HASH" ]]; then
        echo "👀 PRD Changed! Restarting loop..."
        echo "--- PRD CHANGED: RESTARTING LOOP ---" >> "$PROJECT_DIR/progress.log"
        LAST_HASH="$CURRENT_HASH"
        
        # Only reset counter if we are in infinite watch mode
        if [[ "$WATCH_MODE" == "true" ]]; then
            i=1
        fi
    fi

    if [[ "$WATCH_MODE" != "true" ]] && ((i > ITERATIONS)); then
        echo "⏸️  Max iterations reached."
        break
    fi

    PHASE=$(detect_phase)
    ITER_ID=$(printf "iter-%04d" $i)
    ITER_DIR=$(prepare_iteration_context "$ITER_ID")
    
    echo ""
    echo "🔁 Loop $i ($PHASE Phase)"
    echo "   Logs: $ITER_DIR"

    if [ "$PHASE" == "PLAN" ]; then
        MODELS=("${PLAN_MODELS[@]}")
    else
        MODELS=("${BUILD_MODELS[@]}")
    fi

    SUCCESS=false
    for MODEL in "${MODELS[@]}"; do
        echo "   Using: $MODEL"
        set +e
        run_agent "$MODEL" "$PHASE" "$ITER_DIR"
        EXIT_CODE=$?
        set -e
        
        RESPONSE=$(cat "$ITER_DIR/agent_response.txt")
        
        if [ $EXIT_CODE -eq 0 ] && [ -n "$RESPONSE" ]; then
            SUCCESS=true
            echo "$RESPONSE"
            
            if echo "$RESPONSE" | grep -q "<promise>COMPLETE</promise>"; then
                echo "✅ Agent signaled completion."
                if [[ "$WATCH_MODE" != "true" ]]; then 
                    exit 0
                else
                    echo "⏸️  Project Complete. Waiting for changes in prd.md..."
                    while [[ "$(get_current_prd_hash)" == "$LAST_HASH" ]]; do
                        sleep 2
                    done
                    echo "👀 Change detected! Resuming..."
                    i=1
                    continue
                fi
            fi
            break
        else
            echo "   ⚠️  Model $MODEL failed (Exit: $EXIT_CODE). Falling back..."
        fi
    done

    if [ "$SUCCESS" = false ]; then
        echo "❌ All models failed this iteration."
        sleep 2
    fi

    # Sync hash AFTER run so self-edits don't trigger restart
    LAST_HASH=$(get_current_prd_hash)

    ((i++))
    sleep 1
done
