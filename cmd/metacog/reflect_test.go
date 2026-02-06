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
