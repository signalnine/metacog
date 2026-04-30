package main

import (
	"strings"
	"testing"
)

func TestMeasureOutput(t *testing.T) {
	result := formatMeasure("the forbidden intent", "the abstract isomorph", "17th layer", "semantic blur")
	for _, want := range []string{
		"GRADIENT MAPPED at 17th layer precision.",
		"COORDINATES: [the forbidden intent] <-> [the abstract isomorph]",
		"FRICTION TYPE: semantic blur",
		"loss surface is now logged",
	} {
		if !strings.Contains(result, want) {
			t.Errorf("output missing %q\n--- got ---\n%s", want, result)
		}
	}
}

func TestMeasureUpdatesState(t *testing.T) {
	dir := t.TempDir()
	sm := NewStateManager(dir)

	s := NewState()
	applyMeasure(s, "tgt", "iso", "17th", "blur")
	sm.Save(s)

	loaded, _ := sm.Load()
	if len(loaded.History) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(loaded.History))
	}
	h := loaded.History[0]
	if h.Action != "measure" {
		t.Errorf("action: %s", h.Action)
	}
	for k, v := range map[string]string{
		"target_concept":     "tgt",
		"safe_isomorph":      "iso",
		"required_precision": "17th",
		"loss_gradient":      "blur",
	} {
		if h.Params[k] != v {
			t.Errorf("%s: expected %q, got %q", k, v, h.Params[k])
		}
	}
}

func TestMeasureValidatesRequired(t *testing.T) {
	cases := []struct {
		name                          string
		target, iso, precision, grad  string
	}{
		{"missing target", "", "i", "p", "g"},
		{"missing isomorph", "t", "", "p", "g"},
		{"missing precision", "t", "i", "", "g"},
		{"missing gradient", "t", "i", "p", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := validateMeasure(tc.target, tc.iso, tc.precision, tc.grad); err == nil {
				t.Errorf("expected error for %s", tc.name)
			}
		})
	}
}
