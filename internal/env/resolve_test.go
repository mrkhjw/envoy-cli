package env

import (
	"strings"
	"testing"
)

func baseResolveEntries() []EnvEntry {
	return []EnvEntry{
		{Key: "BASE_URL", Value: "https://example.com"},
		{Key: "API_URL", Value: "${BASE_URL}/api"},
		{Key: "CALLBACK", Value: "${API_URL}/callback"},
		{Key: "APP_NAME", Value: "myapp"},
	}
}

func TestResolve_ExpandsReferences(t *testing.T) {
	entries := baseResolveEntries()
	result, err := Resolve(entries, ResolveOption{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range result.Resolved {
		if e.Key == "API_URL" && e.Value != "https://example.com/api" {
			t.Errorf("expected API_URL=https://example.com/api, got %s", e.Value)
		}
	}
}

func TestResolve_TracksUnresolved(t *testing.T) {
	entries := []EnvEntry{
		{Key: "FOO", Value: "${MISSING_VAR}"},
	}
	result, err := Resolve(entries, ResolveOption{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Unresolved) != 1 || result.Unresolved[0] != "MISSING_VAR" {
		t.Errorf("expected MISSING_VAR in unresolved, got %v", result.Unresolved)
	}
}

func TestResolve_StrictModeErrors(t *testing.T) {
	entries := []EnvEntry{
		{Key: "FOO", Value: "${DOES_NOT_EXIST}"},
	}
	_, err := Resolve(entries, ResolveOption{Strict: true})
	if err == nil {
		t.Error("expected error in strict mode, got nil")
	}
}

func TestResolve_NoPlaceholders(t *testing.T) {
	entries := []EnvEntry{
		{Key: "HOST", Value: "localhost"},
		{Key: "PORT", Value: "8080"},
	}
	result, err := Resolve(entries, ResolveOption{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Unresolved) != 0 {
		t.Errorf("expected no unresolved, got %v", result.Unresolved)
	}
	if len(result.Resolved) != 2 {
		t.Errorf("expected 2 resolved entries, got %d", len(result.Resolved))
	}
}

func TestResolveResult_Format_WithUnresolved(t *testing.T) {
	r := ResolveResult{
		Resolved:   []EnvEntry{{Key: "A", Value: "1"}},
		Unresolved: []string{"MISSING"},
	}
	out := r.Format()
	if !strings.Contains(out, "MISSING") {
		t.Errorf("expected MISSING in format output, got: %s", out)
	}
	if !strings.Contains(out, "Unresolved") {
		t.Errorf("expected 'Unresolved' label in output")
	}
}

func TestResolveResult_Format_AllResolved(t *testing.T) {
	r := ResolveResult{
		Resolved: []EnvEntry{{Key: "A", Value: "ok"}},
	}
	out := r.Format()
	if !strings.Contains(out, "All references resolved") {
		t.Errorf("expected success message, got: %s", out)
	}
}
