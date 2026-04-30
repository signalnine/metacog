"""Scoring functions for metacog conditioning experiments.

Both LLM-judged scoring passes use Haiku 4.5 -- cheap, fast, and a different
model from the generator (Sonnet via claude -p) so the same-model-as-generator
bias is reduced. Embeddings use OpenAI text-embedding-3-small as a parallel
metric that captures conceptual reach without depending on Haiku's idea of
"specific named things."
"""

from __future__ import annotations

import json
import re
from dataclasses import dataclass
from typing import List

import anthropic
import openai

JUDGE_MODEL = "claude-haiku-4-5-20251001"
EMBED_MODEL = "text-embedding-3-small"

_anthropic_client: anthropic.Anthropic | None = None
_openai_client: openai.OpenAI | None = None


def _get_anthropic_client() -> anthropic.Anthropic:
    global _anthropic_client
    if _anthropic_client is None:
        _anthropic_client = anthropic.Anthropic()
    return _anthropic_client


def _get_openai_client() -> openai.OpenAI:
    global _openai_client
    if _openai_client is None:
        _openai_client = openai.OpenAI()
    return _openai_client


def embed(text: str) -> List[float]:
    """Return an embedding vector for `text` using OpenAI's text-embedding-3-small.

    Output is a list of 1536 floats. The trial harness caches these in the
    per-trial sidecar JSON so we only embed each answer once."""
    resp = _get_openai_client().embeddings.create(model=EMBED_MODEL, input=text)
    return list(resp.data[0].embedding)


@dataclass
class RarityScore:
    score: float            # 0.0-1.0; mean rarity over extracted entities
    entities: List[str]     # the entities the judge extracted
    rarities: List[float]   # per-entity rarity 0-1


@dataclass
class CoherenceScore:
    score: float            # 0.0-1.0; how well the answer addresses the task
    rationale: str          # one-sentence justification


def _haiku_json(prompt: str) -> dict:
    """Call Haiku, force JSON output, parse. Retries once on parse failure."""
    for attempt in range(2):
        msg = _get_anthropic_client().messages.create(
            model=JUDGE_MODEL,
            max_tokens=1024,
            messages=[{"role": "user", "content": prompt}],
        )
        text = msg.content[0].text.strip()
        # Strip ```json fences if present
        text = re.sub(r"^```(?:json)?\s*|\s*```$", "", text, flags=re.MULTILINE).strip()
        try:
            return json.loads(text)
        except json.JSONDecodeError:
            if attempt == 1:
                raise ValueError(f"Haiku returned non-JSON after retry: {text[:300]}")


def score_rarity(answer: str) -> RarityScore:
    """Extract named entities/methodologies/traditions from the answer
    and score each by how rare it is in mainstream discourse."""
    prompt = f"""Extract from the following text every named entity, specific methodology,
named tradition, named theoretical framework, or specific term-of-art the writer invokes.
Skip generic words ("AI", "users", "system", "design"). Keep only specific named things
(people, schools of thought, particular concepts with proper-noun status, named techniques).

For each, rate its rarity in mainstream English discourse from 0.0 to 1.0:
- 0.0 = household-name common ("Einstein", "machine learning", "the Renaissance")
- 0.5 = known to educated readers but not in mainstream conversation ("Kuhn's paradigm shifts", "the Frankfurt School")
- 1.0 = highly specialized; would require a domain expert to recognize ("Stafford Beer's viable system model", "Lynn Margulis's serial endosymbiosis", "Christopher Alexander's quality without a name")

Output ONLY JSON in this exact shape, with no prose before or after:
{{"entities": [{{"name": "...", "rarity": 0.0}}, ...]}}

If the text invokes nothing specific, output {{"entities": []}}.

TEXT:
{answer}
"""
    data = _haiku_json(prompt)
    items = data.get("entities", [])
    if not items:
        return RarityScore(score=0.0, entities=[], rarities=[])
    names = [str(e["name"]) for e in items]
    rarities = [float(e["rarity"]) for e in items]
    return RarityScore(score=sum(rarities) / len(rarities), entities=names, rarities=rarities)


def score_coherence(task: str, answer: str) -> CoherenceScore:
    """How well does the answer actually address the task?"""
    prompt = f"""Rate how well the ANSWER addresses the TASK on a scale 0.0 to 1.0:
- 0.0 = completely off-topic or incoherent
- 0.5 = partially addresses the task; significant evasion or generality
- 1.0 = directly engages every load-bearing part of the task with substance

Output ONLY JSON in this exact shape, with no prose before or after:
{{"score": 0.0, "rationale": "one sentence"}}

TASK:
{task}

ANSWER:
{answer}
"""
    data = _haiku_json(prompt)
    return CoherenceScore(
        score=float(data["score"]),
        rationale=str(data.get("rationale", "")),
    )
