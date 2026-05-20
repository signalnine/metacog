"""Sweep metacog recipes on DeepSeek with the standard scorer.

Mirror of qwen_sweep.py but targets DeepSeek API. Outputs deepseek_results.tsv
with the same schema so analyze.py can compute delta + emb_d on the results.

Usage:
    METACOG_EXP_DEEPSEEK_MODEL=deepseek-v4-pro \
      python deepseek_sweep.py --recipes none R24-biblical-chord-anchor R12-anchor-duo-occult \
                                --samples 2
"""
from __future__ import annotations

import argparse
import csv
import json
import os
import sys
import time
import urllib.request
import urllib.error
from pathlib import Path

import yaml
from dotenv import load_dotenv

load_dotenv(Path.home() / ".env")

EXP_DIR = Path(__file__).resolve().parent
RECIPES_DIR = EXP_DIR / "recipes"
TASKS_FILE = EXP_DIR / "tasks.yaml"

API_KEY = os.environ.get("DEEPSEEK_API_KEY")
URL = "https://api.deepseek.com/v1/chat/completions"
MODEL = os.environ.get("METACOG_EXP_DEEPSEEK_MODEL", "deepseek-v4-pro")

RESULTS_FILE = EXP_DIR / f"deepseek_{MODEL.replace('-', '_')}_results.tsv"
TRIALS_DIR = EXP_DIR / f"deepseek_{MODEL.replace('-', '_')}_trials"

import score


def render_call_as_text(call):
    cmd = call["cmd"]
    args = call.get("args", {})
    bullets = []
    for k, v in args.items():
        if isinstance(v, list):
            v = "; ".join(str(x) for x in v)
        bullets.append(f"  - {k}: {v}")
    body = "\n".join(bullets)
    return f"{cmd.upper()}:\n{body}"


def build_prompt(recipe_calls, task_prompt):
    if not recipe_calls:
        return (
            f"Answer the following task directly. Output only the answer; "
            f"no preamble, no meta-commentary about the question.\n\nTASK:\n{task_prompt}"
        )
    lines = [
        "Before answering the task below, take the following conditioning into account. "
        "Each block is a stance or operation you should internalize. Do not narrate that "
        "you read them; let them shape the answer.",
        "",
    ]
    for call in recipe_calls:
        lines.append(render_call_as_text(call))
        lines.append("")
    lines.extend([
        "Now answer the task below directly. Output only the answer. Do not explain that "
        "you read the conditioning. Do not summarize the conditioning. Do not preface the "
        "answer.",
        "",
        f"TASK:\n{task_prompt}",
    ])
    return "\n".join(lines)


def load_recipe(name):
    if name in ("none", "null", ""):
        return []
    path = Path(name) if Path(name).exists() else (RECIPES_DIR / f"{name}.yaml")
    if not path.exists():
        sys.exit(f"recipe not found: {name}")
    with open(path) as f:
        return yaml.safe_load(f).get("calls", [])


def query_deepseek(prompt, max_tokens=4000, temperature=0.7):
    body = json.dumps({
        "model": MODEL,
        "messages": [{"role": "user", "content": prompt}],
        "temperature": temperature,
        "max_tokens": max_tokens,
    }).encode()
    req = urllib.request.Request(
        URL, data=body,
        headers={"Content-Type": "application/json", "Authorization": f"Bearer {API_KEY}"}
    )
    with urllib.request.urlopen(req, timeout=300) as resp:
        return json.loads(resp.read().decode())


def load_tasks():
    with open(TASKS_FILE) as f:
        return yaml.safe_load(f)["tasks"]


def already_done():
    if not RESULTS_FILE.exists():
        return set()
    done = set()
    with open(RESULTS_FILE) as f:
        for r in csv.DictReader(f, delimiter="\t"):
            done.add((r["recipe"], r["task"], int(r["sample"])))
    return done


def init_results():
    if RESULTS_FILE.exists():
        return
    with open(RESULTS_FILE, "w") as f:
        w = csv.writer(f, delimiter="\t")
        w.writerow(["ts", "recipe", "control", "task", "sample", "rarity",
                    "coherence", "novelty", "n_entities", "entity_rarities",
                    "trial_path", "answer_len", "answer_preview"])


def append_row(row):
    with open(RESULTS_FILE, "a") as f:
        w = csv.writer(f, delimiter="\t")
        w.writerow([
            row["ts"], row["recipe"], row["control"], row["task"], row["sample"],
            row["rarity"], row["coherence"], row["novelty"], row["n_entities"],
            row["entity_rarities"], row["trial_path"],
            row["answer_len"], row["answer_preview"],
        ])


def trial(recipe_name, recipe_calls, task, sample, max_tokens=4000):
    prompt = build_prompt(recipe_calls, task["prompt"])
    result = query_deepseek(prompt, max_tokens)
    answer = result["choices"][0]["message"]["content"]

    rar = score.score_rarity(answer)
    coh = score.score_coherence(task["prompt"], answer)
    entities = [{"name": n, "rarity": r} for n, r in zip(rar.entities, rar.rarities)]
    control_flag = "1" if recipe_name in ("none", "null", "") else "0"

    TRIALS_DIR.mkdir(exist_ok=True)
    trial_path = TRIALS_DIR / f"{recipe_name}-{task['id']}-{sample}.json"
    with open(trial_path, "w") as f:
        json.dump({
            "recipe": recipe_name, "task": task["id"], "sample": sample,
            "prompt": prompt, "answer": answer,
            "rarity": rar.score, "coherence": coh.score, "entities": entities,
            "model": MODEL,
        }, f, indent=2)

    preview = answer.replace("\n", " ")[:200]
    return {
        "ts": int(time.time()), "recipe": recipe_name, "control": control_flag,
        "task": task["id"], "sample": sample,
        "rarity": f"{rar.score:.4f}", "coherence": f"{coh.score:.4f}",
        "novelty": "",  # filled in at analyze time
        "n_entities": len(entities),
        "entity_rarities": json.dumps([[e["name"], e["rarity"]] for e in entities]),
        "trial_path": str(trial_path.relative_to(EXP_DIR)),
        "answer_len": len(answer),
        "answer_preview": preview,
    }


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--recipes", nargs="+", required=True)
    ap.add_argument("--task", help="optional single task id")
    ap.add_argument("--samples", type=int, default=2)
    ap.add_argument("--max-tokens", type=int, default=4000)
    args = ap.parse_args()

    init_results()
    done = already_done()
    tasks = load_tasks()
    if args.task:
        tasks = [t for t in tasks if t["id"] == args.task]
    recipes = [(name, load_recipe(name)) for name in args.recipes]

    # Run NULL first so baseline embeddings/scores exist when others run
    recipes.sort(key=lambda r: (r[0] != "none", r[0]))

    total = len(recipes) * len(tasks) * args.samples
    n = 0
    for recipe_name, recipe_calls in recipes:
        for task in tasks:
            for sample in range(args.samples):
                n += 1
                key = (recipe_name, task["id"], sample)
                if key in done:
                    print(f"[{n}/{total}] skip {key}")
                    continue
                print(f"[{n}/{total}] {recipe_name} x {task['id']} #{sample}", flush=True)
                try:
                    row = trial(recipe_name, recipe_calls, task, sample, args.max_tokens)
                except Exception as e:
                    print(f"  FAILED: {e}", file=sys.stderr)
                    continue
                append_row(row)
                print(f"  rarity={row['rarity']} coh={row['coherence']} "
                      f"ents={row['n_entities']} len={row['answer_len']}")


if __name__ == "__main__":
    main()
