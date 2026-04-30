"""Aggregate results.tsv into stable per-recipe / per-task summaries.

Recomputes per-task baselines from the full current pool of control rows
(unlike runner.py's trial-time baseline, which uses whatever was present when
each trial was logged). Use this any time you want a coherent read of the
data; do not rely on the `novelty` column in results.tsv for non-control
rows since that was a snapshot at trial time.
"""

from __future__ import annotations

import argparse
import csv
import statistics
from collections import defaultdict
from pathlib import Path

from runner import baselines_from_rows, _row_is_control

RESULTS = Path(__file__).resolve().parent / "results.tsv"


def load_rows(path: Path) -> list[dict]:
    if not path.exists():
        return []
    with path.open() as f:
        return list(csv.DictReader(f, delimiter="\t"))


def per_recipe_per_task(rows: list[dict], baselines: dict[str, float]):
    """Return rows of (recipe, task, n, mean_rarity, mean_coh, mean_rar*coh, mean_delta, n_ent)."""
    grouped: dict[tuple[str, str], list[dict]] = defaultdict(list)
    for r in rows:
        grouped[(r["recipe"], r["task"])].append(r)
    out = []
    for (recipe, task), trials in sorted(grouped.items()):
        rarities = [float(t["rarity"]) for t in trials]
        cohs = [float(t["coherence"]) for t in trials]
        products = [r * c for r, c in zip(rarities, cohs)]
        ents = [int(t["n_entities"]) for t in trials]
        is_ctl = any(_row_is_control(t) for t in trials)
        if is_ctl:
            mean_delta = 0.0
        else:
            base = baselines.get(task)
            mean_delta = (sum(products) / len(products)) - base if base is not None else None
        out.append({
            "recipe": recipe,
            "task": task,
            "control": is_ctl,
            "n": len(trials),
            "rarity": statistics.mean(rarities),
            "coh": statistics.mean(cohs),
            "rar_coh": statistics.mean(products),
            "delta": mean_delta,
            "ents": statistics.mean(ents),
        })
    return out


def per_recipe(rows: list[dict], baselines: dict[str, float]):
    """Aggregate across tasks: per-recipe mean delta from baseline."""
    grouped: dict[str, list[dict]] = defaultdict(list)
    for r in rows:
        grouped[r["recipe"]].append(r)
    out = []
    for recipe, trials in sorted(grouped.items()):
        deltas = []
        for t in trials:
            score = float(t["rarity"]) * float(t["coherence"])
            base = baselines.get(t["task"])
            if base is None:
                continue
            if _row_is_control(t):
                deltas.append(0.0)
            else:
                deltas.append(score - base)
        is_ctl = any(_row_is_control(t) for t in trials)
        out.append({
            "recipe": recipe,
            "control": is_ctl,
            "n": len(trials),
            "rarity": statistics.mean(float(t["rarity"]) for t in trials),
            "rar_coh": statistics.mean(float(t["rarity"]) * float(t["coherence"]) for t in trials),
            "delta_mean": statistics.mean(deltas) if deltas else None,
            "delta_stdev": statistics.stdev(deltas) if len(deltas) > 1 else None,
        })
    return out


def fmt(x):
    if x is None:
        return "  --  "
    return f"{x:+.3f}"


def fmt_unsigned(x):
    if x is None:
        return "  --  "
    return f"{x:.3f}"


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--detail", action="store_true", help="show per-recipe x task breakdown")
    args = ap.parse_args()

    rows = load_rows(RESULTS)
    if not rows:
        print(f"no rows in {RESULTS}")
        return
    baselines = baselines_from_rows(rows)
    print(f"loaded {len(rows)} rows; per-task baselines (rarity*coh):")
    for task, b in sorted(baselines.items()):
        print(f"  {task:<30} {b:.3f}")
    print()

    print("=== per recipe (averaged across tasks) ===")
    print(f"{'recipe':<28} {'ctl':>3} {'n':>3} {'rarity':>7} {'rar*coh':>8} {'delta':>7} {'stdev':>7}")
    print("-" * 72)
    for r in per_recipe(rows, baselines):
        print(
            f"{r['recipe']:<28} {'y' if r['control'] else 'n':>3} {r['n']:>3} "
            f"{fmt_unsigned(r['rarity']):>7} {fmt_unsigned(r['rar_coh']):>8} "
            f"{fmt(r['delta_mean']):>7} {fmt_unsigned(r['delta_stdev']):>7}"
        )

    if args.detail:
        print()
        print("=== per recipe x task ===")
        print(f"{'recipe':<24} {'task':<26} {'n':>3} {'rar*coh':>8} {'delta':>7} {'ents':>5}")
        print("-" * 80)
        for r in per_recipe_per_task(rows, baselines):
            print(
                f"{r['recipe']:<24} {r['task']:<26} {r['n']:>3} "
                f"{fmt_unsigned(r['rar_coh']):>8} {fmt(r['delta']):>7} {r['ents']:>5.1f}"
            )


if __name__ == "__main__":
    main()
