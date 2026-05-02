#!/bin/bash
set -e
cd /home/gabe/metacog/experiments
export METACOG_EXP_BACKEND=codex
export METACOG_EXP_SAMPLES=2
export METACOG_EXP_TIMEOUT=300
export METACOG_EXP_CODEX_REASONING=low
TASKS="0 5 9"

# Test 1: chorus (3 cross-domain becomes + fork + ritual, no disjunction)
for t in $TASKS; do
  echo "=== chorus task $t ==="
  .venv/bin/python runner.py --recipe chorus --task "$t" 2>&1 || echo "FAILED chorus $t"
done

# Test 2: double-extreme (6 extreme becomes + fork + ritual)
for t in $TASKS; do
  echo "=== double-extreme task $t ==="
  .venv/bin/python runner.py --recipe double-extreme --task "$t" 2>&1 || echo "FAILED double-extreme $t"
done

# Test 3: counterpoint-biblical-duo in TEXT mode (rescue test)
export METACOG_EXP_PROMPT_MODE=text-instructions
for t in $TASKS; do
  echo "=== TEXT counterpoint-biblical-duo task $t ==="
  .venv/bin/python runner.py --recipe counterpoint-biblical-duo --task "$t" 2>&1 || echo "FAILED CBD-text $t"
done

echo "ROUND4_DONE"
