package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempArchiveEnv(t *testing.T) []Entry {
	t.Helper()
	return []Entry{
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "SECRET_KEY", Value: "s3cr3t"},
		{Key: "PORT", Value: "8080"},
	}
}

func TestArchive_CreatesFile(t *testing.T) {
	entries := writeTempArchiveEnv(t)
	dest := filepath.Join(t.TempDir(), "archive.json")

	res, err := Archive(entries, dest, "v1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Archived != 3 {
		t.Errorf("expected 3 archived, got %d", res.Archived)
	}
	if res.Label != "v1" {
		t.Errorf("expected label v1, got %s", res.Label)
	}
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		t.Error("expected archive file to exist")
	}
}

func TestLoadArchive_ReturnsEntries(t *testing.T) {
	entries := writeTempArchiveEnv(t)
	dest := filepath.Join(t.TempDir(), "archive.json")

	if _, err := Archive(entries, dest, "test-label"); err != nil {
		t.Fatalf("archive failed: %v", err)
	}

	ae, err := LoadArchive(dest)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if len(ae.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(ae.Entries))
	}
	if ae.Label != "test-label" {
		t.Errorf("expected label test-label, got %s", ae.Label)
	}
}

func TestLoadArchive_MissingFile(t *testing.T) {
	_, err := LoadArchive("/nonexistent/archive.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestArchiveResult_Format(t *testing.T) {
	r := ArchiveResult{Archived: 5, Label: "prod", Path: "/tmp/env.json"}
	out := r.Format()
	if out == "" {
		t.Error("expected non-empty format output")
	}
}
