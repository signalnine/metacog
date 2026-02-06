package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestNewState(t *testing.T) {
	s := NewState()
	if s.Version != StateSchemaVersion {
		t.Errorf("expected version %d, got %d", StateSchemaVersion, s.Version)
	}
	if s.SessionID == "" {
		t.Error("expected non-empty session ID")
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	sm := NewStateManager(dir)

	s := NewState()
	s.Identity = &Identity{Name: "Ada", Lens: "verification", Env: "lab"}

	err := sm.Save(s)
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := sm.Load()
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if loaded.Identity == nil || loaded.Identity.Name != "Ada" {
		t.Error("identity not persisted")
	}
}

func TestAtomicWrite(t *testing.T) {
	dir := t.TempDir()
	sm := NewStateManager(dir)

	s := NewState()
	err := sm.Save(s)
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}

	// Temp file should not exist after save
	tmpPath := filepath.Join(dir, ".state.json.tmp")
	if _, err := os.Stat(tmpPath); !os.IsNotExist(err) {
		t.Error("temp file should not exist after save")
	}

	// State file should exist
	statePath := filepath.Join(dir, "state.json")
	if _, err := os.Stat(statePath); err != nil {
		t.Error("state file should exist after save")
	}
}

func TestLoadMissingFile(t *testing.T) {
	dir := t.TempDir()
	sm := NewStateManager(dir)

	s, err := sm.Load()
	if err != nil {
		t.Fatalf("load of missing file should not error: %v", err)
	}
	if s.Version != StateSchemaVersion {
		t.Error("missing file should return fresh state")
	}
}

func TestLoadCorruptedFile(t *testing.T) {
	dir := t.TempDir()
	statePath := filepath.Join(dir, "state.json")
	os.WriteFile(statePath, []byte("not json{{{"), 0644)

	sm := NewStateManager(dir)
	_, err := sm.Load()
	if err == nil {
		t.Error("corrupted file should return error")
	}
}

func TestLoadFutureVersion(t *testing.T) {
	dir := t.TempDir()
	statePath := filepath.Join(dir, "state.json")
	data, _ := json.Marshal(map[string]any{"version": 999})
	os.WriteFile(statePath, data, 0644)

	sm := NewStateManager(dir)
	_, err := sm.Load()
	if err == nil {
		t.Error("future version should return error")
	}
}

func TestConcurrentAccess(t *testing.T) {
	dir := t.TempDir()
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			sm := NewStateManager(dir)
			s, _ := sm.Load()
			s.Identity = &Identity{Name: "writer", Lens: "concurrent", Env: "test"}
			sm.Save(s)
		}(i)
	}
	wg.Wait()

	// File should be valid JSON after concurrent writes
	sm := NewStateManager(dir)
	s, err := sm.Load()
	if err != nil {
		t.Fatalf("state corrupted after concurrent writes: %v", err)
	}
	if s.Identity == nil {
		t.Error("identity should be set")
	}
}

func TestHistoryRotation(t *testing.T) {
	dir := t.TempDir()
	sm := NewStateManager(dir)

	s := NewState()
	for i := 0; i < 150; i++ {
		s.AddHistory(HistoryEntry{Action: "become", Params: map[string]string{"i": "test"}})
	}

	err := sm.Save(s)
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := sm.Load()
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if len(loaded.History) > 100 {
		t.Errorf("history should cap at 100, got %d", len(loaded.History))
	}
}

func TestRepair(t *testing.T) {
	dir := t.TempDir()
	statePath := filepath.Join(dir, "state.json")
	os.WriteFile(statePath, []byte("corrupt"), 0644)

	sm := NewStateManager(dir)
	err := sm.Repair()
	if err != nil {
		t.Fatalf("repair failed: %v", err)
	}

	s, err := sm.Load()
	if err != nil {
		t.Fatalf("load after repair failed: %v", err)
	}
	if s.Version != StateSchemaVersion {
		t.Error("repaired state should be fresh")
	}
}
