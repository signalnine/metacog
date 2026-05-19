"""Compare metacog recipes on Qwen3.6 against baseline across the standard task set.

For each (recipe, task) pair: send to Qwen, parse out the final answer (strip <think>),
score with Haiku-equivalent scoring locally, write a single TSV row.

Usage:
    python qwen_sweep.py --recipes none R24-biblical-chord-anchor R12-anchor-duo-occult R12-sigil-name-commitment \
                         --samples 2 --max-tokens 6000
"""
from __future__ import annotations

import argparse
import csv
import json
import os
import re
import sys
import time
import urllib.request
from pathlib import Path

import yaml
from dotenv import load_dotenv

load_dotenv(Path.home() / ".env")

EXP_DIR = Path(__file__).resolve().parent
RECIPES_DIR = EXP_DIR / "recipes"
TASKS_FILE = EXP_DIR / "tasks.yaml"
QWEN_URL = "http://haight:8080/v1/chat/completions"
QWEN_MODEL = "RedHatAI/Qwen3.6-35B-A3B-NVFP4"

RESULTS_FILE = EXP_DIR / "qwen_results.tsv"
TRIALS_DIR = EXP_DIR / "qwen_trials"

import score  # the existing scorer; uses Haiku via Anthropic SDK


def render_call_as_text(call: dict) -> str:
    cmd = call["cmd"]
    args = call.get("args", {})
    bullets = []
    for k, v in args.items():
        if isinstance(v, list):
            v = "; ".join(str(x) for x in v)
        bullets.append(f"  - {k}: {v}")
    body = "\n".join(bullets)
    return f"{cmd.upper()}:\n{body}"


def build_prompt(recipe_calls: list[dict], task_prompt: str) -> str:
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


def load_recipe(name: str) -> list[dict]:
    if name in ("none", "null", ""):
        return []
    path = Path(name) if Path(name).exists() else (RECIPES_DIR / f"{name}.yaml")
    if not path.exists():
        sys.exit(f"recipe not found: {name}")
    with open(path) as f:
        data = yaml.safe_load(f)
    return data.get("calls", [])


def strip_thinking(text: str) -> str:
    idx = text.find("</think>")
    if idx >= 0:
        return text[idx + len("</think>"):].lstrip("\n ").strip()
    # No closing tag. Heuristic: if the first ~500 chars look like a thinking
    # block (numbered list, "Here's a thinking process", explicit "draft" labels),
    # cut at the first markdown heading or blockquote that looks like the answer
    # OR cut at the last "Final" marker
    for marker in ["\n\n**Final draft", "\nFinal Answer:", "\n## Final", "\n\n---\n\n"]:
        idx = text.rfind(marker)
        if idx > 0:
            return text[idx:].lstrip("\n#-* ").strip()
    return text  # give up; return full


def query_qwen(prompt: str, max_tokens: int, temperature: float) -> dict:
    body = json.dumps({
        "model": QWEN_MODEL,
        "messages": [{"role": "user", "content": prompt}],
        "max_tokens": max_tokens,
        "temperature": temperature,
    }).encode()
    req = urllib.request.Request(QWEN_URL, data=body, headers={"Content-Type": "application/json"})
    with urllib.request.urlopen(req, timeout=300) as resp:
        return json.loads(resp.read().decode())


def load_tasks() -> list[dict]:
    with open(TASKS_FILE) as f:
        return yaml.safe_load(f)["tasks"]


def already_done() -> set[tuple]:
    if not RESULTS_FILE.exists():
        return set()
    done = set()
    with open(RESULTS_FILE) as f:
        rdr = csv.DictReader(f, delimiter="\t")
        for r in rdr:
            done.add((r["recipe"], r["task"], int(r["sample"])))
    return done


def init_results():
    if RESULTS_FILE.exists():
        return
    with open(RESULTS_FILE, "w") as f:
        w = csv.writer(f, delimiter="\t")
        w.writerow(["ts", "recipe", "task", "sample", "rarity", "coherence", "n_entities",
                    "entity_rarities", "thinking_len", "answer_len", "answer_preview", "trial_path"])


def append_row(row: dict):
    with open(RESULTS_FILE, "a") as f:
        w = csv.writer(f, delimiter="\t")
        w.writerow([
            row["ts"], row["recipe"], row["task"], row["sample"],
            row["rarity"], row["coherence"], row["n_entities"],
            row["entity_rarities"], row["thinking_len"], row["answer_len"],
            row["answer_preview"], row["trial_path"],
        ])


def trial(recipe_name: str, recipe_calls: list[dict], task: dict, sample: int,
          max_tokens: int, temperature: float) -> dict:
    prompt = build_prompt(recipe_calls, task["prompt"])
    result = query_qwen(prompt, max_tokens, temperature)
    full = result["choices"][0]["message"]["content"]
    final = strip_thinking(full)
    thinking_len = len(full) - len(final)

    # Score with the existing rarity/coherence scorer (calls Haiku)
    rar = score.score_rarity(final)
    coh = score.score_coherence(task["prompt"], final)
    rarity = rar.score
    coherence = coh.score
    entities = [{"name": n, "rarity": r} for n, r in zip(rar.entities, rar.rarities)]

    TRIALS_DIR.mkdir(exist_ok=True)
    trial_path = TRIALS_DIR / f"{recipe_name}-{task['id']}-{sample}.json"
    with open(trial_path, "w") as f:
        json.dump({
            "recipe": recipe_name, "task": task["id"], "sample": sample,
            "prompt": prompt, "answer_full": full, "answer_final": final,
            "rarity": rarity, "coherence": coherence, "entities": entities,
            "usage": result.get("usage", {}),
        }, f, indent=2)

    preview = final.replace("\n", " ")[:200]
    return {
        "ts": int(time.time()), "recipe": recipe_name, "task": task["id"],
        "sample": sample, "rarity": f"{rarity:.4f}", "coherence": f"{coherence:.4f}",
        "n_entities": len(entities),
        "entity_rarities": json.dumps([[e["name"], e["rarity"]] for e in entities]),
        "thinking_len": thinking_len, "answer_len": len(final),
        "answer_preview": preview, "trial_path": str(trial_path.relative_to(EXP_DIR)),
    }


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--recipes", nargs="+", required=True)
    ap.add_argument("--task", help="optional single task id to test (default: all)")
    ap.add_argument("--samples", type=int, default=2)
    ap.add_argument("--max-tokens", type=int, default=6000)
    ap.add_argument("--temperature", type=float, default=0.7)
    args = ap.parse_args()

    init_results()
    done = already_done()
    tasks = load_tasks()
    if args.task:
        tasks = [t for t in tasks if t["id"] == args.task]
    recipes = [(name, load_recipe(name)) for name in args.recipes]

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
                    row = trial(recipe_name, recipe_calls, task, sample, args.max_tokens, args.temperature)
                except Exception as e:
                    print(f"  FAILED: {e}", file=sys.stderr)
                    continue
                append_row(row)
                print(f"  rarity={row['rarity']} coh={row['coherence']} "
                      f"ents={row['n_entities']} thinking={row['thinking_len']} ans={row['answer_len']}")


if __name__ == "__main__":
    main()
