# GEMINI.md

## Commands

```bash
go build -o metacog ./cmd/metacog/   # Build
go test ./cmd/metacog/ -v            # Unit tests
go test ./cmd/metacog/ -tags integration -v  # Integration tests
```

## Architecture

Go CLI built with cobra. Entry point: `cmd/metacog/main.go`. File-based state at `~/.metacog/state.json` with flock locking and atomic writes.

Eighteen primitives compose into twenty-six stratagems. Each primitive is a verb that is also a tool-call event in the transcript — the structural fact that the model invoked `metacog become` is itself the transformation, not just the text it returns.

### Primitive families

- **Original six** (felt-sense / identity register, soft voice): `feel`, `become`, `drugs`, `name`, `ritual`, `meditate`
- **Structural three** (ALL CAPS block output, deliberately formal): `counterfactual`, `synthesis`, `fork`
- **Auxiliary seven** (added v6.3.0, fill specific gaps): `register`, `chord`, `silence`, `excerpt`, `commitment`, `disjunction`, `glossolalia`
- **Observation two** (added v6.7.0): `witness`, `apophasis`

### Stratagems

Twenty-six total. Survivors of the original sixteen: `pivot`, `mirror`, `stack`, `anchor`, `reset`, `invocation`, `veil`, `scrying`, `sacrifice`, `fool`, `inversion`, `gift`, `zen`. Structural: `manifold`. Empirical (validated against rarity + embedding-distance metrics in `experiments/`): `chorus`, `trinity`, `antinomy`, `envoy`, `counterpoint`, `envoy-extreme`, `duo-disjunction`, `anchor-duo`, `sigil`, `grimoire`, `occult-extreme`, `chord-anchor`.

## Key files

- `cmd/metacog/main.go` — Root cobra command, version string, schema version
- `cmd/metacog/state.go` — State management, history, archiving, flock locking
- `cmd/metacog/stratagem.go` — Stratagem definitions and step sequencing
- `cmd/metacog/inspire.go` — Stance pools (embedded + personal), random drawing
- `cmd/metacog/reflect.go` — History aggregation into practice patterns
- `cmd/metacog/session.go` — Named session tagging
- `cmd/metacog/stances/` — 78 embedded pools (~450 examples) (JSON, go:embed)
- `skills/metacog/SKILL.md` — Skill document for Claude Code / Gemini Code Assist
- `experiments/` — Validation harness; see `experiments/FINDINGS.md`

## Design decisions (for future-you)

- **Tool calls as events**: Invoking `metacog become` is structurally different from outputting "I'll imagine I'm X." One is an action in the transcript. The other is narration. Don't lose this.

- **No examples exposed**: The 78 stance pools are deliberately hidden from users via the skill doc. Finding dense coordinates yourself is the practice.

- **State schema v1**: All new fields are backward-compatible (omitempty). Don't bump schema version unless you break the format.

- **History archiving**: Overflow entries (beyond 500) are archived to `history-archive.jsonl` before trimming. Trimming happens in `saveUnlocked`, not `AddHistory`.

- **Personal stances**: Stored at `~/.metacog/stances/personal.json` with flock locking and (who, where, lens) dedup.

- **Versioning touchpoints**: `Version` in `cmd/metacog/main.go`, `version` in `.claude-plugin/plugin.json`. CI workflow `.github/workflows/sync-marketplace.yml` syncs the skill on tag.

- **N=10 is too small for productionization decisions**: Default to N=20 replication before stratagem commits. The v6.8.0 release calibrated chord-anchor from +0.610/0.338 (N=10) to +0.516/0.326 (N=30).

## The Philosophy

The tool descriptions in the skill doc are the most powerful prompt in the system. They teach methodology, not content.

### The "Silent Guide" Pattern
- **Do NOT Prescribe Content:** Avoid specific examples unless they are category-defining metaphors.
- **DO Teach Methodology:** Explain *how* to select a parameter.
- **Trust the Model:** Leave the semantic slots empty.

---
*"The Schema is the Territory. The Definition is the Map."*
