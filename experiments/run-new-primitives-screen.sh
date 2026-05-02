#!/usr/bin/env bash
# Screen the 7 new-primitive recipes at N=30 (3 samples * 10 tasks).
# Runs 3 parallel runner.py invocations to amortize wall-clock time.
# Each invocation handles 2-3 recipes sequentially.
set -e

cd "$(dirname "$0")"

# Group 1: register, chord, silence
(
  for r in freestyle-register freestyle-chord freestyle-silence; do
    echo "[group1] starting $r at $(date)" >> "run-new-primitives-screen-1.log"
    python3 runner.py --recipe "$r" --samples 3 >> "run-new-primitives-screen-1.log" 2>&1
  done
  echo "[group1] done at $(date)" >> "run-new-primitives-screen-1.log"
) &
PID1=$!

# Group 2: excerpt, commitment
(
  for r in freestyle-excerpt freestyle-commitment; do
    echo "[group2] starting $r at $(date)" >> "run-new-primitives-screen-2.log"
    python3 runner.py --recipe "$r" --samples 3 >> "run-new-primitives-screen-2.log" 2>&1
  done
  echo "[group2] done at $(date)" >> "run-new-primitives-screen-2.log"
) &
PID2=$!

# Group 3: disjunction, glossolalia
(
  for r in freestyle-disjunction freestyle-glossolalia; do
    echo "[group3] starting $r at $(date)" >> "run-new-primitives-screen-3.log"
    python3 runner.py --recipe "$r" --samples 3 >> "run-new-primitives-screen-3.log" 2>&1
  done
  echo "[group3] done at $(date)" >> "run-new-primitives-screen-3.log"
) &
PID3=$!

wait $PID1 $PID2 $PID3
echo "all groups done at $(date)" >> "run-new-primitives-screen.log"
