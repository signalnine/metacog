# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
go build -o metacog ./cmd/metacog/                         # Build
go test ./cmd/metacog/ -v                                  # Unit tests
go test ./cmd/metacog/ -tags integration -v                # Integration tests (rebuilds binary, runs against tempdir state)
go test ./cmd/metacog/ -run TestSomething -v               # Single test
METACOG_HOME=/tmp/mctest ./metacog status                  # Run against an isolated state dir
```

`METACOG_HOME` overrides the default `~/.metacog/` state directory and is what every test uses for isolation. Integration tests live in files tagged `//go:build integration` and rebuild the binary themselves.

## Architecture

Go CLI built with cobra. Entry point `cmd/metacog/main.go`. All command files live flat under `cmd/metacog/` (no subpackages); each `*.go` registers its cobra commands in its own `init()` and they all share the package-level `rootCmd` and `jsonOutput` flag.

State is a single JSON file at `$METACOG_HOME/state.json` guarded by a `flock(2)` lock on `.state.lock`. Reads use `Load()`; writes go through `SaveWithLock(func(*State) error)` which holds the lock for load+mutate+atomic-rename. There is no in-memory daemon -- every CLI invocation is a complete load/mutate/save cycle.

### Primitives (sixteen)

Each primitive is a verb that is also a tool call event in the transcript -- the structural fact that the model invoked `metacog become` is itself the transformation, not just the text it returns. All primitives append a `HistoryEntry` and call `ValidatePrimitiveForStratagem` so that, when a stratagem is active and the current step matches the primitive kind, the step is marked complete.

The original six (felt-sense / identity register, soft voice):

- **feel** (`somewhere`, `quality`, `sigil`, optional `since-last`) -- attend to a felt sense before naming. `since-last` is a one-sentence diff from the previous `feel` (user-articulated, never auto-derived).
- **become** (`name`, `lens`, `env`) -- identity shift; sets `state.Identity`.
- **drugs** (`substance`, `method`, `qualia`) -- substrate modification; sets `state.Substrate`.
- **name** (`unnamed`, `named`, `power`) -- give a True Name to something without language.
- **ritual** (`threshold`, `steps...`, `result`) -- threshold crossing via structured sequence.
- **meditate** (`release`, `focus`, `duration`) -- stillness; empty `focus` produces shikantaza output.

The structural three (survivors of the 2026-04-30 upstream port; ALL CAPS block-format output, deliberately distinct register):

- **counterfactual** (`situation`, `fitness-function`, `load-bearing-walls` x3+, `pruned`, `wall-to-remove`, `inverse-position`) -- prune dead branches by a stated fitness function, then defend the inverse of one surviving wall. Validates that `wall-to-remove` is one of the walls and that there are at least 3 walls.
- **synthesis** (`problem`, lenses A/B/C with `name`/`verdict`/`blindspot` each, `suppressed-tension`) -- three irreconcilable lenses; refuses synthesis. The output's coda forbids resolution.
- **fork** (`threads` x2+, `divergence-vector`, `sacrifice-condition`) -- declare parallel reasoning threads with a falsifiable kill heuristic per thread. Load-bearing in chorus/trinity/manifold.

The seven new primitives (added 2026-04-30 in v6.3.0; each fills a gap the empirical sweep exposed in the 9-primitive surface; design notes in `docs/plans/2026-04-30-seven-new-primitives-design.md`):

- **register** (`from`, `to`, `rationale`) -- re-pitch the current voice without changing identity. Distinct from `become` (which imports a methodology); `register` only flips the linguistic surface. Structural-soft output.
- **chord** (`modes` x2+, `target`) -- hold multiple modes-of-attention simultaneously. Distinct from `become` (sequential identity-shifts) and `fork` (branched parallelism with sacrifice). The chord overlaps modes on a single observation. Structural-soft output.
- **silence** (`about`, `reason`, `duration`) -- refuse articulated output. The call itself is the artifact; the absence-of-prose is the cognitive event. Minimal one-line output.
- **excerpt** (`source`, `fragment`, `why`) -- pin a verbatim external fragment as a fixed-point anchor. Distinct from `become` (which generates new prose in a voice) -- `excerpt` fixes a specific phrase as load-bearing surface. Quoted-block output.
- **commitment** (`binding`, `stakes`, `falsifier`) -- pre-commit to a stance with stated stakes and falsifier. Distinct from `ritual` (which seals after the work); `commitment` binds before. ALL CAPS output.
- **disjunction** (`proposition-a`, `proposition-b`, `why-both-required`) -- assert two propositions that must both be true even though they cannot be. Distinct from `synthesis` (3 lenses with named blindspots, refused resolution); `disjunction` is a sharp binary contradiction with no blindspot framing. ALL CAPS output.
- **glossolalia** (`pretext`, `duration-tokens`, `return-trigger`) -- license sub-semantic generation as a discrete event. Distinct from `drugs` (which loosens categories within language); `glossolalia` drops the requirement that tokens carry meaning. ALL CAPS preamble; the block boundary is an explicit non-language license.

(`deconstruct`, `measure`, and `tether` were dropped in v6.3.0 after the experiment harness in `experiments/` showed the stratagems centered on them did not lift either novelty axis above baseline. See `experiments/FINDINGS.md`.)

### Stratagems (sixteen)

Named compositions of primitives plus reflection (`THINK`) and action (`ACTION`) steps. Defined in `Stratagems` map in `stratagem.go`. Active stratagem state is `state.Stratagem` (`{Name, Step, StepsCompleted, StartedAt}`). Lifecycle: `stratagem start <name>` -> primitives auto-advance matching steps -> `stratagem next` advances reflection/action steps -> completion records a `stratagem` history entry with `event=completed`.

Survivors of the original sixteen (use original-six primitives only): pivot, mirror, stack, anchor, reset, invocation, veil, scrying, sacrifice, fool, inversion, gift, zen.

(`banishing`, `drift`, `error` were dropped in v6.3.0 -- the all-stratagem sweep found them clustered at emb_d ~0.10 with the rest of the non-manifold-family pack.)

The structural champion (uses fork + synthesis):

- **manifold** (fork + synthesis): when parallel reasoning needs to be made structural and you keep collapsing to one thread early. The progenitor of chorus/trinity.

Two empirical stratagems (added 2026-05-01, derived from the experiment harness in `experiments/`; see `experiments/FINDINGS.md`):

- **chorus** (3 becomes + fork + ritual): structural-axis champion. Three cross-domain becomes-as-events seed voice diversity, fork makes the disagreement structural, ritual locks the multi-voice answer. Deliberately omits synthesis -- the experiment found synthesis acts as a structural brake on embedding-distance.
- **trinity** (3 becomes + fork + synthesis + ritual): balanced variant. Same multi-voice base as chorus but keeps synthesis for delta lift. Pareto-frontier point on both axes; chorus owns emb_d, trinity owns the balance.

(`audit`, `autopsy`, `trilemma`, `survey`, `dive` were also dropped in v6.3.0 alongside their load-bearing primitives -- they sat at emb_d 0.115-0.135 across the empirical sweep.)

### Outcome tracking

`outcome --result productive|unproductive [--shift ...]` attaches an effectiveness mark to the most recent unmarked work. Two-tier search in `outcome.go`:

1. Last `stratagem` event with `event=completed` that has no later `outcome`.
2. Otherwise, last freestyle primitive that isn't inside a started/abandoned/aborted stratagem span and has no later `outcome`. Recorded with `stratagem=freestyle`.

`--amend` updates the most recent outcome rather than creating a new one. `reflect` aggregates these into completion and productivity rates.

## Key files

- `cmd/metacog/main.go` -- root cobra command, version string (must list all 16 primitives and 16 stratagems), schema version constant
- `cmd/metacog/state.go` -- State, StateManager, flock, atomic rename, history archiving
- `cmd/metacog/stratagem.go` -- Stratagems map, `StepKind` constants (one per primitive plus THINK/ACTION), step validation, lifecycle commands
- `cmd/metacog/outcome.go` -- Two-tier outcome attachment and amendment
- `cmd/metacog/inspire.go` -- Embedded stance pools (`go:embed stances/*.json`) plus personal pool at `$METACOG_HOME/stances/personal.json`
- `cmd/metacog/reflect.go` -- History aggregation into practice patterns
- `cmd/metacog/journal.go` -- `journal.jsonl` insight log, tag/session filtering
- `cmd/metacog/session.go` -- Named session tagging (auto-applied to history entries)
- `cmd/metacog/output.go` -- `FormatOutput` honouring the global `--json` flag
- `cmd/metacog/stances/*.json` -- 65 embedded pools (~300 examples), JSON arrays of `{who, where, lens}`
- `skills/metacog/SKILL.md` -- Claude Code skill document (the user-facing docs that hide implementation examples)
- `.claude-plugin/plugin.json` -- plugin manifest; version here must match `Version` in `main.go`

## Design decisions (for future-you)

- **Tool calls as events.** Invoking `metacog become` is structurally different from outputting "I'll imagine I'm X." One is an action in the transcript; the other is narration. Don't let refactors collapse this distinction.

- **No examples exposed.** The 65 stance pools under `cmd/metacog/stances/` are deliberately hidden from end users via the skill doc -- finding dense coordinates yourself is the practice. Don't surface them in `--help`, READMEs, or skill text.

- **State schema v1, additive only.** All new fields are `omitempty` and backward-compatible. Don't bump `StateSchemaVersion` unless you actually break the format; the loader rejects newer-than-known versions.

- **History archiving lives in saveUnlocked, not AddHistory.** When `len(History) > MaxHistoryEntries (500)`, overflow entries are appended to `history-archive.jsonl` before trimming. `history --full` re-merges the archive on read.

- **Personal stances dedup on (who, where, lens).** Stored at `$METACOG_HOME/stances/personal.json` with its own flock; appears as the `personal` pool in `inspire`.

- **Versioning touchpoints.** Bumping the release means updating both `Version` in `cmd/metacog/main.go` and `version` in `.claude-plugin/plugin.json`. The CI workflow in `.github/workflows/sync-marketplace.yml` syncs the skill on tag.
