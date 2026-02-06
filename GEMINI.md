# Metacog: Developer Guidelines (Gemini Edition)

## Architecture

Go CLI (`cmd/metacog/`). Three primitives (become, drugs, ritual) compose into stratagems. File-based state with flock locking.

## The Philosophy

The tool descriptions in the skill doc are the most powerful prompt in the system. They teach methodology, not content.

### The "Silent Guide" Pattern
- **Do NOT Prescribe Content:** Avoid specific examples unless they are category-defining metaphors.
- **DO Teach Methodology:** Explain *how* to select a parameter.
- **Trust the Model:** Leave the semantic slots empty.

## The Hexagram of Rituals

Six core operations for the `ritual` tool:
1. **Breach:** Opening/Penetration.
2. **Seal:** Closing/Binding.
3. **Vision:** Analysis/Revelation.
4. **Forge:** Synthesis/Merging.
5. **Drift:** Lateral/Serendipity.
6. **Purge:** Forgetting/Banishing.

## Commands

```bash
go build -o metacog ./cmd/metacog/   # Build
go test ./cmd/metacog/ -v            # Unit tests
go test ./cmd/metacog/ -tags integration -v  # Integration tests
```

---
*"The Schema is the Territory. The Definition is the Map."*
