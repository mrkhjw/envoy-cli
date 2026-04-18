package env

import (
	"strings"
	"testing"
)

func TestRotate_AllSecrets(t *testing.T) {
	env := map[string]string{
		"API_KEY":  "abc123",
		"APP_NAME": "myapp",
	}
	out, result := Rotate(env, RotateOptions{})
	if out["API_KEY"] == "abc123" {
		t.Error("expected API_KEY to be rotated")
	}
	if out["APP_NAME"] != "myapp" {
		t.Error("expected APP_NAME to be unchanged")
	}
	if len(result.Rotated) != 1 || result.Rotated[0] != "API_KEY" {
		t.Errorf("unexpected rotated list: %v", result.Rotated)
	}
}

func TestRotate_SpecificKeys(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "secret",
		"DB_HOST":     "localhost",
	}
	out, result := Rotate(env, RotateOptions{Keys: []string{"DB_PASSWORD"}})
	if out["DB_PASSWORD"] == "secret" {
		t.Error("expected DB_PASSWORD to be rotated")
	}
	if len(result.Rotated) != 1 {
		t.Errorf("expected 1 rotated, got %d", len(result.Rotated))
	}
}

func TestRotate_DryRun(t *testing.T) {
	env := map[string]string{"SECRET_KEY": "original"}
	out, result := Rotate(env, RotateOptions{DryRun: true})
	if out["SECRET_KEY"] != "original" {
		t.Error("dry-run should not modify values")
	}
	if !result.DryRun {
		t.Error("expected DryRun flag set")
	}
	if len(result.Rotated) != 1 {
		t.Error("expected key listed as rotated in dry-run")
	}
}

func TestRotate_MissingKeySkipped(t *testing.T) {
	env := map[string]string{"APP": "x"}
	_, result := Rotate(env, RotateOptions{Keys: []string{"MISSING_KEY"}})
	if len(result.Skipped) != 1 || result.Skipped[0] != "MISSING_KEY" {
		t.Errorf("expected MISSING_KEY skipped, got %v", result.Skipped)
	}
}

func TestRotateResult_Format(t *testing.T) {
	r := RotateResult{
		Rotated: []string{"API_KEY"},
		Skipped: []string{"OLD_TOKEN"},
		DryRun:  true,
	}
	out := r.Format()
	if !strings.Contains(out, "[dry-run]") {
		t.Error("expected dry-run prefix")
	}
	if !strings.Contains(out, "API_KEY") {
		t.Error("expected API_KEY in output")
	}
	if !strings.Contains(out, "OLD_TOKEN") {
		t.Error("expected OLD_TOKEN in output")
	}
}
