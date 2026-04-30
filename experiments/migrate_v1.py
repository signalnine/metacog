"""One-shot migration: add `control` column to results.tsv and zero out novelty
for control rows (which were previously stored against a rolling baseline).

Idempotent: detects already-migrated TSVs and exits.
"""

import csv
import shutil
from pathlib import Path

from runner import RESULTS_HEADER, _row_is_control

RESULTS = Path(__file__).resolve().parent / "results.tsv"


def main():
    if not RESULTS.exists():
        print("no results.tsv to migrate")
        return
    with RESULTS.open() as f:
        reader = csv.DictReader(f, delimiter="\t")
        old_header = reader.fieldnames
        rows = list(reader)
    if old_header == RESULTS_HEADER:
        print("results.tsv already at target schema; nothing to do")
        return

    backup = RESULTS.with_suffix(".tsv.bak")
    shutil.copy(RESULTS, backup)
    print(f"backed up to {backup}")

    with RESULTS.open("w", newline="") as f:
        w = csv.writer(f, delimiter="\t")
        w.writerow(RESULTS_HEADER)
        for row in rows:
            is_ctl = _row_is_control(row)
            new_row = {
                "ts": row.get("ts", ""),
                "recipe": row.get("recipe", ""),
                "control": "1" if is_ctl else "0",
                "task": row.get("task", ""),
                "sample": row.get("sample", ""),
                "rarity": row.get("rarity", ""),
                "coherence": row.get("coherence", ""),
                # Control rows: novelty is definitionally 0 (delta vs self).
                # Non-control rows: keep the old novelty value as-is even
                # though it was computed against a snapshot of the baseline
                # at trial time. analyze.py recomputes from full pool anyway.
                "novelty": "+0.0000" if is_ctl else row.get("novelty", ""),
                "n_entities": row.get("n_entities", ""),
                "answer_len": row.get("answer_len", ""),
                "answer_preview": row.get("answer_preview", ""),
            }
            w.writerow([new_row[k] for k in RESULTS_HEADER])
    print(f"migrated {len(rows)} rows")


if __name__ == "__main__":
    main()
