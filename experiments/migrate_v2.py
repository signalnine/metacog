"""One-shot migration v2: add `entity_rarities` column to results.tsv.

Old rows lack the per-entity data needed for the alternative metrics
(max/sum/count_high/geo_mean). This migration adds an empty column for
those rows; analyze.py shows "--" on those metrics.

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
    if "entity_rarities" in old_header:
        print("results.tsv already has entity_rarities column; nothing to do")
        return

    backup = RESULTS.with_suffix(".tsv.bak2")
    shutil.copy(RESULTS, backup)
    print(f"backed up to {backup}")

    with RESULTS.open("w", newline="") as f:
        w = csv.writer(f, delimiter="\t")
        w.writerow(RESULTS_HEADER)
        for row in rows:
            new_row = {k: row.get(k, "") for k in RESULTS_HEADER}
            new_row["entity_rarities"] = ""  # not recoverable; empty for legacy rows
            w.writerow([new_row[k] for k in RESULTS_HEADER])
    print(f"migrated {len(rows)} rows")


if __name__ == "__main__":
    main()
