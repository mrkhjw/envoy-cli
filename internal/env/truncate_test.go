package env

import (
	"strings"
	"testing"
)

var baseTruncateEntries = []Entry{
	{Key: "SHORT", Value: "hi"},
	{Key: "API_SECRET", Value: "averylongsecretvaluethatexceedsthemaximumlengthallowed"},
	{Key: "DB_URL", Value: "postgres://user:pass@localhost:5432/mydb"},
	{Key: "NOTE", Value: "just a note"},
}

func TestTruncate_ShortValuesUnchanged(t *testing.T) {
	res := Truncate(baseTruncateEntries, TruncateOptions{MaxLen: 20})
	for _, e := range res.Entries {
		if e.Key == "SHORT" && e.Value != "hi" {
			t.Errorf("expected 'hi', got %q", e.Value)
		}
	}
	if contains(res.Truncated, "SHORT") {
		t.Error("SHORT should not be in truncated list")
	}
}

func TestTruncate_LongValueTruncated(t *testing.T) {
	res := Truncate(baseTruncateEntries, TruncateOptions{MaxLen: 10, Suffix: "…"})
	for _, e := range res.Entries {
		if e.Key == "API_SECRET" {
			if len(e.Value) > 11 {
				t.Errorf("value not truncated: %q", e.Value)
			}
			if !strings.HasSuffix(e.Value, "…") {
				t.Errorf("expected suffix '…', got %q", e.Value)
			}
		}
	}
	if !contains(res.Truncated, "API_SECRET") {
		t.Error("API_SECRET should be in truncated list")
	}
}

func TestTruncate_DryRunDoesNotModify(t *testing.T) {
	orig := baseTruncateEntries[1].Value
	res := Truncate(baseTruncateEntries, TruncateOptions{MaxLen: 10, DryRun: true})
	for _, e := range res.Entries {
		if e.Key == "API_SECRET" && e.Value != orig {
			t.Errorf("dry run modified value: got %q", e.Value)
		}
	}
	if !contains(res.Truncated, "API_SECRET") {
		t.Error("expected API_SECRET in dry-run truncated list")
	}
	if !res.DryRun {
		t.Error("expected DryRun flag set")
	}
}

func TestTruncate_SpecificKeys(t *testing.T) {
	res := Truncate(baseTruncateEntries, TruncateOptions{MaxLen: 10, Keys: []string{"DB_URL"}})
	for _, e := range res.Entries {
		if e.Key == "API_SECRET" && len(e.Value) <= 10 {
			t.Error("API_SECRET should not be truncated when not in Keys")
		}
	}
	if !contains(res.Truncated, "DB_URL") {
		t.Error("expected DB_URL in truncated list")
	}
}

func TestTruncateResult_Format(t *testing.T) {
	res := TruncateResult{Truncated: []string{"A", "B"}}
	out := res.Format()
	if !strings.Contains(out, "2") {
		t.Errorf("expected count in format output, got %q", out)
	}

	empty := TruncateResult{}
	if empty.Format() != "no values truncated" {
		t.Errorf("unexpected empty format: %q", empty.Format())
	}
}

func contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}
