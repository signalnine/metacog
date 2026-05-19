"""Probe Sonnet 4.6 with metacog conditioning. Mirrors qwen_probe.py but uses
the local claude CLI (which is what the rest of the harness uses)."""
import argparse
import json
import os
import subprocess
import sys
from pathlib import Path

import yaml

EXP_DIR = Path(__file__).resolve().parent if "__file__" in dir() else Path("/home/gabe/metacog/experiments")
RECIPES_DIR = EXP_DIR / "recipes"


def render_call_as_text(call):
    cmd = call["cmd"]
    args = call.get("args", {})
    lines = []
    for k, v in args.items():
        if isinstance(v, list):
            v = "; ".join(str(x) for x in v)
        lines.append(f"  - {k}: {v}")
    return f"{cmd.upper()}:\n" + "\n".join(lines)


def build_prompt(calls, task):
    if not calls:
        return task
    parts = [
        "Before answering the task below, take the following conditioning into account. "
        "Each block is a stance or operation you should internalize. Do not narrate that "
        "you read them; let them shape the answer.",
        ""
    ]
    for c in calls:
        parts.append(render_call_as_text(c))
        parts.append("")
    parts.extend([
        "Now answer the task below directly. Output only the answer. Do not explain that "
        "you read the conditioning. Do not summarize the conditioning. Do not preface the answer.",
        "",
        f"TASK:\n{task}",
    ])
    return "\n".join(parts)


def load_recipe(name):
    if name in ("none", ""):
        return []
    p = RECIPES_DIR / f"{name}.yaml"
    with open(p) as f:
        return yaml.safe_load(f).get("calls", [])


def query_claude(prompt, model="claude-sonnet-4-6"):
    env = os.environ.copy()
    env.pop("ANTHROPIC_API_KEY", None)
    cmd = ["claude", "-p", prompt, "--model", model, "--permission-mode",
           "bypassPermissions", "--no-session-persistence"]
    result = subprocess.run(cmd, capture_output=True, text=True, timeout=300, env=env)
    return result.stdout


if __name__ == "__main__":
    ap = argparse.ArgumentParser()
    ap.add_argument("--recipe", required=True)
    ap.add_argument("--task", required=True)
    ap.add_argument("--model", default="claude-sonnet-4-6")
    args = ap.parse_args()
    calls = load_recipe(args.recipe)
    prompt = build_prompt(calls, args.task)
    out = query_claude(prompt, args.model)
    print(out)
