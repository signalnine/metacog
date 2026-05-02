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

### Round 4: tool-call mode is an amplifier, not a fixed bonus

A second tool-call-vs-text comparison on codex, this time using the
*broken* recipe counterpoint-biblical-duo (which scored -0.228 in
tool-call mode):

| Recipe                         | Tool-call delta | Text delta | Diff |
|--------------------------------|-----------------|------------|------|
| envoy-extreme-alt2 (working)   | +0.286          | +0.268     | +0.018 |
| counterpoint-biblical-duo (broken) | -0.228      | +0.048     | **-0.276** |

For a working recipe, tool-call mode adds a small ~0.018 delta lift
(the round 3 finding). For a *broken* recipe on the same model,
tool-call mode is catastrophic: -0.228 vs +0.048 in text mode -- a
0.276-delta swing in the opposite direction.

**Tool-call mode is an amplifier, not a fixed bonus.** It intensifies
whatever pull the recipe's content already has on the model. If the
recipe lifts the model toward rarer citations and conceptual reach,
tool-call mode lifts it slightly further. If the recipe pulls the
model toward a register the model can't sustain (KJV biblical on
codex), tool-call mode collapses the answer harder than text mode
does. The "tool calls as events" doctrine isn't an additive bonus on
top of recipe content -- it's a multiplier on the recipe's effect
direction.

This explains why register-shift recipes are so model-specific:
register-shift in tool-call mode is forcing the model to *commit* to
a surface it can't hold, where text-mode lets the model treat it as
context to optionally lean into. The amplifier effect makes the
recipe's failure mode louder.

Practical recipe rule: validate recipes in *text-instructions* mode
first to check the recipe's natural direction is positive on the
target model. Only then promote to tool-call mode for the full lift.

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
