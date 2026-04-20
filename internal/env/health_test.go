package env

import (
	"strings"
	"testing"
)

func baseHealthEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "API_SECRET", Value: "supersecretvalue"},
		{Key: "DB_HOST", Value: "localhost"},
	}
}

func TestHealth_NoIssues(t *testing.T) {
	result := Health(baseHealthEntries())
	if len(result.Issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(result.Issues))
	}
	if !strings.Contains(result.Format(), "No health issues") {
		t.Error("expected clean message")
	}
}

func TestHealth_EmptyValueWarn(t *testing.T) {
	entries := []Entry{
		{Key: "APP_NAME", Value: ""},
	}
	result := Health(entries)
	if result.Warns != 1 {
		t.Fatalf("expected 1 warning, got %d", result.Warns)
	}
	if result.Issues[0].Severity != "warn" {
		t.Error("expected warn severity")
	}
}

func TestHealth_EmptySecretIsError(t *testing.T) {
	entries := []Entry{
		{Key: "API_SECRET", Value: ""},
	}
	result := Health(entries)
	if result.Errors != 1 {
		t.Fatalf("expected 1 error, got %d", result.Errors)
	}
}

func TestHealth_DuplicateKey(t *testing.T) {
	entries := []Entry{
		{Key: "APP_NAME", Value: "a"},
		{Key: "APP_NAME", Value: "b"},
	}
	result := Health(entries)
	if result.Errors != 1 {
		t.Fatalf("expected 1 error for duplicate, got %d", result.Errors)
	}
}

func TestHealth_ShortSecret(t *testing.T) {
	entries := []Entry{
		{Key: "API_KEY", Value: "short"},
	}
	result := Health(entries)
	if result.Warns != 1 {
		t.Fatalf("expected 1 warn for short secret, got %d", result.Warns)
	}
}

func TestHealth_Format_ContainsIssues(t *testing.T) {
	entries := []Entry{
		{Key: "DB_PASSWORD", Value: ""},
	}
	result := Health(entries)
	out := result.Format()
	if !strings.Contains(out, "DB_PASSWORD") {
		t.Error("expected key in output")
	}
	if !strings.Contains(out, "error") {
		t.Error("expected severity in output")
	}
}
