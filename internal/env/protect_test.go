package env

import (
	"strings"
	"testing"
)

var baseProtectEntries = []Entry{
	{Key: "APP_NAME", Value: "myapp"},
	{Key: "DB_PASSWORD", Value: "s3cr3t"},
	{Key: "API_SECRET", Value: "topsecret"},
	{Key: "PORT", Value: "8080"},
}

func TestProtect_SecretKeysTargeted(t *testing.T) {
	_, result := Protect(baseProtectEntries, ProtectOptions{})
	if len(result.Protected) == 0 {
		t.Fatal("expected secret keys to be protected")
	}
	for _, k := range result.Protected {
		if !isSecret(k) {
			t.Errorf("non-secret key protected: %s", k)
		}
	}
}

func TestProtect_SpecificKeys(t *testing.T) {
	_, result := Protect(baseProtectEntries, ProtectOptions{
		Keys: []string{"APP_NAME", "PORT"},
	})
	if len(result.Protected) != 2 {
		t.Fatalf("expected 2 protected, got %d", len(result.Protected))
	}
}

func TestProtect_DryRunDoesNotTag(t *testing.T) {
	out, result := Protect(baseProtectEntries, ProtectOptions{
		Keys:   []string{"APP_NAME"},
		DryRun: true,
	})
	if len(result.Protected) != 1 {
		t.Fatalf("expected 1 in dry-run protected list")
	}
	for _, e := range out {
		if e.Key == "APP_NAME" && e.Tags != nil && e.Tags["protected"] {
			t.Error("dry-run should not apply tags")
		}
	}
}

func TestProtect_SkipsAlreadyProtected(t *testing.T) {
	entries := []Entry{
		{Key: "DB_PASSWORD", Value: "x", Tags: map[string]bool{"protected": true}},
	}
	_, result := Protect(entries, ProtectOptions{})
	if len(result.Skipped) != 1 {
		t.Fatalf("expected 1 skipped, got %d", len(result.Skipped))
	}
	if len(result.Protected) != 0 {
		t.Fatalf("expected 0 protected, got %d", len(result.Protected))
	}
}

func TestProtectResult_Format(t *testing.T) {
	r := ProtectResult{
		Protected: []string{"DB_PASSWORD", "API_SECRET"},
		Skipped:   []string{"JWT_SECRET"},
	}
	out := r.Format()
	if !strings.Contains(out, "2 protected") {
		t.Errorf("expected count in format, got: %s", out)
	}
	if !strings.Contains(out, "1 skipped") {
		t.Errorf("expected skipped count in format, got: %s", out)
	}
}

func TestProtectResult_Format_NoMatches(t *testing.T) {
	r := ProtectResult{}
	out := r.Format()
	if !strings.Contains(out, "no keys matched") {
		t.Errorf("expected no-match message, got: %s", out)
	}
}
