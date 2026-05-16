# Metacog conditioning experiments

Inspired by [karpathy/autoresearch](https://github.com/karpathy/autoresearch). Instead of mutating a `train.py` to lower `val_bpb`, we mutate **conditioning recipes** (sequences of metacog calls) and measure whether they push the generator into rarer latent space on a fixed task suite.

## Optimization target

"Rare/novel latent space" operationalized as a composite:

- **rarity** -- LLM-judged unusualness of named entities, methodologies, traditions, and specific terms-of-art the answer invokes (matches metacog's True Names doctrine)
- **coherence** -- separately judged "does this answer actually address the task?" Prevents reward-hacking novelty into word salad
- **delta from NULL** -- every (recipe, task) score is reported relative to a NULL (no-conditioning) baseline run on the same task

Both judgments use Haiku. Generator is Sonnet via `claude -p` by default; the harness also supports Opus (`METACOG_EXP_BACKEND=opus METACOG_EXP_GENERATOR=claude-opus-4-7`) and Codex (`METACOG_EXP_BACKEND=codex`). Cross-model judging cuts (some) same-model-as-generator bias.

In addition to rarity and coherence, `analyze.py` also reports **emb_d** -- the mean cosine distance from each task's NULL embedding centroid (OpenAI `text-embedding-3-small`). This captures conceptual reach beyond proper-noun citations and is the second optimization axis alongside delta.

## Architecture

```
experiments/
  runner.py             loops (recipe x task x sample), invokes claude -p / codex, captures, scores
  score.py              rarity + coherence judges via Haiku
  analyze.py            post-hoc analysis with emb_d (recomputes deltas from full pool)
  recipes/*.yaml        one file per conditioning recipe; null.yaml is the control
  tasks.yaml            the task suite -- prompts where novelty has room to vary
  results.tsv           Sonnet/claude trials (one row per trial)
  opus_results.tsv      Opus trials (created when METACOG_EXP_BACKEND=opus)
  codex_results.tsv     codex trials (created when METACOG_EXP_BACKEND=codex)
  FINDINGS.md           comprehensive narrative of all rounds and architectural findings
```

Each trial spins a fresh `METACOG_HOME=$(mktemp -d)` so prior conditioning doesn't leak. `claude -p` actually invokes `metacog` via Bash, so the "tool calls as events" property of the practice is preserved -- the model genuinely emits the calls in its transcript.

## Running

```bash
cd experiments
uv venv && source .venv/bin/activate
uv pip install -r requirements.txt
# ANTHROPIC_API_KEY is loaded automatically from ~/.env via python-dotenv.
# If you keep it elsewhere, export it manually before running.
python runner.py                          # run full suite
python runner.py --recipe pivot           # one recipe, all tasks
python runner.py --recipe null --task 0   # one (recipe, task) pair

python analyze.py                         # post-hoc summary from full pool
python analyze.py --detail                # per-recipe x task breakdown

python -m unittest test_runner.py         # pure-function tests
```

## Reading results

`results.tsv` stores raw measurements (rarity, coherence, n_entities) plus a
trial-time `novelty` snapshot. **Use `analyze.py` for the authoritative
read** -- it recomputes per-task baselines from the full current pool of
control rows, so deltas are stable regardless of trial order. The `novelty`
column in the TSV is preserved for trace-keeping but should be treated as a
snapshot, not a metric to compare across recipes.

Control rows always have novelty `+0.0000` (the delta vs the baseline they
compose is definitionally zero). Non-control rows with no baseline available
at trial time have novelty empty -- run more control trials, then re-derive
deltas with `analyze.py`.

## Iterating

Manual loop (autoresearch-style):
1. Read `analyze.py` output, look at top scores on both axes (delta + emb_d)
2. Hypothesize a recipe variation (swap primitive, ablation, axis-compound, register substitution)
3. Add `recipes/<name>.yaml`
4. Re-run; the runner skips trials already in the results file (keyed by recipe+task+sample) so N can grow incrementally
5. Validate at N=20+ before productionizing as a stratagem (N=10 produces inflated estimates due to regression to mean)

Variations worth trying:
- Swap the `become` stance (78 stance pools are a discrete search space)
- Substitute one structural primitive for another (chord for fork, witness for synthesis, apophasis for silence)
- Ablate scaffolding (drop commitment, drop fork, drop ritual) to find the minimum viable mechanism
- Compose two recipes (anchor + becomes, name + anchors); check for axis interference
- Wrap recipe inside a stratagem vs run as freestyle

## Cross-model probing

Generator backends are switchable via `METACOG_EXP_BACKEND`:
- `claude` (default): Sonnet via `claude -p`
- `opus`: Opus 4.7 (set `METACOG_EXP_GENERATOR=claude-opus-4-7`); writes to `opus_results.tsv`
- `codex`: gpt-5.5 via `codex exec`; writes to `codex_results.tsv`

The `text-instructions` prompt mode (`METACOG_EXP_PROMPT_MODE=text-instructions`) delivers recipe content as plain text rather than tool-calls; useful for validating before committing to tool-call delivery on a new generator. See FINDINGS.md "Tool-call vs text: asymmetric amplifier" for the rule.

## Caveats

- Novelty metrics are biased by the judge model. "Haiku-novel" is not "novel."
- Sample size matters. N=20 is the floor for productionization decisions; N=10 produces inflated estimates.
- Cost: ~$0.05-0.20 per trial on Sonnet, ~$0.10-0.40 on Opus, ~$0.02-0.05 on codex.
- The metric was chosen to detect weirdness along two axes; recipes that win these may not be the recipes you want for any particular downstream task. The stratagems are deliberately optimized for *exploration*, not *task completion*.
