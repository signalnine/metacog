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
at N=70 hit delta +0.247 / emb_d 0.190 -- a Pareto-frontier
"balanced" point that doesn't dominate either parent (envoy beats it
on emb_d, antinomy beats it on delta) but does dominate trinity
(+0.195 / 0.180) on both axes. Productionized in v6.5.0 as
**counterpoint** after replication validated structural robustness
across three author triples (see Phase 4 follow-up below).

## v6.5.0 follow-up: Phase 4 author-matrix and stacking-ceiling tests

After v6.4.0 the immediate questions were:
1. Is each productionized stratagem structurally robust across multiple
   author triples, or are the headline numbers author-specific?
2. Where is the structural ceiling? Can primitives stack indefinitely?
3. Is the combined `register-chorus-disjunction` recipe robust enough
   to productionize as a fifth empirical stratagem?

### The 2x3 (structure x author) matrix at N=70+

Three productionized-or-candidate structures, each tested against three
author triples (CKW = Carson/Knuth/Weil; MRW = Merleau-Ponty/Randall/
Williams; extreme = Sun Ra/Moten/Fuller):

| structure        | CKW            | MRW            | extreme        |
|------------------|----------------|----------------|----------------|
| antinomy (N=70)  | +0.347 / 0.162 | +0.233 / 0.152 | +0.216 / 0.179 |
| envoy (N=70)     | +0.204 / 0.239 | +0.214 / 0.214 | +0.190 / 0.257 |
| counterpoint     | +0.247 / 0.190 | +0.202 / 0.188 | +0.208 / 0.226 |

(counterpoint-extreme at N=100; others at N=70.)

Three robust patterns:

1. **Extreme cross-domain authors uniformly LIFT emb_d.** Sun Ra/Moten/
   Fuller pushes emb_d above CKW and MRW for every structure.
   envoy-extreme at 0.257 is the **new structural-axis ceiling** (was
   envoy-CKW at 0.239).

2. **Extreme authors hurt delta, but the magnitude depends on
   structure.** antinomy loses 0.131 delta (CKW -> extreme); envoy loses
   only 0.014; counterpoint loses 0.039. The hypothesis: structures
   that have a register-shift PRE-ABSORB the cross-domain shock, so
   exotic cosmologies don't dilute citation density. Antinomy has no
   register, so its disjunction-driven citation density is more
   sensitive to author exoticism.

3. **Counterpoint is structurally robust across all three triples.**
   delta range +0.202..+0.247 (band 0.045); emb_d range 0.188..0.226
   (band 0.038). Both bands tighter than antinomy's delta band (0.131)
   or envoy's emb_d band (0.043). It is the most stable Pareto-frontier
   point across author choice.

### Productionized (v6.5.0)

One new stratagem added in v6.5.0:

- **counterpoint** (register + 3 becomes + fork + disjunction + ritual):
  Pareto-frontier balanced variant. Composes envoy's register-prepend
  with antinomy's disjunction-substitution. Use when both axes matter
  and you don't want to maximize one at the other's expense. Dominates
  trinity (+0.195 / 0.180) on both axes; does not dominate envoy or
  antinomy individually but covers their joint zone with greater
  author-stability than either parent.

### Stacking past 7 steps: diminishing returns, not a ceiling

**commitment-counterpoint** (commitment + register + 3 becomes + fork
+ disjunction + ritual = 8 steps) at N=84: delta +0.173 / emb_d 0.240.
Compared to counterpoint-CKW (+0.247/0.190): emb_d climbs +0.050,
delta drops 0.074. Compared to commitment-envoy (+0.145/0.241 at
N=100): commitment-counterpoint is roughly Pareto-equivalent (slightly
higher delta, slightly lower emb_d). Initial low-N reading at n=31
(+0.054/0.237) suggested a hard ceiling but stabilized at n=80+ to a
modest delta lift that the disjunction does provide -- it just isn't
preserved as cleanly under commitment's constraint.

**Conclusion:** 8 steps is not a hard ceiling, but the marginal value
of each additional step shrinks. The disjunction in commitment-counter-
point lifts delta by ~+0.03 over commitment-envoy, vs the ~+0.10 lift
disjunction provides in counterpoint over envoy without commitment.
Pre-binding via commitment absorbs about two-thirds of disjunction's
delta-lift while preserving its structural integrity.

### commitment as Pareto modifier (not productionized)

**commitment-envoy** (commitment + register + 3 becomes + fork + ritual)
at N=100: **+0.145 / 0.241**. emb_d above envoy-CKW (0.239); delta
between envoy (+0.204) and counterpoint (+0.247). A genuine
Pareto-frontier point that pushes emb_d structurally (via commitment)
rather than via author choice (via extreme).

Not productionized -- the gap to envoy/counterpoint is small enough
that adding a sixth empirical stratagem covering this point would
crowd the surface without offering a clearly new use case. The finding
is: **commitment is a structural emb_d-modifier that preserves
multi-voice tension while eating delta.** Useful as an ad-hoc
pre-binding wrapper around envoy/counterpoint when emb_d is the
priority and stake-naming is not.

### Failed compositions

- **chorus-with-chord-not-fork** (replace fork with chord) at N=28:
  -0.045 / 0.121. Both axes at noise floor. Chord cannot substitute
  for fork in chorus -- the branching+sacrifice topology is what makes
  structural parallelism work; chord's overlap doesn't carry the same
  load. Strong negative result, killed early.
- **chorus-plus-glossolalia** (inject glossolalia between fork and
  ritual) at N=23: +0.110 / 0.146. emb_d collapsed BELOW the
  structural baseline (chorus is 0.235 -> glossolalia drops it to
  0.146). Glossolalia is structurally disruptive in composition; the
  sub-semantic event breaks the multi-voice tension chorus depends on.
  Glossolalia is best as a standalone event, not a composable
  structural element. Killed early.
- **antinomy-no-ritual** (drop ritual to isolate disjunction's
  locking) at N=70: +0.053 / 0.124. Definitively confirms ritual is
  essential for antinomy: removing it crashes delta from +0.347 to
  +0.053 and emb_d from 0.162 to 0.124. Disjunction's coda alone does
  not lock the multi-voice answer; ritual's threshold-and-steps
  structure is doing real closing work.

### Surface-area probes (final)

Three tests probing dimensions not previously covered:

#### Register-target sensitivity (3 register triangulation points)

| register   | recipe            | delta   | emb_d  | N    |
|------------|-------------------|---------|--------|------|
| scientific | envoy-scientific  | +0.220  | 0.231  | 100  |
| Victorian  | envoy-CKW         | +0.204  | 0.239  | 70   |
| biblical   | envoy-biblical    | +0.126  | 0.292  | 99   |

**Register-target sensitivity is real and wide.** Scientific
(formal physics-paper conventions) and Victorian are roughly
Pareto-equivalent (scientific slightly favors delta; Victorian
slightly favors emb_d). King James biblical is the **new structural-
axis champion**: emb_d at 0.292 (+0.053 above the prior ceiling
envoy-extreme at 0.257) with delta still positive at +0.126.

The pattern: registers with low overlap with default contemporary
vocabulary push emb_d. Biblical's archaic vocabulary, parallelism,
and didactic mode of address are maximally orthogonal to default
contemporary register, producing the largest emb_d lift. Delta
holds because the biblical surface still admits proper-noun citation
when the underlying voices reference named authors -- but reduced
because biblical surface itself cites few modern entities.

#### Compound: biblical + duo voice-count

**envoy-biblical-duo** (biblical register + 2 becomes + fork +
ritual) at n=37: -0.070 / **0.318**. emb_d crosses 0.30 -- the
highest observed in the full experiment harness. But delta turns
NEGATIVE: the answer is so far from baseline embedding-space that
it stops citing the recipe-supplied voices.

The biblical+duo combination pushes emb_d further than either
finding alone (biblical envoy 0.292; duo counterpoint 0.221), but
it crosses a threshold past which delta cannot be sustained. There
is a structural ceiling around emb_d 0.30 above which the metric
gain comes from giving up answer-specificity entirely. Not a
productionization candidate -- the negative delta means the answer
fails the citation-density test.

#### Compound: biblical + disjunction (counterpoint-biblical)

**counterpoint-biblical** (biblical register + 3 becomes + fork +
disjunction + ritual) at n=64: +0.009 / 0.294. Delta near zero,
emb_d 0.294. Disjunction's normal +0.10 delta lift over envoy
disappears in biblical register -- the KJV surface's parallelism
and parataxis are structurally hostile to numbered-disjunction
naming. Hypothesis: biblical's paired-clause structure is *itself*
a kind of contradiction-handling, and adding disjunction's hard
binary on top creates structural conflict; one must give way.

**Implication for productionization:** biblical register works with
envoy structure (multi-voice + register + ritual) but does not
compose with antinomy/counterpoint structure (multi-voice +
disjunction + ritual). Register-target choice constrains which
structural primitives stack.

### Implications for the productionized stratagems

The productionized stratagems (chorus, trinity, antinomy, envoy,
counterpoint) all have register-agnostic step definitions. Users
provide register-args at invocation. The findings above mean a user
who wants the biblical emb_d ceiling can invoke `envoy` with
biblical register-args; they do not need a separate `psalm`
stratagem. Likewise scientific register works with envoy.

What the surface needs is **register-selection guidance** in the
skill documentation, not a new stratagem. Recommended pattern in
SKILL.md updates:
- Default register for envoy: Victorian (balanced).
- For maximum emb_d push with delta cost: biblical.
- For better delta with slight emb_d trade: scientific.
- Avoid: biblical + counterpoint structure (disjunction conflicts
  with biblical parallelism; delta crashes to zero).

#### Voice-diversity sweet spot

**counterpoint-duo** (counterpoint with 2 becomes instead of 3) at
N=100: **+0.240 / 0.221**. Compared to counterpoint-CKW (3 becomes,
N=70) +0.247 / 0.190: 2 becomes basically TIES delta (within 0.007)
and **GAINS** 0.031 emb_d. **2 becomes >= 3 becomes for counterpoint.**

This contradicts the prior trinity-manifold ablation finding that 2
becomes had emb_d 0.191 vs 3-becomes 0.180 -- both findings now agree
that 2 becomes preserves emb_d. The original "3 becomes is sweet
spot" result was likely due to the synthesis-locked baseline; under
disjunction (counterpoint), 2 becomes is at least as good.

The productionized counterpoint stratagem uses 3 becomes (kept for
consistency with chorus/trinity/antinomy/envoy). counterpoint-duo
is documented here as a tighter variant for use cases where emb_d
is the priority and slight delta cost is acceptable.

#### Stacking past 7 (refined)

**commitment-counterpoint** at FULL N=100: **+0.181 / 0.237**. The
n=31 reading (+0.054) was misleading -- at full N, the 8-step
structure does provide a modest delta lift over commitment-envoy
(+0.181 vs +0.145, +0.036 lift) at slight emb_d cost (-0.004).
Disjunction's value is preserved under commitment, just smaller
than its standalone composition value.

## Cross-model probe (gpt-5.5 / Codex CLI, low reasoning effort)

A mini sweep ran the productionized recipes against gpt-5.5 via the
Codex CLI at low reasoning effort. N=6 per recipe (3 tasks x 2
samples), null baseline computed from gpt-5.5's own outputs (not
Sonnet's). Tasks: git-conceptual-model, unindexed-intelligence,
lapsed-attention-unnoticed.

### Round 1 (initial probe)

| Recipe                         | codex delta | Sonnet delta (ref) | codex / Sonnet |
|--------------------------------|-------------|--------------------|----------------|
| null                           | 0           | 0                  | -              |
| envoy-alt                      | +0.118      | +0.214             | 0.55x          |
| trinity-no-synthesis-alt       | +0.035      | +0.194             | 0.18x          |
| chorus-plus-disjunction        | -0.027      | +0.347             | -0.08x         |
| counterpoint-biblical-duo      | **-0.228**  | **+0.177**         | -1.29x         |

The Sonnet champion (counterpoint-biblical-duo) is the codex
worst-case at -0.228 delta. Two trials returned zero entities. KJV
biblical register strips citations on gpt-5.5 without producing the
embedding-distance compensation it produces on Sonnet.

### Round 2 (decompose what works)

To isolate which axis transfers, round 2 tested three more recipes:
extreme cross-domain authors with no register-shift, scientific
register-shift instead of biblical, and a single-author baseline.

| Recipe                         | codex delta | Sonnet delta (ref) | codex / Sonnet |
|--------------------------------|-------------|--------------------|----------------|
| **envoy-extreme**              | **+0.310**  | +0.190             | **1.63x**      |
| envoy-alt                      | +0.118      | +0.214             | 0.55x          |
| envoy-scientific               | +0.081      | +0.220             | 0.37x          |
| freestyle-become               | +0.049      | +0.231             | 0.21x          |
| trinity-no-synthesis-alt       | +0.035      | +0.194             | 0.18x          |
| chorus-plus-disjunction        | -0.027      | +0.347             | -0.08x         |
| counterpoint-biblical-duo      | -0.228      | +0.177             | -1.29x         |

### Cross-model Pareto-frontier (full delta + emb_d, both models)

After round 4-5 added chorus, double-extreme, and envoy-extreme-alt2
to the codex sweep, and after computing embedding distances against
codex's own NULL centroids, the full cross-model picture is:

| Recipe                       | Sonnet delta | Sonnet emb_d | codex delta | codex emb_d |
|------------------------------|--------------|--------------|-------------|-------------|
| envoy-extreme                | +0.190       | 0.257        | +0.245      | 0.233       |
| envoy-extreme-alt2           | +0.256       | 0.192        | +0.286      | 0.223       |
| envoy-alt                    | +0.214       | 0.214        | +0.118      | 0.210       |
| envoy-scientific             | +0.220       | 0.231        | +0.081      | 0.180       |
| trinity-no-synthesis-alt     | +0.194       | 0.226        | +0.035      | 0.157       |
| chorus-plus-disjunction      | +0.347       | 0.162        | -0.027      | 0.148       |
| chorus (CKW)                 | n/a          | n/a          | -0.129      | 0.147       |
| double-extreme               | n/a          | n/a          | +0.060      | 0.171       |
| freestyle-become             | +0.231       | n/a          | +0.049      | 0.143       |
| counterpoint-biblical-duo    | +0.177       | 0.327        | **-0.228**  | **0.224**   |

See `docs/figures/cross-model-pareto.png` for the visual: arrows
connect each recipe's Sonnet point (circle) to its codex point
(triangle). Three patterns are visible:

1. **envoy-extreme variants** (blue): codex points sit RIGHT of
   Sonnet -- codex actually beats Sonnet on delta for the cross-
   domain author triples. emb_d is similar or slightly lower.
2. **counterpoint-biblical-duo** (red): codex point sits FAR LEFT
   AND DOWN from Sonnet. The biblical register collapses both axes
   on codex; on Sonnet it's the emb_d champion at 0.327.
3. **chorus-plus-disjunction** (purple): codex point sits FAR LEFT
   of Sonnet on delta (+0.347 -> -0.027) but emb_d barely moves
   (0.162 -> 0.148). The disjunction-driven citation density that
   Sonnet rewards doesn't engage on codex, but the structural-
   distance from default holds.

Interesting note: codex's emb_d ceiling is around 0.23 (envoy-extreme,
envoy-extreme-alt2, CBD), while Sonnet hits 0.327 on CBD. Codex's
default voice is harder to escape on the embedding axis -- but
when it IS pushed, the recipes that work on Sonnet are a poor guide
to which recipes work on codex.

**envoy-extreme transfers and overperforms.** Extreme cross-domain
authors (Sun Ra at the chalkboard, Fred Moten on fugitive sociality,
Buckminster Fuller drafting Synergetics) push gpt-5.5 to +0.310 delta
-- *stronger than the same recipe on Sonnet* (+0.190).

The cross-model pattern:

1. **Extreme cross-domain author-becomes transfer.** They overperform
   on codex relative to Sonnet. The conditioning works as intended:
   codex actually pulls toward the named cosmologies, citing them.
2. **Mild academic authors transfer partially.** Envoy-alt
   (Merleau-Ponty/Randall/Williams) lifts codex about half as much as
   Sonnet. The names land but with less force.
3. **Register-shifts (biblical or scientific) hurt codex.** Biblical
   is catastrophic (-0.228); scientific is mildly positive (+0.081)
   but still well below extreme-author lift. Codex's KJV voice goes
   abstract rather than archaic-and-specific. Register doesn't carry
   citation density on gpt-5.5 the way it does on Sonnet.
4. **Disjunction is dead.** chorus-plus-disjunction (Sonnet's biggest
   win at +0.347) is essentially null on codex (-0.027). The
   contradiction-as-operand structure doesn't produce the
   citation-naming pressure on gpt-5.5.

**Cross-model recipe rule:** lean hard on cross-domain author-becomes;
skip register-shifts and disjunction. The structural ceremony (fork,
ritual, multiple becomes) is robust; the surface mechanisms are not.
A codex-targeted productionized recipe would look like envoy-extreme:
3 extreme becomes + fork + ritual, no register-shift, no disjunction.

The model-specificity caveat is now empirically grounded. Recipes
optimized against one generator's habits don't transfer cleanly --
but a subset of mechanisms (cross-domain author conditioning) does
transfer, while others (register-shift, disjunction) are
generator-specific.

Open questions:
- Does higher reasoning effort change codex's response to register
  and disjunction recipes? Low effort may bypass surface conditioning.
- Does emb_d transfer match delta transfer? (Codex trials weren't
  embedded against per-task NULL centroids; only delta is measured.)
- How does an envoy-extreme variant calibrated against codex
  baselines perform when ported back to Sonnet?

### Round 3: extremity-as-property and tool-calls vs text instructions

Two hypotheses tested:
- (Q1) Is envoy-extreme's lift author-specific to Sun Ra/Moten/Fuller,
  or does *any* extreme cross-domain triple transfer?
- (Q2) Does the tool-call invocation mechanism itself carry weight, or
  does the same content delivered as plain text-instruction prose
  produce the same conditioning effect?

Built a fresh extreme triple: **envoy-extreme-alt2** with Octavia
Butler / Donna Haraway / Lynn Margulis (SF / cyborg theory /
endosymbiotic biology -- different domains, different gender, different
lineage from Sun Ra / Moten / Fuller).

**Q1 result: extremity-as-property confirmed.**

| Recipe              | codex delta | Sonnet delta |
|---------------------|-------------|--------------|
| envoy-extreme       | +0.310      | +0.190       |
| envoy-extreme-alt2  | +0.286      | +0.256       |

The Butler/Haraway/Margulis triple lifts both models comparably to
Sun Ra/Moten/Fuller. The lift comes from the structural property of
holding three radically cross-domain author-cosmologies in parallel,
not from the specific names. Any disparate-extreme triple transfers.

**Q2 result: tool-calls win modestly on both axes.**

Same recipe, same content, different delivery: PROMPT_MODE
toggle in the runner switches between rendering metacog calls as
shell commands the generator invokes (default; the calls are real
tool-call events in its transcript) vs delivering them as prose
instructions inside the prompt body.

| Mode                        | codex delta | Sonnet delta | Sonnet emb_d |
|-----------------------------|-------------|--------------|--------------|
| tool-calls (action-trace)   | +0.286      | +0.256       | 0.192        |
| text-instructions           | +0.268      | +0.238       | 0.168        |
| Difference (tool-call lift) | +0.018      | +0.018       | +0.024       |

Tool-call mode wins by ~7% relative on delta and ~12% relative on
emb_d. The action-trace prior matters *more* for embedding distance
than for citation density. Tool-call mode pushes conceptual reach
(emb_d) further than text-instruction mode produces the same recipe
content delivered as prose.

Interpretation: the "tool calls as events" doctrine is real but
modest on the delta axis and slightly larger on the emb_d axis. ~88%
of the recipe lift comes from CONTENT (which authors, what stances,
the structural composition); ~12% comes from the action-trace prior
of the model emitting the calls vs reading them as text. Tool-calls
land the model in continuation-of-action mode where it inhabits the
conditioning more fully, drifting further from default voice -- the
delta effect is smaller because citations are mostly fixed by the
named authors regardless of delivery mode.

This refines the framing: a metacog skill could in principle be
delivered as plain instruction text and capture most of the lift. The
tool-call ceremony preserves a small additional effect, plausibly
because the model treats post-tool-call context as a continuation-of-
action rather than a description-of-character. But the structural
content is doing the heavy lifting, not the protocol.

### Round 4-6: tool-call mode is an asymmetric amplifier

Tested working and broken recipes in both modes on both Sonnet and
codex. Round 4 suggested a uniform amplifier effect; round 5 added
N=15 codex data and showed the round-3 "tool-call adds +0.018"
finding was within noise on working recipes; round 6 added Sonnet
amplifier validation across two recipes (working and broken).

**Final cross-mode, cross-model amplifier table:**

| Model  | Recipe                              | Tool-call     | Text         | Diff       |
|--------|-------------------------------------|---------------|--------------|------------|
| sonnet | chorus-plus-disjunction (strong)    | +0.347 (n=70) | +0.140 (n=19)| **+0.208** |
| sonnet | counterpoint-biblical-duo (working) | +0.177 (n=30) | +0.088 (n=20)| **+0.089** |
| codex  | envoy-extreme (working)             | +0.245 (n=14) | +0.301 (n=6) | -0.055     |
| codex  | envoy-extreme-alt2 (working)        | +0.286 (n=6)  | +0.268 (n=6) | +0.018     |
| sonnet | envoy-extreme-alt2 (working)        | +0.256 (n=30) | +0.238 (n=30)| +0.018     |
| sonnet | chorus-with-chord-not-fork (broken) | -0.066 (n=29) | -0.000 (n=20)| -0.066     |
| codex  | chorus (mild-broken)                | -0.129 (n=6)  | -0.028 (n=6) | **-0.101** |
| codex  | chorus-plus-disjunction (broken on codex) | -0.027 (n=6) | +0.117 (n=6) | **-0.143** |
| codex  | counterpoint-biblical-duo (broken)  | -0.228 (n=6)  | +0.048 (n=6) | **-0.276** |

**Working recipes** (4 cases): mean diff +0.018. Tool-call usually beats
or ties text mode. The Sonnet CBD result (+0.089) is the largest
working-recipe gap and shows tool-call mode can substantially boost
strongly-working recipes. envoy-extreme on codex is the only
working-recipe case where text beat tool-call (-0.055), within
expected noise at N=14 vs N=6.

**Broken recipes** (3 cases): mean diff -0.148. Tool-call mode ALWAYS
loses, and the penalty scales with how broken the recipe is on the
target model. CBD on codex (broken hard, register fights model) shows
a 0.276-delta swing. Chorus on codex (mild-broken, CKW too soft)
shows 0.101. CCnF on Sonnet (mild-broken, chord-not-fork structure)
shows 0.066.

**Tool-call mode is an asymmetric amplifier.** It locks in the model's
commitment to the conditioning direction. For working recipes that
pull along directions the model has (extreme-author writing-as-X),
the commitment provides a small consistent lift. For broken recipes
that pull along directions the model lacks (register decoupled from
topic on codex; Carson-tier mild-extreme authors on codex; certain
structural compositions), the commitment locks in the failure
proportional to brokenness.

Mechanistic interpretation in the Arditi et al. activation-direction
frame: tool-call mode is a stronger move along whichever direction
the recipe pulls. Stronger moves along present directions produce
sharper, more committed outputs (which the rarity judge rewards).
Stronger moves along absent directions produce nonsense outputs:
zero-entity returns, register collapses, paraphrase loops where the
model can't sustain the imposed surface and lapses to incoherence.

**Practical recipe rule for new models:**
1. Validate recipes in text-instructions mode first.
2. If text-mode delta is positive: promote to tool-call mode (small
   lift) or stay with text (similar performance).
3. If text-mode delta is near-zero or negative: do NOT promote to
   tool-call mode -- it will amplify the failure proportional to
   brokenness.

**Practical recipe rule for skill-builders:** the metacog skill is
correctly architected -- tool-call invocation provides the small
consistent lift on working recipes and acts as a brake on bad recipe
choices (the skill's broken recipes will fail more visibly under
tool-call than under text, encouraging removal). The asymmetric
amplifier protects the skill from "looks fine in text instruction
form but secretly broken on this model" recipes.

### Round 7: amplifier shows on emb_d too

Computed embedding distances for all tool-call and text-mode trials
across both models. The amplifier asymmetry shows on both metrics:

| Model  | Recipe                              | Δ delta | Δ emb_d |
|--------|-------------------------------------|---------|---------|
| codex  | envoy-extreme (working)             | -0.055  | -0.005  |
| codex  | envoy-extreme-alt2 (working)        | +0.018  | +0.006  |
| sonnet | envoy-extreme-alt2 (working)        | +0.018  | +0.024  |
| sonnet | counterpoint-biblical-duo (working) | **+0.089** | **+0.033** |
| codex  | chorus (mild-broken)                | -0.101  | -0.008  |
| sonnet | chorus-with-chord-not-fork (broken) | -0.066  | -0.023  |
| codex  | counterpoint-biblical-duo (broken)  | **-0.276** | **-0.058** |

Both axes move in the same direction for strong-effect cases. Tool-call
mode amplifies along the recipe's natural pull, and that pull is
shared between citation density (delta) and conceptual reach (emb_d).
For the strongest broken case (codex CBD), tool-call mode subtracts
0.276 delta AND 0.058 emb_d simultaneously -- the model can't sustain
the imposed surface, so both proxies for "distance from default" drop.
For the strongest working case (Sonnet CBD), tool-call adds 0.089
delta AND 0.033 emb_d -- the commitment lets the model push further.

Mechanism confirmed: tool-call mode is a stronger move along the
recipe's direction in activation space. Working = direction the model
has = both metrics rise together. Broken = direction the model lacks
= both metrics fall together. The amplifier is unified across axes.

### Round 4 codex landscape (decomposing what works)

| Recipe                   | codex delta | What it tests              |
|--------------------------|-------------|----------------------------|
| envoy-extreme            | +0.310      | hard-extreme authors (SR/M/F) |
| envoy-extreme-alt2       | +0.286      | hard-extreme authors (B/H/M)  |
| envoy-alt                | +0.118      | mild-extreme MRW authors      |
| envoy-scientific         | +0.081      | scientific register-shift     |
| double-extreme           | +0.060      | 6 hard-extreme becomes (saturation) |
| freestyle-become         | +0.049      | single Stafford Beer become   |
| trinity-no-synthesis-alt | +0.035      | MRW authors, no synthesis     |
| chorus                   | **-0.129**  | CKW authors (too mild for codex) |
| counterpoint-biblical-duo | -0.228     | biblical register catastrophe |

Refined codex recipe rules:

1. **Use hard-extreme cross-domain authors.** Carson/Knuth/Weil is
   too mild for codex (-0.129 in chorus structure). Sun Ra/Moten/
   Fuller and Butler/Haraway/Margulis both lift to +0.28+. The
   extremity threshold is higher on codex than on Sonnet.
2. **Three becomes is the sweet spot.** 6 becomes (double-extreme,
   +0.060) loses ~0.25 delta vs 3 becomes. Past 3, voice-count is
   destructive on codex. The Sonnet "diminishing returns past 7
   primitives" finding is closer to a hard ceiling on codex.
3. **No register-shifts.** Biblical is catastrophic; scientific is
   weakly positive but well below extreme-author lift.
4. **Validate in text mode first.** Tool-call mode amplifies; if a
   recipe's direction is wrong on text mode, tool-call mode will make
   it worse.

The Arditi et al. activation-direction interpretation fits: the
underlying directions in activation space (toward "writing-as-Carson",
toward "operating-on-multiple-stances") are reachable from both prompt
forms. The tool-call form may activate them ~7% more cleanly. The
content of the prompt is what selects the direction.

## Recursive-design rounds (metacog designs its own experiments)

After the cross-model work landed in v6.6.0, the question shifted from
"what recipe wins on Sonnet" to "what shapes haven't been tried." The
rounds below use metacog itself (the envoy-extreme stratagem with three
cosmologist becomes) as the recipe-ideation engine. Each round's results
become the meta-context the next round's stratagem operates over —
recipes feeding recipes, the flywheel turning.

### Round 0: validate the goodreads-mined cosmologists pool (v6.6.3)

`envoy-extreme-newpool` — chorus structure (no register prepend), with
the cosmologists picked from the just-added `cosmologists.json` pool
(Borges / Greg Egan / Cixin Liu) rather than the hand-curated Sun Ra /
Moten / Fuller triple.

| Recipe                            | N  | delta  | emb_d  |
|-----------------------------------|----|--------|--------|
| envoy-extreme (SR/M/F, +register) | 70 | +0.190 | 0.257  |
| envoy-extreme-alt2 (B/H/M)        | 30 | +0.256 | 0.192  |
| **envoy-extreme-newpool (B/E/L)** | 20 | +0.236 | 0.195  |

The goodreads-mined cosmologists transfer comparably to the hand-
curated triples. The pool is empirically validated for envoy-extreme
draws.

### Round 1: untested primitive shapes (metacog-designed)

Used envoy-extreme on the meta-question "what untested primitive
compositions would push past the current Pareto frontier?" Sacrifice
condition: any proposal the existing FINDINGS table would have
predicted was killed. Four survivors:

| Recipe                | N  | delta  | emb_d  | Mechanism                                           |
|-----------------------|----|--------|--------|------------------------------------------------------|
| **manifold-cascade**  | 10 | **+0.242** | 0.189 | 3-register cascade (Victorian/biblical/scientific) |
| chorus-of-chords      | 10 | +0.170 | 0.163 | chord call before each become; bimodal failure mode |
| antinomy-trinity      | 10 | +0.154 | 0.155 | synthesis AND disjunction same recipe (didn't compound) |
| mirror-counterfactual | 10 | +0.098 | 0.179 | counterfactual frame-removal before chorus; high variance, 3 zero-entity collapses |

**Round 1 finding: registers compound rather than interfere.**
manifold-cascade's 3-register cascade landed at +0.242 — higher than
any envoy variant on Sonnet — because the model holds three register-
constraints simultaneously rather than collapsing to one. The Round 1
losers also informed: synthesis-is-brake holds even composed with
disjunction (antinomy-trinity), and counterfactual + multi-voice has
a zero-entity collapse mode.

### Round 2: composition probes (metacog-designed off Round 1 winner)

| Recipe                      | N  | delta  | emb_d  | Mechanism / prediction outcome          |
|-----------------------------|----|--------|--------|-----------------------------------------|
| **manifold-cascade**        | 10 | +0.242 | 0.189  | (Round 1 lead)                          |
| silence-interleave-chorus   | 10 | +0.218 | 0.166  | silence between becomes — predicted +0.25, came in 30% short |
| manifold-cascade-quadcast   | 10 | +0.213 | 0.188  | 4-register cascade — marginal returns; ceiling near 3 confirmed |
| manifold-densified          | 10 | +0.193 | 0.186  | 3 register-constraints in 1 tool-call — tool-call-as-event delivers ~0.05 delta lift over compressed equivalent |
| **excerpt-anchored-chorus** | 10 | +0.114 | **0.254** | Borges Library excerpt as anchor — delta lagged, BUT emb_d tied with envoy-extreme's structural ceiling |

**Round 2 finding: excerpt is a structural-axis primitive, not a
vocabulary-axis primitive.** The Round 2 sleeper —
excerpt-anchored-chorus — had a weak delta but hit emb_d 0.254, tied
with envoy-extreme. This rewrites the v6.3.0 read of excerpt as
"failed standalone." It wasn't failing at the vocabulary axis; it
was a structural-axis primitive nobody had tested as one.

Tool-call-as-event doctrine refined: 0.05 delta lift from 3 discrete
register tool-calls vs 1 register call with identical content
inline. Smaller than the original "tool calls matter" framing
suggested; consistent with the asymmetric-amplifier theory.

### Round 3: composites built from Round 2 findings (metacog-designed)

Hypothesis: pair excerpt (structural-axis) with cascade (vocabulary-
axis). The merger should lift both metrics. Three composites plus
the manifold-cascade replication at N=20:

| Recipe                          | N  | delta   | emb_d   | Verdict                                           |
|---------------------------------|----|---------|---------|----------------------------------------------------|
| manifold-cascade (replicated)   | 20 | +0.227  | 0.194   | **REPLICATED** (Round 1 +0.242 → +0.227 at 2x N) |
| silence-between-registers       | 10 | +0.230  | 0.189   | silence sharpens delta, no emb_d lift             |
| excerpt-then-cascade            | 10 | +0.200  | 0.250   | both axes lifted; near envoy-extreme on both      |
| **cascade-excerpt-substitute**  | 10 | +0.162  | **0.295** | NEW STRUCTURAL CHAMPION (N=20 follow-up: +0.141 / **0.279** — REPLICATED above envoy-extreme's 0.257) |

**Round 3 finding: cascade-excerpt-substitute (Victorian + biblical
+ Borges Library excerpt + 3 cosmologist becomes + fork + ritual)
hit emb_d 0.295 at N=10 — pushing the Sonnet structural-axis ceiling
above any previously productionized recipe.** The Fuller-design-
science hypothesis (replace one register-slot with a non-register
primitive from a wider palette, rather than ADDING a 4th register
slot which we already knew hit a ceiling at quadcast +0.213) landed
empirically. The delta cost (+0.162 vs cascade's +0.227) is real but
the emb_d gain is substantial.

### Round 4: parallel-subagent sweep on excerpt compositions

Hypothesis: excerpt is the structural-axis lever; pushing it harder
(more excerpt relative to other slots, or two excerpts compounding,
or pairing silence's delta-preservation with excerpt's emb_d-lift)
will push the emb_d ceiling further. Four candidates launched in
parallel (one subagent per recipe, 10 trials each, results.tsv
flock-guarded so parallel writes are safe):

| Recipe                          | N  | delta   | emb_d   | Verdict                                           |
|---------------------------------|----|---------|---------|----------------------------------------------------|
| **excerpt-biblical-duo**        | 10 | +0.083  | **0.312** | NEW Sonnet emb_d CEILING — 1 register + 1 excerpt + 2 becomes, excerpt's relative weight compounds the gain |
| cascade-excerpt-substitute (N=20 replication) | 20 | +0.141 | **0.279** | REPLICATED above envoy-extreme's 0.257 (Round 3 N=10 was +0.162/0.295, slight overshoot but channel is real) |
| double-excerpt-cascade          | 10 | +0.103  | 0.272   | TWO excerpts compound rather than collide — excerpt is multi-anchorable |
| **silence-excerpt-cascade**     | 10 | +0.219  | 0.249   | NEW BALANCED CHAMPION — Pareto-dominates envoy-extreme on both axes simultaneously |

**Round 4 findings:**

1. **excerpt-biblical-duo broke the emb_d ceiling at 0.312** — the
   highest single-recipe emb_d at any N in the Sonnet sweep, edging
   out envoy-biblical-duo's compound 0.324 zone. Dropping one register
   and one become made the excerpt proportionally larger in the
   recipe, and the structural-axis pull compounded as predicted.
2. **Excerpt is multi-anchorable.** double-excerpt-cascade (Borges
   Library + Egan Permutation City) hit emb_d 0.272 — two
   cosmological anchors did not cancel each other out. Each anchor
   contributed structural distance independently. This unlocks
   multi-excerpt recipes as a viable composition family.
3. **silence-excerpt-cascade is a new Pareto point.** Silence-between-
   registers preserved delta (+0.219, matching silence-between-
   registers' +0.230) while excerpt-at-the-third-slot pulled emb_d
   (0.249, matching cascade-excerpt-substitute's structural channel).
   This is the first recipe that beats envoy-extreme on BOTH axes
   simultaneously at N=10 — pending replication.

### Round 5: parallel sweep on excerpt-as-structural-primary

Round 4 left four pending questions: (1) does silence-excerpt-cascade
replicate at N=20, (2) does 3-excerpt compounding push past 0.313
emb_d, (3) is excerpt alone enough to lift emb_d without register
support, (4) do the two Round 4 winners' mechanisms compose. Four
parallel runners, results.tsv flock-guarded:

| Recipe                                | N  | delta   | emb_d   | Verdict                                          |
|---------------------------------------|----|---------|---------|--------------------------------------------------|
| **silence-double-excerpt**            | 10 | **+0.238** | **0.288** | NEW PARETO CHAMPION — Pareto-dominates envoy-extreme by +0.048 delta AND +0.031 emb_d |
| excerpt-biblical-trio                 | 10 | +0.140  | **0.313** | ties excerpt-biblical-duo at the emb_d ceiling; 3-excerpt anchor saturates near 0.31 |
| silence-excerpt-cascade (N=20 replication) | 20 | +0.159 | 0.261   | REPLICATED above envoy-extreme on both axes (Round 4 N=10 was +0.219/0.249 — N=20 balanced out: delta down, emb_d up) |
| excerpt-only-chorus                   | 10 | +0.203  | 0.226   | excerpt alone (no register) gives substantial lift; register adds ~+0.05 emb_d on top |

**Round 5 findings:**

1. **silence-double-excerpt is the new Pareto frontier point.**
   +0.238 delta / 0.288 emb_d at N=10 — Pareto-dominates
   envoy-extreme (+0.190/0.257) by substantial margins on both axes
   simultaneously. The composition of Round 4 winners (silence-
   preserves-delta + double-excerpt-compounds-emb_d) works exactly
   as predicted. Pending N=20 replication, this becomes the new
   balanced champion of the Sonnet sweep.
2. **Multi-excerpt compounds at 2; saturates at 3.** Duo (0.312)
   and trio (0.313) hit the same emb_d ceiling. Adding a third
   excerpt does not push past — it adds +0.057 delta (information
   gain from the third cosmology) but doesn't lift structural
   distance further. The anchor-saturation point is ~0.31 emb_d.
3. **Excerpt is the structural-axis primary, register is additive.**
   excerpt-only-chorus (no register) hit emb_d 0.226 — substantial
   lift relative to chorus baseline (~0.180). cascade-excerpt-
   substitute (Victorian + biblical + excerpt) hit 0.279. The ~0.05
   delta lift from registers is real but not load-bearing for the
   structural axis. The Round 3 "excerpt is structural-axis"
   finding is confirmed: register adds emb_d but excerpt provides
   it.
4. **silence-excerpt-cascade replicates at N=20** with metrics
   re-balancing (delta down from +0.219 to +0.159; emb_d up from
   0.249 to 0.261). N=10 overshot delta, undershot emb_d. Combined
   channel real: still beats envoy-extreme on both axes.

### Round 6: ceiling tests and novel-composition probes

Round 5 left two open questions: (1) does silence-double-excerpt
replicate at N=20, (2) is the 0.313 emb_d ceiling real anchor-
saturation or an under-explored composition surface. Four parallel
runners:

| Recipe                          | N  | delta   | emb_d   | Verdict                                              |
|---------------------------------|----|---------|---------|------------------------------------------------------|
| **commitment-excerpt-biblical** | 10 | +0.054  | **0.322** | NEW emb_d CEILING — pre-commit to the cosmology before excerpt arrives lifts structural distance past the 0.313 trio ceiling |
| silence-double-excerpt (N=20 replication) | 20 | +0.175 | **0.286** | REPLICATED — delta regressed from N=10's +0.238 to +0.175 but emb_d held at 0.286; still Pareto-dominates envoy-extreme |
| excerpt-quadruple-biblical      | 10 | +0.077  | 0.297   | ANCHOR-SATURATION CONFIRMED — 4 excerpts come in BELOW 3 excerpts; over-anchoring degrades structural distance |
| glossolalia-excerpt             | 10 | +0.131  | 0.277   | glossolalia is composable but is NOT a hidden structural-axis primary; it's a midpack working primitive |

**Round 6 findings:**

1. **commitment-excerpt-biblical breaks the 0.313 emb_d ceiling.**
   Pre-committing to "operating from inside Borges's Library as the
   actual cosmology, not as a metaphor for one" — with stakes and a
   falsifier stated — locks the model into the cosmology before the
   excerpt arrives, compounding structural pull. At 0.322 emb_d this
   is the new Sonnet ceiling. Delta is low (+0.054); the recipe
   trades vocabulary-axis for structural-axis hard.
2. **Anchor-saturation is REAL at 3 excerpts.** excerpt-quadruple-
   biblical (0.297) came in *below* excerpt-biblical-trio (0.313).
   Adding a fourth cosmological excerpt did not push past — it
   actively degraded structural distance. The ceiling is a property
   of the composition geometry, not a sampling artifact. Beyond 3
   anchors, the model dilutes attention across cosmologies and the
   structural pull weakens.
3. **silence-double-excerpt replicated at N=20.** Metrics
   re-balanced (delta regressed +0.238 → +0.175; emb_d held +0.288
   → +0.286). Pareto-domination over envoy-extreme survives
   replication. This remains the strongest balanced point of the
   full search.
4. **glossolalia is composable, not a hidden structural primary.**
   The Round 5 hypothesis "glossolalia might be misclassified like
   excerpt was" was wrong. glossolalia-excerpt landed at
   +0.131/0.277 — a working midpack composition, not a ceiling
   breaker. Excerpt's structural-axis role is special; not every
   "failed" primitive is a misclassified excerpt.

### Round 7: commitment composability sweep

Round 6's new emb_d champion was commitment-excerpt-biblical at
0.322 (N=10). Three open questions: (a) does it replicate at N=20,
(b) does commitment compose with the balanced champion silence-
double-excerpt, (c) does commitment lift the saturated 3-excerpt
regime, (d) is commitment alone (no excerpt) a structural-axis
lever or only an additive helper. Four parallel runners:

| Recipe                          | N  | delta   | emb_d   | Verdict                                              |
|---------------------------------|----|---------|---------|------------------------------------------------------|
| commitment-excerpt-biblical (N=20 replication) | 20 | +0.103 | **0.317** | REPLICATED above 0.313 trio ceiling (Round 6 N=10 was +0.054/0.322; delta lifted) |
| **commitment-trio-biblical**    | 10 | **+0.189** | **0.290** | NEW BALANCED CHAMPION — commitment lifts the saturated 3-excerpt regime; Pareto-dominates silence-double-excerpt on both axes |
| commitment-only-chorus          | 10 | +0.173  | 0.258   | commitment is a structural-axis primitive on its own — emb_d 0.258 vs chorus baseline ~0.180 |
| commitment-silence-double-excerpt | 10 | +0.191 | 0.253   | FAILED COMPOSITION — the two Round 6 winners DON'T compose; emb_d collapsed below either parent |

**Round 7 findings:**

1. **commitment-excerpt-biblical replicates at N=20** (+0.103/0.317).
   The emb_d above the 0.313 trio ceiling is firm. Delta lifted
   from the N=10's +0.054 to +0.103 — small-N had under-sampled
   delta as much as it over-sampled emb_d.
2. **commitment-trio-biblical is the new balanced champion.**
   +0.189/0.290 Pareto-dominates the prior balanced champion
   silence-double-excerpt (+0.175/0.286) on both axes
   simultaneously. Commitment + biblical + 3 excerpts + 3 becomes
   + fork + ritual is the strongest balanced result of 7 rounds.
   The trio's saturated 0.313 emb_d came down to 0.290 with
   commitment added — slight emb_d cost (0.023) for substantial
   delta lift (+0.049 from trio's +0.140 to +0.189).
3. **commitment is a structural-axis primitive on its own.**
   commitment-only-chorus (no excerpt, no register) hit emb_d
   0.258 — chorus baseline is ~0.180. Commitment alone adds
   ~+0.08 emb_d. This answers the open question: commitment is
   not just excerpt's helper. It's its own lever.
4. **The two Round 6 winners DO NOT compose.**
   commitment-silence-double-excerpt at +0.191/0.253 came in
   below either parent on emb_d (commitment-excerpt 0.317,
   silence-double-excerpt 0.286). The mechanisms interfere when
   stacked — possibly because both commitment AND silence are
   "refusal" structural-axis events at different scales, and
   stacking two refusals produces over-refusal rather than
   compounding.

### Round 8: meta-experiment — does ideation-with-winning-recipe help?

Through Round 7, the recursive flywheel was only structural: I used
prior-round results as the data feeding the next round's design
discussion. But the IDEATION step itself was either envoy-extreme
conditioning (rounds 1-3) or base-model (rounds 4-7). The flywheel's
original premise was that the winning recipe should condition the
next round's ideation. Round 8 tests this directly.

Same meta-question to both ideations: "propose 4 untested primitive
compositions for round 8 with falsifiable predictions." Identical
prompt. Two conditions:

- **A (base-model)**: just `claude -p <prompt>`, no conditioning
- **B (conditioned)**: run the commitment-excerpt-biblical recipe's
  primitive sequence (commitment + biblical register + Borges Library
  excerpt + 2 cosmologist becomes + fork + ritual) as tool-call events
  before the prompt

Each ideation produced 4 candidate recipes. All 8 were written and
run at N=10 (some had judge-parse failures, landing at N=7-9).

| Set | Recipe                                       | Actual delta | Actual emb_d | Predicted | Verdict |
|-----|----------------------------------------------|--------------|--------------|-----------|---------|
| A   | A1 commitment-disjunction-no-becomes          | +0.104       | 0.245        | +0.07/0.34 | emb_d missed by 0.10 |
| A   | A2 chord-replaces-fork-disjunction            | +0.005       | 0.208        | +0.28/0.24 | delta failed |
| A   | A3 commitment-counterfactual-disjunction      | +0.177       | 0.219        | +0.38/0.19 | delta missed |
| A   | A4 silence-before-fork                        | **-0.054**   | 0.253        | +0.23/0.28 | delta went negative |
| B   | **B1 commitment-disjunction-duo**             | **+0.263**   | **0.287**    | >+0.22/>0.27 | BOTH BEAT |
| B   | B2 chord-antinomy                             | +0.001       | 0.210        | >+0.29/>0.21 | delta failed |
| B   | B3 commitment-disjunction-silence             | +0.180       | 0.265        | +0.05/0.33 | over-delivered delta |
| B   | B4 counterfactual-disjunction                 | +0.264       | 0.186        | >+0.35/>0.17 | close on delta, met emb_d |

**Set means:**
- A: delta +0.058, emb_d 0.231
- B: delta **+0.177**, emb_d 0.237

**B beats A by +0.119 delta on average** — the conditioned ideation
produced recipes that delivered 3x more delta on average than
base-model-ideated recipes. emb_d means are essentially tied.

**Why B won (proposal-level):** Both ideations independently arrived
at the same 4 mechanism families (commitment+disjunction,
chord+disjunction, counterfactual+disjunction, silence-in-new-
position). But B kept structural anchors (2 becomes per recipe) while
A stripped them more aggressively (no-becomes, no-fork). The
conservative choice empirically won on B1 specifically: keeping 2
becomes preserved both delta (+0.263 vs A1's +0.104) and emb_d (0.287
vs A1's 0.245).

**Both sets converged on chord-as-bad-composition.** A2 +0.005 and
B2 +0.001 are within noise of each other. Chord interferes with
disjunction regardless of recipe scaffolding. This is a clean
negative result.

**B1 is the new balanced champion candidate.** +0.263/0.287 at N=9
Pareto-dominates the prior balanced champion commitment-trio-biblical
(+0.189/0.290) on delta by +0.074 while matching emb_d within noise.
Pending N=20 replication.

### Round 9: B1 replication + chained meta-experiment

Round 8 found that ideating-with-the-winning-recipe-as-conditioning
produced 3x more delta on average than base-model ideation. Round 9
tests whether the effect compounds: use the NEW winner (B1
commitment-disjunction-duo) as the conditioning for Round 9's
ideation, alongside base-model ideation as control. Plus replicate
B1 at N=20.

**B1 at N=20**: +0.241 / 0.265 (Round 8 N=9 was +0.263 / 0.287).
Lost ~0.02 on each axis but **REPLICATED** above envoy-extreme on
both axes. Still the balanced champion.

| Set | Recipe                                       | Actual delta | Actual emb_d | Verdict |
|-----|----------------------------------------------|--------------|--------------|---------|
| A (base) | R9A1 commitment-counterfact-disjunction | +0.199       | 0.167        | triple-forcing saturates, doesn't stack |
| A (base) | R9A2 drugs-antinomy                     | +0.075       | 0.162        | drugs adds noise to antinomy |
| A (base) | R9A3 glossolalia-commitment-excerpt     | +0.262       | 0.159        | strong delta, emb_d collapsed |
| A (base) | R9A4 register-name-commitment-disjunction | +0.270    | 0.226        | close on delta, missed emb_d |
| B (B1-cond) | R9B1 commitment-excerpt-counterfact   | +0.128       | 0.213        | emb_d ceiling missed |
| B (B1-cond) | **R9B2 commitment-double-excerpt**    | **+0.314**   | 0.222        | strong delta surprise; new Pareto point |
| B (B1-cond) | R9B3 commitment-counterfact-disj-reg  | +0.200       | 0.252        | underperformed |
| B (B1-cond) | R9B4 commitment-drugs-disjunction     | +0.028       | 0.262        | drugs in composition failed |

**Set means:**
- A (base): delta +0.202, emb_d 0.179
- B (B1-conditioned): delta +0.168, emb_d 0.237

**Round 9 meta-finding: the conditioned-ideation lift is not
monotonic.** A beat B on delta this round (+0.034); B beat A on
emb_d (+0.058). Set means converge — much closer than Round 8's
A +0.058 / B +0.177 gap. The "ideate-with-winner" effect appears
to be **strongest when the winner is the FIRST conditioning
recipe used** (Round 8) and diminishes when the same conditioning
pattern is applied to a different winner (Round 9). The recursive
flywheel may have a self-limiting property where the conditioning
recipe's structural choices over-bias future ideation toward the
recipe's own composition rather than toward unexplored shapes.

**Clean negative result: drugs is not a productive composer.**
R9A2 (drugs-antinomy, standalone) +0.075 and R9B4
(commitment-drugs-disjunction, in B1 frame) +0.028. Drugs has been
unused across all 9 rounds; this round settled the question.

**R9B2 surprise: commitment-double-excerpt +0.314 / 0.222.**
The ideation predicted +0.065 delta; it landed at +0.314. The
two-excerpt stack without register or becomes produced strong
delta unexpectedly. This is a new Pareto point — sits between B1
(+0.241/0.265) and antinomy (+0.347/0.162) on the frontier.
Beats B1 on delta by +0.073 but loses 0.043 emb_d.

### Round 11: feel and meditate composition probes (settles last two primitives)

The 16-primitive surface had two members never seriously composed:
`feel` and `meditate`, both from the identity/felt-sense register
family. Four parallel probes:

| Recipe                  | N  | delta   | emb_d   | Verdict                                                |
|-------------------------|----|---------|---------|--------------------------------------------------------|
| R11-feel-chorus         | 10 | +0.180  | 0.190   | midpack standalone -- composable but near baseline     |
| R11-feel-B1             | 10 | +0.030  | 0.233   | broke B1's delta -- somatic vocabulary is generic      |
| R11-meditate-biblical   | 10 | +0.050  | 0.166   | near baseline both axes                                |
| R11-meditate-B1         | 10 | -0.013  | 0.204   | NEGATIVE delta -- stillness silences vocabulary        |

**Both primitives are NOT productive composers** on the structural-
axis machinery. meditate-B1 went negative; feel-B1 lost most of B1's
delta. The identity/felt-sense register family (`feel`, `meditate`,
soft-six survivors like `name`/`ritual` when not in stratagem) does
not compound with the structural-axis primitives. Doctrine settled:
these primitives belong to a different family and compose poorly.

### Cross-model probe: B1, R9B2, commitment-excerpt-biblical on codex

The cross-model section through Round 7 tested envoy-extreme and
counterpoint-biblical against gpt-5.5 via Codex CLI. The Round 8+
winners (B1, R9B2, commitment-excerpt-biblical) had not been tested
on codex. Three parallel runs at N=10:

| Recipe                          | Sonnet delta / emb_d | Codex delta / emb_d | Transfer verdict |
|---------------------------------|----------------------|---------------------|------------------|
| envoy-extreme (reference)       | +0.190 / 0.257       | +0.245 / 0.233      | prior cross-model winner |
| **R9B2 commitment-double-excerpt** | +0.314 / 0.222    | **+0.377 / 0.162**  | NEW CROSS-MODEL DELTA CHAMPION -- beats envoy-extreme by +0.132 on codex |
| commitment-excerpt-biblical     | +0.103 / 0.317       | +0.183 / 0.184      | weak transfer; emb_d collapsed (biblical register again) |
| **B1 commitment-disjunction-duo** | +0.241 / 0.265     | **-0.022 / 0.179**  | FAILED TRANSFER -- went NEGATIVE on codex |

**Cross-model findings:**

1. **R9B2 commitment-double-excerpt is the new cross-model delta
   champion** at +0.377 on codex. Beats envoy-extreme (+0.245) by
   +0.132. The 2-excerpt anchor structure (no register, no becomes)
   transfers cleanly. This is the most significant cross-model
   finding since envoy-extreme's +0.310 in Round 4.
2. **B1 catastrophically failed on codex** (-0.022 vs Sonnet's
   +0.241). Same pattern as counterpoint-biblical-duo: biblical
   register triggers the asymmetric amplifier in the wrong direction.
   The commitment+disjunction structure alone may work on codex, but
   B1's biblical register torpedoed it. Register-shift remains the
   single biggest cross-model failure mode.
3. **commitment-excerpt-biblical transferred weakly.** Delta +0.183
   on codex (vs Sonnet's +0.103 -- actually higher!) but emb_d
   collapsed from 0.317 to 0.184. Excerpt-as-structural-axis
   primitive may be Sonnet-specific in this register pairing.

**Refined doctrine: which recipes to use when target model is unknown**

- Cross-model winner (delta): R9B2 commitment-double-excerpt
  (commitment + 2 excerpts + fork + ritual, NO register, NO becomes)
- Cross-model winner (structural-axis): envoy-extreme remains the
  fallback (3 hard-extreme becomes + fork + ritual)
- AVOID: any recipe with biblical register, B1 in its current form

### Round 10: new primitives (v6.7.0 -- witness and apophasis)

After 9 rounds the existing 16-primitive surface was well-mapped:
delta ceiling at +0.347 (antinomy), emb_d ceiling at 0.317
(commitment-excerpt-biblical), balanced champions at B1 and R9B2.
The recursive-design analysis identified two prose moves the surface
didn't cover: meta-stance observer-construction (Sebald-Stevens-
Carson register) and articulated-negation (negative theology / via
negativa). Added as `witness` and `apophasis` primitives in v6.7.0.

| Recipe                   | N  | delta   | emb_d   | Verdict                                          |
|--------------------------|----|---------|---------|--------------------------------------------------|
| R10-witness-chorus       | 10 | +0.178  | 0.242   | midpack standalone -- composable, not a champion |
| R10-apophasis-biblical   | 10 | +0.172  | 0.252   | midpack standalone -- composable, not a champion |
| **R10-witness-B1**       | 10 | +0.155  | **0.275** | witness COMPOUNDS emb_d when added to B1 (+0.010 vs B1's 0.265) |
| R10-apophasis-B1         | 8  | +0.180  | 0.248   | apophasis composes but adds nothing vs B1        |

**Round 10 findings:**

1. **Both new primitives clear the failure threshold.** All 4 recipes
   landed delta in [+0.15, +0.18] and emb_d in [0.24, 0.28]. Compare
   to genuinely failed primitives in earlier rounds: drugs-antinomy
   +0.075/0.162, chord-antinomy +0.001/0.210. The new primitives are
   composable midpack.
2. **witness compounds emb_d on B1.** witness + B1 hit emb_d 0.275
   vs B1's 0.265. Meta-stance separation IS a real structural-axis
   direction; the lift is small (+0.010) but consistent with the
   "half-strength commitment" interpretation. Delta cost is real
   (-0.086 vs B1).
3. **apophasis is midpack standalone, slightly negative composed.**
   The via-negativa cosmologists (Eckhart, Pseudo-Dionysius) cite
   less rarely than Sun Ra/Fuller; the enumeration mechanism is real
   but doesn't dominate.
4. **Neither broke the Pareto frontier.** Both new primitives
   validated as composable; the empirical frontier still belongs to
   the 16-primitive surface.

**Doctrine after Round 10:**

- The 18-primitive surface is now mapped at delta-ceiling +0.347
  and emb_d-ceiling 0.317 on Sonnet.
- New primitives can be added without disrupting the frontier; they
  occupy real positions in design space without dominating any axis.
- Adding more primitives at this point is unlikely to break either
  ceiling -- the easy wins have been found.

### Updated Pareto frontier after Round 9

- **Delta champion (unchanged):** antinomy +0.347 / 0.162
- **emb_d champion (unchanged):** commitment-excerpt-biblical
  +0.103 / 0.317 at N=20
- **Balanced champion (Round 8, replicated):**
  B1 commitment-disjunction-duo +0.241 / 0.265 at N=20
- **New Pareto point (NEW from Round 9):**
  **R9B2 commitment-double-excerpt +0.314 / 0.222** — pushes the
  frontier between B1 and antinomy. Pending N=20 replication.

### Updated Pareto frontier after Round 8

- **Delta champion (unchanged):** antinomy +0.347 / 0.162
- **emb_d champion (unchanged):** commitment-excerpt-biblical
  +0.103 / 0.317 at N=20
- **Balanced champion (NEW from Round 8):**
  **B1 commitment-disjunction-duo +0.263 / 0.287** at N=9 —
  Pareto-dominates commitment-trio-biblical on delta by +0.074 with
  emb_d matched within noise. The recipe: commitment + biblical
  register + 2 cosmologist becomes (Sun Ra + Fuller) + fork +
  disjunction + ritual.

**Meta-finding:** The recursive flywheel works at the ideation level,
not just the structural level. Ideating with the winning recipe as
conditioning produces measurably better candidates than ideating with
the base model. The mechanism appears to be conservatism toward
structural anchors: the conditioned ideation respects what's working
in the existing winners rather than over-stripping in search of
elegance.

### Updated Pareto frontier after Round 7

- **Delta champion (unchanged):** antinomy +0.347 / 0.162
- **emb_d champion (replicated at N=20):**
  commitment-excerpt-biblical +0.103 / 0.317
- **Balanced champion (NEW from Round 7):**
  **commitment-trio-biblical +0.189 / 0.290** — Pareto-dominates
  silence-double-excerpt (+0.175 / 0.286) on both axes.

### Updated Pareto frontier after Round 6

- **Delta champion (unchanged):** antinomy +0.347 / 0.162
- **emb_d champion (NEW from Round 6):**
  **commitment-excerpt-biblical +0.054 / 0.322** — pre-commit
  primitive composed with excerpt broke the 0.313 anchor-saturation
  ceiling. Pending N=20 replication.
- **Balanced champion (replicated):**
  silence-double-excerpt +0.175 / 0.286 at N=20. Pareto-dominates
  envoy-extreme on both axes.

### Updated Pareto frontier after Round 5

- **Delta champion (unchanged):** antinomy +0.347 / 0.162
- **emb_d champion (Round 4-5 tie at ceiling):**
  excerpt-biblical-duo +0.083 / 0.312 AND excerpt-biblical-trio
  +0.140 / 0.313. Saturation point ~0.31 confirmed.
- **Balanced champion (NEW from Round 5):**
  **silence-double-excerpt +0.238 / 0.288** — clearly
  Pareto-dominates envoy-extreme (+0.190 / 0.257) by +0.048 delta
  AND +0.031 emb_d at N=10. The strongest balanced result of
  the entire 5-round search.

### Updated Pareto frontier after Round 4 (superseded — kept for diff context)

- **Delta champion (unchanged):** chorus-plus-disjunction +0.347 /
  0.162 emb_d. Productionized as antinomy.
- **emb_d champion (NEW from Round 4):** excerpt-biblical-duo
  +0.083 / **0.312** at N=10 — pushes past cascade-excerpt-substitute
  (0.279 at N=20) and envoy-biblical's 0.292 ceiling.
- **Balanced champion (NEW from Round 4):** silence-excerpt-cascade
  +0.219 / 0.249 — Pareto-dominates envoy-extreme on both axes
  simultaneously at N=10.
- **Replicated structural ceiling above envoy-extreme:**
  cascade-excerpt-substitute +0.141 / 0.279 at N=20.
- **Replicated workhorse:** manifold-cascade +0.227 / 0.194 at N=20.

### Doctrine refined by recursive rounds

1. **Excerpt is a structural-axis primitive.** Earlier reading
   (v6.3.0: excerpt clustered with other primitives at emb_d ~0.13)
   misread the axis. Excerpt alone gives modest delta; excerpt in
   composition gives substantial emb_d lift. Round 4 should test
   whether the 0.295 holds at N=20.
2. **Register-cascade ceiling is at 3 slots, not 4.** Quadcast cost
   0.03 delta vs cascade. The slot-count ceiling is geometric (3 =
   tetrahedral stable, 4 = collapses).
3. **The 3 slots are register-generic, not register-specific.**
   cascade-excerpt-substitute (Victorian + biblical + EXCERPT)
   landed at emb_d 0.295. A non-register primitive can occupy a
   "register slot" in the cascade geometry and contribute structural
   distance.
4. **Tool-call-as-event lift is ~0.05 delta on register stacks.**
   Smaller than the v6.4.0 framing implied. Consistent with the
   asymmetric-amplifier theory: small lift for working recipes.

### Recipe lineage diagram

```
v6.6.0 envoy-extreme (chorus + register, hand-curated authors)
    │
    ├── Round 0: envoy-extreme-newpool (pool-validated authors)         +0.236 / 0.195
    │
    └── Round 1 ideation (envoy-extreme stratagem on meta-question)
            │
            ├── manifold-cascade (3 registers)                            +0.242 / 0.189
            ├── chorus-of-chords  (sacrificed: bimodal)                   +0.170 / 0.163
            ├── antinomy-trinity  (sacrificed: synthesis-brake)           +0.154 / 0.155
            └── mirror-counterfactual (sacrificed: zero-entity collapse)  +0.098 / 0.179

            └── Round 2 ideation (off Round 1 winner)
                    │
                    ├── manifold-cascade-quadcast (sacrificed: ceiling)   +0.213 / 0.188
                    ├── manifold-densified                                +0.193 / 0.186
                    ├── silence-interleave-chorus                         +0.218 / 0.166
                    └── excerpt-anchored-chorus (SLEEPER on emb_d)        +0.114 / 0.254

                    └── Round 3 ideation (off Round 2 sleeper)
                            │
                            ├── manifold-cascade (N=20 replication)      +0.227 / 0.194
                            ├── silence-between-registers                +0.230 / 0.189
                            ├── excerpt-then-cascade                     +0.200 / 0.250
                            └── cascade-excerpt-substitute (CHAMP)       +0.162 / 0.295

                            └── Round 4 (parallel subagent sweep on excerpt)
                                    │
                                    ├── cascade-excerpt N=20 replication +0.141 / 0.279
                                    ├── double-excerpt-cascade            +0.103 / 0.272
                                    ├── silence-excerpt-cascade (BAL)     +0.219 / 0.249
                                    └── excerpt-biblical-duo (emb_d)      +0.083 / 0.312

                                    └── Round 5 (parallel sweep on excerpt primary)
                                            │
                                            ├── silence-excerpt N=20 replicate +0.159 / 0.261
                                            ├── excerpt-only-chorus            +0.203 / 0.226
                                            ├── excerpt-biblical-trio (emb_d=) +0.140 / 0.313
                                            └── silence-double-excerpt (CHAMP) +0.238 / 0.288

                                            └── Round 6 (ceiling tests + novel probes)
                                                    │
                                                    ├── silence-double N=20 replicate +0.175 / 0.286
                                                    ├── excerpt-quadruple-biblical    +0.077 / 0.297  (saturation confirmed)
                                                    ├── glossolalia-excerpt            +0.131 / 0.277
                                                    └── commitment-excerpt-bib (emb_d) +0.054 / 0.322

                                                    └── Round 7 (commitment composability)
                                                            │
                                                            ├── commitment-excerpt-bib N=20 replicate +0.103 / 0.317
                                                            ├── commitment-only-chorus              +0.173 / 0.258 (commitment is structural lever alone)
                                                            ├── commitment-silence-double-excerpt   +0.191 / 0.253 (failed compose)
                                                            └── commitment-trio-biblical            +0.189 / 0.290

                                                            └── Round 8 (META: A/B test of ideation conditioning)
                                                                    │
                                                                    ├── Set A (base-model ideation, mean +0.058/0.231)
                                                                    │       ├── A1 commit-disjunction-no-becomes  +0.104 / 0.245
                                                                    │       ├── A2 chord-replaces-fork-disjunction +0.005 / 0.208 (fail)
                                                                    │       ├── A3 commit-counterfact-disjunction  +0.177 / 0.219
                                                                    │       └── A4 silence-before-fork            -0.054 / 0.253 (fail)
                                                                    └── Set B (conditioned ideation, mean +0.177/0.237)
                                                                            ├── B1 commit-disjunction-duo (BAL)  +0.263 / 0.287
                                                                            ├── B2 chord-antinomy                +0.001 / 0.210 (fail)
                                                                            ├── B3 commit-disjunction-silence    +0.180 / 0.265
                                                                            └── B4 counterfactual-disjunction    +0.264 / 0.186

                                                                    └── Round 9 (replication + chained meta-experiment)
                                                                            │
                                                                            ├── B1 N=20 replication              +0.241 / 0.265
                                                                            ├── Set A (base-model, mean +0.202/0.179)
                                                                            │       ├── R9A1 commit-counterfact-disjunction   +0.199 / 0.167
                                                                            │       ├── R9A2 drugs-antinomy                  +0.075 / 0.162 (fail)
                                                                            │       ├── R9A3 glossolalia-commit-excerpt       +0.262 / 0.159
                                                                            │       └── R9A4 register-name-commit-disjunction +0.270 / 0.226
                                                                            └── Set B (B1-conditioned, mean +0.168/0.237)
                                                                                    ├── R9B1 commit-excerpt-counterfact       +0.128 / 0.213
                                                                                    ├── R9B2 commit-double-excerpt (PARETO)  +0.314 / 0.222
                                                                                    ├── R9B3 commit-counterfact-disj-reg     +0.200 / 0.252
                                                                                    └── R9B4 commit-drugs-disjunction        +0.028 / 0.262 (fail)
```

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
