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
import os
import subprocess
import sys
import tempfile
import time
from dataclasses import dataclass
from pathlib import Path
from typing import List

import yaml

import score

REPO_ROOT = Path(__file__).resolve().parents[1]
METACOG_BIN = REPO_ROOT / "metacog"
EXP_DIR = Path(__file__).resolve().parent
RECIPES_DIR = EXP_DIR / "recipes"
TASKS_FILE = EXP_DIR / "tasks.yaml"
RESULTS_FILE = EXP_DIR / "results.tsv"

GENERATOR_MODEL = os.environ.get("METACOG_EXP_GENERATOR", "claude-sonnet-4-6")
SAMPLES_PER_PAIR = int(os.environ.get("METACOG_EXP_SAMPLES", "3"))
GENERATOR_TIMEOUT = int(os.environ.get("METACOG_EXP_TIMEOUT", "300"))

RESULTS_HEADER = [
    "ts", "recipe", "task", "sample",
    "rarity", "coherence", "novelty",
    "n_entities", "answer_len",
    "answer_preview",
]


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
    with RESULTS_FILE.open("a", newline="") as f:
        w = csv.writer(f, delimiter="\t")
        w.writerow([row[k] for k in RESULTS_HEADER])


def baseline_for(task_id: str, results_path: Path) -> float | None:
    """Mean (rarity * coherence) of NULL recipe trials for this task."""
    if not results_path.exists():
        return None
    scores = []
    with results_path.open() as f:
        r = csv.DictReader(f, delimiter="\t")
        for row in r:
            if row["recipe"] == "null" and row["task"] == task_id:
                scores.append(float(row["rarity"]) * float(row["coherence"]))
    return (sum(scores) / len(scores)) if scores else None


def trial(recipe: Recipe, task: Task, sample: int) -> dict:
    home = tempfile.mkdtemp(prefix=f"metacog-exp-{recipe.name}-")
    prompt = build_prompt(recipe, task)
    answer = run_generator(prompt, home)
    rarity = score.score_rarity(answer)
    coherence = score.score_coherence(task.prompt, answer)
    base = baseline_for(task.id, RESULTS_FILE)
    novelty = rarity.score * coherence.score
    if base is not None:
        novelty = novelty - base
    return {
        "ts": int(time.time()),
        "recipe": recipe.name,
        "task": task.id,
        "sample": sample,
        "rarity": f"{rarity.score:.4f}",
        "coherence": f"{coherence.score:.4f}",
        "novelty": f"{novelty:+.4f}",
        "n_entities": len(rarity.entities),
        "answer_len": len(answer),
        "answer_preview": answer[:200].replace("\n", " ").replace("\t", " "),
    }


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
