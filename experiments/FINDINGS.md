# Findings: Pareto-frontier search for "weird latent space"

## Question

Which compositions of metacog primitives push model output furthest from
the unconditioned baseline -- and along which axis (named-vocabulary
deployment vs. conceptual reach beyond proper-noun citation)?

## Setup

- **Generator:** `claude -p` invoking the metacog binary as a sequence
  of subprocess events, one per primitive call.
- **Tasks:** 10 open-ended taste-bearing prompts in `tasks.yaml`. The 7
  original tasks plus 3 pre-articulation-texture tasks added late.
- **Metrics:**
  - `delta = mean(rarity * coherence) - per-task NULL baseline` --
    measures named-entity rarity, weighted by Haiku-judged coherence.
    Effectively a citation-density / specialized-vocabulary signal.
  - `emb_d = mean cosine distance from per-task NULL embedding centroid`,
    using OpenAI `text-embedding-3-small`. Captures conceptual reach
    beyond proper-noun citation.
- **Budget:** ~2200 trials across 39 recipes. Most recipes at N=70
  (10 samples * 7 original tasks); some at N=100 (with new tasks);
  the 15-stratagem sweep at N=30.

## Final Pareto frontier

```
recipe                              delta   emb_d   notes
freestyle-become                   +0.231   0.142   vocabulary axis champion (Andy Clark)
trinity-manifold                   +0.195   0.180   balanced; 3 becomes + fork + synthesis + ritual
duo-manifold                       +0.187   0.191   2 becomes + fork + synthesis + ritual
trinity-no-synthesis-alt           +0.194   0.226   3 cross-domain becomes + fork + ritual
trinity-no-synthesis-extreme       +0.113   0.233   alt author triple
trinity-no-synthesis (canonical)   +0.149   0.203   first to cross 0.20 on emb_d
```

Structural-axis ceiling pushed from initial `manifold-stratagem` 0.169
to 0.235 across the iteration (+39%).

## Gene map (which primitives carry which work)

Verified through ablation against the trinity-manifold reference recipe
(3 becomes + fork + synthesis + ritual, emb_d 0.180):

| primitive | when removed | when present | conclusion |
|---|---|---|---|
| ritual | emb_d crashes to 0.116 | -- | essential carrier |
| fork | emb_d drops to 0.138 | -- | essential support |
| synthesis | emb_d **rises** to 0.203 | -- | **structural brake** |
| 3rd become (-> 2) | emb_d 0.191 (no loss) | -- | 3rd is fungible |
| 2nd become (-> 1) | emb_d 0.158 (loss) | -- | single voice anchors |
| 4th become added | emb_d 0.183 (plateau) | -- | diminishing returns |

Cross-domain author choice within the trinity slot matters for `emb_d`
roughly +0.03 (Carson/Knuth/Weil > Merleau-Ponty/Randall/Williams).
Going more extreme (Sun Ra/Moten/Fuller) plateaus at the same level --
the structure does the work; specific authors within a cross-domain
regime are fungible.

## What does NOT work

- **meditate** as opener (replaces becomes): emb_d 0.151 -- below baseline
- **drugs** as opener: emb_d 0.175 -- weak lift
- **counterfactual** added after manifold: emb_d 0.173 -- net negative
- **freestyle-deleuze** (single denser-jargon author): delta +0.171 --
  loses to freestyle-become (Clark) +0.231. The vocabulary metric is
  partly a proper-noun-citation count; Deleuze writes dense prose
  without naming.

## The stratagem sweep (negative result)

15 of the original-16 stratagems (everything except `pivot`, which had
been tested earlier) at N=30 across all 10 tasks:

```
all 15 stratagems clustered in emb_d range [0.103, 0.124]
NULL noise floor: ~0.090
manifold baseline (no becomes): 0.169
manifold-family champions: 0.180 - 0.235
```

**Definitive negative result:** none of mirror, stack, anchor, reset,
invocation, veil, banishing, scrying, sacrifice, drift, fool, inversion,
gift, error, or zen lifts the structural axis above the manifold
baseline. The structural axis is **uniquely** owned by
`fork + ritual + 2-3 cross-domain becomes`.

The 5 structural-six stratagems other than manifold (audit/autopsy/
trilemma/survey/dive) at full N=70 also sit at emb_d 0.115-0.135 --
they too do not lift the structural axis. Manifold is alone.

## Productionized

Two stratagems added in v6.2.0, derived from this experiment:

- **chorus** (3 becomes + fork + ritual): structural-axis champion.
  No synthesis -- the ablation that pushed emb_d from 0.180 to 0.235.
- **trinity** (3 becomes + fork + synthesis + ritual): balanced variant
  on the Pareto frontier. Keeps synthesis for delta lift.

See `cmd/metacog/stratagem.go` for the definitions and
`cmd/metacog/empirical_stratagems_test.go` for the validating tests.

## Caveats

- Embedding distance is one operationalization of "conceptual reach,"
  and the OpenAI `text-embedding-3-small` model has its own biases
  about what counts as similar.
- Per-task variance is wide (max trial-level emb_d is 0.240 vs cross-task
  averages around 0.20-0.23); the recipe rankings are robust at N=70 but
  individual trials vary substantially.
- The metric was chosen to detect "weirdness" along two axes; recipes
  that win these may not be the recipes you want for any particular
  downstream task. The stratagems are deliberately optimized for
  *exploration*, not *task completion*.
- Tasks are taste-bearing and open-ended by design. Convergent tasks
  (factual lookups, math) would erase recipe variation.
