package env

import (
	"testing"
)

var baseCompactEntries = []Entry{
	{Key: "APP_NAME", Value: "myapp"},
	{Key: "", Value: "", Comment: true, Raw: "# this is a comment"},
	{Key: "DEBUG", Value: ""},
	{Key: "SECRET_KEY", Value: "abc123"},
	{Key: "", Value: "", Comment: true, Raw: "# another comment"},
	{Key: "PORT", Value: "8080"},
}

func TestCompact_RemovesComments(t *testing.T) {
	result := Compact(baseCompactEntries, CompactOptions{RemoveComments: true})
	for _, e := range result.Compacted {
		if e.Comment {
			t.Errorf("expected no comments, got: %s", e.Raw)
		}
	}
	if result.Removed != 2 {
		t.Errorf("expected 2 removed, got %d", result.Removed)
	}
}

func TestCompact_RemovesEmptyValues(t *testing.T) {
	result := Compact(baseCompactEntries, CompactOptions{RemoveEmpty: true})
	for _, e := range result.Compacted {
		if !e.Comment && e.Value == "" {
			t.Errorf("expected no empty values, got key: %s", e.Key)
		}
	}
	if result.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", result.Removed)
	}
}

func TestCompact_RemovesBoth(t *testing.T) {
	result := Compact(baseCompactEntries, CompactOptions{RemoveComments: true, RemoveEmpty: true})
	if result.Removed != 3 {
		t.Errorf("expected 3 removed, got %d", result.Removed)
	}
	if len(result.Compacted) != 3 {
		t.Errorf("expected 3 remaining entries, got %d", len(result.Compacted))
	}
}

func TestCompact_DryRunDoesNotChange(t *testing.T) {
	result := Compact(baseCompactEntries, CompactOptions{RemoveComments: true, DryRun: true})
	if len(result.Compacted) != len(baseCompactEntries) {
		t.Errorf("dry run should not modify entries")
	}
	if result.Removed != 2 {
		t.Errorf("expected removed count 2, got %d", result.Removed)
	}
}

func TestCompactResult_Format_NothingRemoved(t *testing.T) {
	r := CompactResult{Removed: 0}
	if r.Format() != "nothing to compact" {
		t.Errorf("unexpected format: %s", r.Format())
	}
}

func TestCompactResult_Format_WithDryRun(t *testing.T) {
	r := CompactResult{Removed: 3, DryRun: true}
	out := r.Format()
	if out != "compacted 3 entries (dry run)" {
		t.Errorf("unexpected format: %s", out)
	}
}
