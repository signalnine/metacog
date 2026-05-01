# Seven new primitives for v6.3.0

Date: 2026-04-30
Source of inspiration: experiment harness in `experiments/`. After the
all-stratagem sweep showed the structural-novelty axis is uniquely owned
by manifold-family recipes, we want primitives that supply *new levers*
not yet present in the 9-primitive surface (feel, become, drugs, name,
ritual, meditate, counterfactual, synthesis, fork).

## What's missing in the 9?

| Domain | Current primitive | Gap |
|---|---|---|
| Identity / voice | become | identity-shift only; no register-flip without identity-flip |
| Branches | fork | sequential branches with sacrifice; no simultaneous-hold |
| Output | (all primitives produce) | no deliberate-silence event |
| Anchor | counterfactual | only structural-walls; no external-fragment fixed-point |
| Reasoning | (none) | no pre-commitment / binding-before-exploration |
| Lens-conflict | synthesis | 3-lens refused-resolution; no 2-prop hard-contradiction |
| Substrate | drugs | category-loosening; no sub-semantic dropout |

These are the seven gaps; one primitive per gap. The names are short,
verb-first, and avoid overlap with existing primitive names.

## The seven

### register
Flip the linguistic register of the current voice without changing
identity. Where `become` imports a methodology, `register` re-pitches
the same speaker — academic to vernacular, oracular to technical,
earnest to arch. Empirically motivated by the observation that the
emb_d axis depends partly on rare vocabulary which is register-bound
within an author.

Flags: `--from`, `--to`, `--rationale`
Output: structural-soft block, single shift line + coda.
State: appends history; does not touch `state.Identity`.

### chord
Hold multiple modes-of-attention simultaneously without alternating.
Distinct from `become` (sequential identity-shifts) and `fork`
(branched parallelism with sacrifice). The chord is a single event in
which several attentional modes overlap on the same observation. Where
chorus uses three discrete `become` events to seed voice-diversity at
the IDENTITY level, `chord` adds attentional polyphony at the
ATTENTION level — orthogonal to identity.

Flags: `--modes` (array, min 2), `--target`
Output: structural-soft block listing the modes plus a coda forbidding
mode-switching during the held window.
State: history-only.

### silence
Refuse to produce articulated output about something. The call itself
is the artifact: it marks a held question that articulation would
falsify, or that requires a longer ripening than the current move.
Strongest expression of "tool calls as events" — the call is structurally
visible; the absence-of-prose is the cognitive event.

Flags: `--about`, `--reason`, `--duration` (free-form)
Output: a single-line ack and nothing more.
State: history-only.

### excerpt
Cite a verbatim fragment from outside the current generative frame as a
fixed point. Where `become` channels a voice and generates new prose in
that voice, `excerpt` pins a specific phrase that interrupts the model's
own continuation. The fragment is treated as load-bearing surface, not
a stylistic prompt.

Flags: `--source`, `--fragment`, `--why`
Output: structural block displaying the fragment as a quoted anchor
plus attribution and load-bearing rationale.
State: history-only.

### commitment
Pre-commit to a binding stance, prediction, or position before
exploration. Stakes-bearing reasoning: the cost of motivated reasoning
goes up because the binding is now in the transcript. Distinct from
`ritual` (which seals after the work) — `commitment` binds before.

Flags: `--binding`, `--stakes`, `--falsifier`
Output: ALL-CAPS structural block with the three slots and a coda
declaring the constraint active.
State: history-only.

### disjunction
Assert two propositions that must both be true even though they cannot
be. Distinct from `synthesis` (3 lenses with named blindspots, refused
resolution): `disjunction` is a sharp binary contradiction with no
blindspot framing — the contradiction itself is the cognitive operand,
not a meta-comment on lens-conflict.

Flags: `--proposition-a`, `--proposition-b`, `--why-both-required`
Output: ALL-CAPS structural block with both propositions stated as
fact and a coda forbidding resolution.
State: history-only.

### glossolalia
Deliberately leave the domain of semantic language. Where `drugs`
loosens categories within language, `glossolalia` drops the requirement
that tokens carry meaning. The schema is tight precisely because the
output is loose — the call exists to license sub-semantic generation
*as a discrete event* rather than letting it bleed into surrounding
prose.

Flags: `--pretext`, `--duration-tokens` (int, hint not enforced),
`--return-trigger`
Output: structural ALL-CAPS preamble followed by a block boundary the
caller is licensed to fill with non-meaning-bearing tokens, then a
return marker.
State: history-only.

## Per-primitive flag mapping

| Primitive | Required flags |
|---|---|
| register | `--from`, `--to`, `--rationale` |
| chord | `--modes` (min 2), `--target` |
| silence | `--about`, `--reason`, `--duration` |
| excerpt | `--source`, `--fragment`, `--why` |
| commitment | `--binding`, `--stakes`, `--falsifier` |
| disjunction | `--proposition-a`, `--proposition-b`, `--why-both-required` |
| glossolalia | `--pretext`, `--duration-tokens`, `--return-trigger` |

Validation:
- `chord.--modes` must contain at least 2 entries.
- `glossolalia.--duration-tokens` must be a positive integer.
- All other required flags must be non-empty.

## Output registers

- **register, chord**: structural-soft (block format, sentence-case headers).
- **silence**: minimal — single ack line.
- **excerpt**: quoted-fragment block with attribution.
- **commitment, disjunction, glossolalia**: structural ALL-CAPS (the
  most binding / most extreme cognitive moves get the most formal
  register).

## Stratagem integration

Each primitive gets a `StepKind` constant in `stratagem.go`. Each is
added to the switch in `outcome.go::findLastPrimitive`. No new
stratagems on this pass — wiring keeps the option open and we want the
experiment harness to recommend stratagem compositions, not pre-specify
them.

## File layout

```
cmd/metacog/
  register.go     + _test.go
  chord.go        + _test.go
  silence.go      + _test.go
  excerpt.go      + _test.go
  commitment.go   + _test.go
  disjunction.go  + _test.go
  glossolalia.go  + _test.go
  stratagem.go    + 7 StepKind constants
  outcome.go      + 7 cases in findLastPrimitive
  main.go         + extend version-string primitive list (9 -> 16)
```

## Testing

For each primitive:
- Missing required flag returns a usage error.
- History entry contains all expected params (arrays joined with `"; "`
  per the existing convention).
- Output contains expected section headers and coda.
- `outcome_test.go` cross-cutting: each new primitive can be attached
  as a freestyle outcome.
- Integration test: each primitive can satisfy a stratagem step when
  its `StepKind` is the current expected step.

## Verification gate

```
go build -o metacog ./cmd/metacog/
go test ./cmd/metacog/ -v
go test ./cmd/metacog/ -tags integration -v
```

All three must pass with zero failures before commit.

## Versioning

The drops + adds together are v6.3.0. After this commit:
- 9 primitives -> 16 primitives
- 16 stratagems unchanged from Phase 1
- main.go Version stays "6.3.0"
- plugin.json description updates to "Sixteen primitives compose into
  sixteen transformation stratagems."

## Out of scope (intentionally)

- New stratagems composing the new primitives. The experiment harness
  will tell us which compositions are worth productionizing. Premature
  stratagem-design risks repeating the audit/autopsy/trilemma/survey/
  dive over-fit.
- Recipes in `experiments/recipes/` for the new primitives. Will be
  added in a separate commit after the primitives land.
- Changes to `feel.--since-last`, the existing structural-soft pattern,
  or the `--json` global flag.
