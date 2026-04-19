package env

import (
	"strings"
	"testing"
)

func TestDedupe_NoDuplicates(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", Value: "bar"},
		{Key: "BAZ", Value: "qux"},
	}
	res := Dedupe(entries)
	if res.Dupes != 0 {
		t.Errorf("expected 0 dupes, got %d", res.Dupes)
	}
	if len(res.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(res.Entries))
	}
}

func TestDedupe_RemovesDuplicateKeepLast(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", Value: "first"},
		{Key: "BAR", Value: "keep"},
		{Key: "FOO", Value: "last"},
	}
	res := Dedupe(entries)
	if res.Dupes != 1 {
		t.Errorf("expected 1 dupe, got %d", res.Dupes)
	}
	if len(res.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(res.Entries))
	}
	for _, e := range res.Entries {
		if e.Key == "FOO" && e.Value != "last" {
			t.Errorf("expected last value 'last', got %s", e.Value)
		}
	}
}

func TestDedupe_MultipleDuplicates(t *testing.T) {
	entries := []Entry{
		{Key: "A", Value: "1"},
		{Key: "A", Value: "2"},
		{Key: "B", Value: "x"},
		{Key: "B", Value: "y"},
		{Key: "C", Value: "z"},
	}
	res := Dedupe(entries)
	if res.Dupes != 2 {
		t.Errorf("expected 2 dupes, got %d", res.Dupes)
	}
	if len(res.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(res.Entries))
	}
}

func TestDedupeResult_Format_NoDupes(t *testing.T) {
	res := DedupeResult{Total: 3, Dupes: 0}
	out := res.Format()
	if !strings.Contains(out, "no duplicates") {
		t.Errorf("expected 'no duplicates' in output, got: %s", out)
	}
}

func TestDedupeResult_Format_WithDupes(t *testing.T) {
	res := DedupeResult{
		Total:   4,
		Dupes:   2,
		Removed: []string{"FOO", "BAR"},
	}
	out := res.Format()
	if !strings.Contains(out, "FOO") || !strings.Contains(out, "BAR") {
		t.Errorf("expected removed keys in output, got: %s", out)
	}
}
