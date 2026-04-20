package env

import (
	"strings"
	"testing"
)

func baseWrapEntries() []Entry {
	return []Entry{
		{Key: "SHORT", Value: "hi"},
		{Key: "LONG_VALUE", Value: strings.Repeat("a", 120)},
		{Key: "SECRET_TOKEN", Value: strings.Repeat("b", 100)},
		{Key: "#comment", Value: "", Comment: true},
	}
}

func TestWrap_TruncatesLongValues(t *testing.T) {
	entries := baseWrapEntries()
	res := Wrap(entries, WrapOptions{MaxLength: 80})
	for _, e := range res.Wrapped {
		if !e.Comment && len(e.Value) > 80 {
			t.Errorf("key %s: value length %d exceeds max 80", e.Key, len(e.Value))
		}
	}
	if res.Modified != 2 {
		t.Errorf("expected 2 modified, got %d", res.Modified)
	}
}

func TestWrap_ShortValuesUnchanged(t *testing.T) {
	entries := baseWrapEntries()
	res := Wrap(entries, WrapOptions{MaxLength: 80})
	for _, e := range res.Wrapped {
		if e.Key == "SHORT" && e.Value != "hi" {
			t.Errorf("expected SHORT unchanged, got %q", e.Value)
		}
	}
}

func TestWrap_QuoteOption(t *testing.T) {
	entries := []Entry{{Key: "MSG", Value: "hello world"}}
	res := Wrap(entries, WrapOptions{MaxLength: 80, Quote: true})
	if !strings.HasPrefix(res.Wrapped[0].Value, "\"") {
		t.Errorf("expected quoted value, got %q", res.Wrapped[0].Value)
	}
	if res.Modified != 1 {
		t.Errorf("expected 1 modified, got %d", res.Modified)
	}
}

func TestWrap_SpecificKeys(t *testing.T) {
	entries := baseWrapEntries()
	res := Wrap(entries, WrapOptions{MaxLength: 80, Keys: []string{"LONG_VALUE"}})
	for _, e := range res.Wrapped {
		if e.Key == "SECRET_TOKEN" && len(e.Value) <= 80 {
			t.Errorf("SECRET_TOKEN should not have been truncated")
		}
	}
	if res.Modified != 1 {
		t.Errorf("expected 1 modified, got %d", res.Modified)
	}
}

func TestWrap_DryRunDoesNotChange(t *testing.T) {
	original := strings.Repeat("x", 120)
	entries := []Entry{{Key: "BIG", Value: original}}
	res := Wrap(entries, WrapOptions{MaxLength: 80, DryRun: true})
	if res.Wrapped[0].Value != original {
		t.Errorf("dry run should not modify value")
	}
	if res.Modified != 1 {
		t.Errorf("expected modified count 1 even in dry run, got %d", res.Modified)
	}
}

func TestWrapResult_Format(t *testing.T) {
	r := WrapResult{Modified: 3, DryRun: false}
	if !strings.Contains(r.Format(), "3") {
		t.Errorf("format should mention count: %s", r.Format())
	}
	r.DryRun = true
	if !strings.Contains(r.Format(), "dry run") {
		t.Errorf("format should mention dry run: %s", r.Format())
	}
}
