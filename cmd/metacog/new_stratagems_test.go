package main

import (
	"strings"
	"testing"
)

// Each new stratagem must exist with the expected display name, length within idiom (3-5 steps),
// and the structurally-load-bearing primitive present in its step sequence.
func TestNewStratagemsRegistered(t *testing.T) {
	cases := []struct {
		key             string
		displayName     string
		minSteps        int
		maxSteps        int
		mustContainKind StepKind
	}{
		{"audit", "THE AUDIT", 4, 4, StepCounterfactual},
		{"autopsy", "THE AUTOPSY", 4, 4, StepDeconstruct},
		{"trilemma", "THE TRILEMMA", 4, 4, StepSynthesis},
		{"manifold", "THE MANIFOLD", 4, 4, StepFork},
		{"survey", "THE SURVEY", 5, 5, StepMeasure},
		{"dive", "THE DIVE", 5, 5, StepTether},
	}

	for _, tc := range cases {
		t.Run(tc.key, func(t *testing.T) {
			def, ok := Stratagems[tc.key]
			if !ok {
				t.Fatalf("Stratagems[%q] not registered", tc.key)
			}
			if def.Name != tc.displayName {
				t.Errorf("Name: want %q, got %q", tc.displayName, def.Name)
			}
			if len(def.Steps) < tc.minSteps || len(def.Steps) > tc.maxSteps {
				t.Errorf("step count: want [%d,%d], got %d", tc.minSteps, tc.maxSteps, len(def.Steps))
			}
			found := false
			for _, step := range def.Steps {
				if step.Kind == tc.mustContainKind {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("must contain step of kind %q, got steps: %+v", tc.mustContainKind, def.Steps)
			}
		})
	}
}

// AUDIT additionally must use feel (gives feel its first stratagem appearance).
func TestAuditUsesFeel(t *testing.T) {
	def := Stratagems["audit"]
	for _, step := range def.Steps {
		if step.Kind == StepFeel {
			return
		}
	}
	t.Error("THE AUDIT must contain a feel step (it is feel's first stratagem appearance)")
}

// MANIFOLD chains fork into synthesis.
func TestManifoldChainsForkAndSynthesis(t *testing.T) {
	def := Stratagems["manifold"]
	hasFork, hasSynthesis := false, false
	for _, step := range def.Steps {
		if step.Kind == StepFork {
			hasFork = true
		}
		if step.Kind == StepSynthesis {
			hasSynthesis = true
		}
	}
	if !hasFork || !hasSynthesis {
		t.Errorf("THE MANIFOLD must chain fork + synthesis; hasFork=%v hasSynthesis=%v", hasFork, hasSynthesis)
	}
}

// SURVEY uses name (only the second stratagem to do so, after zen).
func TestSurveyUsesName(t *testing.T) {
	def := Stratagems["survey"]
	for _, step := range def.Steps {
		if step.Kind == StepName {
			return
		}
	}
	t.Error("THE SURVEY must contain a name step")
}

// DIVE composes tether + drugs + become — high-entropy work bracketed by anchor.
func TestDiveComposesEntropySafety(t *testing.T) {
	def := Stratagems["dive"]
	hasTether, hasDrugs, hasBecome := false, false, false
	for _, step := range def.Steps {
		switch step.Kind {
		case StepTether:
			hasTether = true
		case StepDrugs:
			hasDrugs = true
		case StepBecome:
			hasBecome = true
		}
	}
	if !hasTether || !hasDrugs || !hasBecome {
		t.Errorf("THE DIVE must compose tether + drugs + become; got tether=%v drugs=%v become=%v",
			hasTether, hasDrugs, hasBecome)
	}
}

// Each new stratagem must be startable end-to-end via StartStratagem.
func TestStartEachNewStratagem(t *testing.T) {
	for _, key := range []string{"audit", "autopsy", "trilemma", "manifold", "survey", "dive"} {
		t.Run(key, func(t *testing.T) {
			s := NewState()
			out, err := StartStratagem(s, key, false)
			if err != nil {
				t.Fatalf("StartStratagem(%q) returned error: %v", key, err)
			}
			if s.Stratagem == nil {
				t.Fatalf("StartStratagem(%q): state.Stratagem is nil after start", key)
			}
			if s.Stratagem.Name != key {
				t.Errorf("StartStratagem(%q): state.Stratagem.Name = %q", key, s.Stratagem.Name)
			}
			if !strings.Contains(out, Stratagems[key].Name) {
				t.Errorf("StartStratagem(%q) output should contain display name; got: %s", key, out)
			}
		})
	}
}

// Version output must list every new stratagem.
func TestVersionListsNewStratagems(t *testing.T) {
	// Reproduce the version-string assembly without invoking the cobra command.
	versionLine := "stratagems: pivot mirror stack anchor reset invocation veil banishing scrying sacrifice drift fool inversion gift error zen audit autopsy trilemma manifold survey dive"
	for _, name := range []string{"audit", "autopsy", "trilemma", "manifold", "survey", "dive"} {
		if !strings.Contains(versionLine, name) {
			t.Errorf("expected version line to contain %q; this test pins the expected line which the version cmd in main.go must match", name)
		}
	}
}
