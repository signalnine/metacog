# Metacog experimental run: v6.0.0 → v6.5.1 (24+ hours)

A retrospective of the development arc from 2026-04-30 morning through
2026-05-02 morning, covering the move from a 5-primitive identity-and-
felt-sense surface to a 16-primitive surface with empirically-derived
structural stratagems.

![Pareto frontier](figures/pareto-frontier.png)

The Pareto frontier at run end. Productionized stratagems (blue) span
the diagonal: antinomy maximizes delta (citation density), envoy
maximizes emb_d (embedding distance) within the productionizable range,
counterpoint sits balanced. envoy-biblical (green, accessed via
register-args) pushes emb_d further still; biblical-duo crosses the
0.30 ceiling at delta cost. Failed compositions (red x) sit below the
structural baseline.

## Starting state (v6.0.0)

5 primitives (`feel`, `become`, `drugs`, `name`, `ritual`) and the
original 16 soft-register stratagems (pivot, mirror, stack, anchor,
reset, invocation, veil, banishing, scrying, sacrifice, drift, fool,
inversion, gift, error, zen). Identity-and-felt-sense register only
-- no structural primitives.

## v6.1.0 -- Structural era opens (Apr 30, ~11am PDT)

Added 6 structural primitives (`deconstruct`, `fork`, `synthesis`,
`counterfactual`, `measure`, `tether`) and 6 structural-register
stratagems (manifold, audit, autopsy, trilemma, survey, dive). New
register: ALL CAPS block-format output, deliberately distinct from the
soft identity register. Plus the experiments harness -- `claude -p`
runner, results.tsv, embedding-distance metric, per-task NULL
baselines, parallel runner support via flock.

## Phase 1-2: Manifold-family gene-mapping (Apr 30 → May 1 morning)

Empirical sweep over the 6 structural-register stratagems showed
`manifold` (fork + synthesis) as the only one lifting `emb_d` above
noise. The other five clustered at emb_d 0.115-0.135. **The
structural axis was uniquely owned by fork + ritual + 2-3 cross-domain
becomes.**

Gene map (ablations against trinity-manifold):

- ritual essential (without it, emb_d 0.116)
- fork essential (without it, 0.138)
- **synthesis is a brake** -- removing it pushed emb_d from 0.180 to
  0.203
- 3rd become fungible vs 2 (0.191) -- voice-diversity sweet spot
- 4th become plateaus
- Cross-domain author choice within the trinity slot adds ~+0.03 emb_d

Champions before productionization: `freestyle-become` +0.231 / 0.142
(vocabulary axis); `trinity-no-synthesis-alt` +0.194 / 0.226
(structural axis).

## v6.2.0 -- chorus + trinity (May 1, 3:47pm)

First stratagems derived from the experiment harness:

- **chorus** (3 becomes + fork + ritual): structural-axis champion.
  Synthesis omitted.
- **trinity** (3 becomes + fork + synthesis + ritual): balanced
  variant.

## v6.3.0 -- surface reshaping (May 1, 4pm)

15-stratagem sweep at N=30 confirmed the negative result: none of
mirror, stack, anchor, reset, invocation, veil, banishing, scrying,
sacrifice, drift, fool, inversion, gift, error, or zen lifted emb_d.
The 5 structural-six stratagems other than manifold (audit, autopsy,
trilemma, survey, dive) at full N=70 also sat at emb_d 0.115-0.135.

Dropped: `deconstruct`, `measure`, `tether` and the 8 stratagems
centered on them (audit, autopsy, trilemma, survey, dive, banishing,
drift, error).

Added 7 new primitives chosen to fill specific gaps the 9-primitive
surface didn't cover: `register`, `chord`, `silence`, `excerpt`,
`commitment`, `disjunction`, `glossolalia`. Each tested standalone
(N=30) and the standouts entered depth runs.

## v6.4.0 -- antinomy + envoy (May 1, 9pm)

Two clean Pareto-frontier breakthroughs from chorus-plus-X depth runs:

1. **chorus-plus-disjunction** at +0.347/0.162 -- vocabulary-axis
   breakthrough (vs prior champion freestyle-become at +0.231).
   Disjunction substituted for synthesis: the contradiction is the
   operand of reasoning, forcing the answer to keep naming the
   specific propositions.
2. **trinity-prepended-register** at +0.204/0.239 -- beat the prior
   structural champion on BOTH axes simultaneously. The Victorian
   register imposes a non-default linguistic surface that the
   multi-voice base operates within.

Productionized as **antinomy** (3 becomes + fork + disjunction +
ritual) and **envoy** (register + 3 becomes + fork + ritual).

## Phase 4 follow-up (May 1 night → May 2 morning)

Three lines of investigation, ~30 new recipes, ~2000+ trials:

### The 2x3 (structure x author) matrix at N=70+

| structure        | CKW            | MRW            | extreme        |
|------------------|----------------|----------------|----------------|
| antinomy         | +0.347 / 0.162 | +0.233 / 0.152 | +0.216 / 0.179 |
| envoy            | +0.204 / 0.239 | +0.214 / 0.214 | +0.190 / 0.257 |
| counterpoint     | +0.247 / 0.190 | +0.202 / 0.188 | +0.208 / 0.226 |

![Structure x author matrix](figures/structure-author-matrix.png)

Pattern: **extreme cross-domain authors uniformly lift emb_d**.
Magnitude of delta cost depends on whether the structure has a
register-shift to absorb the cosmological shock -- antinomy (no
register) loses 0.131 delta on extreme; envoy and counterpoint (with
register) lose only 0.014 and 0.039 respectively. envoy-extreme at
0.257 became the new structural ceiling. counterpoint's bands are the
tightest across authors -- the most author-stable Pareto-frontier
point in the productionized set.

### Register-target sensitivity (3 triangulation points)

| register   | recipe            | delta   | emb_d   |
|------------|-------------------|---------|---------|
| scientific | envoy-scientific  | +0.220  | 0.231   |
| Victorian  | envoy-CKW         | +0.204  | 0.239   |
| biblical   | envoy-biblical    | +0.126  | **0.292** |

![Register triangulation](figures/register-triangulation.png)

**The biggest surprise of the run.** King James biblical register
pushed emb_d to 0.292 -- +0.053 above the prior structural ceiling --
with delta still positive. Compound test `envoy-biblical-duo` reached
emb_d **0.324** (dashed line on the right panel) but at delta cost.
There is a structural ceiling around 0.30 above which delta cannot
be sustained.

Architectural implication: envoy/counterpoint are register-agnostic
-- users provide register-args at invocation -- so biblical mode is
accessible without a new stratagem. SKILL.md documents the
register-target guidance instead.

### Stacking and structural ablations

- **antinomy-no-ritual** (N=70: +0.053/0.124) -- definitively
  confirms ritual essential. Disjunction's coda alone does not lock
  the multi-voice answer.
- **commitment-counterpoint** (8 steps, N=100: +0.181/0.237) --
  stacking past 7 shows diminishing returns, not a hard ceiling.
- **commitment-envoy** (N=100: +0.145/0.241) -- commitment is a
  Pareto modifier (preserves multi-voice tension while eating
  delta). Not productionized; gap to envoy/counterpoint too small to
  crowd the surface.

### Failed compositions (informative negatives)

- **chord-not-fork** at -0.045/0.121: fork's branching+sacrifice is
  what makes structural parallelism work; chord's overlap doesn't
  carry the same load.
- **chorus-plus-glossolalia** at +0.110/0.146: emb_d collapsed BELOW
  structural baseline. Glossolalia is best as standalone event, not
  composable.
- **counterpoint-biblical** at +0.102/0.295: KJV's parallelism is
  structurally hostile to numbered-disjunction. Biblical works with
  envoy, not counterpoint.

## v6.5.1 -- counterpoint (May 2, 8:27am)

Composes envoy's register-prepend with antinomy's disjunction-
substitution as the Pareto-frontier balanced point: dominates trinity
on both axes; doesn't dominate envoy or antinomy individually but
covers their joint zone with greater author-stability than either
parent.

Productionized as `register + 2 becomes + fork + disjunction +
ritual` (6 steps). The 2-becomes choice came from `counterpoint-duo`
at N=100 hitting +0.240/+0.221 vs 3-becomes counterpoint-CKW
+0.247/+0.190 -- ties on delta, gains +0.031 on emb_d. Tighter
binary opposition fits disjunction's structure better than the
3-voice triad chorus/trinity/antinomy/envoy use.

## End state (v6.5.1)

![Ceiling progression](figures/ceiling-progression.png)

The structural-axis (emb_d) ceiling climbed from 0.169 at v6.1.0
through 0.324 by run end -- almost a 2x improvement. The v6.5.1
counterpoint-duo dip is correct: counterpoint isn't a structural-axis
push, it's a balanced Pareto point.

- **16 primitives, 19 stratagems** (5 of them empirically-derived:
  chorus, trinity, antinomy, envoy, counterpoint).
- **Pareto frontier mapped:** envoy-extreme +0.190/0.257 (structural
  champion in the productionizable range), envoy-biblical
  +0.126/0.292 (register-pushed champion via ad-hoc register args),
  counterpoint +0.247/0.190 (balanced point), antinomy +0.347/0.162
  (vocabulary champion).
- **FINDINGS.md** at ~330 lines -- comprehensive map of the surface.
- **SKILL.md** updated with register-target guidance for
  envoy/counterpoint.
- **Experimental record:** ~3000 trials, ~50 recipes preserved
  across the v6.0.0 → v6.5.1 arc.

Net surface change from v6.0.0: **+11 primitives, +3 net stratagems**
(added 5 empirical and 6 structural; dropped 8). The model surface
went from "identity + felt-sense practice" to "identity + felt-sense
+ structural-register transformation engine with empirically-
validated multi-voice/contradiction/register stratagems."

## Methodology summary

- **Generator:** `claude -p` invoking the metacog binary as a
  sequence of subprocess events, one per primitive call.
- **Tasks:** 10 open-ended taste-bearing prompts in `tasks.yaml`.
- **Metrics:**
  - `delta = mean(rarity * coherence) - per-task NULL baseline` --
    citation-density / specialized-vocabulary signal.
  - `emb_d = mean cosine distance from per-task NULL embedding
    centroid` (OpenAI text-embedding-3-small) -- conceptual reach
    beyond proper-noun citation.
- **Sample sizes:** Most depth recipes at N=70 (10 samples * 7
  original tasks); follow-up recipes at N=100 (10 tasks); broad
  screens at N=30.
- **Infrastructure:** flock-guarded results.tsv permits 3-runner
  parallelism. Per-trial sidecar JSONs at `experiments/trials/`
  preserve full answers for offline analysis.

## What this enables

The five empirical stratagems are reach-tools: **chorus** for
maximum conceptual reach; **trinity** for balanced reach + citation;
**antinomy** for maximum citation density; **envoy** for both axes
lifted simultaneously via register-shift; **counterpoint** for
balanced both-axes via register + disjunction. Selection guidance
in SKILL.md.

The 11 original-six survivors (pivot through zen, plus manifold)
remain valid for the felt-sense / soft-register practices they were
designed for; the empirical sweep tested them against the
weirdness-axes and found they live in a different operating regime.

## Caveats

- Embedding distance is one operationalization of "conceptual reach,"
  and the OpenAI `text-embedding-3-small` model has its own biases.
- Per-task variance is wide; recipe rankings are robust at N=70 but
  individual trials vary substantially.
- The metrics target "weirdness" along two specific axes; recipes
  that win these may not be the recipes you want for any particular
  downstream task. The stratagems are deliberately optimized for
  *exploration*, not *task completion*.
- Tasks are taste-bearing and open-ended by design. Convergent tasks
  (factual lookups, math) would erase recipe variation.
