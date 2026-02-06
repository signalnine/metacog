package main

import (
	"testing"
)

func TestFormatStatus(t *testing.T) {
	s := NewState()
	s.Identity = &Identity{Name: "Ada", Lens: "verification", Env: "lab"}
	s.Substrate = &Substrate{Substance: "caffeine", Method: "antagonism", Qualia: "sharp"}

	output := FormatStatus(s)
	if output == "" {
		t.Error("expected non-empty status")
	}
}

func TestFormatStatusEmpty(t *testing.T) {
	s := NewState()
	output := FormatStatus(s)
	if output == "" {
		t.Error("expected non-empty status even for empty state")
	}
}

func TestFormatHistory(t *testing.T) {
	s := NewState()
	s.AddHistory(HistoryEntry{Action: "become", Params: map[string]string{"name": "Ada"}})
	s.AddHistory(HistoryEntry{Action: "drugs", Params: map[string]string{"substance": "caffeine"}})

	output := FormatHistory(s)
	if output == "" {
		t.Error("expected non-empty history")
	}
}

func TestFormatHistoryEmpty(t *testing.T) {
	s := NewState()
	output := FormatHistory(s)
	if output == "" {
		t.Error("expected non-empty output for empty history")
	}
}
