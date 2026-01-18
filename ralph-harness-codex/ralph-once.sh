#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CODEX_BIN="${CODEX_BIN:-codex}"
CODEX_MODEL="${CODEX_MODEL:-}"

cd "$ROOT"

if [[ ! -f "prd.json" ]]; then
  echo "Missing prd.json in $ROOT" >&2
  exit 1
fi

touch progress.txt

PROMPT_BASE="$(cat "$ROOT/ralph-prompt.md")"
PRD_CONTENT="$(cat "$ROOT/prd.json")"
PROGRESS_CONTENT="$(cat "$ROOT/progress.txt")"

PROMPT="${PROMPT_BASE}

<PRD>
${PRD_CONTENT}
</PRD>

<PROGRESS>
${PROGRESS_CONTENT}
</PROGRESS>
"

if [[ -n "$CODEX_MODEL" ]]; then
  "$CODEX_BIN" exec --full-auto --skip-git-repo-check -m "$CODEX_MODEL" -C "$ROOT" "$PROMPT"
else
  "$CODEX_BIN" exec --full-auto --skip-git-repo-check -C "$ROOT" "$PROMPT"
fi
