package env

import (
	"strings"
	"testing"
)

var baseSplitEntries = []Entry{
	{Key: "APP_NAME", Value: "myapp"},
	{Key: "DB_PASSWORD", Value: "secret"},
	{Key: "API_KEY", Value: "key123"},
	{Key: "PORT", Value: "8080"},
}

func TestSplit_ByKeys(t *testing.T) {
	result := Split(baseSplitEntries, SplitOptions{Keys: []string{"APP_NAME", "PORT"}})
	if len(result.Matched) != 2 {
		t.Errorf("expected 2 matched, got %d", len(result.Matched))
	}
	if len(result.Remainder) != 2 {
		t.Errorf("expected 2 remainder, got %d", len(result.Remainder))
	}
}

func TestSplit_Inverted(t *testing.T) {
	result := Split(baseSplitEntries, SplitOptions{Keys: []string{"APP_NAME"}, Invert: true})
	if len(result.Matched) != 3 {
		t.Errorf("expected 3 matched (inverted), got %d", len(result.Matched))
	}
	if len(result.Remainder) != 1 {
		t.Errorf("expected 1 remainder, got %d", len(result.Remainder))
	}
}

func TestSplit_NoKeys(t *testing.T) {
	result := Split(baseSplitEntries, SplitOptions{})
	if len(result.Matched) != 0 {
		t.Errorf("expected 0 matched, got %d", len(result.Matched))
	}
	if len(result.Remainder) != 4 {
		t.Errorf("expected 4 remainder, got %d", len(result.Remainder))
	}
}

func TestSplit_DryRun(t *testing.T) {
	result := Split(baseSplitEntries, SplitOptions{Keys: []string{"PORT"}, DryRun: true})
	if !result.DryRun {
		t.Error("expected DryRun to be true")
	}
}

func TestSplitResult_Format(t *testing.T) {
	result := SplitResult{Matched: baseSplitEntries[:2], Remainder: baseSplitEntries[2:], DryRun: true}
	out := result.Format()
	if !strings.Contains(out, "matched: 2") {
		t.Errorf("expected matched count in format, got: %s", out)
	}
	if !strings.Contains(out, "dry run") {
		t.Errorf("expected dry run in format, got: %s", out)
	}
}
