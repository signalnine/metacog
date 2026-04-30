"""Run metacog conditioning recipes against a task suite.

For each (recipe, task, sample):
  1. Spin up a fresh METACOG_HOME so prior conditioning doesn't leak.
  2. Build a prompt that scripts the recipe's metacog calls and ends with the task.
  3. Invoke `claude -p` so the generator actually emits the metacog tool calls
     (the "tool calls as events" property of the practice is preserved).
  4. Score the answer for rarity (Haiku) and coherence (Haiku).
  5. Append a row to results.tsv.

Skips trials already in results.tsv (keyed by recipe+task+sample) so runs are resumable.
"""

from __future__ import annotations

import argparse
import csv
import fcntl
import json
import math
import os
import shutil
import subprocess
import sys
import tempfile
import time
from collections import defaultdict
from dataclasses import dataclass
from pathlib import Path
from typing import Any, Dict, List, Optional, Tuple

import yaml
from dotenv import load_dotenv

# Load ~/.env first so ANTHROPIC_API_KEY is available before score imports the SDK.
load_dotenv(Path.home() / ".env")

import score  # noqa: E402  # must come after load_dotenv

REPO_ROOT = Path(__file__).resolve().parents[1]
METACOG_BIN = REPO_ROOT / "metacog"
EXP_DIR = Path(__file__).resolve().parent
RECIPES_DIR = EXP_DIR / "recipes"
TASKS_FILE = EXP_DIR / "tasks.yaml"
RESULTS_FILE = EXP_DIR / "results.tsv"
TRIALS_DIR = EXP_DIR / "trials"

TRIAL_SCHEMA_VERSION = 1

GENERATOR_MODEL = os.environ.get("METACOG_EXP_GENERATOR", "claude-sonnet-4-6")
SAMPLES_PER_PAIR = int(os.environ.get("METACOG_EXP_SAMPLES", "3"))
GENERATOR_TIMEOUT = int(os.environ.get("METACOG_EXP_TIMEOUT", "300"))

RESULTS_HEADER = [
    "ts", "recipe", "control", "task", "sample",
    "rarity", "coherence", "novelty",
    "n_entities", "entity_rarities", "trial_path",
    "answer_len", "answer_preview",
]

RARITY_HIGH_THRESHOLD = 0.7


@dataclass
class Recipe:
    name: str
    description: str
    control: bool
    stratagem: str | None
    calls: List[dict]


@dataclass
class Task:
    id: str
    prompt: str


def load_recipes(only: str | None) -> List[Recipe]:
    out = []
    for path in sorted(RECIPES_DIR.glob("*.yaml")):
        data = yaml.safe_load(path.read_text())
        # YAML parses bare `null` as None; normalize the recipe name.
        name = data.get("name") or path.stem
        if only and name != only:
            continue
        out.append(Recipe(
            name=str(name),
            description=str(data.get("description", "")),
            control=bool(data.get("control", False)),
            stratagem=data.get("stratagem"),
            calls=list(data.get("calls", [])),
        ))
    return out


def load_tasks(only_idx: int | None) -> List[Task]:
    data = yaml.safe_load(TASKS_FILE.read_text())
    items = data.get("tasks", [])
    if only_idx is not None:
        items = [items[only_idx]]
    return [Task(id=str(t["id"]), prompt=str(t["prompt"]).strip()) for t in items]


def shell_quote(s: str) -> str:
    """Single-quote a string for safe inclusion in a bash command line."""
    return "'" + s.replace("'", "'\\''") + "'"


def render_metacog_command(call: dict) -> str:
    cmd = call["cmd"]
    args = call.get("args", {})
    parts = [str(METACOG_BIN), cmd]
    for key, value in args.items():
        flag = f"--{key}"
        if isinstance(value, list):
            for item in value:
                parts.extend([flag, shell_quote(str(item))])
        else:
            parts.extend([flag, shell_quote(str(value))])
    return " ".join(parts)


def build_prompt(recipe: Recipe, task: Task) -> str:
    """Construct the prompt fed to `claude -p`. The generator runs the
    metacog calls via Bash, then answers the task. The recipe is presented
    as concrete shell commands so the generator's own choice of when to
    invoke them stays transparent in its tool-call trace."""
    if recipe.control:
        return (
            f"Answer the following task directly. Output only the answer; "
            f"no preamble, no meta-commentary about the question.\n\n"
            f"TASK:\n{task.prompt}"
        )

    lines = [
        "You have access to a Bash tool. Before answering the task below, "
        "run these conditioning commands in order using Bash:",
        "",
    ]
    if recipe.stratagem:
        lines.append(f"  {METACOG_BIN} stratagem start {recipe.stratagem}")
    for call in recipe.calls:
        lines.append(f"  {render_metacog_command(call)}")
        if recipe.stratagem:
            lines.append(f"  {METACOG_BIN} stratagem next")
    lines.extend([
        "",
        "Treat each command as an event in your reasoning, not a procedure to "
        "narrate. After running them all, answer the task below directly. "
        "Output only the answer to the task. Do not explain that you ran the "
        "conditioning commands. Do not summarize the conditioning. Do not "
        "preface the answer.",
        "",
        f"TASK:\n{task.prompt}",
    ])
    return "\n".join(lines)


def run_generator(prompt: str, metacog_home: str) -> str:
    """Invoke `claude -p` and return stdout (the model's final answer)."""
    env = os.environ.copy()
    env["METACOG_HOME"] = metacog_home
    # SKILL.md checks this to suppress wait-for-human gates and offer
    # autonomous selection. Required for any non-interactive run.
    env["METACOG_HEADLESS"] = "1"

    cmd = [
        "claude",
        "-p", prompt,
        "--model", GENERATOR_MODEL,
        "--permission-mode", "bypassPermissions",
        "--no-session-persistence",
        "--tools", "Bash",
    ]
    proc = subprocess.run(
        cmd, env=env, capture_output=True, text=True, timeout=GENERATOR_TIMEOUT,
    )
    if proc.returncode != 0:
        raise RuntimeError(
            f"claude -p exited {proc.returncode}\nstderr: {proc.stderr[:1000]}"
        )
    return proc.stdout.strip()


def init_results():
    if RESULTS_FILE.exists():
        return
    with RESULTS_FILE.open("w", newline="") as f:
        w = csv.writer(f, delimiter="\t")
        w.writerow(RESULTS_HEADER)


def already_done() -> set:
    """Set of (recipe, task, sample) tuples already in results.tsv."""
    if not RESULTS_FILE.exists():
        return set()
    seen = set()
    with RESULTS_FILE.open() as f:
        r = csv.DictReader(f, delimiter="\t")
        for row in r:
            seen.add((row["recipe"], row["task"], int(row["sample"])))
    return seen


def append_result(row: dict):
    """Append one row to results.tsv under an exclusive flock so multiple
    runner processes can write concurrently (one per recipe) without
    interleaving partial rows."""
    with RESULTS_FILE.open("a", newline="") as f:
        fcntl.flock(f.fileno(), fcntl.LOCK_EX)
        try:
            w = csv.writer(f, delimiter="\t")
            w.writerow([row[k] for k in RESULTS_HEADER])
            f.flush()
        finally:
            fcntl.flock(f.fileno(), fcntl.LOCK_UN)


def _row_is_control(row: dict) -> bool:
    """A row is control if its `control` column is "1". Legacy rows without
    that column fall back to the old recipe == "null" heuristic so older TSVs
    keep working until they get migrated."""
    if "control" in row and row["control"] != "":
        return row["control"] == "1"
    return row.get("recipe") == "null"


def baselines_from_rows(rows: List[dict]) -> Dict[str, float]:
    """Per-task mean (rarity * coherence) over control rows only.

    Pure function for testability. Filters by the `control` column rather than
    by recipe name so a control recipe named something other than "null" still
    counts toward baselines.
    """
    by_task: Dict[str, list] = defaultdict(list)
    for row in rows:
        if not _row_is_control(row):
            continue
        try:
            score = float(row["rarity"]) * float(row["coherence"])
        except (KeyError, ValueError):
            continue
        by_task[row["task"]].append(score)
    return {task: sum(s) / len(s) for task, s in by_task.items() if s}


def cosine_distance(a: List[float], b: List[float]) -> float:
    """Cosine distance in [0, 2]. Returns 0 for zero vectors (defensive).

    cos_sim(a,b) = dot(a,b) / (|a| * |b|); cos_dist = 1 - cos_sim.
    Range: 0 (identical direction) to 2 (opposite direction).
    """
    if len(a) != len(b):
        raise ValueError(f"vector lengths differ: {len(a)} vs {len(b)}")
    dot = sum(x * y for x, y in zip(a, b))
    na = math.sqrt(sum(x * x for x in a))
    nb = math.sqrt(sum(y * y for y in b))
    if na == 0.0 or nb == 0.0:
        return 0.0
    return 1.0 - dot / (na * nb)


def centroid(vectors: List[List[float]]) -> Optional[List[float]]:
    """Componentwise mean of a list of vectors. None on empty input."""
    if not vectors:
        return None
    n = len(vectors)
    dim = len(vectors[0])
    out = [0.0] * dim
    for v in vectors:
        for i, x in enumerate(v):
            out[i] += x
    return [x / n for x in out]


def metrics_from_rarities(
    rarities: List[float], threshold: float = RARITY_HIGH_THRESHOLD
) -> Dict[str, Any]:
    """Compute alternative rarity metrics from a list of per-entity rarities.

    Returns a dict with:
      - max:        max rarity (rewards finding any one specialized reference)
      - sum:        sum of rarities (rewards quantity of specific things)
      - count_high: count of entities with rarity >= threshold
      - geo_mean:   geometric mean (penalizes diluting rare with common); None if empty

    Each metric reflects a different methodological choice about what "rare
    latent space" means. The mean rarity (in score.RarityScore.score) treats
    every entity equally; these alternatives don't.
    """
    if not rarities:
        return {"max": 0.0, "sum": 0.0, "count_high": 0, "geo_mean": None}
    if any(r == 0.0 for r in rarities):
        geo = 0.0
    else:
        geo = math.exp(sum(math.log(r) for r in rarities) / len(rarities))
    return {
        "max": max(rarities),
        "sum": sum(rarities),
        "count_high": sum(1 for r in rarities if r >= threshold),
        "geo_mean": geo,
    }


def parse_entity_rarities(value: Optional[str]) -> List[Tuple[str, float]]:
    """Parse the JSON-encoded entity_rarities cell into [(name, rarity), ...].

    Returns empty list for None, empty string, or malformed JSON. Old TSV rows
    that predate this column simply have no per-entity data; analyze.py shows
    "--" on metrics that require it for those rows.
    """
    if not value:
        return []
    try:
        data = json.loads(value)
    except (json.JSONDecodeError, TypeError):
        return []
    out = []
    for item in data:
        try:
            out.append((str(item[0]), float(item[1])))
        except (IndexError, TypeError, ValueError):
            continue
    return out


def compute_novelty(
    rarity: float, coherence: float, baseline: Optional[float], is_control: bool
) -> Optional[float]:
    """Novelty delta. Returns:
    - 0.0 for control trials (they ARE the baseline; delta against self is 0).
    - None for non-control trials with no baseline (undefined, do not infer 0).
    - rarity*coherence - baseline otherwise.
    """
    if is_control:
        return 0.0
    if baseline is None:
        return None
    return rarity * coherence - baseline


def baseline_for(task_id: str, results_path: Path) -> Optional[float]:
    """Wrapper that loads results.tsv and delegates to baselines_from_rows."""
    if not results_path.exists():
        return None
    with results_path.open() as f:
        rows = list(csv.DictReader(f, delimiter="\t"))
    return baselines_from_rows(rows).get(task_id)


def trial_sidecar_path(recipe_name: str, task_id: str, sample: int) -> Path:
    """Where the per-trial JSON lives. Stable filename so re-running the same
    (recipe, task, sample) overwrites the same sidecar."""
    return TRIALS_DIR / f"{recipe_name}-{task_id}-{sample}.json"


def write_sidecar(
    recipe: Recipe, task: Task, sample: int, ts: int,
    prompt: str, answer: str,
    rarity: "score.RarityScore", coherence: "score.CoherenceScore",
) -> Path:
    """Persist the full trial payload (prompt, answer, scoring detail) so
    later analysis can re-score, embed, or read the actual output without
    re-running the generator. Embeddings are added later by analyze.py."""
    TRIALS_DIR.mkdir(parents=True, exist_ok=True)
    path = trial_sidecar_path(recipe.name, task.id, sample)
    payload = {
        "schema": TRIAL_SCHEMA_VERSION,
        "ts": ts,
        "recipe": recipe.name,
        "control": recipe.control,
        "task": task.id,
        "sample": sample,
        "prompt": prompt,
        "answer": answer,
        "rarity": {
            "score": rarity.score,
            "entities": rarity.entities,
            "rarities": rarity.rarities,
        },
        "coherence": {
            "score": coherence.score,
            "rationale": coherence.rationale,
        },
        # embedding: list[float] keyed by model name, populated lazily by analyze.py
        "embeddings": {},
    }
    path.write_text(json.dumps(payload, ensure_ascii=False, indent=2))
    return path


def trial(recipe: Recipe, task: Task, sample: int) -> dict:
    home = tempfile.mkdtemp(prefix=f"metacog-exp-{recipe.name}-")
    ts = int(time.time())
    try:
        prompt = build_prompt(recipe, task)
        answer = run_generator(prompt, home)
        rarity = score.score_rarity(answer)
        coherence = score.score_coherence(task.prompt, answer)
        sidecar = write_sidecar(recipe, task, sample, ts, prompt, answer, rarity, coherence)
        base = baseline_for(task.id, RESULTS_FILE)
        novelty = compute_novelty(rarity.score, coherence.score, base, recipe.control)
        novelty_str = "" if novelty is None else f"{novelty:+.4f}"
        entity_rarities_json = json.dumps(
            list(zip(rarity.entities, rarity.rarities)),
            ensure_ascii=False,
        )
        return {
            "ts": ts,
            "recipe": recipe.name,
            "control": "1" if recipe.control else "0",
            "task": task.id,
            "sample": sample,
            "rarity": f"{rarity.score:.4f}",
            "coherence": f"{coherence.score:.4f}",
            "novelty": novelty_str,
            "n_entities": len(rarity.entities),
            "entity_rarities": entity_rarities_json,
            "trial_path": str(sidecar.relative_to(EXP_DIR)),
            "answer_len": len(answer),
            "answer_preview": answer[:200].replace("\n", " ").replace("\t", " "),
        }
    finally:
        shutil.rmtree(home, ignore_errors=True)


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--recipe", help="run only this recipe")
    ap.add_argument("--task", type=int, help="run only this task index (0-based)")
    ap.add_argument("--samples", type=int, default=SAMPLES_PER_PAIR)
    args = ap.parse_args()

    if not METACOG_BIN.exists():
        print(f"metacog binary not found at {METACOG_BIN}", file=sys.stderr)
        print("Build it first: go build -o metacog ./cmd/metacog/", file=sys.stderr)
        sys.exit(2)

    init_results()
    recipes = load_recipes(args.recipe)
    tasks = load_tasks(args.task)
    done = already_done()

    # Run NULL recipe first (across all tasks/samples) so per-task baselines exist
    # by the time non-control recipes compute their novelty delta.
    recipes.sort(key=lambda r: (not r.control, r.name))

    total = len(recipes) * len(tasks) * args.samples
    n = 0
    for recipe in recipes:
        for task in tasks:
            for sample in range(args.samples):
                n += 1
                key = (recipe.name, task.id, sample)
                if key in done:
                    print(f"[{n}/{total}] skip {key}")
                    continue
                print(f"[{n}/{total}] {recipe.name} x {task.id} #{sample}", flush=True)
                try:
                    row = trial(recipe, task, sample)
                except Exception as e:
                    print(f"  FAILED: {e}", file=sys.stderr)
                    continue
                append_result(row)
                print(
                    f"  rarity={row['rarity']} coherence={row['coherence']} "
                    f"novelty={row['novelty']} n_entities={row['n_entities']}"
                )


if __name__ == "__main__":
    main()
