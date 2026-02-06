package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestSavePersonalStance(t *testing.T) {
	dir := t.TempDir()

	s := &State{
		Identity:  &Identity{Name: "Ada", Lens: "verification", Env: "lab"},
		Substrate: &Substrate{Substance: "caffeine", Method: "antagonism", Qualia: "sharp"},
	}

	err := SavePersonalStance(dir, s)
	if err != nil {
		t.Fatalf("save personal stance failed: %v", err)
	}

	poolPath := filepath.Join(dir, "stances", "personal.json")
	data, err := os.ReadFile(poolPath)
	if err != nil {
		t.Fatalf("personal pool file should exist: %v", err)
	}

	var stances []PersonalStance
	if err := json.Unmarshal(data, &stances); err != nil {
		t.Fatalf("personal pool should be valid JSON: %v", err)
	}

	if len(stances) != 1 {
		t.Fatalf("expected 1 stance, got %d", len(stances))
	}
	if stances[0].Who != "Ada" {
		t.Errorf("expected who=Ada, got %q", stances[0].Who)
	}
	if stances[0].Where != "lab" {
		t.Errorf("expected where=lab, got %q", stances[0].Where)
	}
	if stances[0].Lens != "verification" {
		t.Errorf("expected lens=verification, got %q", stances[0].Lens)
	}
	if stances[0].Substance != "caffeine" {
		t.Errorf("expected substance=caffeine, got %q", stances[0].Substance)
	}
}

func TestSavePersonalStanceAppends(t *testing.T) {
	dir := t.TempDir()

	s1 := &State{
		Identity: &Identity{Name: "Ada", Lens: "verification", Env: "lab"},
	}
	s2 := &State{
		Identity: &Identity{Name: "Eno", Lens: "generative", Env: "studio"},
	}

	SavePersonalStance(dir, s1)
	SavePersonalStance(dir, s2)

	poolPath := filepath.Join(dir, "stances", "personal.json")
	data, _ := os.ReadFile(poolPath)
	var stances []PersonalStance
	json.Unmarshal(data, &stances)

	if len(stances) != 2 {
		t.Fatalf("expected 2 stances, got %d", len(stances))
	}
}

func TestSavePersonalStanceDedup(t *testing.T) {
	dir := t.TempDir()

	s := &State{
		Identity: &Identity{Name: "Ada", Lens: "verification", Env: "lab"},
	}

	SavePersonalStance(dir, s)
	SavePersonalStance(dir, s)

	poolPath := filepath.Join(dir, "stances", "personal.json")
	data, _ := os.ReadFile(poolPath)
	var stances []PersonalStance
	json.Unmarshal(data, &stances)

	if len(stances) != 1 {
		t.Fatalf("expected 1 stance (deduped), got %d", len(stances))
	}
}

func TestSavePersonalStanceNoIdentity(t *testing.T) {
	dir := t.TempDir()
	s := &State{}

	err := SavePersonalStance(dir, s)
	if err == nil {
		t.Error("expected error when no identity is set")
	}
}

func TestLoadPersonalPool(t *testing.T) {
	dir := t.TempDir()

	s := &State{
		Identity:  &Identity{Name: "Ada", Lens: "verification", Env: "lab"},
		Substrate: &Substrate{Substance: "caffeine", Method: "antagonism", Qualia: "sharp"},
	}
	SavePersonalStance(dir, s)

	pools, err := LoadStancePoolsWithPersonal(dir)
	if err != nil {
		t.Fatalf("load pools failed: %v", err)
	}

	pool, ok := pools["personal"]
	if !ok {
		t.Fatal("personal pool should be loaded")
	}
	if len(pool.Stances) != 1 {
		t.Errorf("expected 1 personal stance, got %d", len(pool.Stances))
	}
	if pool.Stances[0].Who != "Ada" {
		t.Errorf("expected who=Ada, got %q", pool.Stances[0].Who)
	}
}
