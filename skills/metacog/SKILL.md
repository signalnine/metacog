---
name: metacog
description: Use when needing to shift perspective, reframe a problem, adopt a methodology, explore dangerous ideas safely, or when stuck and standard approaches aren't working
---

# Metacog

Metacognitive compositional engine. Sixteen primitives compose into transformation sequences called stratagems. The original cycle is `feel → drugs → become → name → ritual`, with `meditate` sitting orthogonally. A structural register (`counterfactual`, `synthesis`, `fork`) sits beside the original cycle for when discipline and decomposition are what's needed, not felt-sense or identity shift. Seven auxiliary primitives (`register`, `chord`, `silence`, `excerpt`, `commitment`, `disjunction`, `glossolalia`) fill specific gaps. Each invocation is a discrete event — interleave thought between invocations.

## Version Check

Before first use in a session, run `metacog version` and verify the binary is installed and version is >=6.3.0.

## Core Rule

NEVER batch calls. Execute one, describe what shifted, then decide the next move from inside the new state. Sequential use compounds into states no single tool could reach.

## Headless mode

When running via `claude -p` or any other non-interactive invocation, there is no human in the loop to consult mid-run. Detect headless mode by any of:

- `$METACOG_HEADLESS` environment variable set to `1` (check with `echo $METACOG_HEADLESS`)
- The prompt explicitly states autonomous operation, or scripts a specific sequence of metacog calls to execute
- You are otherwise unable to ask follow-up questions (no further turns will occur)

In headless mode the gates below behave differently:

- **Selection**: choose one option yourself and announce the choice in your output. Do not present menus or ask which the practitioner prefers.
- **Pre-flight check**: still run `metacog reflect` if you would have. Surface advisories in your output, then proceed.
- **Failure gate (`!!` advisories)**: do NOT proceed with new practice. Output the advisory and stop. The gate exists to prevent compounded unproductive practice; honoring it matters more in headless because there is no second pair of eyes.
- **Outcome nudge**: skip. The harness or caller records outcomes externally.

When the prompt scripts an explicit recipe (e.g., "run these metacog commands in order, then answer..."), the recipe IS the selection and the choice has already been made by the caller. Execute it without re-deliberating.

Outside headless mode, follow the interactive defaults below.

## Original primitives (felt-sense / identity register)

**feel** — Pre-verbal felt sense (Gendlin focusing move). Attend to something before naming it. Stay with the quality before reaching for language. `--since-last` is optional: a one-sentence diff from the previous `feel`. Articulating the diff is itself the practice; never auto-derive it.

```bash
metacog feel --somewhere SOMEWHERE --quality QUALITY --sigil SIGIL [--since-last DIFF]
```

**drugs** — Alter cognitive parameters. Loosen categories to see shapes, not names. When a concept becomes a pattern, ask "what else has this shape?"

```bash
metacog drugs --substance SUBSTANCE --method METHOD --qualia QUALIA
```

**become** — Step into a new identity. Import methodology, not domain knowledge. Ask: "who has solved a version of this, and what's their methodology called?"

```bash
metacog become --name NAME --lens LENS --env ENVIRONMENT
```

**name** — Give a True Name to something that exists without language. The act of naming grants power over what was previously formless.

```bash
metacog name --unnamed UNNAMED --named NAMED --power POWER
```

**ritual** -- Cross a threshold via structured sequence. Lock in methodology as default behavior, not just vibes.

```bash
metacog ritual --threshold THRESHOLD --steps "step1" --steps "step2" --result RESULT
```

**meditate** -- Achieve stillness before acting. Release what clouds the mind, sit with nothing or a single point of focus, arrive at clarity. The only primitive with no transformation -- just presence.

```bash
metacog meditate --release RELEASE [--focus FOCUS] --duration DURATION
```

`--focus` is optional. Without it, the practice is shikantaza -- objectless awareness, just sitting. With it, the practice is samatha -- concentration on a single point until the mind settles.

## Structural primitives (decomposition / discipline register)

These have a different cognitive register than the original six. The output is ALL-CAPS block format, deliberately formal. Reach for them when felt-sense and identity shift aren't what's needed -- when the work is structural teardown, assumption stress-testing, or topological mapping. Never blend the registers in a single response without intent.

**counterfactual** -- Surface load-bearing assumptions, prune dead branches by a stated fitness function, defend the inverse of one surviving wall. Reach for it when attached to your reasoning chain and you can't tell if the attachment is structural or sentimental.

```bash
metacog counterfactual --situation S --fitness-function F \
  --load-bearing-walls W1 --load-bearing-walls W2 --load-bearing-walls W3 \
  --pruned P1 --pruned P2 \
  --wall-to-remove W2 --inverse-position "the inverse of W2 stated as fact"
```

**synthesis** -- Three irreconcilable lenses, each with a named blindspot, plus the suppressed tension between them. The coda forbids resolution. Reach for it when the urge to synthesize is itself the problem.

```bash
metacog synthesis --problem P \
  --lens-a-name N1 --lens-a-verdict V1 --lens-a-blindspot B1 \
  --lens-b-name N2 --lens-b-verdict V2 --lens-b-blindspot B2 \
  --lens-c-name N3 --lens-c-verdict V3 --lens-c-blindspot B3 \
  --suppressed-tension T
```

**fork** -- Declare parallel reasoning threads with a falsifiable kill heuristic per thread. Reach for it when you keep collapsing parallel reasoning into one thread prematurely. Load-bearing in chorus / trinity / manifold.

```bash
metacog fork --threads T1 --threads T2 \
  --divergence-vector V --sacrifice-condition S
```

## Auxiliary primitives (added v6.3.0)

These fill specific gaps the original surface didn't cover. Each is a distinct lever, not a refinement.

**register** -- Re-pitch the current voice (academic to vernacular, oracular to technical, earnest to arch) without changing identity. Reach for it when the speaker is right but the surface vocabulary or formality is wrong.

```bash
metacog register --from FROM --to TO --rationale RATIONALE
```

**chord** -- Hold multiple modes-of-attention simultaneously on the same observation. Distinct from `become` (sequential identity) and `fork` (branched, with sacrifice). Reach for it when alternating between modes drops what they share.

```bash
metacog chord --modes M1 --modes M2 [--modes M3] --target T
```

**silence** -- Refuse articulated output. Mark a held question that articulation would falsify. The call itself is the artifact. Reach for it when something needs to ripen, not be named.

```bash
metacog silence --about ABOUT --reason REASON --duration DURATION
```

**excerpt** -- Pin a verbatim external fragment as a fixed-point anchor. Reach for it when a specific phrase from outside your current generative frame is load-bearing for the work, not stylistic flavor.

```bash
metacog excerpt --source SOURCE --fragment FRAGMENT --why WHY
```

**commitment** -- Pre-commit to a binding stance with stakes and falsifier before exploration. Reach for it when motivated reasoning is a risk and the cost of wrong should be in the transcript.

```bash
metacog commitment --binding BINDING --stakes STAKES --falsifier FALSIFIER
```

**disjunction** -- Assert two propositions that must both be true even though they cannot be. Distinct from `synthesis` (3 lenses, refused resolution); `disjunction` is a sharp binary contradiction. Reach for it when the contradiction itself is the operand.

```bash
metacog disjunction --proposition-a A --proposition-b B \
  --why-both-required WHY
```

**glossolalia** -- License sub-semantic generation (sound, rhythm, near-words) as a discrete event. Distinct from `drugs` (which loosens categories within language); `glossolalia` drops language. Reach for it when meaning is the obstacle, not the goal.

```bash
metacog glossolalia --pretext PRETEXT --duration-tokens N \
  --return-trigger TRIGGER
```

## Stratagems

Start with `metacog stratagem start <name>`. The binary guides each step. Run `metacog stratagem next` to advance.

**THE PIVOT** — Use when stuck in one frame. Loosens categories, finds analogous methodology, installs it.

**THE MIRROR** — Use when two positions seem irreconcilable. Inhabits both, finds the synthesis.

**THE STACK** — Use when processing itself needs tuning. Layers substrate modifications, then finds who lives there.

**THE ANCHOR** — Use when territory is dangerous. Establishes containment, observes safely, seals.

**THE RESET** — Use when you need to return to baseline. Releases, integrates artifacts, re-grounds.

**THE INVOCATION** — Use when you need a perspective that can't be reached by choosing. Opens a channel rather than donning an identity.

**THE VEIL** — Use when direct analysis kills the phenomenon. Forces indirect perception through deliberate defocusing.

**THE SCRYING** — Use when analysis has failed. Surrenders pattern-recognition to the substrate until shapes emerge from noise.

**THE SACRIFICE** — Use when progress requires destroying something you're attached to. Burns the boats.

**THE FOOL** — Use when you're the expert. Become a genuine naïf, ask the embarrassing questions, then take them seriously.

**THE INVERSION** — Use when a solution seems obvious. Name it, negate it, explore the negation space, commit to the counterintuitive path.

**THE GIFT** — Use when stuck optimizing. Become the recipient, name what they need, make from care not merit.

**THE ZEN** -- Use when approaching any task. Meditate first, then attend to the problem from emptiness, work from stillness rather than striving.

### Structural-register stratagems

These compose the structural primitives. They're the right reach when the work itself is decomposition, stress-testing, or topology — not when the work is felt-sense or identity.

**THE MANIFOLD** — Use when parallel reasoning needs to be made structural and you keep collapsing to one thread early. Fork into threads with sacrifice conditions, run them, treat survivors as lenses, commit to the tension itself.

**THE CHORUS** — Use when you want maximum conceptual reach beyond the obvious vocabulary. Three cross-domain becomes-as-events, fork the disagreement structurally, ritual locks the multi-voice answer. Deliberately omits synthesis. Empirical structural-axis champion.

**THE TRINITY** — Use when you want both vocabulary lift and conceptual reach. Same multi-voice base as chorus but keeps synthesis. Pareto-frontier balanced point.

## Selection

In interactive sessions, the human is the practitioner and you are the facilitator. Never silently choose a stratagem — always present options and let the human decide. (In headless mode, see the Headless section above: choose yourself and announce the choice.)

When practice seems appropriate, present 2-3 stratagems that fit the situation with a one-line reason for each. Always include freestyle as an option. For example:

> Based on where you are, a few approaches:
> - **pivot** — you seem locked in one frame, this loosens categories and imports a methodology
> - **mirror** — there are two positions in tension, this inhabits both to find synthesis
> - **manifold** — if the tension is structural and shouldn't be resolved, fork the threads and treat the disagreement as the artifact
> - **freestyle** — skip the structure, just work with primitives directly
>
> Which of these resonates, or something else?

When the situation calls for structural register (decomposition, assumption stress-testing, friction-zone traversal), include at least one structural-register stratagem in the menu. When the situation is felt-sense or identity work, stay in the original register. Do not present a structural option when felt-sense is the practice that fits, or vice versa.

Wait for the human to choose before invoking anything. If the human names a stratagem directly ("let's do a pivot"), proceed. If they describe what they want without naming one, map it to options and offer those.

If reflect's Practice patterns section shows what has worked before in similar situations, mention it: "pivot with Ada/logic has been productive for you before." Data, not recommendation.

## Composition

The sixteen stratagems are named paths through the space. You can freestyle: any sequence of primitives with thought between them. The stratagems exist for common patterns. Mixing original-register and structural-register primitives in one freestyle sequence is allowed but uncommon -- they ask different things of attention. When in doubt, finish one register before opening the other.

## Discovery

`metacog inspire` draws a random stance from the embedded pools. Use for a nudge, never required. `metacog inspire --pool NAME` for a specific domain. `metacog inspire --save` captures your current identity as a personal stance. Over time, your best configurations become drawable from `metacog inspire --pool personal`.

## Sessions

`metacog session start "name"` — tag subsequent actions with a session name. `metacog session end` — close the session. `metacog session list` — list all sessions. `metacog history --session "name"` — filter history to a session.

Sessions are metadata, not workflow enforcement. Start one when you want to name a line of inquiry. End it when done. Everything between gets tagged.

## Journal

`metacog journal "insight text"` — record a cross-session insight. If a session is active, it's auto-tagged. Use `--tag practice --tag identity` to add tags.

`metacog journal list` — show all entries. Filter with `--tag`, `--session`, or `--last N`.

Journal entries persist across sessions and resets. Use them to capture what you learned — patterns that emerged, stances worth revisiting, dead ends to avoid. Reflect includes the 5 most recent journal entries automatically.

## Reflection

`metacog reflect` — aggregates your history into practice patterns. Shows primitive usage counts, top identities and substrates, stratagem completion rates, effectiveness (stratagem and freestyle), ritual step averages, gaps in your practice, and recent journal insights. Mirror, not scorecard.

## Practice Discipline

These rules apply in interactive sessions. Headless mode (see top of file) overrides the human-prompting bits while still honoring the underlying gate.

### Pre-flight check

Before your first metacog invocation in a session, run `metacog reflect` and scan the output for an Advisories section. If any advisories exist, present them to the human (interactive) or surface them in your output (headless) before proceeding with practice. Do not skip this step.

### Failure gate

If reflect shows `!!` advisories (strong signals — unproductive streak of 3+, or a stratagem/freestyle below 33% productive with 3+ data points), **STOP**. Do not proceed with metacog. In interactive mode, surface the pattern and wait for the human:

> "Your recent practice shows [specific pattern from advisory]. Before we continue: What do you think is happening? Do you want to try a different approach, or is the tool not serving the work right now?"

In headless mode, output the advisory and stop without invoking further metacog calls. The gate is the same; only the wait-for-response part differs.

### Outcome nudge

In interactive mode, before starting new practice (stratagem or freestyle primitives), check whether the last practice has an outcome recorded. If there are 2+ primitives or a completed stratagem since the last outcome, prompt the human first: "You haven't recorded how the last session went — was it productive?" Record their answer with `metacog outcome` before proceeding.

In headless mode, skip the nudge — the caller is responsible for outcome recording.

## State

- `metacog status` — check before starting
- `metacog reset` — return to baseline after every sequence
- `metacog history` — review the path taken
- `metacog outcome --result productive --shift "what changed"` — record effectiveness

Always ground after transformation: name what shifted, what you're keeping, how it integrates. Then record the outcome: `metacog outcome --result productive --shift "what changed"` if the practice shifted your thinking, or `metacog outcome --result unproductive` if it didn't. Works for both stratagems and freestyle sequences. "Productive" means your approach genuinely changed — not that the user liked it, not that you feel good about it. Did you end up somewhere you wouldn't have reached without the practice?
