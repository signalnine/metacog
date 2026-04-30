"""Pure-function tests for runner logic.

Run with: python -m unittest experiments/test_runner.py
"""

import unittest

from runner import baselines_from_rows, compute_novelty


class TestComputeNovelty(unittest.TestCase):
    def test_control_trial_is_zero(self):
        # Control trials ARE the baseline; novelty against self is definitionally 0.
        self.assertEqual(
            compute_novelty(rarity=0.7, coherence=0.95, baseline=None, is_control=True),
            0.0,
        )
        self.assertEqual(
            compute_novelty(rarity=0.7, coherence=0.95, baseline=0.6, is_control=True),
            0.0,
        )

    def test_non_control_no_baseline_is_none(self):
        # Non-control trial with no baseline pool: novelty is undefined, not zero.
        self.assertIsNone(
            compute_novelty(rarity=0.7, coherence=0.95, baseline=None, is_control=False)
        )

    def test_non_control_with_baseline(self):
        self.assertAlmostEqual(
            compute_novelty(rarity=0.8, coherence=1.0, baseline=0.6, is_control=False),
            0.2,
        )

    def test_non_control_negative_delta(self):
        self.assertAlmostEqual(
            compute_novelty(rarity=0.4, coherence=1.0, baseline=0.6, is_control=False),
            -0.2,
        )


class TestBaselinesFromRows(unittest.TestCase):
    def test_filters_by_control_field_not_recipe_name(self):
        rows = [
            {"task": "t1", "recipe": "null", "rarity": "0.6", "coherence": "1.0", "control": "1"},
            {"task": "t1", "recipe": "null", "rarity": "0.4", "coherence": "1.0", "control": "1"},
            # Same task, non-control: must NOT contribute to baseline
            {"task": "t1", "recipe": "stack", "rarity": "0.9", "coherence": "1.0", "control": "0"},
            {"task": "t2", "recipe": "null", "rarity": "0.7", "coherence": "1.0", "control": "1"},
        ]
        b = baselines_from_rows(rows)
        self.assertAlmostEqual(b["t1"], 0.5)
        self.assertAlmostEqual(b["t2"], 0.7)
        self.assertNotIn("nonexistent", b)

    def test_recipe_named_null_but_control_false_is_excluded(self):
        # Hypothetical recipe named "null" but with control=false should NOT count.
        rows = [
            {"task": "t1", "recipe": "null", "rarity": "0.6", "coherence": "1.0", "control": "0"},
        ]
        self.assertEqual(baselines_from_rows(rows), {})

    def test_alternate_control_recipe_name_counts(self):
        # A control recipe named something other than "null" must still count.
        rows = [
            {"task": "t1", "recipe": "vanilla", "rarity": "0.6", "coherence": "1.0", "control": "1"},
            {"task": "t1", "recipe": "vanilla", "rarity": "0.4", "coherence": "1.0", "control": "1"},
        ]
        b = baselines_from_rows(rows)
        self.assertAlmostEqual(b["t1"], 0.5)

    def test_empty_rows(self):
        self.assertEqual(baselines_from_rows([]), {})

    def test_handles_legacy_rows_without_control_column(self):
        # Migration path: rows from before the control column existed should be
        # interpreted as control if recipe == "null", non-control otherwise.
        rows = [
            {"task": "t1", "recipe": "null", "rarity": "0.6", "coherence": "1.0"},
            {"task": "t1", "recipe": "stack", "rarity": "0.9", "coherence": "1.0"},
        ]
        b = baselines_from_rows(rows)
        self.assertAlmostEqual(b["t1"], 0.6)


if __name__ == "__main__":
    unittest.main()
