//go:build integration

package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func buildBinary(t *testing.T) string {
	t.Helper()
	binary := filepath.Join(t.TempDir(), "metacog")
	cmd := exec.Command("go", "build", "-o", binary, "./cmd/metacog")
	// Run from repo root
	cmd.Dir = filepath.Join(".", "..", "..")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	return binary
}

func runMetacog(t *testing.T, binary, stateDir string, args ...string) (string, error) {
	t.Helper()
	cmd := exec.Command(binary, args...)
	cmd.Env = append(os.Environ(), "METACOG_HOME="+stateDir)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func TestIntegrationFullPivot(t *testing.T) {
	binary := buildBinary(t)
	stateDir := t.TempDir()

	// Start pivot
	out, err := runMetacog(t, binary, stateDir, "stratagem", "start", "pivot")
	if err != nil {
		t.Fatalf("start pivot: %v\n%s", err, out)
	}

	// Call drugs
	_, err = runMetacog(t, binary, stateDir, "drugs", "--substance", "test", "--method", "test", "--qualia", "test")
	if err != nil {
		t.Fatalf("drugs: %v", err)
	}

	// Advance
	_, err = runMetacog(t, binary, stateDir, "stratagem", "next")
	if err != nil {
		t.Fatalf("next after drugs: %v", err)
	}

	// THINK step — advance freely
	_, err = runMetacog(t, binary, stateDir, "stratagem", "next")
	if err != nil {
		t.Fatalf("next after THINK: %v", err)
	}

	// Call become
	_, err = runMetacog(t, binary, stateDir, "become", "--name", "test", "--lens", "test", "--env", "test")
	if err != nil {
		t.Fatalf("become: %v", err)
	}

	// Advance
	_, err = runMetacog(t, binary, stateDir, "stratagem", "next")
	if err != nil {
		t.Fatalf("next after become: %v", err)
	}

	// THINK step
	_, err = runMetacog(t, binary, stateDir, "stratagem", "next")
	if err != nil {
		t.Fatalf("next after THINK 2: %v", err)
	}

	// Call ritual
	_, err = runMetacog(t, binary, stateDir, "ritual", "--threshold", "test", "--steps", "s1", "--result", "test")
	if err != nil {
		t.Fatalf("ritual: %v", err)
	}

	// Final advance — should complete
	out, err = runMetacog(t, binary, stateDir, "stratagem", "next")
	if err != nil {
		t.Fatalf("final next: %v\n%s", err, out)
	}

	// Status should show no active stratagem
	out, err = runMetacog(t, binary, stateDir, "status")
	if err != nil {
		t.Fatalf("status: %v\n%s", err, out)
	}
}

func TestIntegrationStatePersistence(t *testing.T) {
	binary := buildBinary(t)
	stateDir := t.TempDir()

	runMetacog(t, binary, stateDir, "become", "--name", "Ada", "--lens", "logic", "--env", "lab")

	out, err := runMetacog(t, binary, stateDir, "status")
	if err != nil {
		t.Fatalf("status: %v\n%s", err, out)
	}
	if !strings.Contains(out, "Ada") {
		t.Error("status should show Ada")
	}
}

func TestIntegrationReset(t *testing.T) {
	binary := buildBinary(t)
	stateDir := t.TempDir()

	runMetacog(t, binary, stateDir, "become", "--name", "Ada", "--lens", "logic", "--env", "lab")
	runMetacog(t, binary, stateDir, "reset")

	out, _ := runMetacog(t, binary, stateDir, "status")
	if strings.Contains(out, "Ada") {
		t.Error("reset should clear identity")
	}
}

func TestIntegrationRepair(t *testing.T) {
	binary := buildBinary(t)
	stateDir := t.TempDir()

	// Corrupt state file
	os.MkdirAll(stateDir, 0755)
	os.WriteFile(filepath.Join(stateDir, "state.json"), []byte("corrupt{{{"), 0644)

	out, err := runMetacog(t, binary, stateDir, "repair")
	if err != nil {
		t.Fatalf("repair: %v\n%s", err, out)
	}

	// Should work now
	_, err = runMetacog(t, binary, stateDir, "status")
	if err != nil {
		t.Fatalf("status after repair: %v", err)
	}
}

func TestIntegrationJSONOutput(t *testing.T) {
	binary := buildBinary(t)
	stateDir := t.TempDir()

	out, err := runMetacog(t, binary, stateDir, "version", "--json")
	if err != nil {
		t.Fatalf("version --json: %v\n%s", err, out)
	}

	var parsed map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &parsed); err != nil {
		t.Errorf("not valid JSON: %v\noutput: %s", err, out)
	}
}

func TestIntegrationReflect(t *testing.T) {
	binary := buildBinary(t)
	stateDir := t.TempDir()

	runMetacog(t, binary, stateDir, "become", "--name", "Ada", "--lens", "logic", "--env", "lab")
	runMetacog(t, binary, stateDir, "drugs", "--substance", "caffeine", "--method", "antagonism", "--qualia", "sharp")
	runMetacog(t, binary, stateDir, "ritual", "--threshold", "test", "--steps", "s1", "--steps", "s2", "--result", "done")

	out, err := runMetacog(t, binary, stateDir, "reflect")
	if err != nil {
		t.Fatalf("reflect: %v\n%s", err, out)
	}
	if !strings.Contains(out, "become: 1") {
		t.Errorf("reflect should show become count:\n%s", out)
	}
}

func TestIntegrationSession(t *testing.T) {
	binary := buildBinary(t)
	stateDir := t.TempDir()

	out, err := runMetacog(t, binary, stateDir, "session", "start", "test-session")
	if err != nil {
		t.Fatalf("session start: %v\n%s", err, out)
	}

	runMetacog(t, binary, stateDir, "become", "--name", "Ada", "--lens", "logic", "--env", "lab")

	out, err = runMetacog(t, binary, stateDir, "session", "end")
	if err != nil {
		t.Fatalf("session end: %v\n%s", err, out)
	}

	out, err = runMetacog(t, binary, stateDir, "session", "list")
	if err != nil {
		t.Fatalf("session list: %v\n%s", err, out)
	}
	if !strings.Contains(out, "test-session") {
		t.Errorf("session list should contain 'test-session':\n%s", out)
	}

	out, err = runMetacog(t, binary, stateDir, "history", "--session", "test-session")
	if err != nil {
		t.Fatalf("history --session: %v\n%s", err, out)
	}
	if !strings.Contains(out, "Ada") {
		t.Errorf("filtered history should contain Ada:\n%s", out)
	}
}

func TestIntegrationInspireSave(t *testing.T) {
	binary := buildBinary(t)
	stateDir := t.TempDir()

	runMetacog(t, binary, stateDir, "become", "--name", "Ada", "--lens", "logic", "--env", "lab")

	out, err := runMetacog(t, binary, stateDir, "inspire", "--save")
	if err != nil {
		t.Fatalf("inspire --save: %v\n%s", err, out)
	}
	if !strings.Contains(out, "Ada") {
		t.Errorf("save output should mention identity:\n%s", out)
	}

	out, err = runMetacog(t, binary, stateDir, "inspire", "--pool", "personal")
	if err != nil {
		t.Fatalf("inspire --pool personal: %v\n%s", err, out)
	}
	if !strings.Contains(out, "Ada") {
		t.Errorf("personal pool should contain Ada:\n%s", out)
	}
}

func TestIntegrationExitCodes(t *testing.T) {
	binary := buildBinary(t)
	stateDir := t.TempDir()

	// Missing required flag — exit code 1
	_, err := runMetacog(t, binary, stateDir, "become", "--name", "Ada")
	if err == nil {
		t.Error("expected error for missing flags")
	}

	// Stratagem error — no active stratagem
	_, err = runMetacog(t, binary, stateDir, "stratagem", "next")
	if err == nil {
		t.Error("expected error for next with no stratagem")
	}
}
