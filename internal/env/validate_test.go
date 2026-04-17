package env

import (
	"testing"
)

func TestValidate_ValidFile(t *testing.T) {
	lines := []string{
		"APP_NAME=myapp",
		"PORT=8080",
		"# a comment",
		"",
		"SECRET_KEY=abc123",
	}
	result := Validate(nil, lines)
	if !result.Valid() {
		t.Errorf("expected valid, got errors: %s", result.Summary())
	}
}

func TestValidate_DuplicateKey(t *testing.T) {
	lines := []string{
		"APP_NAME=myapp",
		"APP_NAME=other",
	}
	result := Validate(nil, lines)
	if result.Valid() {
		t.Error("expected invalid due to duplicate key")
	}
	if len(result.Errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(result.Errors))
	}
}

func TestValidate_EmptyValue(t *testing.T) {
	lines := []string{"API_KEY="}
	result := Validate(nil, lines)
	if result.Valid() {
		t.Error("expected invalid due to empty value")
	}
}

func TestValidate_InvalidFormat(t *testing.T) {
	lines := []string{"NOTAVALIDLINE"}
	result := Validate(nil, lines)
	if result.Valid() {
		t.Error("expected invalid due to bad format")
	}
}

func TestValidate_EmptyKey(t *testing.T) {
	lines := []string{"=value"}
	result := Validate(nil, lines)
	if result.Valid() {
		t.Error("expected invalid due to empty key")
	}
}

func TestValidationResult_Summary_Valid(t *testing.T) {
	r := &ValidationResult{}
	s := r.Summary()
	if s != "✔ No validation errors found." {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestValidationResult_Summary_Invalid(t *testing.T) {
	r := &ValidationResult{
		Errors: []ValidationError{
			{Line: 1, Key: "FOO", Message: "empty value"},
		},
	}
	s := r.Summary()
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
