#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [[ $# -lt 1 ]]; then
  echo "Usage: $0 <iterations>" >&2
  exit 1
fi

ITERATIONS="$1"

for ((i=1; i<=ITERATIONS; i++)); do
  echo "Ralph iteration $i/$ITERATIONS"
  RESULT="$("$ROOT/ralph-once.sh" 2>&1)"

  echo "$RESULT"

  if [[ "$RESULT" == *"<promise>COMPLETE</promise>"* ]]; then
    echo "PRD complete, exiting."
    exit 0
  fi
done

echo "Reached iteration cap."
