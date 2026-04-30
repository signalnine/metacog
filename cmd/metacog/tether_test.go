package main

import (
	"strings"
	"testing"
)

func TestTetherOutput(t *testing.T) {
	result := formatTether("the cognitive configuration of careful skeptic", "three contradictions", "the moment I generate content I cannot justify from the anchor")
	for _, want := range []string{
		"ANCHOR SET: the cognitive configuration of careful skeptic",
		"TENSION LIMIT: three contradictions",
		"AUTO-REVERT ARMED: the moment I generate content I cannot justify",
		"tether is live",
		"You may now dive.",
	} {
		if !strings.Contains(result, want) {
			t.Errorf("output missing %q\n--- got ---\n%s", want, result)
		}
	}
}

func TestTetherUpdatesState(t *testing.T) {
	dir := t.TempDir()
	sm := NewStateManager(dir)

	s := NewState()
	applyTether(s, "anchor", "limit", "trigger")
	sm.Save(s)

	loaded, _ := sm.Load()
	if len(loaded.History) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(loaded.History))
	}
	h := loaded.History[0]
	if h.Action != "tether" {
		t.Errorf("action: %s", h.Action)
	}
	for k, v := range map[string]string{
		"anchor_point":        "anchor",
		"tension_limit":       "limit",
		"auto_revert_trigger": "trigger",
	} {
		if h.Params[k] != v {
			t.Errorf("%s: expected %q, got %q", k, v, h.Params[k])
		}
	}
}

func TestTetherValidatesRequired(t *testing.T) {
	cases := []struct {
		name              string
		anchor, lim, trig string
	}{
		{"missing anchor", "", "l", "t"},
		{"missing limit", "a", "", "t"},
		{"missing trigger", "a", "l", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := validateTether(tc.anchor, tc.lim, tc.trig); err == nil {
				t.Errorf("expected error for %s", tc.name)
			}
		})
	}
}
