"""Probe Qwen3.6 at haight:8080 with metacog conditioning in text-instructions mode.

Usage:
    python qwen_probe.py --recipe psalter-chord --task "your question"
    python qwen_probe.py --recipe none --task "baseline question"
    python qwen_probe.py --recipe path/to/recipe.yaml --task "..." --json
"""
from __future__ import annotations

import argparse
import json
import re
import sys
import urllib.request
from pathlib import Path

import yaml

EXP_DIR = Path(__file__).resolve().parent
RECIPES_DIR = EXP_DIR / "recipes"
QWEN_URL = "http://haight:8080/v1/chat/completions"
QWEN_MODEL = "RedHatAI/Qwen3.6-35B-A3B-NVFP4"


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


def build_text_mode_prompt(recipe_calls: list[dict], task: str) -> str:
    if not recipe_calls:
        return task
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
        f"TASK:\n{task}",
    ])
    return "\n".join(lines)


def load_recipe(recipe_arg: str) -> list[dict]:
    if recipe_arg in ("none", "null", ""):
        return []
    # accept either a recipe name (in recipes/) or a path
    path = Path(recipe_arg)
    if not path.exists():
        path = RECIPES_DIR / f"{recipe_arg}.yaml"
    if not path.exists():
        sys.exit(f"recipe not found: {recipe_arg}")
    with open(path) as f:
        data = yaml.safe_load(f)
    return data.get("calls", [])


THINK_BLOCK_RE = re.compile(r"^(?:Here's a thinking process:|<think>).*?(?:</think>|(?=\n\nFinal answer:|\Z))",
                            re.DOTALL | re.IGNORECASE | re.MULTILINE)


def strip_thinking(text: str) -> str:
    # Qwen3 thinking section ends with </think>. Sometimes it opens with <think>,
    # sometimes with "Here's a thinking process:" plain prose. Cut everything up
    # to and including </think>.
    idx = text.find("</think>")
    if idx >= 0:
        return text[idx + len("</think>"):].lstrip("\n ").strip()
    return text


def query_qwen(prompt: str, max_tokens: int = 1500, temperature: float = 0.7) -> dict:
    body = json.dumps({
        "model": QWEN_MODEL,
        "messages": [{"role": "user", "content": prompt}],
        "max_tokens": max_tokens,
        "temperature": temperature,
    }).encode()
    req = urllib.request.Request(QWEN_URL, data=body, headers={"Content-Type": "application/json"})
    with urllib.request.urlopen(req, timeout=180) as resp:
        return json.loads(resp.read().decode())


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--recipe", required=True, help="recipe name (e.g. psalter-chord) or 'none' for baseline")
    ap.add_argument("--task", required=True, help="the task/question to ask")
    ap.add_argument("--max-tokens", type=int, default=1500)
    ap.add_argument("--temperature", type=float, default=0.7)
    ap.add_argument("--json", action="store_true", help="emit full JSON instead of just the answer")
    ap.add_argument("--keep-thinking", action="store_true", help="don't strip Qwen's thinking section")
    args = ap.parse_args()

    calls = load_recipe(args.recipe)
    prompt = build_text_mode_prompt(calls, args.task)
    result = query_qwen(prompt, args.max_tokens, args.temperature)

    answer = result["choices"][0]["message"]["content"]
    final = answer if args.keep_thinking else strip_thinking(answer)

    if args.json:
        print(json.dumps({
            "recipe": args.recipe,
            "task": args.task,
            "prompt_len": len(prompt),
            "answer_full": answer,
            "answer_final": final,
            "usage": result.get("usage", {}),
        }, indent=2))
    else:
        print(final)


if __name__ == "__main__":
    main()
