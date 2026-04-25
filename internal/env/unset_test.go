package env

import (
	"strings"
	"testing"
)

var baseUnsetEntries = []Entry{
	{Key: "APP_NAME", Value: "myapp"},
	{Key: "DB_PASSWORD", Value: "secret"},
	{Key: "PORT", Value: "8080"},
	{Key: "DEBUG", Value: "true"},
}

func TestUnset_RemovesKeys(t *testing.T) {
	res := Unset(baseUnsetEntries, []string{"APP_NAME", "PORT"}, false)
	if len(res.Removed) != 2 {
		t.Fatalf("expected 2 removed, got %d", len(res.Removed))
	}
	for _, e := range res.Entries {
		if e.Key == "APP_NAME" || e.Key == "PORT" {
			t.Errorf("key %q should have been removed", e.Key)
		}
	}
}

func TestUnset_SkipsMissingKeys(t *testing.T) {
	res := Unset(baseUnsetEntries, []string{"NONEXISTENT"}, false)
	if len(res.Skipped) != 1 || res.Skipped[0] != "NONEXISTENT" {
		t.Errorf("expected NONEXISTENT in skipped, got %v", res.Skipped)
	}
	if len(res.Entries) != len(baseUnsetEntries) {
		t.Errorf("entries should be unchanged")
	}
}

func TestUnset_DryRunDoesNotRemove(t *testing.T) {
	res := Unset(baseUnsetEntries, []string{"DEBUG"}, true)
	if !res.DryRun {
		t.Error("expected DryRun to be true")
	}
	if len(res.Entries) != len(baseUnsetEntries) {
		t.Errorf("dry-run should not modify entries")
	}
	if len(res.Removed) != 1 || res.Removed[0] != "DEBUG" {
		t.Errorf("expected DEBUG in removed list")
	}
}

func TestUnset_CaseInsensitiveKeys(t *testing.T) {
	res := Unset(baseUnsetEntries, []string{"app_name"}, false)
	if len(res.Removed) != 1 {
		t.Fatalf("expected 1 removed, got %d", len(res.Removed))
	}
}

func TestUnsetResult_Format(t *testing.T) {
	res := Unset(baseUnsetEntries, []string{"APP_NAME", "MISSING"}, false)
	out := res.Format()
	if !strings.Contains(out, "removed: 1") {
		t.Errorf("expected removed count in format output, got: %s", out)
	}
	if !strings.Contains(out, "MISSING (not found)") {
		t.Errorf("expected skipped key in format output, got: %s", out)
	}
}

func TestUnsetResult_Format_DryRun(t *testing.T) {
	res := Unset(baseUnsetEntries, []string{"PORT"}, true)
	out := res.Format()
	if !strings.Contains(out, "[dry-run]") {
		t.Errorf("expected dry-run label in output, got: %s", out)
	}
}
