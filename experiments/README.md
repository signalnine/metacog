# Metacog conditioning experiments

Inspired by [karpathy/autoresearch](https://github.com/karpathy/autoresearch). Instead of mutating a `train.py` to lower `val_bpb`, we mutate **conditioning recipes** (sequences of metacog calls) and measure whether they push the generator into rarer latent space on a fixed task suite.

## Optimization target

"Rare/novel latent space" operationalized as a composite:

- **rarity** -- LLM-judged unusualness of named entities, methodologies, traditions, and specific terms-of-art the answer invokes (matches metacog's True Names doctrine)
- **coherence** -- separately judged "does this answer actually address the task?" Prevents reward-hacking novelty into word salad
- **delta from NULL** -- every (recipe, task) score is reported relative to a NULL (no-conditioning) baseline run on the same task

Both judgments use Haiku. Generator is Sonnet via `claude -p`. Cross-model judging cuts (some) same-model-as-generator bias.

## Architecture

```
experiments/
  runner.py         loops (recipe x task x sample), invokes claude -p, captures, scores
  score.py          rarity + coherence judges via Haiku
  recipes/*.yaml    one file per conditioning recipe; null.yaml is the control
  tasks.yaml        the task suite -- prompts where novelty has room to vary
  results.tsv       autoresearch-style log: one row per trial
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
1. Read `results.tsv`, look at top scores
2. Hypothesize a recipe variation (swap stance pool, prepend a primitive, wrap in a stratagem)
3. Add `recipes/<name>.yaml`
4. Re-run; the runner skips trials already in `results.tsv` (keyed by recipe+task+sample)

Variations worth seeding the agent with:
- Swap the `become` stance (your 64 stance pools are already a discrete search space)
- Prepend `feel` for felt-sense register, or `deconstruct` for structural register
- Wrap recipe inside a stratagem vs run as freestyle
- Vary parameter density: specific named methodologies (Bourdieu's habitus) vs generic descriptions (social conditioning)

## Caveats

- Novelty metrics are biased by the judge model. "Haiku-novel" is not "novel."
- Sample size matters. ~3 resamples per (recipe, task) is the floor for any signal.
- Cost: ~$0.05-0.20 per trial. 5 recipes x 5 tasks x 3 samples = ~$4-15.
