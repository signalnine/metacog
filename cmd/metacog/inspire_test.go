package main

import (
	"strings"
	"testing"
)

func TestLoadStances(t *testing.T) {
	pools, err := LoadStancePools()
	if err != nil {
		t.Fatalf("LoadStancePools failed: %v", err)
	}
	if len(pools) == 0 {
		t.Fatal("expected at least one pool")
	}
	for name, pool := range pools {
		if len(pool.Stances) == 0 {
			t.Errorf("pool %q has no stances", name)
		}
		for i, s := range pool.Stances {
			if s.Who == "" || s.Where == "" || s.Lens == "" {
				t.Errorf("pool %q stance %d has empty field: %+v", name, i, s)
			}
		}
	}
}

func TestListPools(t *testing.T) {
	pools, err := LoadStancePools()
	if err != nil {
		t.Fatalf("LoadStancePools failed: %v", err)
	}
	names := ListPoolNames(pools)
	if len(names) == 0 {
		t.Fatal("expected at least one pool name")
	}
	// Verify sorted
	for i := 1; i < len(names); i++ {
		if names[i] < names[i-1] {
			t.Errorf("pool names not sorted: %q before %q", names[i-1], names[i])
		}
	}
}

func TestRandomStance(t *testing.T) {
	pools, err := LoadStancePools()
	if err != nil {
		t.Fatalf("LoadStancePools failed: %v", err)
	}
	stance, pool, err := RandomStance(pools, "")
	if err != nil {
		t.Fatalf("RandomStance failed: %v", err)
	}
	if stance == nil {
		t.Fatal("expected a stance, got nil")
	}
	if pool == "" {
		t.Fatal("expected a pool name, got empty string")
	}
	if stance.Who == "" || stance.Where == "" || stance.Lens == "" {
		t.Errorf("stance has empty field: %+v", stance)
	}
}

func TestRandomStanceFromPool(t *testing.T) {
	pools, err := LoadStancePools()
	if err != nil {
		t.Fatalf("LoadStancePools failed: %v", err)
	}
	stance, pool, err := RandomStance(pools, "philosophy")
	if err != nil {
		t.Fatalf("RandomStance from pool failed: %v", err)
	}
	if stance == nil {
		t.Fatal("expected a stance, got nil")
	}
	if pool != "philosophy" {
		t.Errorf("expected pool 'philosophy', got %q", pool)
	}
	if stance.Who == "" || stance.Where == "" || stance.Lens == "" {
		t.Errorf("stance has empty field: %+v", stance)
	}
}

func TestRandomStanceUnknownPool(t *testing.T) {
	pools, err := LoadStancePools()
	if err != nil {
		t.Fatalf("LoadStancePools failed: %v", err)
	}
	_, _, err = RandomStance(pools, "nonexistent-pool-xyz")
	if err == nil {
		t.Fatal("expected error for unknown pool, got nil")
	}
}

func TestRandomStanceEmptyNamedPool(t *testing.T) {
	pools := map[string]StancePool{
		"empty": {Name: "empty", Stances: []Stance{}},
	}
	_, _, err := RandomStance(pools, "empty")
	if err == nil {
		t.Fatal("expected error for empty named pool, got nil")
	}
	if !strings.Contains(err.Error(), "no stances") {
		t.Errorf("expected error mentioning 'no stances', got: %v", err)
	}
}

func TestRandomStanceAllPoolsEmpty(t *testing.T) {
	pools := map[string]StancePool{
		"a": {Name: "a", Stances: []Stance{}},
		"b": {Name: "b", Stances: []Stance{}},
	}
	_, _, err := RandomStance(pools, "")
	if err == nil {
		t.Fatal("expected error when all pools empty, got nil")
	}
}

func TestRandomStanceUnnamedSkipsEmptyPools(t *testing.T) {
	pools := map[string]StancePool{
		"empty": {Name: "empty", Stances: []Stance{}},
		"good":  {Name: "good", Stances: []Stance{{Who: "w", Where: "x", Lens: "y"}}},
	}
	for i := 0; i < 50; i++ {
		stance, pool, err := RandomStance(pools, "")
		if err != nil {
			t.Fatalf("RandomStance failed on iteration %d: %v", i, err)
		}
		if pool != "good" {
			t.Errorf("expected pool 'good', got %q", pool)
		}
		if stance == nil || stance.Who != "w" {
			t.Errorf("unexpected stance: %+v", stance)
		}
	}
}
