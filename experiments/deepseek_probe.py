"""Probe DeepSeek with metacog conditioning via their OpenAI-compatible API."""
import argparse
import json
import os
import urllib.request
import urllib.error
from pathlib import Path

import yaml
from dotenv import load_dotenv

load_dotenv(Path.home() / ".env")

EXP_DIR = Path(__file__).resolve().parent
RECIPES_DIR = EXP_DIR / "recipes"

API_KEY = os.environ.get("DEEPSEEK_API_KEY")
URL = "https://api.deepseek.com/v1/chat/completions"


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


def query_deepseek(prompt, model="deepseek-v4-pro"):
    body = json.dumps({
        "model": model,
        "messages": [{"role": "user", "content": prompt}],
        "temperature": 0.7,
        "max_tokens": 4000,
    }).encode()
    req = urllib.request.Request(
        URL, data=body,
        headers={"Content-Type": "application/json", "Authorization": f"Bearer {API_KEY}"}
    )
    try:
        with urllib.request.urlopen(req, timeout=180) as resp:
            data = json.loads(resp.read().decode())
        return data["choices"][0]["message"]["content"]
    except urllib.error.HTTPError as e:
        return f"[HTTP {e.code}: {e.read().decode()[:300]}]"


if __name__ == "__main__":
    ap = argparse.ArgumentParser()
    ap.add_argument("--recipe", required=True)
    ap.add_argument("--task", required=True)
    ap.add_argument("--model", default="deepseek-v4-pro")
    args = ap.parse_args()
    calls = load_recipe(args.recipe)
    prompt = build_prompt(calls, args.task)
    print(query_deepseek(prompt, args.model))
