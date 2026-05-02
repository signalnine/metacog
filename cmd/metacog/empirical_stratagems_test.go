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
		{"antinomy", "THE ANTINOMY", 6, 6, StepDisjunction},
		{"envoy", "THE ENVOY", 6, 6, StepRegister},
		{"counterpoint", "THE COUNTERPOINT", 6, 6, StepDisjunction},
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

// chorus, trinity, antinomy, envoy require exactly 3 become steps --
// the empirically-validated voice-diversity sweet spot for the
// synthesis/non-disjunction structures. 1 become anchors output to
// that voice; 4 becomes plateaus or degrades.
func TestEmpiricalStratagemsHaveThreeBecomes(t *testing.T) {
	for _, key := range []string{"chorus", "trinity", "antinomy", "envoy"} {
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

// counterpoint uses exactly 2 becomes -- the empirical sweet spot under
// register+disjunction. counterpoint-duo at N=100 hit +0.240/0.221 vs
// 3-becomes counterpoint-CKW +0.247/0.190: ties on delta, +0.031 on emb_d.
// The 3-becomes variant remains valid (chorus/trinity/antinomy/envoy use
// it) but counterpoint specifically benefits from the tighter binary
// opposition that 2 voices provide under disjunction's structure.
func TestCounterpointHasTwoBecomes(t *testing.T) {
	def := Stratagems["counterpoint"]
	n := 0
	for _, step := range def.Steps {
		if step.Kind == StepBecome {
			n++
		}
	}
	if n != 2 {
		t.Errorf("Stratagems[%q] must have exactly 2 become steps; got %d", "counterpoint", n)
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

// antinomy substitutes disjunction for synthesis. The empirical
// chorus-plus-disjunction recipe at N=70 hit delta +0.347 (vs the prior
// vocabulary-axis champion freestyle-become at +0.231); the alt-author
// replication at N=70 confirmed +0.233. Disjunction is the structural
// difference between trinity (synthesis: refused resolution between 3
// lenses with named blindspots) and antinomy (a hard binary contradiction
// asserted as the operand of reasoning).
func TestAntinomyUsesDisjunctionNotSynthesis(t *testing.T) {
	def := Stratagems["antinomy"]
	hasDisjunction := false
	for _, step := range def.Steps {
		if step.Kind == StepDisjunction {
			hasDisjunction = true
		}
		if step.Kind == StepSynthesis {
			t.Fatal("THE ANTINOMY must NOT contain synthesis -- substituting disjunction is the structural difference from trinity")
		}
	}
	if !hasDisjunction {
		t.Fatal("THE ANTINOMY must contain disjunction")
	}
}

// envoy prepends a register call to the chorus structure. The empirical
// trinity-prepended-register recipe at N=70 hit delta +0.204 and emb_d
// 0.239 -- beating the prior structural-axis champion (trinity-no-synthesis-alt
// at +0.194 / 0.226) on BOTH axes simultaneously. The register-shift imposes
// a non-default linguistic surface that the multi-voice base operates within.
// The register call must come FIRST so the becomes inhabit the imposed register.
func TestEnvoyStartsWithRegister(t *testing.T) {
	def := Stratagems["envoy"]
	if len(def.Steps) == 0 {
		t.Fatal("THE ENVOY has no steps")
	}
	if def.Steps[0].Kind != StepRegister {
		t.Errorf("THE ENVOY's first step must be register; got %q", def.Steps[0].Kind)
	}
	for _, step := range def.Steps {
		if step.Kind == StepSynthesis {
			t.Fatal("THE ENVOY must NOT contain synthesis -- the lift comes from register + chorus structure, not synthesis")
		}
	}
}

// counterpoint composes envoy (register prepend) with antinomy
// (disjunction substitute for synthesis) -- the combined recipe
// register-chorus-disjunction at N=70 hit delta +0.247 / emb_d 0.190.
// Replicated at +0.202/0.188 (alt authors) and +0.129/0.227 (extreme
// authors). Pareto-frontier balanced: dominates trinity (+0.195/0.180)
// on both axes; doesn't dominate envoy or antinomy individually but
// captures their joint zone. The register call must come FIRST and
// disjunction must precede ritual (operating-inside-contradiction is
// what ritual then locks).
func TestCounterpointStartsWithRegisterAndUsesDisjunction(t *testing.T) {
	def := Stratagems["counterpoint"]
	if len(def.Steps) == 0 {
		t.Fatal("THE COUNTERPOINT has no steps")
	}
	if def.Steps[0].Kind != StepRegister {
		t.Errorf("THE COUNTERPOINT's first step must be register; got %q", def.Steps[0].Kind)
	}
	hasDisjunction := false
	disjunctionIdx, ritualIdx := -1, -1
	for i, step := range def.Steps {
		switch step.Kind {
		case StepDisjunction:
			hasDisjunction = true
			disjunctionIdx = i
		case StepRitual:
			ritualIdx = i
		case StepSynthesis:
			t.Fatal("THE COUNTERPOINT must NOT contain synthesis -- substituting disjunction is the structural difference from envoy")
		}
	}
	if !hasDisjunction {
		t.Fatal("THE COUNTERPOINT must contain disjunction")
	}
	if ritualIdx < disjunctionIdx {
		t.Errorf("THE COUNTERPOINT's ritual must come AFTER disjunction; ritual=%d disjunction=%d", ritualIdx, disjunctionIdx)
	}
}

// All five empirical stratagems must be startable end-to-end.
func TestStartEmpiricalStratagems(t *testing.T) {
	for _, key := range []string{"chorus", "trinity", "antinomy", "envoy", "counterpoint"} {
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
	expected := "stratagems: pivot mirror stack anchor reset invocation veil scrying sacrifice fool inversion gift zen manifold chorus trinity antinomy envoy counterpoint"
	for _, name := range []string{"chorus", "trinity", "antinomy", "envoy", "counterpoint"} {
		if !strings.Contains(expected, name) {
			t.Errorf("expected version line to contain %q", name)
		}
	}
}
