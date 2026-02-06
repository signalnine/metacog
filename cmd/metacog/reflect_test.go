package main

import (
	"strings"
	"testing"
)

func TestReflectEmpty(t *testing.T) {
	s := NewState()
	output := FormatReflection(s)
	if !strings.Contains(output, "No history") {
		t.Error("empty history should say so")
	}
}

func TestReflectPrimitiveCounts(t *testing.T) {
	s := NewState()
	for i := 0; i < 5; i++ {
		s.AddHistory(HistoryEntry{Action: "become", Params: map[string]string{"name": "Ada"}})
	}
	for i := 0; i < 3; i++ {
		s.AddHistory(HistoryEntry{Action: "drugs", Params: map[string]string{"substance": "caffeine"}})
	}
	s.AddHistory(HistoryEntry{Action: "ritual", Params: map[string]string{"threshold": "test"}})

	output := FormatReflection(s)
	if !strings.Contains(output, "become: 5") {
		t.Errorf("expected become: 5 in output:\n%s", output)
	}
	if !strings.Contains(output, "drugs: 3") {
		t.Errorf("expected drugs: 3 in output:\n%s", output)
	}
	if !strings.Contains(output, "ritual: 1") {
		t.Errorf("expected ritual: 1 in output:\n%s", output)
	}
}

func TestReflectTopIdentities(t *testing.T) {
	s := NewState()
	for i := 0; i < 4; i++ {
		s.AddHistory(HistoryEntry{Action: "become", Params: map[string]string{"name": "Ada"}})
	}
	for i := 0; i < 2; i++ {
		s.AddHistory(HistoryEntry{Action: "become", Params: map[string]string{"name": "Eno"}})
	}
	s.AddHistory(HistoryEntry{Action: "become", Params: map[string]string{"name": "Doepfer"}})

	output := FormatReflection(s)
	if !strings.Contains(output, "Ada") {
		t.Errorf("expected Ada in top identities:\n%s", output)
	}
}

func TestReflectStratagemUsage(t *testing.T) {
	s := NewState()
	s.AddHistory(HistoryEntry{Action: "stratagem", Params: map[string]string{"name": "pivot", "event": "completed"}})
	s.AddHistory(HistoryEntry{Action: "stratagem", Params: map[string]string{"name": "pivot", "event": "completed"}})
	s.AddHistory(HistoryEntry{Action: "stratagem", Params: map[string]string{"name": "pivot", "event": "completed"}})
	s.AddHistory(HistoryEntry{Action: "stratagem", Params: map[string]string{"name": "stack", "event": "completed"}})
	s.AddHistory(HistoryEntry{Action: "stratagem", Params: map[string]string{"name": "mirror", "event": "started"}})

	output := FormatReflection(s)
	if !strings.Contains(output, "pivot: 3") {
		t.Errorf("expected 'pivot: 3' in output:\n%s", output)
	}
	if !strings.Contains(output, "stack: 1") {
		t.Errorf("expected 'stack: 1' in output:\n%s", output)
	}
	if !strings.Contains(output, "anchor") {
		t.Errorf("expected 'anchor' in never-completed list:\n%s", output)
	}
	if !strings.Contains(output, "mirror") {
		t.Errorf("expected 'mirror' in never-completed list:\n%s", output)
	}
}

func TestReflectRitualStepAverage(t *testing.T) {
	s := NewState()
	s.AddHistory(HistoryEntry{Action: "ritual", Params: map[string]string{"steps": "a; b; c"}})
	s.AddHistory(HistoryEntry{Action: "ritual", Params: map[string]string{"steps": "a; b; c; d; e"}})

	output := FormatReflection(s)
	if !strings.Contains(output, "4.0") {
		t.Errorf("expected average 4.0 steps in output:\n%s", output)
	}
}

func TestReflectEffectiveness(t *testing.T) {
	s := NewState()

	// pivot: 2 productive, 1 unproductive = 66%
	s.AddHistory(HistoryEntry{Action: "stratagem", Params: map[string]string{"name": "pivot", "event": "completed"}})
	s.AddHistory(HistoryEntry{Action: "outcome", Params: map[string]string{"result": "productive", "stratagem": "pivot", "shift": "reframed"}})
	s.AddHistory(HistoryEntry{Action: "stratagem", Params: map[string]string{"name": "pivot", "event": "completed"}})
	s.AddHistory(HistoryEntry{Action: "outcome", Params: map[string]string{"result": "productive", "stratagem": "pivot"}})
	s.AddHistory(HistoryEntry{Action: "stratagem", Params: map[string]string{"name": "pivot", "event": "completed"}})
	s.AddHistory(HistoryEntry{Action: "outcome", Params: map[string]string{"result": "unproductive", "stratagem": "pivot"}})

	// stack: 1 productive = 100% [provisional]
	s.AddHistory(HistoryEntry{Action: "stratagem", Params: map[string]string{"name": "stack", "event": "completed"}})
	s.AddHistory(HistoryEntry{Action: "outcome", Params: map[string]string{"result": "productive", "stratagem": "stack"}})

	// mirror: completed but no outcome = unmeasured
	s.AddHistory(HistoryEntry{Action: "stratagem", Params: map[string]string{"name": "mirror", "event": "completed"}})

	output := FormatReflection(s)

	if !strings.Contains(output, "self-reported") {
		t.Errorf("expected 'self-reported' framing in output:\n%s", output)
	}
	if !strings.Contains(output, "pivot: 66%") && !strings.Contains(output, "pivot: 67%") {
		t.Errorf("expected 'pivot: 66%%' or 'pivot: 67%%' in output:\n%s", output)
	}
	if !strings.Contains(output, "stack: 100%") {
		t.Errorf("expected 'stack: 100%%' in output:\n%s", output)
	}
	if !strings.Contains(output, "[provisional]") {
		t.Errorf("expected '[provisional]' tag for stack:\n%s", output)
	}
	if !strings.Contains(output, "unmeasured") {
		t.Errorf("expected 'unmeasured' for mirror:\n%s", output)
	}
	if !strings.Contains(output, "Overall:") {
		t.Errorf("expected Overall rate in output:\n%s", output)
	}
}

func TestReflectNoOutcomes(t *testing.T) {
	s := NewState()
	s.AddHistory(HistoryEntry{Action: "become", Params: map[string]string{"name": "test"}})

	output := FormatReflection(s)
	// Should not contain effectiveness section when no outcomes exist
	if strings.Contains(output, "effectiveness") {
		t.Errorf("should not show effectiveness with no outcomes:\n%s", output)
	}
}
