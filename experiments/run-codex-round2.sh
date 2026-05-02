#!/bin/bash
set -e
cd /home/gabe/metacog/experiments
export METACOG_EXP_BACKEND=codex
export METACOG_EXP_SAMPLES=2
export METACOG_EXP_TIMEOUT=240
export METACOG_EXP_CODEX_REASONING=low
RECIPES="envoy-extreme envoy-scientific freestyle-become"
TASKS="0 5 9"
for r in $RECIPES; do
  for t in $TASKS; do
    echo "=== $r task $t ==="
    .venv/bin/python runner.py --recipe "$r" --task "$t" 2>&1 || echo "FAILED $r task $t"
  done
done
echo "ROUND2_DONE"
