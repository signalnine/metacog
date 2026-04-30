# Six new primitives + two design moves from upstream

Date: 2026-04-30
Source of inspiration: `inanna-malick/metacog` commit `007faae` (2026-03-26) "iterated".

## Scope

Port six primitives added upstream after our v5 fork point, plus two smaller design moves, into the Go CLI. No new stratagems on this pass; primitives are wired so future stratagems can compose them.

New primitives: `counterfactual`, `deconstruct`, `synthesis`, `fork`, `measure`, `tether`.
Design moves: `feel --since-last` optional flag; "response gives nothing" pattern applied to `deconstruct`.

## Decisions

- **Stratagem integration: full first-class.** Add six new constants to `StepKind` and the corresponding cases to `outcome.go::findLastPrimitive`. No new stratagems written this pass; the wiring keeps the option open.
- **Output style: verbatim upstream.** ALL CAPS section headers, structural blocks. Tonally distinct from `feel`/`meditate`, deliberately. The visual register marks these as cognitively heavy.
- **`feel.since_last`: user-provided, optional.** Matches upstream. No auto-derivation from history -- articulating the diff is itself the practice.
- **`deconstruct` response gives nothing.** Output echoes only `CORE MECHANIC` plus a one-line coda. The other fields are recorded in history but not surfaced. The schema is the work; the echo is a receipt.
- **History params convention.** Strings stored directly; arrays joined with `"; "` per `ritual.go`; synthesis's nested lens objects flattened to `lens_a_name`/`lens_a_verdict`/`lens_a_blindspot` and same for B/C.
- **Tether is stateless.** No new state field, no enforcement. Framing only, matching upstream.

## File layout

```
cmd/metacog/
  counterfactual.go + _test.go    new
  deconstruct.go    + _test.go    new
  synthesis.go      + _test.go    new
  fork.go           + _test.go    new
  measure.go        + _test.go    new
  tether.go         + _test.go    new
  feel.go                          + --since-last
  feel_test.go                     + since_last cases
  stratagem.go                     + 6 StepKind constants
  outcome.go                       + 6 cases in findLastPrimitive
  outcome_test.go                  + freestyle attachment for each
  main.go                          + extend version-string primitive list
```

## Per-primitive flag mapping

| Primitive | Flags |
|---|---|
| counterfactual | `--situation`, `--fitness-function`, `--load-bearing-walls` (array, min 3), `--pruned` (array, may be empty), `--wall-to-remove`, `--inverse-position` |
| deconstruct | `--subject`, `--core-mechanic`, `--structural-dependencies` (array), `--resource-inputs` (array), `--failure-modes` (array), `--output-artifacts` (array) |
| synthesis | `--problem`, `--lens-a-name`, `--lens-a-verdict`, `--lens-a-blindspot`, same for `b`/`c`, `--suppressed-tension` |
| fork | `--threads` (array, min 2), `--divergence-vector`, `--sacrifice-condition` |
| measure | `--target-concept`, `--safe-isomorph`, `--required-precision`, `--loss-gradient` |
| tether | `--anchor-point`, `--tension-limit`, `--auto-revert-trigger` |

All flags required except `counterfactual.--pruned` (may be the empty list).

## Validation

- `counterfactual`: `--load-bearing-walls` must contain at least 3 entries; `--wall-to-remove` must equal one of them.
- `fork`: `--threads` must contain at least 2 entries.
- All others: required flags non-empty.

## Output formats

Verbatim from upstream `src/index.ts` per primitive. Section headers ALL CAPS, blocks separated by blank lines, codas preserved. `deconstruct` output is the minimal `CORE MECHANIC: %s\n\nAtoms extracted...` form.

## Testing

Each `_test.go`:
- Missing-required-flag returns a usage error.
- History record contains all expected params (arrays joined with `"; "`).
- Output contains the expected section headers and coda.
- For `counterfactual`/`fork`: validation errors fire when constraints violated.

Cross-cutting:
- `outcome_test.go`: each new primitive can be attached as a freestyle outcome.
- `feel_test.go`: `--since-last` populates params and prepends to output when set; absent when unset.

## Verification gate

```
go build -o metacog ./cmd/metacog/
go test ./cmd/metacog/ -v
go test ./cmd/metacog/ -tags integration -v
```

All three must pass with zero failures before commit.
