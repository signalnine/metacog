#!/bin/bash
set -e
cd /home/gabe/metacog/experiments
export METACOG_EXP_PROMPT_MODE=text-instructions
export METACOG_EXP_TIMEOUT=300

# Codex first (faster, fewer trials)
export METACOG_EXP_BACKEND=codex
export METACOG_EXP_CODEX_REASONING=low
export METACOG_EXP_SAMPLES=2
for t in 0 5 9; do
  echo "=== TEXT codex chorus-plus-disjunction task $t ==="
  .venv/bin/python runner.py --recipe chorus-plus-disjunction --task "$t" 2>&1 || echo "FAILED"
done

# Sonnet next (full sweep)
unset METACOG_EXP_BACKEND
unset METACOG_EXP_CODEX_REASONING
export METACOG_EXP_SAMPLES=2
echo "=== TEXT sonnet chorus-plus-disjunction ==="
.venv/bin/python runner.py --recipe chorus-plus-disjunction 2>&1 || echo "FAILED"
echo "AMP_CPD_DONE"
