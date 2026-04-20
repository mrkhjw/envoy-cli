package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempCompareEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env: %v", err)
	}
	return path
}

func TestCompareFile_Identical(t *testing.T) {
	p1 := writeTempCompareEnv(t, "APP=prod\nDEBUG=false\n")
	p2 := writeTempCompareEnv(t, "APP=prod\nDEBUG=false\n")

	result, err := CompareFile(p1, p2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.OnlyIn1) != 0 || len(result.OnlyIn2) != 0 || len(result.Conflicts) != 0 {
		t.Errorf("expected no differences, got %+v", result)
	}
}

func TestCompareFile_Differences(t *testing.T) {
	p1 := writeTempCompareEnv(t, "APP=prod\nONLY_IN_1=yes\n")
	p2 := writeTempCompareEnv(t, "APP=staging\nONLY_IN_2=yes\n")

	result, err := CompareFile(p1, p2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %d", len(result.Conflicts))
	}
	if len(result.OnlyIn1) != 1 {
		t.Errorf("expected 1 key only in file1, got %d", len(result.OnlyIn1))
	}
	if len(result.OnlyIn2) != 1 {
		t.Errorf("expected 1 key only in file2, got %d", len(result.OnlyIn2))
	}
}

func TestCompareFile_MissingFile(t *testing.T) {
	p1 := writeTempCompareEnv(t, "APP=prod\n")
	_, err := CompareFile(p1, "/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
