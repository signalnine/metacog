#!/bin/bash
set -e
cd /home/gabe/metacog/experiments
export METACOG_EXP_BACKEND=codex
export METACOG_EXP_TIMEOUT=300
export METACOG_EXP_CODEX_REASONING=low

# Test 1: chorus in TEXT mode (amplifier test on broken recipe)
export METACOG_EXP_PROMPT_MODE=text-instructions
export METACOG_EXP_SAMPLES=2
TASKS="0 5 9"
for t in $TASKS; do
  echo "=== TEXT chorus task $t ==="
  .venv/bin/python runner.py --recipe chorus --task "$t" 2>&1 || echo "FAILED chorus-text $t"
done

# Test 2: envoy-extreme in TEXT mode (amplifier test on working recipe)
for t in $TASKS; do
  echo "=== TEXT envoy-extreme task $t ==="
  .venv/bin/python runner.py --recipe envoy-extreme --task "$t" 2>&1 || echo "FAILED envoy-extreme-text $t"
done

# Test 3: envoy-extreme TOOL-CALL at higher N (5 samples per task = 15 trials)
unset METACOG_EXP_PROMPT_MODE
export METACOG_EXP_SAMPLES=5
for t in $TASKS; do
  echo "=== HiN envoy-extreme task $t ==="
  .venv/bin/python runner.py --recipe envoy-extreme --task "$t" 2>&1 || echo "FAILED envoy-extreme-hiN $t"
done

echo "ROUND5_DONE"
