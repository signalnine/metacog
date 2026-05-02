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

## Productionized (v6.2.0)

Two stratagems added in v6.2.0, derived from this experiment:

- **chorus** (3 becomes + fork + ritual): structural-axis champion.
  No synthesis -- the ablation that pushed emb_d from 0.180 to 0.235.
- **trinity** (3 becomes + fork + synthesis + ritual): balanced variant
  on the Pareto frontier. Keeps synthesis for delta lift.

See `cmd/metacog/stratagem.go` for the definitions and
`cmd/metacog/empirical_stratagems_test.go` for the validating tests.

## v6.3.0 follow-up: seven new primitives + Pareto-frontier extensions

After v6.2.0, the surface was reshaped: dropped `deconstruct`,
`measure`, `tether` and the 8 stratagems centered on them (audit,
autopsy, trilemma, survey, dive, banishing, drift, error) -- all sat at
emb_d ~0.10-0.13 in the all-stratagem sweep. Added 7 new primitives
(`register`, `chord`, `silence`, `excerpt`, `commitment`,
`disjunction`, `glossolalia`) chosen to fill specific gaps the
9-primitive surface didn't cover.

### Single-primitive screen (N=30 each)

| primitive | delta | emb_d | comment |
|---|---|---|---|
| register | -0.088 | **0.153** | highest emb_d among new singles; citation-stripping artifact suspected |
| excerpt | -0.018 | 0.142 | second-highest emb_d; bimodal (+0.76 and 0.0) |
| commitment | +0.000 | 0.128 | flat |
| disjunction | **+0.105** | 0.124 | best delta among new singles |
| glossolalia | -0.037 | 0.118 | weak |
| chord | -0.069 | 0.116 | weak |
| silence | -0.082 | 0.110 | weak |

Five of seven beat the dropped-stratagem cluster (0.10-0.13) on at
least one axis. Two (register, excerpt) reached above all original-16
stratagems on emb_d. Standalone delta winner is disjunction.

### Composition with chorus / trinity (N=70 each)

| recipe | delta | emb_d | structure |
|---|---|---|---|
| trinity-prepended-register | **+0.204** | **0.239** | register + 3 becomes + fork + ritual (Carson/Knuth/Weil) |
| chorus-plus-disjunction | **+0.347** | 0.162 | 3 becomes + fork + disjunction + ritual (Carson/Knuth/Weil) |
| chorus-plus-disjunction-alt | +0.233 | 0.152 | same structure, Merleau-Ponty/Randall/Williams |
| chorus-plus-register | +0.177 | 0.173 | register + 3 becomes + fork + ritual (additional path) |
| chorus-plus-excerpt | +0.120 | 0.159 | excerpt + 3 becomes + fork + ritual |

Two clean Pareto-frontier breakthroughs:

1. **trinity-prepended-register beats prior structural champion on
   BOTH axes simultaneously** (delta +0.204 vs +0.194; emb_d 0.239 vs
   0.226). The register-prepend is genuinely orthogonal to the
   trinity-no-synthesis structure -- not a citation-stripping artifact.
   The Victorian register imposes a non-default linguistic surface
   that the multi-voice base then operates within, producing both
   citation density (via the becomes) AND rare embedding-space contour
   (via the register).

2. **chorus-plus-disjunction is the new vocabulary-axis champion**
   (delta +0.347; was +0.231 with freestyle-become). Disjunction
   substitutes for synthesis in the chorus structure: where synthesis
   refuses resolution between 3 lenses with named blindspots,
   disjunction asserts a hard binary contradiction as the operand of
   reasoning. The contradiction lifts citation density dramatically
   because operating-inside-contradiction requires the answer to keep
   naming the specific propositions.

Author choice mattered for the original chorus-plus-disjunction
(+0.347) but the alt-author replication confirmed +0.233 -- still the
delta floor for this structure. Carson/Knuth/Weil amplifies, but the
+0.20+ delta lift is the structural floor.

### Failed compositions (negative results)

- chorus-plus-excerpt at +0.120 / 0.159: excerpt's standalone emb_d
  0.142 was largely a citation-stripping artifact. When chorus's
  becomes restore citations, excerpt's emb_d advantage dilutes.
- chord, silence, glossolalia, commitment as compositions: not
  attempted in depth after the screen results clustered them at
  emb_d ~0.11-0.13. The screen findings stand: not enough lift to
  justify depth runs.

### Productionized (v6.4.0)

Two new stratagems added in v6.4.0, derived from the chorus-plus-X
depth runs:

- **antinomy** (3 becomes + fork + disjunction + ritual): vocabulary-
  axis champion. Substitutes disjunction for synthesis in the chorus
  structure. At N=70: delta +0.347 (Carson/Knuth/Weil) /
  +0.233 (Merleau-Ponty/Randall/Williams). Author choice amplifies but
  the +0.20+ delta lift is the structural floor.
- **envoy** (register + 3 becomes + fork + ritual): both-axes
  champion. Prepends a register-shift to the chorus structure. At
  N=70: delta +0.204, emb_d 0.239. Beats the prior structural champion
  (trinity-no-synthesis-alt at +0.194 / 0.226) on both axes
  simultaneously.

See `cmd/metacog/stratagem.go` for the definitions and
`cmd/metacog/empirical_stratagems_test.go` for the validating tests.

The combined `register-chorus-disjunction` recipe (envoy + antinomy)
was tested but did NOT dominate either parent -- the register's
citation-stripping interacts with disjunction's citation-amplification
to produce a balanced but undominant result. The two findings are
more useful as separate stratagems than composed.

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
