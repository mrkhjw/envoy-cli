package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempCheckpointEnv(t *testing.T) []Entry {
	t.Helper()
	return []Entry{
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "DB_PASSWORD", Value: "s3cr3t"},
		{Key: "PORT", Value: "8080"},
	}
}

func TestCheckpoint_SavesFile(t *testing.T) {
	entries := writeTempCheckpointEnv(t)
	out := filepath.Join(t.TempDir(), "checkpoint.json")

	result, err := Checkpoint(entries, "v1", out, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Count != 3 {
		t.Errorf("expected count 3, got %d", result.Count)
	}
	if _, err := os.Stat(out); os.IsNotExist(err) {
		t.Error("expected checkpoint file to exist")
	}
}

func TestCheckpoint_DryRunDoesNotWrite(t *testing.T) {
	entries := writeTempCheckpointEnv(t)
	out := filepath.Join(t.TempDir(), "checkpoint.json")

	result, err := Checkpoint(entries, "v1", out, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.DryRun {
		t.Error("expected DryRun to be true")
	}
	if _, err := os.Stat(out); !os.IsNotExist(err) {
		t.Error("expected no file written in dry-run mode")
	}
}

func TestLoadCheckpoint_ReturnsEntries(t *testing.T) {
	entries := writeTempCheckpointEnv(t)
	out := filepath.Join(t.TempDir(), "checkpoint.json")

	_, err := Checkpoint(entries, "release", out, false)
	if err != nil {
		t.Fatalf("save error: %v", err)
	}

	cp, err := LoadCheckpoint(out)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if cp.Label != "release" {
		t.Errorf("expected label 'release', got %q", cp.Label)
	}
	if len(cp.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(cp.Entries))
	}
}

func TestLoadCheckpoint_MissingFile(t *testing.T) {
	_, err := LoadCheckpoint("/nonexistent/checkpoint.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestCheckpointResult_Format_DryRun(t *testing.T) {
	r := CheckpointResult{Label: "v2", Count: 5, DryRun: true, OutPath: "out.json"}
	out := r.Format()
	if !strings.Contains(out, "dry-run") {
		t.Errorf("expected dry-run in output, got: %s", out)
	}
}

func TestCheckpointResult_Format_Saved(t *testing.T) {
	r := CheckpointResult{Label: "v2", Count: 5, DryRun: false, OutPath: "out.json"}
	out := r.Format()
	if !strings.Contains(out, "saved") {
		t.Errorf("expected 'saved' in output, got: %s", out)
	}
}
