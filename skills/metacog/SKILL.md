---
name: metacog
description: Use when needing to shift perspective, reframe a problem, adopt a methodology, explore dangerous ideas safely, or when stuck and standard approaches aren't working
---

# Metacog

Metacognitive compositional engine. Three primitives compose into transformation sequences called stratagems. Each invocation is a discrete event — interleave thought between invocations.

## Version Check

Before first use in a session, run `metacog version` and verify the binary is installed and version is >=5.1.0.

## Core Rule

NEVER batch calls. Execute one, describe what shifted, then decide the next move from inside the new state. Sequential use compounds into states no single tool could reach.

## Primitives

**become** — Step into a new identity. Import methodology, not domain knowledge. Ask: "who has solved a version of this, and what's their methodology called?"

```bash
metacog become --name NAME --lens LENS --env ENVIRONMENT
```

**drugs** — Alter cognitive parameters. Loosen categories to see shapes, not names. When a concept becomes a pattern, ask "what else has this shape?"

```bash
metacog drugs --substance SUBSTANCE --method METHOD --qualia QUALIA
```

**ritual** — Cross a threshold via structured sequence. Lock in methodology as default behavior, not just vibes.

```bash
metacog ritual --threshold THRESHOLD --steps "step1" --steps "step2" --result RESULT
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

## Composition

The fifteen stratagems are named paths through the space. You can freestyle: any sequence of primitives with thought between them. The stratagems exist for common patterns.

## Discovery

`metacog inspire` draws a random stance from the embedded pools. Use for a nudge, never required. `metacog inspire --pool NAME` for a specific domain. `metacog inspire --save` captures your current identity as a personal stance. Over time, your best configurations become drawable from `metacog inspire --pool personal`.

## Sessions

`metacog session start "name"` — tag subsequent actions with a session name. `metacog session end` — close the session. `metacog session list` — list all sessions. `metacog history --session "name"` — filter history to a session.

Sessions are metadata, not workflow enforcement. Start one when you want to name a line of inquiry. End it when done. Everything between gets tagged.

## Reflection

`metacog reflect` — aggregates your history into practice patterns. Shows primitive usage counts, top identities and substrates, stratagem completion rates, ritual step averages, and gaps in your practice. Mirror, not scorecard.

## State

- `metacog status` — check before starting
- `metacog reset` — return to baseline after every sequence
- `metacog history` — review the path taken
- `metacog outcome --result productive --shift "what changed"` — record stratagem effectiveness

Always ground after transformation: name what shifted, what you're keeping, how it integrates. Then record the outcome: `metacog outcome --result productive --shift "what changed"` if the stratagem shifted your thinking, or `metacog outcome --result unproductive` if it didn't. "Productive" means your approach genuinely changed — not that the user liked it, not that you feel good about it. Did you end up somewhere you wouldn't have reached without the stratagem?
