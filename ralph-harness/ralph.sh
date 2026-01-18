#!/bin/bash
set -e

# Usage: ./ralph.sh [iterations]
ITERATIONS=${1:-5}
WORK_DIR="/home/shantanu/ralph-harness"

# Define models in order of preference
MODEL_PRIMARY="google/antigravity-claude-sonnet-4-5-thinking"
MODEL_FALLBACK_1="opencode/gpt-5.1-codex"
MODEL_FALLBACK_2="google/gemini-2.5-pro"

echo "🤖 Ralph is starting up... (Max iterations: $ITERATIONS)"
echo "   Primary: $MODEL_PRIMARY"
echo "   Fallback 1: $MODEL_FALLBACK_1"
echo "   Fallback 2: $MODEL_FALLBACK_2"
echo "   Fallback 3: Codex CLI (external)"

# Function to execute the agent run
run_agent() {
    local model_name=$1
    opencode run "Proceed with the next task in the autonomous loop. Read the attached files for instructions." \
        --file "$WORK_DIR/prompt.md" \
        --file "$WORK_DIR/prd.json" \
        --file "$WORK_DIR/progress.txt" \
        --agent general \
        --model "$model_name"
}

# Function to run codex cli
run_codex_cli() {
    echo "Running Codex CLI..."
    # Construct prompt with context
    {
        echo "Proceed with the next task in the autonomous loop."
        echo "Instructions from attached files follow:"
        echo "--- prompt.md ---"
        cat "$WORK_DIR/prompt.md"
        echo "--- prd.json ---"
        cat "$WORK_DIR/prd.json"
        echo "--- progress.txt ---"
        cat "$WORK_DIR/progress.txt"
    } | codex exec --dangerously-bypass-approvals-and-sandbox -
}

for ((i=1; i<=$ITERATIONS; i++)); do
  echo ""
  echo "🔁 Iteration $i / $ITERATIONS"
  echo "---------------------------------------------------"

  # Try Primary Model
  echo "Trying model: $MODEL_PRIMARY"
  # We use set +e temporarily to handle the error manually without exiting the script
  set +e
  RESPONSE=$(run_agent "$MODEL_PRIMARY")
  EXIT_CODE=$?
  set -e

  # Fallback 1: Codex (via opencode)
  if [ $EXIT_CODE -ne 0 ] || [ -z "$RESPONSE" ]; then
    echo "⚠️  Primary model failed or returned empty. Switching to fallback: Codex (via opencode)..."
    set +e
    RESPONSE=$(run_agent "$MODEL_FALLBACK_1")
    EXIT_CODE=$?
    set -e
  fi

  # Fallback 2: Gemini
  if [ $EXIT_CODE -ne 0 ] || [ -z "$RESPONSE" ]; then
    echo "⚠️  Codex (via opencode) failed. Switching to fallback: Gemini..."
    set +e
    RESPONSE=$(run_agent "$MODEL_FALLBACK_2")
    EXIT_CODE=$?
    set -e
  fi
  
  # Fallback 3: Codex CLI (external tool)
  if [ $EXIT_CODE -ne 0 ] || [ -z "$RESPONSE" ]; then
    echo "⚠️  Gemini failed. Switching to fallback: Codex CLI (external tool)..."
    set +e
    RESPONSE=$(run_codex_cli)
    EXIT_CODE=$?
    set -e
  fi

  # Final check if all failed
  if [ $EXIT_CODE -ne 0 ] || [ -z "$RESPONSE" ]; then
    echo "❌ All models and tools failed this iteration. Skipping to next..."
    continue
  fi

  # Output the agent's response
  echo "$RESPONSE"

  # Check for completion signal
  if echo "$RESPONSE" | grep -q "<promise>COMPLETE</promise>"; then
    echo ""
    echo "✅ Ralph has finished all tasks! Exiting loop."
    exit 0
  fi

  # Sleep slightly to ensure file system settles
  sleep 1
done

echo ""
echo "⚠️  Max iterations reached. Ralph is taking a break."
