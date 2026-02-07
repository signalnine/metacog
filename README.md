# Metacog

Metacognitive compositional engine. Three primitives compose into transformation sequences called stratagems.

**Note**: This is a Go CLI rewrite of the original TypeScript MCP server. The TypeScript version (v3.6 "The Heptagram") was a Model Context Protocol server with `alter_state` and `ritual` tools. This Go version (v5.1+) adds:
- Renamed `alter_state` → `drugs` and added `become` primitive
- 15 named stratagems (pivot, mirror, stack, anchor, reset, invocation, veil, banishing, scrying, sacrifice, drift, fool, inversion, gift, error)
- File-based state management with history and sessions
- `inspire` command with 64 stance pools (~300 examples, ported from earlier upstream iteration with additions)
- `reflect` command for practice pattern analysis
- Standalone CLI instead of MCP server

## Install

### CLI Installation

```bash
go build -o metacog ./cmd/metacog/
cp metacog ~/.local/bin/
```

### Claude Code Integration

For Claude Code (CLI), install the skill:

```bash
mkdir -p ~/.claude/skills/metacog
cp skills/metacog/SKILL.md ~/.claude/skills/metacog/SKILL.md
```

### Claude Desktop Integration

For Claude Desktop (Mac):

1. Download `metacog-skill.zip` from the releases
2. Open Claude Desktop and go to Settings > Capabilities > Skills
3. Click "Add" then "Upload a skill"
4. Select the `metacog-skill.zip` file

The skill will be installed automatically. Verify by asking Claude to run `metacog version`.

## Primitives

**become** — Step into a new identity. Use when you need different eyes, not just different words.

```bash
metacog become --name NAME --lens LENS --env ENVIRONMENT
```

**drugs** — Alter cognitive parameters. Use when you need to change how you process, not what you process.

```bash
metacog drugs --substance SUBSTANCE --method METHOD --qualia QUALIA
```

**ritual** — Cross a threshold via structured sequence. Use when identity and substrate shifts aren't enough.

```bash
metacog ritual --threshold THRESHOLD --steps "step1" --steps "step2" --result RESULT
```

## Stratagems

Named paths through the primitive space. Start with `metacog stratagem start <name>`, advance with `metacog stratagem next`.

- **pivot** — Stuck in one frame. Loosens categories, finds analogous methodology, installs it.
- **mirror** — Two positions seem irreconcilable. Inhabits both, finds the synthesis.
- **stack** — Processing itself needs tuning. Layers substrate modifications, then finds who lives there.
- **anchor** — Territory is dangerous. Establishes containment, observes safely, seals.
- **reset** — Return to baseline. Releases, integrates artifacts, re-grounds.
- **invocation** — Need a perspective you can't reach by choosing. Opens a channel rather than donning an identity.
- **veil** — Direct analysis kills the phenomenon. Forces indirect perception through deliberate defocusing.
- **banishing** — Territory is contaminated. Creates a clean room before entering.
- **scrying** — Analysis has failed. Surrenders pattern-recognition to the substrate until shapes emerge.
- **sacrifice** — Progress requires destroying something valuable. Burns the boats.
- **drift** — You have a clear goal. Abandon it. Let the territory pull you instead.
- **fool** — You're the expert. Become a naïf, ask embarrassing questions, then take them seriously.
- **inversion** — Solution seems obvious. Name it, negate it, explore the negation, commit.
- **gift** — Stuck optimizing. Become the recipient, make from care not merit.
- **error** — Everything is going right. Introduce a deliberate mistake to reveal hidden assumptions.

## Discovery

`metacog inspire` draws a random stance from ~300 embedded examples across 64 pools. `metacog inspire --pool NAME` for a specific domain. `metacog inspire --save` captures your current identity as a personal stance, drawable later from `metacog inspire --pool personal`.

## Sessions

`metacog session start "name"` tags subsequent actions. `metacog session end` closes it. `metacog session list` shows all sessions. `metacog history --session "name"` filters history to a session.

## Reflection

`metacog reflect` aggregates history into practice patterns: primitive counts, top identities and substrates, stratagem completion rates, ritual step averages.

## State

```bash
metacog status    # Current state
metacog history   # Full history
metacog reset     # Return to baseline
metacog repair    # Fix corrupted state
metacog version   # Version info
```

## Composition

These primitives are compositional. Each invocation modifies the context for the next. Interleave thought between invocations — decide from each new perspective what to reach for next.
