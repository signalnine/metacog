#!/usr/bin/env bash
# Run the 3 chorus+new-primitive depth recipes at samples=7 (= N=70 across 10 tasks).
# These probe whether the strongest standalone-primitive signals compose
# with chorus or plateau at the manifold ceiling.
set -e

cd "$(dirname "$0")"

# Group A: chorus-plus-register (highest standalone emb_d)
(
  echo "[A] starting chorus-plus-register at $(date)" >> "run-depth-chorus-register.log"
  python3 runner.py --recipe chorus-plus-register --samples 7 >> "run-depth-chorus-register.log" 2>&1
  echo "[A] done at $(date)" >> "run-depth-chorus-register.log"
) &
PID_A=$!

# Group B: chorus-plus-excerpt (second-highest standalone emb_d, bimodal)
(
  echo "[B] starting chorus-plus-excerpt at $(date)" >> "run-depth-chorus-excerpt.log"
  python3 runner.py --recipe chorus-plus-excerpt --samples 7 >> "run-depth-chorus-excerpt.log" 2>&1
  echo "[B] done at $(date)" >> "run-depth-chorus-excerpt.log"
) &
PID_B=$!

# Group C: chorus-plus-disjunction (best standalone delta)
(
  echo "[C] starting chorus-plus-disjunction at $(date)" >> "run-depth-chorus-disjunction.log"
  python3 runner.py --recipe chorus-plus-disjunction --samples 7 >> "run-depth-chorus-disjunction.log" 2>&1
  echo "[C] done at $(date)" >> "run-depth-chorus-disjunction.log"
) &
PID_C=$!

wait $PID_A $PID_B $PID_C
echo "all depth groups done at $(date)" >> "run-depth-chorus-compositions.log"
