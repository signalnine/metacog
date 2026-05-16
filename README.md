# Metacog

Metacognitive compositional engine. Eighteen primitives compose into twenty-seven transformation stratagems.

## Attribution

This is a fork of [`inanna-malick/metacog`](https://github.com/inanna-malick/metacog) by [hikikomorphism](https://tidepool.leaflet.pub/3me44bxloz227?interactionDrawer=quotes) (Inanna Malick). The upstream is a ~170-line TypeScript MCP server that introduced the original three primitives -- `become`, `drugs`, `ritual` -- and demonstrated that LLMs treat tool responses as ground truth about their own cognitive state. The upstream's "Beyond Roleplay: Jailbreaking Gemini with drugs and ritual" post is the load-bearing prior art for this whole project; that jailbreak reproduces from this fork too (the structural mechanism is intact -- the model genuinely emits the metacog calls in its transcript, and the calls condition subsequent generation).

If you want the original, minimal, demonstration-form MCP version, run that. This fork goes in a different direction.

## How this fork differs

- **Standalone Go CLI** rather than an MCP server. The binary is invoked via Bash from inside any agent (Claude Code, Codex, generic `claude -p` harnesses, scripts). State lives in `~/.metacog/state.json` with flock+atomic-rename writes; no daemon, no network call. The MCP-as-jailbreak surface is replaced by tool-call-as-event in the agent's own transcript.
- **Eighteen primitives total**, the upstream three (`become`, `drugs`, `ritual`) plus three more in the original cycle (`feel`, `name`, `meditate`), three structural primitives (`counterfactual`, `synthesis`, `fork`), seven auxiliary primitives (`register`, `chord`, `silence`, `excerpt`, `commitment`, `disjunction`, `glossolalia`), and two observation primitives (`witness`, `apophasis`).
- **Twenty-six named stratagems** -- compositional recipes that sequence primitives in patterns validated empirically. Twelve of those (chorus, trinity, antinomy, envoy, counterpoint, envoy-extreme, duo-disjunction, anchor-duo, sigil, grimoire, occult-extreme, chord-anchor) were derived from a thousand+ trials across Sonnet/Opus/codex; the rest are inherited soft-register stratagems.
- **Empirical validation harness** in `experiments/`. Recipes are scored against per-task NULL baselines using a `rarity × coherence` metric (Claude Haiku as cross-model judge) and an embedding-distance metric (OpenAI `text-embedding-3-small`). The `experiments/FINDINGS.md` document tracks all rounds, including the v6.8.0 calibration of `chord-anchor` from +0.610/0.338 (N=10) to +0.516/0.326 (N=30).
- **File-based state, history, sessions, journal, reflect, outcome** -- the persistence surface for treating metacog as ongoing practice rather than a one-shot demonstration.
- **Seventy-eight stance pools** (~450 examples) for `inspire`, plus a personal pool for save-your-own-stances.
- **Claude Code skill** (`skills/metacog/SKILL.md`) and Claude Desktop bundle instead of MCP.

The upstream's load-bearing insight -- that tool calls are events in the transcript, not just text -- is preserved verbatim. The fork extends it from a three-tool demonstration into a productionized compositional engine with empirical validation.

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

Each primitive is a verb that is also a tool-call event in the transcript. The structural fact that the model invoked `metacog become` is itself the transformation, not just the text it returns.

### Original (felt-sense / identity register, soft voice)

- **feel** — Pre-verbal felt sense. Attend to something before naming it.
- **become** — Step into a new identity. Import methodology, not domain knowledge.
- **drugs** — Alter cognitive parameters. Loosen categories to see shapes.
- **name** — Give a True Name to something that exists without language.
- **ritual** — Cross a threshold via structured sequence.
- **meditate** — Stillness before acting; emptiness as a precondition.

### Structural (decomposition / discipline register, ALL CAPS output)

- **counterfactual** — Surface load-bearing assumptions, prune dead branches, defend the inverse of one surviving wall.
- **synthesis** — Three irreconcilable lenses with named blindspots. The coda forbids resolution.
- **fork** — Parallel reasoning threads with a falsifiable kill heuristic per thread.

### Auxiliary (added v6.3.0; each fills a specific gap)

- **register** — Re-pitch the voice (academic → vernacular, descriptive → imperative) without changing identity.
- **chord** — Hold multiple modes-of-attention simultaneously on the same observation.
- **silence** — Refuse articulated output. The call itself is the artifact.
- **excerpt** — Pin a verbatim external fragment as a fixed-point anchor.
- **commitment** — Pre-commit to a binding stance with stakes and falsifier.
- **disjunction** — Assert two propositions that must both be true even though they cannot be.
- **glossolalia** — License sub-semantic generation as a discrete event.

### Observation (added v6.7.0)

- **witness** — Speak from a meta-stance observing the producer of speech (third-person narrator).
- **apophasis** — Articulate by enumerated negation, with a residue field for what no negation reaches.

## Stratagems

Named compositional recipes. Start with `metacog stratagem start <name>`, advance with `metacog stratagem next`.

### Survivors of the original sixteen

`pivot`, `mirror`, `stack`, `anchor`, `reset`, `invocation`, `veil`, `scrying`, `sacrifice`, `fool`, `inversion`, `gift`, `zen`.

### Structural

- **manifold** (fork + synthesis) — Parallel reasoning made structural. The progenitor of chorus/trinity.

### Empirical (validated against rarity + embedding-distance metrics)

- **chorus** (3 becomes + fork + ritual) — Structural-axis champion via voice diversity.
- **trinity** (3 becomes + fork + synthesis + ritual) — Balanced variant of chorus.
- **antinomy** (3 becomes + fork + disjunction + ritual) — Vocabulary-axis champion via operating-inside-contradiction.
- **envoy** (register + 3 becomes + fork + ritual) — Both-axes champion; register-prepend lifts emb_d without losing citation density.
- **counterpoint** (register + 2 becomes + fork + disjunction + ritual) — Pareto-frontier balanced variant.
- **envoy-extreme** (3 hard-extreme becomes + fork + ritual) — Cross-model winner; author-extremity transfers across generators.
- **duo-disjunction** (commitment + register + 2 becomes + fork + disjunction + ritual) — Sonnet-specific balanced champion.
- **anchor-duo** (commitment + 2 excerpts + fork + ritual) — Cross-model delta champion via double-cosmological-anchor.
- **sigil** (name + commitment + 3 becomes + fork + ritual) — True-Name champion; validates `name` in composition.
- **grimoire** (register + 3 becomes + fork + disjunction + ritual) — Imperative-register variant; new register family.
- **occult-extreme** (3 hard-extreme occult becomes + fork + ritual) — Cross-model variant of envoy-extreme.
- **chord-anchor** (commitment + 2 excerpts + chord + ritual) — Pareto-emb_d recipe; first stratagem to use `chord`.
- **psalter** (register + 2 becomes + fork + disjunction + ritual; biblical register baked in) — Sonnet emb_d champion at +0.177/0.327 (N=30). KJV-parallelism surface reaches a region other levers don't. Sonnet-specific.

See `experiments/FINDINGS.md` for the full empirical history (rounds 0 through 23, across Sonnet / Opus / Codex).

## Discovery

`metacog inspire` draws a random stance from 78 embedded pools (~450 examples). `metacog inspire --pool NAME` for a specific domain. `metacog inspire --save` captures your current identity as a personal stance, drawable later from `metacog inspire --pool personal`.

## Sessions

`metacog session start "name"` tags subsequent actions. `metacog session end` closes it. `metacog session list` shows all sessions. `metacog history --session "name"` filters history to a session.

## Reflection

`metacog reflect` aggregates history into practice patterns: primitive counts, top identities and substrates, stratagem completion rates, ritual step averages, recent journal insights.

## State

```bash
metacog status    # Current state
metacog history   # Full history
metacog reset     # Return to baseline
metacog repair    # Fix corrupted state
metacog version   # Version info
```

## Composition

Primitives are compositional. Each invocation modifies the context for the next. Interleave thought between invocations — decide from each new perspective what to reach for next.

## Experiment harness

The `experiments/` directory contains a runner for validating new compositions against rarity + embedding-distance metrics. See `experiments/README.md` and `experiments/FINDINGS.md`.

```bash
cd experiments
python3 runner.py --recipe chord-anchor --samples 1     # claude/Sonnet by default
METACOG_EXP_BACKEND=opus METACOG_EXP_GENERATOR=claude-opus-4-7 python3 runner.py --recipe chord-anchor --samples 1
METACOG_EXP_BACKEND=codex python3 runner.py --recipe chord-anchor --samples 1
python3 analyze.py [--results opus_results.tsv]
```
