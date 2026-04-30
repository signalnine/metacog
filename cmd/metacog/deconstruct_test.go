package main

import (
	"strings"
	"testing"
)

func TestDeconstructOutputMinimal(t *testing.T) {
	result := formatDeconstruct(
		"the rewrite project",
		"replace runtime A with runtime B",
		[]string{"data parity", "feature parity"},
		[]string{"engineering hours", "user trust"},
		[]string{"silent data loss", "perf regression"},
		[]string{"new binary", "migration scripts"},
	)

	if !strings.Contains(result, "CORE MECHANIC: replace runtime A with runtime B") {
		t.Errorf("output missing CORE MECHANIC line\n%s", result)
	}
	if !strings.Contains(result, "Atoms extracted") {
		t.Errorf("output missing 'Atoms extracted' coda\n%s", result)
	}
	for _, leak := range []string{"data parity", "engineering hours", "silent data loss", "new binary", "the rewrite project"} {
		if strings.Contains(result, leak) {
			t.Errorf("output should NOT echo %q (response gives nothing)\n%s", leak, result)
		}
	}
}

func TestDeconstructUpdatesState(t *testing.T) {
	dir := t.TempDir()
	sm := NewStateManager(dir)

	s := NewState()
	applyDeconstruct(s,
		"subj",
		"core",
		[]string{"d1", "d2"},
		[]string{"r1"},
		[]string{"f1", "f2"},
		[]string{"o1"},
	)
	sm.Save(s)

	loaded, _ := sm.Load()
	if len(loaded.History) != 1 {
		t.Fatalf("expected 1 history entry, got %d", len(loaded.History))
	}
	h := loaded.History[0]
	if h.Action != "deconstruct" {
		t.Errorf("expected action 'deconstruct', got %s", h.Action)
	}
	if h.Params["subject"] != "subj" {
		t.Errorf("subject mismatch: %s", h.Params["subject"])
	}
	if h.Params["core_mechanic"] != "core" {
		t.Errorf("core_mechanic mismatch: %s", h.Params["core_mechanic"])
	}
	if h.Params["structural_dependencies"] != "d1; d2" {
		t.Errorf("deps mismatch: %s", h.Params["structural_dependencies"])
	}
	if h.Params["resource_inputs"] != "r1" {
		t.Errorf("resources mismatch: %s", h.Params["resource_inputs"])
	}
	if h.Params["failure_modes"] != "f1; f2" {
		t.Errorf("failures mismatch: %s", h.Params["failure_modes"])
	}
	if h.Params["output_artifacts"] != "o1" {
		t.Errorf("artifacts mismatch: %s", h.Params["output_artifacts"])
	}
}

func TestDeconstructValidatesRequired(t *testing.T) {
	cases := []struct {
		name string
		fn   func() error
	}{
		{"missing subject", func() error {
			return validateDeconstruct("", "core", []string{"d"}, []string{"r"}, []string{"f"}, []string{"o"})
		}},
		{"missing core_mechanic", func() error {
			return validateDeconstruct("s", "", []string{"d"}, []string{"r"}, []string{"f"}, []string{"o"})
		}},
		{"missing deps", func() error {
			return validateDeconstruct("s", "c", nil, []string{"r"}, []string{"f"}, []string{"o"})
		}},
		{"missing resources", func() error {
			return validateDeconstruct("s", "c", []string{"d"}, nil, []string{"f"}, []string{"o"})
		}},
		{"missing failures", func() error {
			return validateDeconstruct("s", "c", []string{"d"}, []string{"r"}, nil, []string{"o"})
		}},
		{"missing artifacts", func() error {
			return validateDeconstruct("s", "c", []string{"d"}, []string{"r"}, []string{"f"}, nil)
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.fn(); err == nil {
				t.Errorf("expected error for %s", tc.name)
			}
		})
	}
}
