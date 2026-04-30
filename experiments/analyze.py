"""Aggregate results.tsv into stable per-recipe / per-task summaries.

Recomputes per-task baselines from the full current pool of control rows
(unlike runner.py's trial-time baseline, which uses whatever was present when
each trial was logged). Use this any time you want a coherent read of the
data; do not rely on the `novelty` column in results.tsv for non-control
rows since that was a snapshot at trial time.

Reports the mean rarity metric (rarity*coherence delta vs baseline) plus
four alternative metrics that disagree with mean rarity in informative ways:

  max:  per-trial max entity rarity (rewards finding ANY one rare reference)
  sum:  per-trial sum of entity rarities (rewards quantity of specific things)
  hi:   per-trial count of entities at rarity >= 0.7 (density at the rare end)
  geo:  per-trial geometric mean (penalizes diluting rare with common)

Old rows from before the entity_rarities column existed show "--" on the
alternative metrics; rerun them to fill in.
"""

from __future__ import annotations

import argparse
import csv
import json
import statistics
from collections import defaultdict
from pathlib import Path
from typing import Dict, List, Optional

from dotenv import load_dotenv

# Embedding step needs OPENAI_API_KEY; load before score is imported.
load_dotenv(Path.home() / ".env")

from runner import (
    EXP_DIR,
    _row_is_control,
    baselines_from_rows,
    centroid,
    cosine_distance,
    metrics_from_rarities,
    parse_entity_rarities,
)
import score

RESULTS = Path(__file__).resolve().parent / "results.tsv"


def load_rows(path: Path) -> list[dict]:
    if not path.exists():
        return []
    with path.open() as f:
        return list(csv.DictReader(f, delimiter="\t"))


def load_sidecar(row: dict) -> Optional[dict]:
    """Load and return the per-trial sidecar JSON, or None if no path or
    file missing. Legacy rows have empty trial_path."""
    path_str = row.get("trial_path") or ""
    if not path_str:
        return None
    path = EXP_DIR / path_str
    if not path.exists():
        return None
    try:
        return json.loads(path.read_text())
    except (json.JSONDecodeError, OSError):
        return None


def ensure_embedding(row: dict, sidecar: dict, model: str = score.EMBED_MODEL) -> Optional[List[float]]:
    """Return the embedding for this trial, computing and caching it on the
    sidecar if not already present. Returns None if the answer is unavailable."""
    if "embeddings" not in sidecar:
        sidecar["embeddings"] = {}
    cached = sidecar["embeddings"].get(model)
    if cached is not None:
        return cached
    answer = sidecar.get("answer")
    if not answer:
        return None
    vec = score.embed(answer)
    sidecar["embeddings"][model] = vec
    # Persist back to disk so we only embed once.
    path = EXP_DIR / row["trial_path"]
    path.write_text(json.dumps(sidecar, ensure_ascii=False, indent=2))
    return vec


def embedding_centroids_per_task(rows: List[dict]) -> Dict[str, List[float]]:
    """Compute per-task centroids of NULL/control trial embeddings.

    Embeds answers as needed, caching to sidecars. Tasks with no control
    embeddings available get no entry.
    """
    by_task: Dict[str, List[List[float]]] = defaultdict(list)
    for row in rows:
        if not _row_is_control(row):
            continue
        sidecar = load_sidecar(row)
        if sidecar is None:
            continue
        vec = ensure_embedding(row, sidecar)
        if vec is None:
            continue
        by_task[row["task"]].append(vec)
    return {t: centroid(v) for t, v in by_task.items() if v}


def embedding_distance_for(row: dict, centroids: Dict[str, List[float]]) -> Optional[float]:
    """Cosine distance from this trial's embedding to its task's NULL centroid.

    Returns None if the trial lacks a sidecar/embedding, or its task has no
    centroid (no control embeddings yet)."""
    centroid_vec = centroids.get(row["task"])
    if centroid_vec is None:
        return None
    sidecar = load_sidecar(row)
    if sidecar is None:
        return None
    vec = ensure_embedding(row, sidecar)
    if vec is None:
        return None
    return cosine_distance(vec, centroid_vec)


def trial_metrics(row: dict) -> Optional[dict]:
    """Extract alternative metrics for one row. Returns None if per-entity
    data is missing or inconsistent (legacy row)."""
    parsed = parse_entity_rarities(row.get("entity_rarities"))
    try:
        n_expected = int(row.get("n_entities") or 0)
    except ValueError:
        n_expected = 0
    if not parsed and n_expected > 0:
        # Legacy row: had entities at trial time but didn't persist them.
        return None
    rarities = [r for _, r in parsed]
    return metrics_from_rarities(rarities)


def safe_mean(xs: List[float]) -> Optional[float]:
    return statistics.mean(xs) if xs else None


def per_recipe(rows: List[dict], baselines: dict[str, float],
               embed_centroids: Dict[str, List[float]]):
    grouped: dict[str, list[dict]] = defaultdict(list)
    for r in rows:
        grouped[r["recipe"]].append(r)

    out = []
    for recipe, trials in sorted(grouped.items()):
        rarities = [float(t["rarity"]) for t in trials]
        cohs = [float(t["coherence"]) for t in trials]
        products = [r * c for r, c in zip(rarities, cohs)]
        ents = [int(t["n_entities"]) for t in trials]

        deltas = []
        for t in trials:
            base = baselines.get(t["task"])
            if base is None:
                continue
            sc = float(t["rarity"]) * float(t["coherence"])
            deltas.append(0.0 if _row_is_control(t) else sc - base)

        # Alternative entity metrics: only over trials with per-entity data
        max_vals, sum_vals, hi_vals, geo_vals = [], [], [], []
        for t in trials:
            m = trial_metrics(t)
            if m is None:
                continue
            max_vals.append(m["max"])
            sum_vals.append(m["sum"])
            hi_vals.append(m["count_high"])
            if m["geo_mean"] is not None:
                geo_vals.append(m["geo_mean"])

        # Embedding distance: cosine distance from this task's NULL centroid.
        # Control trials by definition have ~0 expected mean (vs themselves),
        # but stdev within control is the natural noise floor.
        emb_dists = []
        for t in trials:
            d = embedding_distance_for(t, embed_centroids)
            if d is not None:
                emb_dists.append(d)

        is_ctl = any(_row_is_control(t) for t in trials)
        out.append({
            "recipe": recipe,
            "control": is_ctl,
            "n": len(trials),
            "n_metrics": len(max_vals),
            "n_emb": len(emb_dists),
            "rarity": statistics.mean(rarities),
            "rar_coh": statistics.mean(products),
            "delta": safe_mean(deltas),
            "delta_stdev": statistics.stdev(deltas) if len(deltas) > 1 else None,
            "max": safe_mean(max_vals),
            "sum": safe_mean(sum_vals),
            "hi": safe_mean(hi_vals),
            "geo": safe_mean(geo_vals),
            "emb_dist": safe_mean(emb_dists),
            "emb_stdev": statistics.stdev(emb_dists) if len(emb_dists) > 1 else None,
            "ents": statistics.mean(ents),
        })
    return out


def per_recipe_per_task(rows: List[dict], baselines: dict[str, float]):
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

        max_vals, sum_vals, hi_vals = [], [], []
        for t in trials:
            m = trial_metrics(t)
            if m is None:
                continue
            max_vals.append(m["max"])
            sum_vals.append(m["sum"])
            hi_vals.append(m["count_high"])

        out.append({
            "recipe": recipe,
            "task": task,
            "control": is_ctl,
            "n": len(trials),
            "rar_coh": statistics.mean(products),
            "delta": mean_delta,
            "max": safe_mean(max_vals),
            "sum": safe_mean(sum_vals),
            "hi": safe_mean(hi_vals),
            "ents": statistics.mean(ents),
        })
    return out


def fmt_delta(x):
    return "  --  " if x is None else f"{x:+.3f}"


def fmt_pos(x, width=6, prec=3):
    if x is None:
        return f"{'  --  ':>{width}}"
    return f"{x:>{width}.{prec}f}"


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

    print("computing embedding centroids from control trials (may embed if uncached)...")
    embed_centroids = embedding_centroids_per_task(rows)
    if embed_centroids:
        print(f"  centroids ready for {len(embed_centroids)} task(s): "
              f"{', '.join(sorted(embed_centroids))}")
    else:
        print("  no centroids available (no control sidecars yet -- rerun NULL baseline)")
    print()

    print("=== per recipe (averaged across tasks) ===")
    print(f"{'recipe':<26} {'ctl':>3} {'n':>3} {'rar*coh':>7} {'delta':>7}  "
          f"{'max':>5} {'sum':>5} {'hi':>4} {'geo':>5}  "
          f"{'emb_d':>6}  {'ents':>4}  {'(n)':>4} {'(e)':>4}")
    print("-" * 105)
    for r in per_recipe(rows, baselines, embed_centroids):
        print(
            f"{r['recipe']:<26} {'y' if r['control'] else 'n':>3} {r['n']:>3} "
            f"{fmt_pos(r['rar_coh'], 7)} {fmt_delta(r['delta']):>7}  "
            f"{fmt_pos(r['max'], 5)} {fmt_pos(r['sum'], 5, 2)} "
            f"{fmt_pos(r['hi'], 4, 1)} {fmt_pos(r['geo'], 5)}  "
            f"{fmt_pos(r['emb_dist'], 6)}  "
            f"{r['ents']:>4.1f}  {r['n_metrics']:>4} {r['n_emb']:>4}"
        )
    print()
    print("Column key:")
    print("  rar*coh = mean of rarity * coherence (the original metric)")
    print("  delta   = mean (rar*coh) - per-task baseline (control rows define baseline)")
    print("  max     = mean of per-trial max(entity_rarity); rewards finding ONE specialized name")
    print("  sum     = mean of per-trial sum(entity_rarity); rewards quantity of specifics")
    print("  hi      = mean of per-trial count(rarity >= 0.7); density at the rare end)")
    print("  geo     = mean of per-trial geometric mean of rarities; penalizes dilution")
    print("  emb_d   = mean cosine distance from this task's NULL embedding centroid")
    print("            (captures conceptual reach beyond proper-noun citations)")
    print("  ents    = mean entity count per trial")
    print("  (n)     = trials with per-entity data available (legacy rows excluded from max/sum/hi/geo)")
    print("  (e)     = trials with embeddings available (legacy rows excluded from emb_d)")

    if args.detail:
        print()
        print("=== per recipe x task ===")
        print(f"{'recipe':<24} {'task':<26} {'n':>3} {'rar*coh':>7} {'delta':>7}  "
              f"{'max':>5} {'sum':>5} {'hi':>4}  {'ents':>4}")
        print("-" * 92)
        for r in per_recipe_per_task(rows, baselines):
            print(
                f"{r['recipe']:<24} {r['task']:<26} {r['n']:>3} "
                f"{fmt_pos(r['rar_coh'], 7)} {fmt_delta(r['delta']):>7}  "
                f"{fmt_pos(r['max'], 5)} {fmt_pos(r['sum'], 5, 2)} "
                f"{fmt_pos(r['hi'], 4, 1)}  {r['ents']:>4.1f}"
            )


if __name__ == "__main__":
    main()
