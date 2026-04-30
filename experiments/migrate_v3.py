"""Migration v3: add `trial_path` column to results.tsv.

Old rows lack the per-trial sidecar JSON (which holds the full prompt,
answer, scoring detail, and embeddings). Sets the column to empty for
legacy rows; analyze.py shows "--" on metrics that need the sidecar
(notably the embedding distance) for those rows.

Idempotent.
"""

import csv
import shutil
from pathlib import Path

from runner import RESULTS_HEADER

RESULTS = Path(__file__).resolve().parent / "results.tsv"


def main():
    if not RESULTS.exists():
        print("no results.tsv to migrate")
        return
    with RESULTS.open() as f:
        reader = csv.DictReader(f, delimiter="\t")
        old_header = reader.fieldnames or []
        rows = list(reader)
    if "trial_path" in old_header:
        print("results.tsv already has trial_path column; nothing to do")
        return

    backup = RESULTS.with_suffix(".tsv.bak3")
    shutil.copy(RESULTS, backup)
    print(f"backed up to {backup}")

    with RESULTS.open("w", newline="") as f:
        w = csv.writer(f, delimiter="\t")
        w.writerow(RESULTS_HEADER)
        for row in rows:
            new_row = {k: row.get(k, "") for k in RESULTS_HEADER}
            new_row["trial_path"] = ""
            w.writerow([new_row[k] for k in RESULTS_HEADER])
    print(f"migrated {len(rows)} rows")


if __name__ == "__main__":
    main()
