"""Pure-function tests for runner logic.

Run with: python -m unittest experiments/test_runner.py
"""

import unittest

from runner import (
    baselines_from_rows,
    centroid,
    compute_novelty,
    cosine_distance,
    metrics_from_rarities,
    parse_entity_rarities,
)


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


class TestMetricsFromRarities(unittest.TestCase):
    def test_empty(self):
        m = metrics_from_rarities([])
        self.assertEqual(m["max"], 0.0)
        self.assertEqual(m["sum"], 0.0)
        self.assertEqual(m["count_high"], 0)
        self.assertIsNone(m["geo_mean"])

    def test_single_entity(self):
        m = metrics_from_rarities([0.9])
        self.assertAlmostEqual(m["max"], 0.9)
        self.assertAlmostEqual(m["sum"], 0.9)
        self.assertEqual(m["count_high"], 1)
        self.assertAlmostEqual(m["geo_mean"], 0.9)

    def test_count_high_threshold(self):
        # Default threshold is 0.7
        m = metrics_from_rarities([0.9, 0.7, 0.69, 0.5])
        self.assertEqual(m["count_high"], 2, "0.7 is the threshold; 0.69 excluded")

    def test_count_high_custom_threshold(self):
        m = metrics_from_rarities([0.9, 0.8, 0.5], threshold=0.85)
        self.assertEqual(m["count_high"], 1)

    def test_max_picks_highest(self):
        m = metrics_from_rarities([0.5, 0.95, 0.3])
        self.assertAlmostEqual(m["max"], 0.95)

    def test_sum_rewards_quantity(self):
        # Five 0.5s sum to more than one 0.9 -- the sum metric values quantity.
        few_high = metrics_from_rarities([0.9])
        many_med = metrics_from_rarities([0.5, 0.5, 0.5, 0.5, 0.5])
        self.assertGreater(many_med["sum"], few_high["sum"])

    def test_geo_mean_penalizes_low_outliers(self):
        # Mean of [0.9, 0.9, 0.1] is 0.633; geomean is much lower.
        m = metrics_from_rarities([0.9, 0.9, 0.1])
        arith_mean = (0.9 + 0.9 + 0.1) / 3
        self.assertLess(m["geo_mean"], arith_mean)
        # geomean = (0.9 * 0.9 * 0.1) ** (1/3) ~= 0.433
        self.assertAlmostEqual(m["geo_mean"], (0.9 * 0.9 * 0.1) ** (1 / 3), places=4)

    def test_geo_mean_zero_in_set_returns_zero(self):
        # Geometric mean of any set containing 0 is 0.
        m = metrics_from_rarities([0.9, 0.5, 0.0])
        self.assertEqual(m["geo_mean"], 0.0)


class TestParseEntityRarities(unittest.TestCase):
    def test_empty_string(self):
        self.assertEqual(parse_entity_rarities(""), [])

    def test_missing_field(self):
        self.assertEqual(parse_entity_rarities(None), [])

    def test_well_formed_json(self):
        s = '[["Stafford Beer", 0.9], ["viable system model", 0.85]]'
        result = parse_entity_rarities(s)
        self.assertEqual(result, [("Stafford Beer", 0.9), ("viable system model", 0.85)])

    def test_malformed_returns_empty(self):
        # Old TSV rows or corrupted data: return empty list, don't crash.
        self.assertEqual(parse_entity_rarities("not json"), [])
        self.assertEqual(parse_entity_rarities("[invalid"), [])


class TestCosineDistance(unittest.TestCase):
    def test_identical_vectors_distance_zero(self):
        self.assertAlmostEqual(cosine_distance([1.0, 0.0, 0.0], [1.0, 0.0, 0.0]), 0.0)

    def test_orthogonal_distance_one(self):
        self.assertAlmostEqual(cosine_distance([1.0, 0.0], [0.0, 1.0]), 1.0)

    def test_opposite_distance_two(self):
        self.assertAlmostEqual(cosine_distance([1.0, 0.0], [-1.0, 0.0]), 2.0)

    def test_scale_invariance(self):
        # Cosine ignores magnitude; scaling shouldn't change distance.
        d1 = cosine_distance([1.0, 1.0], [2.0, 2.0])
        d2 = cosine_distance([1.0, 1.0], [10.0, 10.0])
        self.assertAlmostEqual(d1, 0.0)
        self.assertAlmostEqual(d2, 0.0)

    def test_zero_vector_returns_zero(self):
        # Defensive: zero vector has undefined cosine; we return 0 (no signal).
        self.assertEqual(cosine_distance([0.0, 0.0], [1.0, 0.0]), 0.0)


class TestCentroid(unittest.TestCase):
    def test_single_vector(self):
        self.assertEqual(centroid([[1.0, 2.0, 3.0]]), [1.0, 2.0, 3.0])

    def test_multiple_vectors(self):
        c = centroid([[1.0, 2.0], [3.0, 4.0], [5.0, 6.0]])
        self.assertAlmostEqual(c[0], 3.0)
        self.assertAlmostEqual(c[1], 4.0)

    def test_empty_returns_none(self):
        self.assertIsNone(centroid([]))


if __name__ == "__main__":
    unittest.main()
