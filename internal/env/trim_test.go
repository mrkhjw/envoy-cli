package env

import (
	"testing"
)

func TestTrim_AllValues(t *testing.T) {
	entries := []Entry{
		{Key: "HOST", Value: "  localhost  "},
		{Key: "PORT", Value: " 8080 "},
	}
	out, result := Trim(entries, TrimOptions{})
	if len(result.Trimmed) != 2 {
		t.Fatalf("expected 2 trimmed, got %d", len(result.Trimmed))
	}
	if out[0].Value != "localhost" {
		t.Errorf("expected 'localhost', got %q", out[0].Value)
	}
	if out[1].Value != "8080" {
		t.Errorf("expected '8080', got %q", out[1].Value)
	}
}

func TestTrim_SpecificKeys(t *testing.T) {
	entries := []Entry{
		{Key: "HOST", Value: "  localhost  "},
		{Key: "PORT", Value: " 8080 "},
	}
	out, result := Trim(entries, TrimOptions{Keys: []string{"HOST"}})
	if len(result.Trimmed) != 1 {
		t.Fatalf("expected 1 trimmed, got %d", len(result.Trimmed))
	}
	if out[0].Value != "localhost" {
		t.Errorf("expected 'localhost', got %q", out[0].Value)
	}
	if out[1].Value != " 8080 " {
		t.Errorf("PORT should be untouched, got %q", out[1].Value)
	}
}

func TestTrim_DryRun(t *testing.T) {
	entries := []Entry{
		{Key: "HOST", Value: "  localhost  "},
	}
	out, result := Trim(entries, TrimOptions{DryRun: true})
	if len(result.Trimmed) != 1 {
		t.Fatalf("expected 1 in trimmed list, got %d", len(result.Trimmed))
	}
	if out[0].Value != "  localhost  " {
		t.Errorf("dry run should not modify value, got %q", out[0].Value)
	}
}

func TestTrim_TrimLeftOnly(t *testing.T) {
	entries := []Entry{
		{Key: "KEY", Value: "  value  "},
	}
	out, _ := Trim(entries, TrimOptions{TrimLeft: true})
	if out[0].Value != "value  " {
		t.Errorf("expected 'value  ', got %q", out[0].Value)
	}
}

func TestTrim_NoWhitespace(t *testing.T) {
	entries := []Entry{
		{Key: "CLEAN", Value: "already"},
	}
	_, result := Trim(entries, TrimOptions{})
	if len(result.Trimmed) != 0 {
		t.Errorf("expected 0 trimmed, got %d", len(result.Trimmed))
	}
}

func TestTrimResult_Format(t *testing.T) {
	r := TrimResult{Trimmed: []string{"HOST", "PORT"}, Skipped: []string{"DB"}}
	f := r.Format()
	if f == "" {
		t.Error("expected non-empty format output")
	}
}
