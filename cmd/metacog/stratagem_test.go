package main

import (
	"testing"
)

func TestStratagemDefinitions(t *testing.T) {
	names := []string{"pivot", "mirror", "stack", "anchor", "reset", "invocation", "veil", "banishing", "scrying", "sacrifice"}
	for _, name := range names {
		def, ok := Stratagems[name]
		if !ok {
			t.Errorf("stratagem %q not defined", name)
			continue
		}
		if len(def.Steps) == 0 {
			t.Errorf("stratagem %q has no steps", name)
		}
	}
}

func TestPivotStepSequence(t *testing.T) {
	def := Stratagems["pivot"]
	expected := []StepKind{StepDrugs, StepThink, StepBecome, StepThink, StepRitual}
	if len(def.Steps) != len(expected) {
		t.Fatalf("pivot: expected %d steps, got %d", len(expected), len(def.Steps))
	}
	for i, step := range def.Steps {
		if step.Kind != expected[i] {
			t.Errorf("pivot step %d: expected %v, got %v", i+1, expected[i], step.Kind)
		}
	}
}

func TestStartStratagem(t *testing.T) {
	s := NewState()
	output, err := StartStratagem(s, "pivot", false)
	if err != nil {
		t.Fatalf("start pivot failed: %v", err)
	}
	if s.Stratagem == nil {
		t.Fatal("stratagem should be active")
	}
	if s.Stratagem.Name != "pivot" {
		t.Errorf("expected pivot, got %s", s.Stratagem.Name)
	}
	if s.Stratagem.Step != 0 {
		t.Errorf("expected step 0, got %d", s.Stratagem.Step)
	}
	if output == "" {
		t.Error("expected non-empty output")
	}
}

func TestStartStratagemWhileActive(t *testing.T) {
	s := NewState()
	StartStratagem(s, "pivot", false)

	_, err := StartStratagem(s, "mirror", false)
	if err == nil {
		t.Error("expected error starting second stratagem without force")
	}
}

func TestStartStratagemForce(t *testing.T) {
	s := NewState()
	StartStratagem(s, "pivot", false)

	_, err := StartStratagem(s, "mirror", true)
	if err != nil {
		t.Fatalf("force start failed: %v", err)
	}
	if s.Stratagem.Name != "mirror" {
		t.Error("expected mirror after force")
	}
	// Check pivot was recorded as abandoned
	found := false
	for _, h := range s.History {
		if h.Action == "stratagem" && h.Status == "abandoned" {
			found = true
			break
		}
	}
	if !found {
		t.Error("abandoned stratagem should be in history")
	}
}

func TestAdvanceThinkStep(t *testing.T) {
	s := NewState()
	StartStratagem(s, "pivot", false)

	// Step 0 is drugs — simulate calling drugs
	s.Stratagem.StepsCompleted = append(s.Stratagem.StepsCompleted, "drugs")

	// Advance past drugs step
	_, err := AdvanceStratagem(s)
	if err != nil {
		t.Fatalf("advance past drugs failed: %v", err)
	}

	// Now at step 1 (THINK) — should advance freely
	_, err = AdvanceStratagem(s)
	if err != nil {
		t.Fatalf("advance past THINK failed: %v", err)
	}

	// Now at step 2 (become)
	if s.Stratagem.Step != 2 {
		t.Errorf("expected step 2, got %d", s.Stratagem.Step)
	}
}

func TestAdvanceRequiresPrimitive(t *testing.T) {
	s := NewState()
	StartStratagem(s, "pivot", false)

	// Step 0 is drugs — try to advance without calling drugs
	_, err := AdvanceStratagem(s)
	if err == nil {
		t.Error("expected error advancing without calling drugs first")
	}
}

func TestCompleteStratagem(t *testing.T) {
	s := NewState()
	StartStratagem(s, "reset", false)

	// Reset: ritual, THINK, ritual
	// Step 0: ritual
	s.Stratagem.StepsCompleted = append(s.Stratagem.StepsCompleted, "ritual")
	AdvanceStratagem(s)

	// Step 1: THINK — advances freely
	AdvanceStratagem(s)

	// Step 2: ritual
	s.Stratagem.StepsCompleted = append(s.Stratagem.StepsCompleted, "ritual")
	output, err := AdvanceStratagem(s)
	if err != nil {
		t.Fatalf("final advance failed: %v", err)
	}

	if s.Stratagem != nil {
		t.Error("stratagem should be cleared after completion")
	}
	if output == "" {
		t.Error("completion should produce output")
	}
}

func TestAbortStratagem(t *testing.T) {
	s := NewState()
	StartStratagem(s, "pivot", false)

	err := AbortStratagem(s)
	if err != nil {
		t.Fatalf("abort failed: %v", err)
	}
	if s.Stratagem != nil {
		t.Error("stratagem should be nil after abort")
	}
}

func TestAbortNoStratagem(t *testing.T) {
	s := NewState()
	err := AbortStratagem(s)
	if err == nil {
		t.Error("expected error aborting with no active stratagem")
	}
}

func TestStratagemStatus(t *testing.T) {
	s := NewState()
	StartStratagem(s, "pivot", false)

	output := StratagemStatus(s)
	if output == "" {
		t.Error("expected non-empty status output")
	}
}

func TestUnknownStratagem(t *testing.T) {
	s := NewState()
	_, err := StartStratagem(s, "nonexistent", false)
	if err == nil {
		t.Error("expected error for unknown stratagem")
	}
}
