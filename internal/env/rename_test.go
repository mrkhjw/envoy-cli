package env

import (
	"os"
	"testing"
)

func writeTempRenameEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "rename_test_*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestRename_Success(t *testing.T) {
	entries := []EnvEntry{{Key: "OLD_KEY", Value: "value1"}, {Key: "OTHER", Value: "value2"}}
	updated, result, err := Rename(entries, "OLD_KEY", "NEW_KEY", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Renamed {
		t.Error("expected Renamed=true")
	}
	if updated[0].Key != "NEW_KEY" {
		t.Errorf("expected NEW_KEY, got %s", updated[0].Key)
	}
}

func TestRename_KeyNotFound(t *testing.T) {
	entries := []EnvEntry{{Key: "FOO", Value: "bar"}}
	_, result, err := Rename(entries, "MISSING", "NEW", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Renamed {
		t.Error("expected Renamed=false")
	}
}

func TestRename_NewKeyExistsNoOverwrite(t *testing.T) {
	entries := []EnvEntry{{Key: "OLD", Value: "v1"}, {Key: "NEW", Value: "v2"}}
	_, _, err := Rename(entries, "OLD", "NEW", false)
	if err == nil {
		t.Error("expected error when new key exists and overwrite=false")
	}
}

func TestRename_NewKeyExistsWithOverwrite(t *testing.T) {
	entries := []EnvEntry{{Key: "OLD", Value: "v1"}, {Key: "NEW", Value: "v2"}}
	updated, result, err := Rename(entries, "OLD", "NEW", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Renamed {
		t.Error("expected Renamed=true")
	}
	if len(updated) != 1 {
		t.Errorf("expected 1 entry, got %d", len(updated))
	}
	if updated[0].Value != "v1" {
		t.Errorf("expected value v1, got %s", updated[0].Value)
	}
}

func TestRenameFile_Valid(t *testing.T) {
	path := writeTempRenameEnv(t, "API_KEY=secret\nDEBUG=true\n")
	result, err := RenameFile(path, "DEBUG", "VERBOSE", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Renamed {
		t.Error("expected file rename to succeed")
	}
	entries, _ := ParseFile(path)
	for _, e := range entries {
		if e.Key == "DEBUG" {
			t.Error("old key DEBUG should not exist")
		}
	}
}

func TestRenameResult_Format(t *testing.T) {
	r := RenameResult{OldKey: "FOO", NewKey: "BAR", Renamed: true}
	if r.Format() == "" {
		t.Error("expected non-empty format")
	}
	r2 := RenameResult{OldKey: "MISSING", NewKey: "BAR", Renamed: false}
	if r2.Format() == "" {
		t.Error("expected non-empty format for not-found case")
	}
}
