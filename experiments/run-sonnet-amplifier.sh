#!/bin/bash
set -e
cd /home/gabe/metacog/experiments
export METACOG_EXP_PROMPT_MODE=text-instructions
# Sonnet runs full task suite by default (10 tasks); 2 samples each = 20 trials
export METACOG_EXP_SAMPLES=2
# Working recipe in text mode -- predict: convergent with tool-call
echo "=== TEXT counterpoint-biblical-duo (sonnet, working) ==="
.venv/bin/python runner.py --recipe counterpoint-biblical-duo 2>&1 || echo "FAILED CBD"
# Broken recipe in text mode -- predict: less negative than tool-call
echo "=== TEXT chorus-with-chord-not-fork (sonnet, broken) ==="
.venv/bin/python runner.py --recipe chorus-with-chord-not-fork 2>&1 || echo "FAILED CCnF"
echo "AMP_DONE"
