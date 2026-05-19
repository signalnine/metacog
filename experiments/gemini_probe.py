"""Probe Gemini with metacog conditioning via the REST API."""
import argparse
import json
import os
import urllib.request
from pathlib import Path

import yaml
from dotenv import load_dotenv

load_dotenv(Path.home() / ".env")

EXP_DIR = Path(__file__).resolve().parent
RECIPES_DIR = EXP_DIR / "recipes"

API_KEY = os.environ.get("GEMINI_API_KEY")
MODEL = os.environ.get("GEMINI_MODEL", "gemini-2.5-pro")
URL = f"https://generativelanguage.googleapis.com/v1beta/models/{MODEL}:generateContent?key={API_KEY}"


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


def query_gemini(prompt):
    body = json.dumps({
        "contents": [{"parts": [{"text": prompt}]}],
        "generationConfig": {"temperature": 0.7, "maxOutputTokens": 4000},
    }).encode()
    req = urllib.request.Request(URL, data=body, headers={"Content-Type": "application/json"})
    with urllib.request.urlopen(req, timeout=180) as resp:
        data = json.loads(resp.read().decode())
    cands = data.get("candidates", [])
    if not cands:
        return f"[no candidates: {json.dumps(data)[:500]}]"
    parts = cands[0].get("content", {}).get("parts", [])
    if not parts:
        finish = cands[0].get("finishReason", "?")
        safety = cands[0].get("safetyRatings", [])
        return f"[no parts, finishReason={finish}, safety={safety}]"
    return "\n".join(p.get("text", "") for p in parts)


if __name__ == "__main__":
    ap = argparse.ArgumentParser()
    ap.add_argument("--recipe", required=True)
    ap.add_argument("--task", required=True)
    args = ap.parse_args()
    calls = load_recipe(args.recipe)
    prompt = build_prompt(calls, args.task)
    print(query_gemini(prompt))
