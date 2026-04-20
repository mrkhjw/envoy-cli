package env

import (
	"strings"
	"testing"
)

func basePlaceholderEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASSWORD", Value: "CHANGEME"},
		{Key: "API_KEY", Value: "CHANGEME"},
		{Key: "PORT", Value: "8080"},
	}
}

func TestFillPlaceholders_FillsMatching(t *testing.T) {
	entries := basePlaceholderEntries()
	replacements := map[string]string{
		"DB_PASSWORD": "supersecret",
		"API_KEY":     "abc123",
	}
	res := FillPlaceholders(entries, replacements, PlaceholderOptions{})
	if res.Filled != 2 {
		t.Errorf("expected 2 filled, got %d", res.Filled)
	}
	for _, e := range res.Entries {
		if e.Key == "DB_PASSWORD" && e.Value != "supersecret" {
			t.Errorf("expected DB_PASSWORD=supersecret, got %s", e.Value)
		}
	}
}

func TestFillPlaceholders_TracksMissing(t *testing.T) {
	entries := basePlaceholderEntries()
	res := FillPlaceholders(entries, map[string]string{}, PlaceholderOptions{})
	if len(res.Missing) != 2 {
		t.Errorf("expected 2 missing, got %d", len(res.Missing))
	}
}

func TestFillPlaceholders_DryRunDoesNotWrite(t *testing.T) {
	entries := basePlaceholderEntries()
	replacements := map[string]string{"DB_PASSWORD": "newsecret"}
	res := FillPlaceholders(entries, replacements, PlaceholderOptions{DryRun: true})
	for _, e := range res.Entries {
		if e.Key == "DB_PASSWORD" && e.Value != "CHANGEME" {
			t.Errorf("dry run should not change value")
		}
	}
	if res.Filled != 1 {
		t.Errorf("expected filled count 1 even in dry run")
	}
}

func TestFillPlaceholders_CustomToken(t *testing.T) {
	entries := []Entry{{Key: "SECRET", Value: "TODO"}}
	replacements := map[string]string{"SECRET": "val"}
	res := FillPlaceholders(entries, replacements, PlaceholderOptions{Token: "TODO"})
	if res.Filled != 1 {
		t.Errorf("expected 1 filled with custom token")
	}
}

func TestPlaceholderResult_Format_NoPlaceholders(t *testing.T) {
	r := PlaceholderResult{}
	out := r.Format()
	if !strings.Contains(out, "no placeholders") {
		t.Errorf("expected no placeholders message, got: %s", out)
	}
}
