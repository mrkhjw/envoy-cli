package env

import (
	"strings"
	"testing"
)

func TestInterpolate_BasicSubstitution(t *testing.T) {
	env := map[string]string{
		"HOST":     "localhost",
		"PORT":     "5432",
		"DATABASE_URL": "postgres://${HOST}:${PORT}/db",
	}
	result := Interpolate(env)
	if got := result.Resolved["DATABASE_URL"]; got != "postgres://localhost:5432/db" {
		t.Errorf("expected expanded URL, got %q", got)
	}
	if len(result.Unresolved) != 0 {
		t.Errorf("expected no unresolved, got %v", result.Unresolved)
	}
}

func TestInterpolate_UnresolvedReference(t *testing.T) {
	env := map[string]string{
		"URL": "http://${MISSING_HOST}/path",
	}
	result := Interpolate(env)
	if len(result.Unresolved) != 1 || result.Unresolved[0] != "MISSING_HOST" {
		t.Errorf("expected MISSING_HOST unresolved, got %v", result.Unresolved)
	}
	if result.Resolved["URL"] != "http://${MISSING_HOST}/path" {
		t.Errorf("expected original placeholder preserved")
	}
}

func TestInterpolate_NoPlaceholders(t *testing.T) {
	env := map[string]string{
		"KEY": "plainvalue",
	}
	result := Interpolate(env)
	if result.Resolved["KEY"] != "plainvalue" {
		t.Errorf("expected unchanged value")
	}
	if len(result.Unresolved) != 0 {
		t.Errorf("expected no unresolved")
	}
}

func TestInterpolate_EmptyMap(t *testing.T) {
	result := Interpolate(map[string]string{})
	if len(result.Resolved) != 0 {
		t.Errorf("expected empty resolved map")
	}
}

func TestInterpolateResult_Format(t *testing.T) {
	r := InterpolateResult{
		Resolved:   map[string]string{"A": "1", "B": "2"},
		Unresolved: []string{"MISSING"},
	}
	out := r.Format()
	if !strings.Contains(out, "Resolved: 2") {
		t.Errorf("expected resolved count in format output")
	}
	if !strings.Contains(out, "MISSING") {
		t.Errorf("expected unresolved key in format output")
	}
}
