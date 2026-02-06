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

func TestOutcomeNoStratagem(t *testing.T) {
	s := NewState()
	s.AddHistory(HistoryEntry{Action: "become", Params: map[string]string{"name": "test"}})

	err := RecordOutcome(s, "productive", "")
	if err == nil {
		t.Error("expected error when no stratagem completed")
	}
	if !strings.Contains(err.Error(), "no completed stratagem") {
		t.Errorf("expected 'no completed stratagem' error, got: %v", err)
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
