package env

import (
	"strings"
	"testing"
)

func baseEntries() []Entry {
	return []Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_PASSWORD", Value: "s3cr3t"},
		{Key: "PORT", Value: "8080"},
		{Key: "API_SECRET", Value: "topsecret"},
	}
}

func TestPin_AllKeys(t *testing.T) {
	entries := baseEntries()
	result := Pin(entries, nil, false)
	if len(result.Pinned) != 4 {
		t.Errorf("expected 4 pinned, got %d", len(result.Pinned))
	}
	if len(result.Skipped) != 0 {
		t.Errorf("expected 0 skipped, got %d", len(result.Skipped))
	}
}

func TestPin_SpecificKeys(t *testing.T) {
	entries := baseEntries()
	result := Pin(entries, []string{"APP_ENV", "PORT"}, false)
	if len(result.Pinned) != 2 {
		t.Errorf("expected 2 pinned, got %d", len(result.Pinned))
	}
	if result.Pinned["APP_ENV"] != "production" {
		t.Errorf("unexpected value for APP_ENV")
	}
	if len(result.Skipped) != 2 {
		t.Errorf("expected 2 skipped, got %d", len(result.Skipped))
	}
}

func TestPin_DryRun(t *testing.T) {
	entries := baseEntries()
	result := Pin(entries, nil, true)
	if !result.DryRun {
		t.Error("expected DryRun to be true")
	}
	formatted := result.Format()
	if !strings.Contains(formatted, "[dry-run]") {
		t.Error("expected [dry-run] in output")
	}
}

func TestPin_MasksSecrets(t *testing.T) {
	entries := baseEntries()
	result := Pin(entries, nil, false)
	formatted := result.Format()
	if strings.Contains(formatted, "s3cr3t") {
		t.Error("expected secret value to be masked")
	}
	if strings.Contains(formatted, "topsecret") {
		t.Error("expected API_SECRET value to be masked")
	}
}

func TestPin_Format_ShowsCount(t *testing.T) {
	entries := baseEntries()
	result := Pin(entries, []string{"PORT"}, false)
	formatted := result.Format()
	if !strings.Contains(formatted, "Pinned 1 key(s)") {
		t.Errorf("unexpected format output: %s", formatted)
	}
}
