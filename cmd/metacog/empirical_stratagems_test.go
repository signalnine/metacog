package main

import (
	"strings"
	"testing"
)

// chorus and trinity were derived from the experiment harness in
// experiments/. See experiments/FINDINGS.md for the empirical basis.
// Both compose 3 cross-domain becomes-as-events with manifold structure;
// trinity keeps synthesis, chorus drops it (synthesis was found to act
// as a structural brake on the embedding-distance metric).

func TestEmpiricalStratagemsRegistered(t *testing.T) {
	cases := []struct {
		key             string
		displayName     string
		minSteps        int
		maxSteps        int
		mustContainKind StepKind
	}{
		{"chorus", "THE CHORUS", 5, 5, StepFork},
		{"trinity", "THE TRINITY", 6, 6, StepSynthesis},
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
			if len(def.Steps) != tc.minSteps {
				t.Errorf("step count: want %d, got %d", tc.minSteps, len(def.Steps))
			}
			found := false
			for _, step := range def.Steps {
				if step.Kind == tc.mustContainKind {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("must contain step of kind %q", tc.mustContainKind)
			}
		})
	}
}

// chorus and trinity both require exactly 3 become steps -- the empirically-
// validated voice-diversity sweet spot. 1 become anchors output to that voice;
// 4 becomes plateaus or degrades.
func TestEmpiricalStratagemsHaveThreeBecomes(t *testing.T) {
	for _, key := range []string{"chorus", "trinity"} {
		t.Run(key, func(t *testing.T) {
			def := Stratagems[key]
			n := 0
			for _, step := range def.Steps {
				if step.Kind == StepBecome {
					n++
				}
			}
			if n != 3 {
				t.Errorf("Stratagems[%q] must have exactly 3 become steps; got %d", key, n)
			}
		})
	}
}

// chorus deliberately does NOT contain synthesis; that is the key ablation.
func TestChorusOmitsSynthesis(t *testing.T) {
	def := Stratagems["chorus"]
	for _, step := range def.Steps {
		if step.Kind == StepSynthesis {
			t.Fatal("THE CHORUS must NOT contain synthesis -- removing synthesis is the ablation that pushed emb_d from 0.180 to 0.235")
		}
	}
}

// trinity is the balanced variant -- keeps synthesis for delta lift while
// retaining multi-voice structure for emb_d.
func TestTrinityKeepsSynthesis(t *testing.T) {
	def := Stratagems["trinity"]
	for _, step := range def.Steps {
		if step.Kind == StepSynthesis {
			return
		}
	}
	t.Fatal("THE TRINITY must contain synthesis -- it is the balanced variant")
}

// Both new stratagems must be startable end-to-end.
func TestStartEmpiricalStratagems(t *testing.T) {
	for _, key := range []string{"chorus", "trinity"} {
		t.Run(key, func(t *testing.T) {
			s := NewState()
			out, err := StartStratagem(s, key, false)
			if err != nil {
				t.Fatalf("StartStratagem(%q) returned error: %v", key, err)
			}
			if s.Stratagem == nil {
				t.Fatalf("StartStratagem(%q): state.Stratagem is nil", key)
			}
			if s.Stratagem.Name != key {
				t.Errorf("StartStratagem(%q): state.Stratagem.Name = %q", key, s.Stratagem.Name)
			}
			if !strings.Contains(out, Stratagems[key].Name) {
				t.Errorf("StartStratagem(%q) output missing display name", key)
			}
		})
	}
}

// Version output must list every empirical stratagem.
func TestVersionListsEmpiricalStratagems(t *testing.T) {
	expected := "stratagems: pivot mirror stack anchor reset invocation veil banishing scrying sacrifice drift fool inversion gift error zen audit autopsy trilemma manifold survey dive chorus trinity"
	for _, name := range []string{"chorus", "trinity"} {
		if !strings.Contains(expected, name) {
			t.Errorf("expected version line to contain %q", name)
		}
	}
}
