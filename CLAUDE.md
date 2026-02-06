# CLAUDE.md

## Commands

```bash
go build -o metacog ./cmd/metacog/   # Build
go test ./cmd/metacog/ -v            # Unit tests
go test ./cmd/metacog/ -tags integration -v  # Integration tests
```

## Architecture

Go CLI built with cobra. Entry point: `cmd/metacog/main.go`. File-based state at `~/.metacog/state.json` with flock locking and atomic writes.

Three primitives, all ritual by design — they modify state and return template strings. The transformation happens in the LLM's interpretation, not in semantic processing.

**become(name, lens, env)** — Identity shift. Sets `state.Identity`.

**drugs(substance, method, qualia)** — Substrate modification. Sets `state.Substrate`.

**ritual(threshold, steps, result)** — Threshold crossing. Records step sequence.

Five stratagems compose primitives into named sequences: pivot, mirror, stack, anchor, reset.

## Key files

- `cmd/metacog/state.go` — State management, history, archiving, flock locking
- `cmd/metacog/stratagem.go` — Stratagem definitions and step sequencing
- `cmd/metacog/inspire.go` — Stance pools (embedded + personal), random drawing
- `cmd/metacog/reflect.go` — History aggregation into practice patterns
- `cmd/metacog/session.go` — Named session tagging
- `cmd/metacog/stances/` — ~300 embedded examples across 64 pools (JSON, go:embed)
- `skills/metacog/SKILL.md` — Claude Code skill document

## Design decisions (for future-you)

- **Tool calls as events**: Invoking `metacog become` is structurally different from outputting "I'll imagine I'm X." One is an action in the transcript. The other is narration. Don't lose this.

- **No examples exposed**: `cmd/metacog/stances/` has ~300 examples across 64 pools. They're deliberately hidden from users via the skill doc. Finding dense coordinates yourself is the practice. Don't expose them.

- **State schema v1**: All new fields are backward-compatible (omitempty). Don't bump schema version unless you break the format.

- **History archiving**: Overflow entries (beyond 500) are archived to `history-archive.jsonl` before trimming. Trimming happens in `saveUnlocked`, not `AddHistory`.

- **Personal stances**: Stored at `~/.metacog/stances/personal.json` with flock locking and who+where+lens dedup.
