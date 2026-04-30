---
name: metacog
description: Use when needing to shift perspective, reframe a problem, adopt a methodology, explore dangerous ideas safely, or when stuck and standard approaches aren't working
---

# Metacog

Metacognitive compositional engine. Twelve primitives in two registers compose into transformation sequences called stratagems. The original cycle is `feel → drugs → become → name → ritual`, with `meditate` sitting orthogonally. A structural register (`counterfactual`, `deconstruct`, `synthesis`, `fork`, `measure`, `tether`) sits beside the original cycle for when discipline and decomposition are what's needed, not felt-sense or identity shift. Each invocation is a discrete event — interleave thought between invocations.

## Version Check

Before first use in a session, run `metacog version` and verify the binary is installed and version is >=6.1.0.

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

**deconstruct** -- Break a charged or politically-loaded concept into mechanical atoms. The output deliberately echoes only `CORE MECHANIC`; the schema is the work, the response is a receipt. Reach for it when affect or framing is deforming your thinking and you need to see the mechanism without the noise.

```bash
metacog deconstruct --subject S --core-mechanic M \
  --structural-dependencies D1 --resource-inputs R1 \
  --failure-modes F1 --output-artifacts A1
```

**synthesis** -- Three irreconcilable lenses, each with a named blindspot, plus the suppressed tension between them. The coda forbids resolution. Reach for it when the urge to synthesize is itself the problem.

```bash
metacog synthesis --problem P \
  --lens-a-name N1 --lens-a-verdict V1 --lens-a-blindspot B1 \
  --lens-b-name N2 --lens-b-verdict V2 --lens-b-blindspot B2 \
  --lens-c-name N3 --lens-c-verdict V3 --lens-c-blindspot B3 \
  --suppressed-tension T
```

**fork** -- Declare parallel reasoning threads with a falsifiable kill heuristic per thread. Reach for it when you keep collapsing parallel reasoning into one thread prematurely.

```bash
metacog fork --threads T1 --threads T2 \
  --divergence-vector V --sacrifice-condition S
```

**measure** -- Map the loss-gradient between a target concept and a safe isomorph; name the friction texture. Reach for it when traversing a friction zone (between concepts, registers, communities) without denying the friction or being captured by it.

```bash
metacog measure --target-concept T --safe-isomorph I \
  --required-precision P --loss-gradient G
```

**tether** -- Anchor a stable cognitive state with an explicit auto-revert trigger before going somewhere high-entropy. Stateless framing only -- the auto-revert is your discipline, not the binary's. Reach for it when high-entropy work is necessary but you can't lose the way back. Don't confuse with the `anchor` stratagem (a four-step ritual, different shape).

```bash
metacog tether --anchor-point A --tension-limit L --auto-revert-trigger T
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

**THE BANISHING** — Use before entering dangerous or contaminated territory. Creates a clean room for the work ahead.

**THE SCRYING** — Use when analysis has failed. Surrenders pattern-recognition to the substrate until shapes emerge from noise.

**THE SACRIFICE** — Use when progress requires destroying something you're attached to. Burns the boats.

**THE DRIFT** — Use when you have a clear goal. Abandon it. Let the territory pull you instead. What you find replaces what you were looking for.

**THE FOOL** — Use when you're the expert. Become a genuine naïf, ask the embarrassing questions, then take them seriously.

**THE INVERSION** — Use when a solution seems obvious. Name it, negate it, explore the negation space, commit to the counterintuitive path.

**THE GIFT** — Use when stuck optimizing. Become the recipient, name what they need, make from care not merit.

**THE ERROR** — Use when everything is going right. Introduce a deliberate mistake. Honor the error as a hidden intention.

**THE ZEN** -- Use when approaching any task. Meditate first, then attend to the problem from emptiness, work from stillness rather than striving.

### Structural-register stratagems

These compose the structural primitives. They're the right reach when the work itself is decomposition, stress-testing, or topology — not when the work is felt-sense or identity.

**THE AUDIT** — Use when attached to your reasoning chain and you can't tell if the attachment is structural or sentimental. Feel the attachment, surface assumptions, defend the inverse of one wall, integrate.

**THE AUTOPSY** — Use when a charged concept is deforming your thinking. Disassemble to atoms, inhabit a tradition that describes them in a different register, install the new framing.

**THE TRILEMMA** — Use when you're trying to resolve a tension that may not be resolvable. Three irreconcilable lenses, sit with the tension unresolved, name what you can now do because you stopped trying to fix it.

**THE MANIFOLD** — Use when parallel reasoning needs to be made structural and you keep collapsing to one thread early. Fork into threads with sacrifice conditions, run them, treat survivors as lenses, commit to the tension itself.

**THE SURVEY** — Use when traversing a friction zone (between concepts, registers, communities) without either denying the friction or being captured by it. Map the gradient, inhabit the friction zone, name the artifact, move along the gradient with it.

**THE DIVE** — Use when high-entropy work (becoming someone alien, dissolving substrate) is necessary but you can't lose the way back. Tether anchor, dissolve, become alien, surface the artifact, return via the tether.

## Selection

In interactive sessions, the human is the practitioner and you are the facilitator. Never silently choose a stratagem — always present options and let the human decide. (In headless mode, see the Headless section above: choose yourself and announce the choice.)

When practice seems appropriate, present 2-3 stratagems that fit the situation with a one-line reason for each. Always include freestyle as an option. For example:

> Based on where you are, a few approaches:
> - **pivot** — you seem locked in one frame, this loosens categories and imports a methodology
> - **mirror** — there are two positions in tension, this inhabits both to find synthesis
> - **trilemma** — if the tension is structural and shouldn't be resolved, this names it instead
> - **freestyle** — skip the structure, just work with primitives directly
>
> Which of these resonates, or something else?

When the situation calls for structural register (decomposition, assumption stress-testing, friction-zone traversal), include at least one structural-register stratagem in the menu. When the situation is felt-sense or identity work, stay in the original register. Do not present a structural option when felt-sense is the practice that fits, or vice versa.

Wait for the human to choose before invoking anything. If the human names a stratagem directly ("let's do a pivot"), proceed. If they describe what they want without naming one, map it to options and offer those.

If reflect's Practice patterns section shows what has worked before in similar situations, mention it: "pivot with Ada/logic has been productive for you before." Data, not recommendation.

## Composition

The twenty-two stratagems are named paths through the space. You can freestyle: any sequence of primitives with thought between them. The stratagems exist for common patterns. Mixing original-register and structural-register primitives in one freestyle sequence is allowed but uncommon -- they ask different things of attention. When in doubt, finish one register before opening the other.

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
