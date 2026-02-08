package main

import (
	"strings"
	"testing"
)

func TestOutcomeRecordsResult(t *testing.T) {
	s := NewState()
	s.AddHistory(HistoryEntry{Action: "stratagem", Params: map[string]string{"name": "pivot", "event": "completed"}})

	err := RecordOutcome(s, "productive", "reframed the problem")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Find outcome entry
	found := false
	for _, h := range s.History {
		if h.Action == "outcome" {
			found = true
			if h.Params["result"] != "productive" {
				t.Errorf("expected result=productive, got %s", h.Params["result"])
			}
			if h.Params["shift"] != "reframed the problem" {
				t.Errorf("expected shift text, got %s", h.Params["shift"])
			}
			if h.Params["stratagem"] != "pivot" {
				t.Errorf("expected stratagem=pivot, got %s", h.Params["stratagem"])
			}
		}
	}
	if !found {
		t.Error("outcome entry not found in history")
	}
}

func TestOutcomeUnproductive(t *testing.T) {
	s := NewState()
	s.AddHistory(HistoryEntry{Action: "stratagem", Params: map[string]string{"name": "drift", "event": "completed"}})

	err := RecordOutcome(s, "unproductive", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, h := range s.History {
		if h.Action == "outcome" {
			if h.Params["result"] != "unproductive" {
				t.Errorf("expected unproductive, got %s", h.Params["result"])
			}
			if h.Params["shift"] != "" {
				t.Errorf("expected empty shift, got %s", h.Params["shift"])
			}
			if h.Params["stratagem"] != "drift" {
				t.Errorf("expected stratagem=drift, got %s", h.Params["stratagem"])
			}
		}
	}
}

func TestOutcomeAutoCaptures(t *testing.T) {
	s := NewState()
	// Multiple stratagems — should capture the most recent
	s.AddHistory(HistoryEntry{Action: "stratagem", Params: map[string]string{"name": "pivot", "event": "completed"}})
	s.AddHistory(HistoryEntry{Action: "become", Params: map[string]string{"name": "test"}})
	s.AddHistory(HistoryEntry{Action: "stratagem", Params: map[string]string{"name": "mirror", "event": "completed"}})

	err := RecordOutcome(s, "productive", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, h := range s.History {
		if h.Action == "outcome" {
			if h.Params["stratagem"] != "mirror" {
				t.Errorf("expected mirror (most recent), got %s", h.Params["stratagem"])
			}
		}
	}
}

func TestOutcomeFreestyle(t *testing.T) {
	s := NewState()
	// Freestyle: primitives without a stratagem
	s.AddHistory(HistoryEntry{Action: "become", Params: map[string]string{"name": "Ada"}})
	s.AddHistory(HistoryEntry{Action: "drugs", Params: map[string]string{"substance": "caffeine"}})

	err := RecordOutcome(s, "productive", "reframed via freestyle")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, h := range s.History {
		if h.Action == "outcome" {
			found = true
			if h.Params["stratagem"] != "freestyle" {
				t.Errorf("expected stratagem=freestyle, got %s", h.Params["stratagem"])
			}
			if h.Params["result"] != "productive" {
				t.Errorf("expected result=productive, got %s", h.Params["result"])
			}
			if h.Params["shift"] != "reframed via freestyle" {
				t.Errorf("expected shift text, got %s", h.Params["shift"])
			}
		}
	}
	if !found {
		t.Error("outcome entry not found in history")
	}
}

func TestOutcomeFreestyleAfterStratagem(t *testing.T) {
	s := NewState()
	// Completed stratagem with outcome already recorded
	s.AddHistory(HistoryEntry{Action: "stratagem", Params: map[string]string{"name": "pivot", "event": "started"}})
	s.AddHistory(HistoryEntry{Action: "drugs", Params: map[string]string{"substance": "test"}})
	s.AddHistory(HistoryEntry{Action: "stratagem", Params: map[string]string{"name": "pivot", "event": "completed"}})
	s.AddHistory(HistoryEntry{Action: "outcome", Params: map[string]string{"result": "productive", "stratagem": "pivot"}})
	// Then freestyle primitives
	s.AddHistory(HistoryEntry{Action: "become", Params: map[string]string{"name": "Ada"}})

	err := RecordOutcome(s, "unproductive", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should be a freestyle outcome, not duplicate stratagem
	outcomeCount := 0
	for _, h := range s.History {
		if h.Action == "outcome" {
			outcomeCount++
			if outcomeCount == 2 {
				if h.Params["stratagem"] != "freestyle" {
					t.Errorf("expected second outcome as freestyle, got %s", h.Params["stratagem"])
				}
			}
		}
	}
	if outcomeCount != 2 {
		t.Errorf("expected 2 outcomes, got %d", outcomeCount)
	}
}

func TestOutcomeNoPrimitives(t *testing.T) {
	s := NewState()
	// Empty history — no primitives at all

	err := RecordOutcome(s, "productive", "")
	if err == nil {
		t.Error("expected error when no primitives or stratagems")
	}
}

func TestOutcomeFreestyleIgnoresStratagemPrimitives(t *testing.T) {
	s := NewState()
	// Primitives inside a stratagem span should not count as freestyle
	s.AddHistory(HistoryEntry{Action: "stratagem", Params: map[string]string{"name": "pivot", "event": "started"}})
	s.AddHistory(HistoryEntry{Action: "drugs", Params: map[string]string{"substance": "test"}})
	s.AddHistory(HistoryEntry{Action: "become", Params: map[string]string{"name": "test"}})
	s.AddHistory(HistoryEntry{Action: "stratagem", Params: map[string]string{"name": "pivot", "event": "completed"}})
	s.AddHistory(HistoryEntry{Action: "outcome", Params: map[string]string{"result": "productive", "stratagem": "pivot"}})

	err := RecordOutcome(s, "productive", "")
	if err == nil {
		t.Error("expected error — no freestyle primitives outside stratagem")
	}
}

func TestOutcomeDuplicate(t *testing.T) {
	s := NewState()
	s.AddHistory(HistoryEntry{Action: "stratagem", Params: map[string]string{"name": "pivot", "event": "completed"}})

	RecordOutcome(s, "productive", "first")

	err := RecordOutcome(s, "unproductive", "second")
	if err == nil {
		t.Error("expected error on duplicate outcome")
	}
	if !strings.Contains(err.Error(), "already recorded") {
		t.Errorf("expected 'already recorded' error, got: %v", err)
	}
}

func TestOutcomeAmend(t *testing.T) {
	s := NewState()
	s.AddHistory(HistoryEntry{Action: "stratagem", Params: map[string]string{"name": "pivot", "event": "completed"}})
	RecordOutcome(s, "productive", "initial assessment")

	err := AmendOutcome(s, "unproductive", "on reflection, no real shift")
	if err != nil {
		t.Fatalf("amend failed: %v", err)
	}

	// Find the outcome — should be updated
	for _, h := range s.History {
		if h.Action == "outcome" {
			if h.Params["result"] != "unproductive" {
				t.Errorf("expected amended result=unproductive, got %s", h.Params["result"])
			}
			if h.Params["shift"] != "on reflection, no real shift" {
				t.Errorf("expected amended shift, got %s", h.Params["shift"])
			}
		}
	}
}

func TestOutcomeAmendNoExisting(t *testing.T) {
	s := NewState()
	err := AmendOutcome(s, "productive", "nothing to amend")
	if err == nil {
		t.Error("expected error when no outcome to amend")
	}
}

func TestOutcomeInvalidResult(t *testing.T) {
	s := NewState()
	s.AddHistory(HistoryEntry{Action: "stratagem", Params: map[string]string{"name": "pivot", "event": "completed"}})

	err := RecordOutcome(s, "maybe", "")
	if err == nil {
		t.Error("expected error for invalid result value")
	}
}
