#!/bin/bash
set -e
cd /home/gabe/metacog/experiments
export METACOG_EXP_BACKEND=codex
export METACOG_EXP_SAMPLES=2
export METACOG_EXP_TIMEOUT=240
export METACOG_EXP_CODEX_REASONING=low
RECIPE=envoy-extreme-alt2
TASKS="0 5 9"
for t in $TASKS; do
  echo "=== $RECIPE task $t ==="
  .venv/bin/python runner.py --recipe "$RECIPE" --task "$t" 2>&1 || echo "FAILED $RECIPE task $t"
done
echo "ROUND3_DONE"
