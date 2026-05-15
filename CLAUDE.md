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

### Primitives (eighteen)

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

The two newest primitives (added 2026-05-14 in v6.7.0 after the recursive-design rounds in `experiments/FINDINGS.md` identified two prose moves the existing surface didn't cover; both ALL CAPS structural-surface output):

- **witness** (`position`, `observed`, `distance`) -- speak from a meta-stance observing the producer of speech. Distinct from `become` (which adopts an identity); `witness` is the third-person observer construction (Sebald's narrator, late Stevens, Carson's *Plainwater*). Holds structural separation between speaker and content; collapsing to first-person breaks the witness.
- **apophasis** (`subject`, `negation` x3+, `residue`) -- articulate by enumerated negation. Distinct from `silence` (which refuses output) and `disjunction` (which asserts binary contradiction); `apophasis` enumerates what something is NOT as the load-bearing articulation, with a residue field for what no negation reaches. The negative-theology register (Pseudo-Dionysius, Eckhart, Mahayana via negativa).

### Stratagems (twenty-six)

Named compositions of primitives plus reflection (`THINK`) and action (`ACTION`) steps. Defined in `Stratagems` map in `stratagem.go`. Active stratagem state is `state.Stratagem` (`{Name, Step, StepsCompleted, StartedAt}`). Lifecycle: `stratagem start <name>` -> primitives auto-advance matching steps -> `stratagem next` advances reflection/action steps -> completion records a `stratagem` history entry with `event=completed`.

Survivors of the original sixteen (use original-six primitives only): pivot, mirror, stack, anchor, reset, invocation, veil, scrying, sacrifice, fool, inversion, gift, zen.

(`banishing`, `drift`, `error` were dropped in v6.3.0 -- the all-stratagem sweep found them clustered at emb_d ~0.10 with the rest of the non-manifold-family pack.)

The structural champion (uses fork + synthesis):

- **manifold** (fork + synthesis): when parallel reasoning needs to be made structural and you keep collapsing to one thread early. The progenitor of chorus/trinity.

Twelve empirical stratagems (first two added 2026-05-01 in v6.2.0; the next two added 2026-05-02 in v6.4.0 after the seven-new-primitives experiment harness; the fifth added 2026-05-02 in v6.5.0 after the 2x3 (structure x author) matrix validated the combined recipe across three author triples; the sixth added 2026-05-02 in v6.6.0 after the cross-model probe found register-shifts are generator-specific while extreme-author-becomes transfer; the seventh added 2026-05-15 in v6.7.1 after the recursive-design rounds; the eighth added 2026-05-15 in v6.7.2 after the cross-model probe found a new cross-model delta champion; the three new in v6.7.3 added 2026-05-15 after the round-12 occult-anchor sweep validated three new compositions; the twelfth added 2026-05-15 in v6.8.0 after rounds 16-20 of Opus validation produced the first Pareto-breakthrough on both axes simultaneously; see `experiments/FINDINGS.md`):

- **chorus** (3 becomes + fork + ritual): structural-axis champion. Three cross-domain becomes-as-events seed voice diversity, fork makes the disagreement structural, ritual locks the multi-voice answer. Deliberately omits synthesis -- the experiment found synthesis acts as a structural brake on embedding-distance.
- **trinity** (3 becomes + fork + synthesis + ritual): balanced variant. Same multi-voice base as chorus but keeps synthesis for delta lift.
- **antinomy** (3 becomes + fork + disjunction + ritual): vocabulary-axis champion. Substitutes disjunction (hard binary contradiction asserted as the operand of reasoning) for synthesis (refused-resolution between 3 lenses). At N=70 hit delta +0.347 (alt-author replication +0.233) -- the prior vocabulary-axis champion, freestyle-become, was at +0.231. Operating-inside-contradiction lifts citation density dramatically because the answer must keep naming the specific propositions.
- **envoy** (register + 3 becomes + fork + ritual): both-axes champion. Prepends a register-shift to the chorus structure, imposing a non-default linguistic surface that the multi-voice base then operates within. At N=70 hit delta +0.204 / emb_d 0.239 -- beats the prior structural champion (trinity-no-synthesis-alt at +0.194 / 0.226) on BOTH axes simultaneously. The register isn't a citation-stripping artifact: composing it with the trinity base preserves citations while pushing emb_d. Author-pattern result: extreme cross-domain authors push emb_d to 0.257 (envoy-extreme N=70).
- **counterpoint** (register + 2 becomes + fork + disjunction + ritual): Pareto-frontier balanced variant. Composes envoy's register-prepend with antinomy's disjunction-substitution. The 3-becomes variant hit delta +0.247 / emb_d 0.190 at N=70 (replicated at +0.202/0.188 with alt authors); the 2-becomes variant (counterpoint-duo) at N=100 hit +0.240/0.221, basically tying on delta and gaining +0.031 on emb_d. Productionized as 2-becomes in v6.5.1. Use when both axes matter and you don't want to maximize one at the other's expense. Unlike chorus/trinity/antinomy/envoy (3 becomes), counterpoint specifically benefits from the tighter binary opposition under disjunction's structure.
- **envoy-extreme** (3 becomes + fork + ritual, no register): cross-model winner. Same structure as chorus, but step prose explicitly demands HARD-extreme cross-domain authors (Sun Ra/Octavia Butler/Hilma af Klint-tier cosmologists, NOT Carson/Knuth-tier mild-academic-essayists). Empirically validated against gpt-5.5 via Codex CLI (round 4-5): hit delta +0.310 on codex (vs Sonnet's +0.190 on the same recipe). Author-extremity transfers cleanly across models; register-shifts (envoy/counterpoint's register step) are generator-specific. Use when target generator is unknown, when the recipe must work outside Sonnet, or when register-shift attempts have failed on the target.
- **duo-disjunction** (commitment + register + 2 becomes + fork + disjunction + ritual): balanced champion from the recursive-design rounds. Pre-commits to holding a binary contradiction, imposes register, then two hard-extreme cross-domain becomes carry the contradiction's poles. At N=20 hit delta +0.241 / emb_d 0.265 -- Pareto-dominates envoy-extreme (+0.190/0.257) on both axes simultaneously on Sonnet. Distinct from antinomy (which uses 3 becomes and no register/commitment): duo-disjunction commits FIRST so the contradiction is pre-locked when becomes import, then the imposed register holds across both poles. Sonnet-specific: biblical register torpedoes it on codex (-0.022). Use when both axes matter on Sonnet specifically.
- **anchor-duo** (commitment + 2 excerpts + fork + ritual): cross-model delta champion. Pre-commits to operating from two cosmological excerpts as a single composite cosmos; both fixed-points are load-bearing architecture, not metaphor. No register, no becomes -- just commitment + double-anchor. On Sonnet hit +0.314/0.222 (N=10); on codex hit **+0.377**/0.162 (N=10) -- beats envoy-extreme's cross-model +0.245 by +0.132. The 2-excerpt anchor structure transfers cleanly across models where register-shifts collapse. Use when the target model is unknown OR when you want the strongest delta lift available cross-model.
- **sigil** (name + commitment + 3 becomes + fork + ritual): True-Name champion. Validates `name` primitive in composition. Coins a sigil (Zos-Kia-Aleph-tier — a name with no prior corpus presence), charges it via commitment, then three sigil-magick / occult-cosmology author-becomes operate from inside the charged sigil. At N=10 on Sonnet hit delta +0.327 / emb_d 0.201 -- beats chorus on delta (+0.327 vs ~+0.18) and matches it on emb_d. The coined sigil propagates as load-bearing citation because the answer keeps re-citing the name. Distinct from anchor-duo (which uses external excerpts as anchors): sigil generates its own anchor in-place. Sonnet-specific; cross-model untested. Use when you want delta lift via coined-vocabulary rather than borrowed cosmology.
- **grimoire** (register + 3 becomes + fork + disjunction + ritual): imperative-register variant. Re-pitches the surface to grimoire / ritual imperative register (verb-first, second-person address to the operator, "Let the magus now…", "thus it is done") -- a register family structurally distinct from biblical parallelism, Victorian judgment, scientific hedging. Three occult-method author-becomes (disciplined-system / chaos-meta-belief / cut-up-procedural) then disjunction holds two-of-three as load-bearing. At N=10 on Sonnet hit +0.307 / emb_d 0.238 -- the highest emb_d of the round-12 family. The imperative register pushes a previously-untested region of register space. Sonnet-specific.
- **chord-anchor** (commitment + 2 excerpts + chord + ritual): Pareto-EMB_D recipe on both Sonnet and Opus. Substitutes `chord` (modes held simultaneously, single composite attention) for `fork` (parallel threads with sacrifice) in anchor-duo. Calibrated numbers at N=20-30 (initial N=10 was inflated): Opus +0.516 / 0.326 (vs anchor-duo's +0.596/0.278 -- chord trades -0.080 delta for +0.048 emb_d). Sonnet +0.328 / 0.251 (beats anchor-duo-occult +0.277/0.191 on BOTH axes; the bare variant without commitment hit +0.369 / 0.236 at N=10). The chord primitive is the structural counterpart to fork — every sentence attends through both anchored cosmoses simultaneously rather than alternating threads, producing more structural distance per token at the same citation density. First productionized stratagem using the `chord` primitive (was unused in compositions until v6.8.0). Use when both axes matter and you want simultaneous-attention binding instead of fork's threaded binding. **Does NOT transfer to codex** — all chord recipes on codex hit -0.005 to -0.103 delta at N=10-20; codex's generator doesn't execute the chord-attention pattern. Use occult-extreme on codex instead.
- **occult-extreme** (3 becomes + fork + ritual, no register, occult anchors): cross-model variant of envoy-extreme. Same 5-step structure as envoy-extreme but step prose explicitly demands HARD-extreme occult cosmologists (Crowley/Liber AL-tier, Spare/Zos Kia-tier, Dee/Enochian-tier, Carroll/chaos-tier, P-Orridge/TOPY-tier, Bruno/De Umbris Idearum-tier, Hermetic-tier, neoplatonic-magical-tier). At N=10 on codex hit delta +0.281 (beat Sonnet's +0.199 on the same recipe) -- occult anchors specifically transfer cross-model where mild-academic-essayist authors don't. The author-extremity transfer rule (envoy-extreme finding) compounds with anchor-domain specificity. Use when target generator is unknown OR when you want the envoy-extreme structure with explicitly occult anchors.

### Tool-call vs text: asymmetric amplifier (round 4-7 finding)

Cross-model probe testing recipe content delivered as tool-calls (default) vs plain text instructions found tool-call invocation is an asymmetric amplifier on the recipe's natural pull, not a fixed bonus:
- Working recipes (positive delta): tool-call mode adds small lift (+0.018 average across 4 cases).
- Broken recipes (negative delta): tool-call mode amplifies failure proportional to brokenness (-0.066 to -0.276 across 3 cases).

Practical rule for porting recipes to new generators: validate in text-instructions mode first via `METACOG_EXP_PROMPT_MODE=text-instructions` in the experiment harness. If text-mode delta is positive, tool-call mode is safe. If negative, do NOT promote -- tool-call will amplify the failure. Mechanistic story (Arditi et al. activation-direction frame): tool-call invocation is a stronger move along whichever direction the recipe pulls; stronger moves along *present* directions sharpen outputs, stronger moves along *absent* directions produce nonsense.

(`audit`, `autopsy`, `trilemma`, `survey`, `dive` were also dropped in v6.3.0 alongside their load-bearing primitives -- they sat at emb_d 0.115-0.135 across the empirical sweep.)

### Outcome tracking

`outcome --result productive|unproductive [--shift ...]` attaches an effectiveness mark to the most recent unmarked work. Two-tier search in `outcome.go`:

1. Last `stratagem` event with `event=completed` that has no later `outcome`.
2. Otherwise, last freestyle primitive that isn't inside a started/abandoned/aborted stratagem span and has no later `outcome`. Recorded with `stratagem=freestyle`.

`--amend` updates the most recent outcome rather than creating a new one. `reflect` aggregates these into completion and productivity rates.

## Key files

- `cmd/metacog/main.go` -- root cobra command, version string (must list all 18 primitives and 22 stratagems), schema version constant
- `cmd/metacog/state.go` -- State, StateManager, flock, atomic rename, history archiving
- `cmd/metacog/stratagem.go` -- Stratagems map, `StepKind` constants (one per primitive plus THINK/ACTION), step validation, lifecycle commands
- `cmd/metacog/outcome.go` -- Two-tier outcome attachment and amendment
- `cmd/metacog/inspire.go` -- Embedded stance pools (`go:embed stances/*.json`) plus personal pool at `$METACOG_HOME/stances/personal.json`
- `cmd/metacog/reflect.go` -- History aggregation into practice patterns
- `cmd/metacog/journal.go` -- `journal.jsonl` insight log, tag/session filtering
- `cmd/metacog/session.go` -- Named session tagging (auto-applied to history entries)
- `cmd/metacog/output.go` -- `FormatOutput` honouring the global `--json` flag
- `cmd/metacog/stances/*.json` -- 78 embedded pools (~450 examples), JSON arrays of `{who, where, lens}`
- `skills/metacog/SKILL.md` -- Claude Code skill document (the user-facing docs that hide implementation examples)
- `.claude-plugin/plugin.json` -- plugin manifest; version here must match `Version` in `main.go`

## Design decisions (for future-you)

- **Tool calls as events.** Invoking `metacog become` is structurally different from outputting "I'll imagine I'm X." One is an action in the transcript; the other is narration. Don't let refactors collapse this distinction.

- **No examples exposed.** The 65 stance pools under `cmd/metacog/stances/` are deliberately hidden from end users via the skill doc -- finding dense coordinates yourself is the practice. Don't surface them in `--help`, READMEs, or skill text.

- **State schema v1, additive only.** All new fields are `omitempty` and backward-compatible. Don't bump `StateSchemaVersion` unless you actually break the format; the loader rejects newer-than-known versions.

- **History archiving lives in saveUnlocked, not AddHistory.** When `len(History) > MaxHistoryEntries (500)`, overflow entries are appended to `history-archive.jsonl` before trimming. `history --full` re-merges the archive on read.

- **Personal stances dedup on (who, where, lens).** Stored at `$METACOG_HOME/stances/personal.json` with its own flock; appears as the `personal` pool in `inspire`.

- **Versioning touchpoints.** Bumping the release means updating both `Version` in `cmd/metacog/main.go` and `version` in `.claude-plugin/plugin.json`. The CI workflow in `.github/workflows/sync-marketplace.yml` syncs the skill on tag.
