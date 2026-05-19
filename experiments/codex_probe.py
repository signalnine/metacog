"""Probe Codex (gpt-5.5) with metacog conditioning. Same shape as sonnet_probe."""
import argparse
import json
import os
import subprocess
import sys
from pathlib import Path

import yaml

EXP_DIR = Path(__file__).resolve().parent
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


def query_codex(prompt):
    env = os.environ.copy()
    env.pop("ANTHROPIC_API_KEY", None)
    cmd = [
        "codex", "exec",
        "--skip-git-repo-check",
        "--dangerously-bypass-approvals-and-sandbox",
        "-c", 'model_reasoning_effort="low"',
        "--json",
        prompt,
    ]
    proc = subprocess.run(cmd, env=env, capture_output=True, text=True, timeout=300)
    answer = ""
    for line in proc.stdout.splitlines():
        line = line.strip()
        if not line:
            continue
        try:
            evt = json.loads(line)
        except json.JSONDecodeError:
            continue
        if evt.get("type") == "item.completed":
            item = evt.get("item", {})
            if item.get("type") == "agent_message" and item.get("text"):
                answer = item["text"]
    return answer or f"[no agent_message; stderr head: {proc.stderr[:300]}]"


if __name__ == "__main__":
    ap = argparse.ArgumentParser()
    ap.add_argument("--recipe", required=True)
    ap.add_argument("--task", required=True)
    args = ap.parse_args()
    calls = load_recipe(args.recipe)
    prompt = build_prompt(calls, args.task)
    print(query_codex(prompt))
