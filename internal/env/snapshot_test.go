package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempSnapshotEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "snap-*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestTakeSnapshot_CreatesFile(t *testing.T) {
	src := writeTempSnapshotEnv(t, "APP_NAME=envoy\nSECRET_KEY=abc123\n")
	dest := filepath.Join(t.TempDir(), "snap.json")

	result, err := TakeSnapshot(src, dest)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.KeyCount != 2 {
		t.Errorf("expected 2 keys, got %d", result.KeyCount)
	}
	if result.Path != dest {
		t.Errorf("expected path %s, got %s", dest, result.Path)
	}
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		t.Error("snapshot file was not created")
	}
}

func TestLoadSnapshot_ReturnsEntries(t *testing.T) {
	src := writeTempSnapshotEnv(t, "DB_HOST=localhost\nDB_PORT=5432\n")
	dest := filepath.Join(t.TempDir(), "snap.json")

	if _, err := TakeSnapshot(src, dest); err != nil {
		t.Fatalf("snapshot failed: %v", err)
	}

	snap, err := LoadSnapshot(dest)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if snap.Entries["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %s", snap.Entries["DB_HOST"])
	}
	if snap.Source != src {
		t.Errorf("expected source %s, got %s", src, snap.Source)
	}
}

func TestTakeSnapshot_MissingSource(t *testing.T) {
	_, err := TakeSnapshot("/nonexistent/.env", "/tmp/out.json")
	if err == nil {
		t.Error("expected error for missing source file")
	}
}

func TestSnapshotResult_Format(t *testing.T) {
	src := writeTempSnapshotEnv(t, "KEY=val\n")
	dest := filepath.Join(t.TempDir(), "snap.json")
	result, err := TakeSnapshot(src, dest)
	if err != nil {
		t.Fatal(err)
	}
	out := result.Format()
	if out == "" {
		t.Error("expected non-empty format output")
	}
}
