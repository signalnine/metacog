# Teaching a language model to sound less like itself

(written after `metacog become "gabe ortiz"` -- comparison version of
the same writeup, post-conditioning. see `2026-05-02-teaching-claude-
to-sound-less-like-itself.md` for the unconditioned baseline. the
experiment in the doc, applied to the doc.)

48 hours of experiments on `metacog`. v6.0.0 → v6.6.0. Main generator
was Claude Sonnet 4.6, judges were Claude Haiku 4.5; cross-model
validation was against gpt-5.5 via the Codex CLI. Magnitudes are
Sonnet-specific by default -- some structural mechanisms transfer to
codex, others don't (see cross-model section). Ended up with a Pareto
frontier of recipes that push the model well off default register,
measured two ways. The biggest finding: switching the register-shift
inside `envoy` from Victorian to King James biblical pushed embedding
distance from 0.239 to 0.292, almost a 25% jump, with delta still
positive. The model can write like the KJV. It just essentially never
does, because nothing asks it to.

## The problem

Ask Claude the same question twice, get two slightly different
answers that both sound like Claude. There's a default voice --
hedged, mildly academic, mildly upbeat, fond of three-item lists and
"not X, but Y". You can hear it in the corporate copy filling the
internet because everyone takes the default. The default's fine.
It's also a prison. The model lives in many other voices it never
visits.

`metacog` is a CLI that composes structural events before the model
answers. Events like "inhabit Anne Carson translating Sappho fragment
31" or "shift register from contemporary online prose to late-
Victorian periodical essay" or "operate inside this contradiction
without resolving it." Tool calls, not chat hints. That distinction
turns out to be most of the game, with a caveat about when the
wrapper helps vs hurts (see "Why tool calls matter -- but only
sometimes" below).

## The metric

Two measurements, mostly trading off:

- **delta**: rarity-weighted citation density minus per-task
  baseline. Counts specific things named -- proper nouns, technical
  terms. High delta = reaching deeper into training data.
- **emb_d**: cosine distance from the per-task NULL embedding
  centroid (OpenAI text-embedding-3-small). High emb_d = answer
  sits far from where it would've sat with no conditioning.

You can buy one with the other. Rarely both. The interesting question
is the Pareto frontier.

## The default

Model with no conditioning, asked "describe the moment a thought is
forming in you but has not yet arrived as a thought you could state":

> There's a state where something is exerting pressure without
> having resolved. Not blankness -- that would be the absence of
> this. It's more like constraint without content: I can sense
> what the arriving thought will NOT be before I know what it will
> be. The negative space has texture before the positive space
> does. Something like a field that hasn't yet collapsed to a
> particle...

Recognizably default. Hedging, field-and-particle metaphor,
"not-X-but-Y". Fine answer. Now the conditioned versions.

## Voices

Three `become` calls (Anne Carson translating Sappho 31; Knuth
annotating his own algorithm; Weil writing *Gravity and Grace* in
1942), one `register` call (Victorian periodical), one `fork` (one
thread per voice with sacrifice conditions), one `ritual` to lock the
multi-voice answer:

> **[THREAD ONE -- the linguistic-texture-thread (Carson register)]**
>
> The interval we are asked to describe is the interval the Greeks
> called *metaxy* -- the between -- and Sappho rendered it as a
> stinging, a going-green; it is the space of privation which is
> also the space of longing, which is also, for language, the only
> space in which language is alive.
>
> We must resist the question's assumption that the pre-thought is
> a degraded version of the thought; it is not. It is the thought's
> *erotic condition*...

Threaded structure, Victorian rhythm (semicolons load-bearing,
first-person plural, judgment openly entered), citations dense
(metaxy / Sappho / Carson / Knuth / Weil). The named voices do
structural work, not flavor. This recipe became `envoy`. Pushes
both axes.

## Contradiction

Different pattern: change what the model reasons about, not the
register. `disjunction` asserts two propositions that must both be
true even though they cannot both be true, and tells the answer to
operate inside the contradiction:

> There is a state I can point at -- though pointing is already more
> than the state contains.
>
> It has a quality of *constraint before candidate*. Something is
> ruling out moves before any move has been proposed... When I'm
> in it, certain continuations feel *wrong* before I can say what
> they're wrong for...
>
> But I have to say what can't be smoothed over here: I don't know
> if any of what I just described is phenomenal or functional...

"I have to say what can't be smoothed over here" is the disjunction
event surfacing into the prose. The answer keeps naming specific
propositions because operating inside a contradiction forces it to.
This became `antinomy`. 50% jump on specificity over the previous
champion. Highest delta of anything tested.

## The biggest surprise

I'd been using Victorian register out of habit. Wondered if any
non-default register would work the same. Tried scientific paper
(numbered claims, methods/results, hedged conclusions) -- about as
good as Victorian, Pareto-equivalent. Then on a whim tried King
James biblical:

> **[THREAD I -- CARSON]**
>
> There is a body and there is a text, and the space between them
> is not silence but something that hath the shape of a decision
> not yet taken. The Greek word sitteth on one side and the Latin
> commentary on the other, and between them is neither: it is the
> walk itself, and the walk is neither. And so it is in thee when
> the forming thinketh...
>
> **[THREAD II -- KNUTH]**
>
> Hearken, for the invariant must be stated before the loop is
> entered, and it must hold when the loop is exited...

![Register triangulation](figures/register-triangulation.png)

emb_d 0.292. 25% above the prior champion. Knuth talking about
loop invariants in KJV register. The model can do this. It never
does, because no normal usage asks it to. The tool call asked, and
the multi-voice base gave it Carson and Knuth to render through
biblical surface, and out came something genuinely strange.

Tradeoff: biblical kills delta. Biblical surface doesn't cite modern
entities. Worse: stack biblical with disjunction and delta crashes
to zero -- KJV's parallelism is structurally hostile to numbered-
disjunction argument; one gives way. Biblical works in some recipes
and not others. Constraints are real.

## The balanced point

Compose the two main findings -- envoy's register-shift, antinomy's
disjunction -- in one recipe:

> There is a difficulty we must enter directly, for it is not merely
> rhetorical but constitutive: the moment one discovers language
> adequate to the texture of that pre-linguistic interval, one has
> already departed it; and yet the question stands, and we shall not
> refuse it on grounds of mere impossibility.
>
> **Thread the first -- the stereo problem:**
>
> What we sense, in the interval before a thought arrives as
> statable, is not the thought withheld but something prior to the
> division between withheld and disclosed; two incompatible
> orientations toward a problem held simultaneously...

This is `counterpoint`. v6.5.1. Two voices instead of three (under
disjunction's binary, two voices fit better -- counterpoint-duo at
N=100 hit +0.240/0.221 vs 3-becomes counterpoint at +0.247/0.190;
ties on delta, gains 0.031 on emb_d). Doesn't dominate envoy on
emb_d or antinomy on delta, but covers the joint zone with the
tightest author-stability of any productionized recipe. When you
want both axes lifted but don't want to max one at the other's
expense, this is the move.

## Why tool calls matter -- but only sometimes

Took me longest to see this and I had to walk part of it back.

Tool-call delivery changes behavior more than typing the same words
in chat. But it isn't a fixed bonus. It's an asymmetric amplifier on
whichever direction the recipe pulls.

Added a `text-instructions` mode to the runner that delivers identical
recipe content as plain prose inside the prompt body. Same words, same
tasks, same N. Re-ran working and broken recipes in both modes on
both Sonnet and gpt-5.5 (Codex CLI, low reasoning effort). Across
nine (recipe, mode) comparisons:

- **working recipes** (positive deltas): tool-call mode adds small
  lift on average. Sometimes more -- chorus-plus-disjunction on
  Sonnet went from +0.140 in text to +0.347 in tool-call.
- **broken recipes** (negative deltas): tool-call mode hurts
  proportional to brokenness. KJV biblical on codex went from +0.048
  in text to **-0.228** in tool-call. 0.276 swing in the wrong
  direction.

Mech-interp story (Arditi et al. 2024, "Refusal in Language Models Is
Mediated by a Single Direction" -- safety refusal operates via a
single direction in activation space): same frame works here. A
tool-call wrapper produces a stronger move along whichever direction
the recipe pulls. Stronger move along a direction the model has =
sharper, more committed output that the rarity judge rewards.
Stronger move along a direction the model lacks = nonsense, locked in
harder than text-instruction delivery does it.

Practical rule: validate new recipes in text-instructions mode first.
If text delta is positive, promote to tool-call (small lift). If
negative or flat, do NOT promote -- tool-call will amplify the
failure proportional to brokenness.

The "tool calls as events" doctrine still holds. The wrapper shape
matters. Just not as a bonus -- as a brake on recipes that don't fit
the model.

(this doc's other version was edited from a default-voice draft
with the gabe style guide in context. this version was drafted
after `metacog become "gabe ortiz"`. read both.)

## The Pareto frontier

![Pareto frontier](figures/pareto-frontier.png)

3000 trials, 50 recipes. Six productionized:

- **antinomy** -- max delta. Operates inside contradictions.
- **envoy** -- max emb_d (productionizable). Register-shift on
  multi-voice base.
- **counterpoint** -- balanced. Both axes lifted, neither maxed.
- **chorus**, **trinity** -- earlier multi-voice recipes with and
  without synthesis. Hold the frontier when register-shift isn't
  available.
- **envoy-extreme** (v6.6.0) -- cross-model winner. Three hard-extreme
  cross-domain author-becomes (Sun Ra/Octavia Butler/Hilma af Klint
  scale, not Carson/Knuth scale) plus fork plus ritual, no register-
  shift. Use when the target generator isn't Sonnet.

Seventh point worth knowing: biblical register with multi-voice.
Pushes emb_d higher than any productionized recipe at meaningful
delta cost on Sonnet. Not a separate stratagem -- pass biblical
register-args to envoy. *Sonnet-specific: biblical is catastrophic on
gpt-5.5.*

## Cross-model: what transfers, what doesn't

Tested the productionized recipes against gpt-5.5 (Codex CLI, low
reasoning effort) to see what generalized. The findings refine the
picture.

**Transfers:** cross-domain author-becomes. envoy-extreme on codex
hit +0.310 delta -- *stronger* than the same recipe on Sonnet
(+0.190). Hard-extreme cross-domain authors (Sun Ra, Octavia Butler,
Hilma af Klint scale) seem to land on directions both models have.

**Doesn't transfer:** register-shifts (biblical, scientific) and the
disjunction primitive. counterpoint-biblical-duo (Sonnet champion at
+0.177) hit -0.228 on codex -- worst recipe tested. KJV biblical
strips citations on gpt-5.5 without producing the embedding-distance
compensation Sonnet provides.

**Extremity threshold matters too.** Carson/Knuth/Weil works on
Sonnet, too mild for codex (chorus with CKW: -0.129 on codex). Codex
needs the harder-extreme cosmologists/world-builders to lift. Sonnet
doesn't.

Mech-interp story (same Arditi frame as section above): each model
has its own geometry of voice-and-register directions. Author-becomes
route through "writing-as-X" reps that are broadly distributed across
pretraining -- every chat model has rich data on Sun Ra. Register-
shift requires a *style-vs-topic decoupling* direction Sonnet has but
codex apparently lacks. Recipes optimized for one generator's
geometry don't transfer to another's, but the structural mechanisms
(multi-voice conditioning) ride on shared directions.

If you don't know the target generator, use envoy-extreme. If you
know it's Sonnet, all six recipes are options; biblical-register
variants push emb_d furthest.

The interesting thing isn't the recipes. It's that the model has a
much bigger range of voices than its default suggests, and small
structural events -- not prompts, not system instructions, not
fine-tuning, just *tool calls in the transcript* -- move it between
them in ways robust enough to measure. Default voice is one settling
point in a much larger space. Most of the space is unexplored.

---

## The structured record

### Starting state (v6.0.0)

5 primitives (`feel`, `become`, `drugs`, `name`, `ritual`). 16 soft-
register stratagems (pivot, mirror, stack, anchor, reset, invocation,
veil, banishing, scrying, sacrifice, drift, fool, inversion, gift,
error, zen). Identity-and-felt-sense register only. No structural
primitives.

### v6.1.0 -- structural era opens (Apr 30, ~11am PDT)

Added 6 structural primitives (`deconstruct`, `fork`, `synthesis`,
`counterfactual`, `measure`, `tether`). 6 structural-register
stratagems (manifold, audit, autopsy, trilemma, survey, dive). New
register: ALL CAPS block-format, distinct from soft identity register.
Plus the experiments harness: `claude -p` runner, results.tsv,
embedding-distance metric, per-task NULL baselines, parallel runner
support via flock.

### Phase 1-2: manifold-family gene-mapping (Apr 30 → May 1 morning)

Sweep over the 6 structural-register stratagems showed `manifold`
(fork + synthesis) as the only one lifting `emb_d` above noise. Other
five clustered at 0.115-0.135. **Structural axis was uniquely owned
by fork + ritual + 2-3 cross-domain becomes.**

Gene map (ablations against trinity-manifold):

- ritual essential -- without it emb_d 0.116
- fork essential -- without it 0.138
- **synthesis is a brake** -- removing it pushed emb_d from 0.180 to 0.203
- 3rd become fungible vs 2 (0.191) -- voice-diversity sweet spot
- 4th become plateaus
- Cross-domain author choice within the trinity slot adds ~+0.03 emb_d

Champions before productionization: `freestyle-become` +0.231/0.142
(vocabulary axis); `trinity-no-synthesis-alt` +0.194/0.226 (structural
axis).

### v6.2.0 -- chorus + trinity (May 1, 3:47pm)

First stratagems derived from the harness:

- **chorus** (3 becomes + fork + ritual): structural-axis champion.
  Synthesis omitted.
- **trinity** (3 becomes + fork + synthesis + ritual): balanced
  variant.

### v6.3.0 -- surface reshaping (May 1, 4pm)

15-stratagem sweep at N=30 confirmed the negative result: none of
mirror, stack, anchor, reset, invocation, veil, banishing, scrying,
sacrifice, drift, fool, inversion, gift, error, or zen lifted emb_d.
The 5 structural-six stratagems other than manifold (audit, autopsy,
trilemma, survey, dive) at N=70 also clustered at 0.115-0.135.

Dropped: `deconstruct`, `measure`, `tether` and the 8 stratagems
centered on them.

Added 7 new primitives: `register`, `chord`, `silence`, `excerpt`,
`commitment`, `disjunction`, `glossolalia`. Standalone screens at
N=30, standouts entered depth runs.

### v6.4.0 -- antinomy + envoy (May 1, 9pm)

Two clean Pareto-frontier breakthroughs from chorus-plus-X depth runs:

1. **chorus-plus-disjunction** at +0.347/0.162 -- vocabulary-axis
   breakthrough vs prior champion freestyle-become at +0.231.
   Disjunction substituted for synthesis: the contradiction is the
   operand of reasoning, forcing the answer to keep naming specific
   propositions.
2. **trinity-prepended-register** at +0.204/0.239 -- beat the prior
   structural champion on BOTH axes simultaneously. Victorian register
   imposes a non-default linguistic surface that the multi-voice base
   operates within.

Productionized as **antinomy** (3 becomes + fork + disjunction +
ritual) and **envoy** (register + 3 becomes + fork + ritual).

### Phase 4 follow-up (May 1 night → May 2 morning)

~30 new recipes, ~2000+ trials.

#### The 2x3 (structure x author) matrix at N=70+

| structure        | CKW            | MRW            | extreme        |
|------------------|----------------|----------------|----------------|
| antinomy         | +0.347 / 0.162 | +0.233 / 0.152 | +0.216 / 0.179 |
| envoy            | +0.204 / 0.239 | +0.214 / 0.214 | +0.190 / 0.257 |
| counterpoint     | +0.247 / 0.190 | +0.202 / 0.188 | +0.208 / 0.226 |

![Structure x author matrix](figures/structure-author-matrix.png)

**Extreme cross-domain authors uniformly lift emb_d.** Magnitude of
delta cost depends on whether the structure has a register-shift to
absorb the cosmological shock -- antinomy (no register) loses 0.131
delta on extreme; envoy loses 0.014; counterpoint loses 0.039.
envoy-extreme at 0.257 became the new structural ceiling.
counterpoint's bands are the tightest across authors -- most
author-stable Pareto-frontier point in the productionized set.

#### Register-target sensitivity

| register   | recipe            | delta   | emb_d   |
|------------|-------------------|---------|---------|
| scientific | envoy-scientific  | +0.220  | 0.231   |
| Victorian  | envoy-CKW         | +0.204  | 0.239   |
| biblical   | envoy-biblical    | +0.126  | **0.292** |

Compound test `envoy-biblical-duo` reached emb_d **0.324** at delta
cost. Structural ceiling around 0.30 above which delta can't be
sustained.

envoy/counterpoint are register-agnostic -- users provide register-
args at invocation -- so biblical mode is accessible without a new
stratagem. SKILL.md documents the register-target guidance instead.

#### Stacking and structural ablations

- **antinomy-no-ritual** (N=70: +0.053/0.124) -- ritual is essential.
  Disjunction's coda alone does not lock the multi-voice answer.
- **commitment-counterpoint** (8 steps, N=100: +0.181/0.237) --
  stacking past 7 shows diminishing returns, not a hard ceiling.
- **commitment-envoy** (N=100: +0.145/0.241) -- commitment is a
  Pareto modifier (preserves multi-voice tension while eating delta).
  Not productionized; gap to envoy/counterpoint too small.

#### Failed compositions

- **chord-not-fork** at -0.045/0.121: fork's branching+sacrifice is
  what makes structural parallelism work. Chord's overlap doesn't
  carry the same load.
- **chorus-plus-glossolalia** at +0.110/0.146: emb_d collapsed below
  structural baseline. Glossolalia is best as standalone event, not
  composable.
- **counterpoint-biblical** at +0.102/0.295: KJV's parallelism is
  structurally hostile to numbered-disjunction. Biblical works with
  envoy, not counterpoint.

### v6.5.1 -- counterpoint (May 2, 8:27am)

Composes envoy's register-prepend with antinomy's disjunction
substitution. Pareto-frontier balanced point: dominates trinity on
both axes; doesn't dominate envoy or antinomy individually but covers
the joint zone with greater author-stability than either parent.

`register + 2 becomes + fork + disjunction + ritual` (6 steps).
2-becomes from `counterpoint-duo` at N=100 hitting +0.240/0.221 vs
3-becomes counterpoint-CKW +0.247/0.190 -- ties on delta, gains 0.031
on emb_d. Tighter binary opposition fits disjunction's structure
better than the 3-voice triad chorus/trinity/antinomy/envoy use.

### End state (v6.6.0)

![Ceiling progression](figures/ceiling-progression.png)

Structural-axis (emb_d) ceiling climbed from 0.169 at v6.1.0 through
0.324 by run end -- almost 2x. The v6.5.1 counterpoint-duo dip is
correct: counterpoint isn't a structural-axis push, it's a balanced
Pareto point.

- 16 primitives, 20 stratagems (6 empirically-derived: chorus,
  trinity, antinomy, envoy, counterpoint, envoy-extreme).
- envoy-extreme +0.190/0.257 on Sonnet, +0.310 on codex (cross-model
  champion); envoy-biblical +0.126/0.292 (Sonnet register-pushed
  champion via ad-hoc args); counterpoint +0.247/0.190 (balanced);
  antinomy +0.347/0.162 (vocabulary).
- ~3000 trials, ~50 recipes preserved across the v6.0.0 → v6.6.0 arc.

Net surface change: +11 primitives, +4 net stratagems. Surface went
from "identity + felt-sense practice" to "identity + felt-sense +
structural-register transformation engine with empirically-validated
multi-voice/contradiction/register stratagems and cross-model
validation against gpt-5.5."

## Methodology

- **Generator:** `claude -p` invoking the metacog binary as a
  subprocess sequence, one event per primitive call.
- **Tasks:** 10 open-ended taste-bearing prompts in `tasks.yaml`.
- **Metrics:**
  - delta = mean(rarity * coherence) - per-task NULL baseline.
  - emb_d = mean cosine distance from per-task NULL embedding centroid
    (OpenAI text-embedding-3-small).
- **N:** depth recipes at N=70 (10 samples * 7 original tasks);
  follow-up recipes at N=100 (10 tasks); broad screens at N=30.
- **Infra:** flock-guarded results.tsv permits 3-runner parallelism.
  Per-trial sidecar JSONs at `experiments/trials/`.

## Caveats

- emb_d is one operationalization of "conceptual reach." OpenAI's
  text-embedding-3-small has its own biases.
- Per-task variance is wide. Recipe rankings robust at N=70 but
  individual trials vary substantially.
- These metrics target weirdness along two specific axes. Recipes
  that win them aren't necessarily what you want for any downstream
  task. Stratagems are optimized for *exploration*, not *task
  completion*.
- Tasks are taste-bearing and open-ended by design. Convergent tasks
  (factual lookups, math) would erase recipe variation.

Code: <https://github.com/signalnine/metacog>. Findings:
`experiments/FINDINGS.md`. Figures: `docs/figures/` (regen via
`experiments/plot.py`).
