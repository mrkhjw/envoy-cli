package env

import (
	"strings"
	"testing"
)

func TestSanitize_TrimValues(t *testing.T) {
	entries := []Entry{
		{Key: "APP_NAME", Value: "  myapp  "},
		{Key: "PORT", Value: "8080"},
	}
	result := Sanitize(entries, SanitizeOptions{TrimValues: true})
	if result.Entries[0].Value != "myapp" {
		t.Errorf("expected 'myapp', got %q", result.Entries[0].Value)
	}
	if result.Cleaned != 1 {
		t.Errorf("expected 1 cleaned, got %d", result.Cleaned)
	}
	if result.Unchanged != 1 {
		t.Errorf("expected 1 unchanged, got %d", result.Unchanged)
	}
}

func TestSanitize_TrimKeys(t *testing.T) {
	entries := []Entry{
		{Key: "  DB_HOST  ", Value: "localhost"},
	}
	result := Sanitize(entries, SanitizeOptions{TrimKeys: true})
	if result.Entries[0].Key != "DB_HOST" {
		t.Errorf("expected 'DB_HOST', got %q", result.Entries[0].Key)
	}
	if result.Cleaned != 1 {
		t.Errorf("expected 1 cleaned, got %d", result.Cleaned)
	}
}

func TestSanitize_RemoveNullBytes(t *testing.T) {
	entries := []Entry{
		{Key: "API_KEY", Value: "abc\x00def"},
	}
	result := Sanitize(entries, SanitizeOptions{RemoveNullBytes: true})
	if result.Entries[0].Value != "abcdef" {
		t.Errorf("expected 'abcdef', got %q", result.Entries[0].Value)
	}
	if result.Cleaned != 1 {
		t.Errorf("expected 1 cleaned, got %d", result.Cleaned)
	}
}

func TestSanitize_NormalizeLineEndings(t *testing.T) {
	entries := []Entry{
		{Key: "NOTES", Value: "line1\r\nline2\rline3"},
	}
	result := Sanitize(entries, SanitizeOptions{NormalizeLineEndings: true})
	expected := "line1\nline2\nline3"
	if result.Entries[0].Value != expected {
		t.Errorf("expected %q, got %q", expected, result.Entries[0].Value)
	}
}

func TestSanitize_StripControlChars(t *testing.T) {
	entries := []Entry{
		{Key: "MSG", Value: "hello\x01\x02world"},
	}
	result := Sanitize(entries, SanitizeOptions{StripControlChars: true})
	if result.Entries[0].Value != "helloworld" {
		t.Errorf("expected 'helloworld', got %q", result.Entries[0].Value)
	}
}

func TestSanitize_NoOptions(t *testing.T) {
	entries := []Entry{
		{Key: "KEY", Value: "  value  "},
	}
	result := Sanitize(entries, SanitizeOptions{})
	if result.Entries[0].Value != "  value  " {
		t.Errorf("expected unchanged value, got %q", result.Entries[0].Value)
	}
	if result.Unchanged != 1 {
		t.Errorf("expected 1 unchanged, got %d", result.Unchanged)
	}
}

func TestSanitizeResult_Format_MasksSecrets(t *testing.T) {
	result := SanitizeResult{
		Entries:   []Entry{{Key: "SECRET_KEY", Value: "s3cr3t"}, {Key: "APP", Value: "myapp"}},
		Cleaned:   1,
		Unchanged: 1,
	}
	out := result.Format()
	if !strings.Contains(out, "SECRET_KEY=***") {
		t.Errorf("expected secret to be masked, got: %s", out)
	}
	if !strings.Contains(out, "APP=myapp") {
		t.Errorf("expected APP value visible, got: %s", out)
	}
	if !strings.Contains(out, "sanitized: 1") {
		t.Errorf("expected summary line, got: %s", out)
	}
}
