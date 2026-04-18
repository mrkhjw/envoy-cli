package env

import (
	"os"
	"testing"
)

func writeTempStripEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "strip-*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestStrip_RemovesSpecifiedKeys(t *testing.T) {
	entries := []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "SECRET_KEY", Value: "abc123"},
		{Key: "PORT", Value: "8080"},
	}
	kept, result := Strip(entries, []string{"SECRET_KEY"})
	if len(result.Removed) != 1 || result.Removed[0] != "SECRET_KEY" {
		t.Errorf("expected SECRET_KEY removed, got %v", result.Removed)
	}
	if len(kept) != 2 {
		t.Errorf("expected 2 kept entries, got %d", len(kept))
	}
}

func TestStrip_NoMatchingKeys(t *testing.T) {
	entries := []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "PORT", Value: "8080"},
	}
	kept, result := Strip(entries, []string{"MISSING"})
	if len(result.Removed) != 0 {
		t.Errorf("expected no removals, got %v", result.Removed)
	}
	if len(kept) != 2 {
		t.Errorf("expected 2 kept entries, got %d", len(kept))
	}
}

func TestStrip_CaseInsensitiveKeys(t *testing.T) {
	entries := []Entry{
		{Key: "SECRET_KEY", Value: "abc"},
		{Key: "PORT", Value: "9000"},
	}
	_, result := Strip(entries, []string{"secret_key"})
	if len(result.Removed) != 1 {
		t.Errorf("expected case-insensitive match, got %v", result.Removed)
	}
}

func TestStripFile_RemovesAndWrites(t *testing.T) {
	path := writeTempStripEnv(t, "APP_NAME=myapp\nSECRET_KEY=abc123\nPORT=8080\n")
	defer os.Remove(path)

	result, err := StripFile(path, []string{"SECRET_KEY"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Removed) != 1 || result.Removed[0] != "SECRET_KEY" {
		t.Errorf("expected SECRET_KEY removed, got %v", result.Removed)
	}

	entries, _ := ParseFile(path)
	for _, e := range entries {
		if e.Key == "SECRET_KEY" {
			t.Error("SECRET_KEY should have been removed from file")
		}
	}
}

func TestStripResult_Format_NoRemovals(t *testing.T) {
	r := StripResult{Removed: nil, Kept: 3}
	out := r.Format()
	if out != "No keys removed.\n" {
		t.Errorf("unexpected format output: %q", out)
	}
}
