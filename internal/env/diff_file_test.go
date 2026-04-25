package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempDiffEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env: %v", err)
	}
	return p
}

func TestDiffFile_Added(t *testing.T) {
	a := writeTempDiffEnv(t, "FOO=bar\n")
	b := writeTempDiffEnv(t, "FOO=bar\nBAZ=qux\n")

	result, err := DiffFile(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Added) != 1 || result.Added[0].Key != "BAZ" {
		t.Errorf("expected BAZ added, got %+v", result.Added)
	}
}

func TestDiffFile_Removed(t *testing.T) {
	a := writeTempDiffEnv(t, "FOO=bar\nBAZ=qux\n")
	b := writeTempDiffEnv(t, "FOO=bar\n")

	result, err := DiffFile(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Removed) != 1 || result.Removed[0].Key != "BAZ" {
		t.Errorf("expected BAZ removed, got %+v", result.Removed)
	}
}

func TestDiffFile_Changed(t *testing.T) {
	a := writeTempDiffEnv(t, "FOO=bar\n")
	b := writeTempDiffEnv(t, "FOO=newval\n")

	result, err := DiffFile(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Changed) != 1 || result.Changed[0].Key != "FOO" {
		t.Errorf("expected FOO changed, got %+v", result.Changed)
	}
}

func TestDiffFile_MissingFileA(t *testing.T) {
	b := writeTempDiffEnv(t, "FOO=bar\n")
	_, err := DiffFile("/nonexistent/.env", b)
	if err == nil {
		t.Error("expected error for missing file A")
	}
}

func TestDiffFile_MissingFileB(t *testing.T) {
	a := writeTempDiffEnv(t, "FOO=bar\n")
	_, err := DiffFile(a, "/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file B")
	}
}
