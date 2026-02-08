# Practice Patterns in Reflect

## Context

Two opportunity areas from the metacog review:
1. **Learning/surfacing** — Reflect shows data but doesn't adapt. Surface productive patterns, flag underused primitives.
2. **Cross-session composition** — Insights and transformation products die with the conversation.

This design addresses both in a single new section in `metacog reflect` output. No new commands, no new state, no new data capture. Pure aggregation of existing history.

## Design

### New function: `FormatPracticePatterns(s *State) string`

**File:** `cmd/metacog/reflect.go`

Returns a "Practice patterns" section with two halves:

### What Worked

Scan history for outcome entries where `result=productive`. For each, walk backward to find the nearest `become` entry (name + lens) and nearest `drugs` entry (substance). Extract the stratagem and shift from the outcome itself.

Show the last 5 productive outcomes in reverse chronological order. If more exist, append `(N more productive outcomes in history)`.

Format: `stratagem — name/lens + substance — "shift"`
- Omit identity if no become found
- Omit substance if no drugs found
- Show `(no config)` if neither found
- Omit quoted shift if no shift text recorded

### Underused

Count become, drugs, ritual calls from history. If total >= 5, flag any primitive at <20% of total. Each flag names the primitive's function:
- become: "identity shifting"
- drugs: "substrate modification"
- ritual: "threshold-crossing"

Skip section entirely if all three are >=20%.

### Example Output

```
Practice patterns:
  What worked (last 5 of 8):
    pivot — Ada/logic + caffeine — "reframed auth as state machine"
    freestyle — Eno/ambient + psilocybin — "found the shape underneath"
    mirror — Ada/logic — "synthesis: both assumed single-user"
    freestyle — (no config) — "shifted from debugging to redesigning"
    pivot — Doepfer/modular + caffeine — "saw the feedback loop"
    (3 more productive outcomes in history)

  Underused:
    ritual is 8% of your practice (1 of 12 primitives) — threshold-crossing is available
    drugs is 17% of your practice (2 of 12 primitives) — substrate modification is available
```

## Wiring

Insert into reflectCmd output after FormatRecentInsights, before FormatAdvisories:
1. FormatReflection (existing)
2. FormatRecentInsights (existing)
3. FormatPracticePatterns (new)
4. FormatAdvisories (existing)

## SKILL.md

No changes. The existing pre-flight check instruction already tells the LLM to run `metacog reflect` and surface what it finds. The new section will be visible automatically.

## Tests

Unit tests in `cmd/metacog/reflect_test.go`:
- `TestPracticePatternsWhatWorked` — 3 productive outcomes with configs shows all 3
- `TestPracticePatternsOverflow` — 7 productive outcomes shows last 5 + overflow note
- `TestPracticePatternsNoConfig` — productive outcome with no become/drugs shows `(no config)`
- `TestPracticePatternsNoShift` — productive outcome without shift text omits quote
- `TestPracticePatternsUnderused` — 10 becomes, 1 drugs, 0 rituals flags drugs and ritual
- `TestPracticePatternsBalanced` — roughly even usage, no underused section
- `TestPracticePatternsEmpty` — no productive outcomes + <5 primitives returns empty string

Integration test in `cmd/metacog/integration_test.go`:
- `TestPracticePatternsIntegration` — build productive outcomes via binary, verify section appears in reflect

## Verification

```bash
go test ./cmd/metacog/ -v
go test ./cmd/metacog/ -tags integration -v
go build -o metacog ./cmd/metacog/
```
