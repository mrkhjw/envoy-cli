package env

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestParseFile_Basic(t *testing.T) {
	path := writeTempEnv(t, "APP_NAME=envoy\nDB_PASSWORD=secret123\n# comment\n\nPORT=8080\n")
	defer os.Remove(path)

	entries, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[1].Secret != true {
		t.Errorf("DB_PASSWORD should be marked as secret")
	}
	if entries[0].Secret != false {
		t.Errorf("APP_NAME should not be marked as secret")
	}
}

func TestParseFile_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `API_KEY="my-quoted-key"`+"\n")
	defer os.Remove(path)

	entries, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[0].Value != "my-quoted-key" {
		t.Errorf("expected unquoted value, got %q", entries[0].Value)
	}
}

func TestIsSecret(t *testing.T) {
	cases := map[string]bool{
		"DB_PASSWORD":   true,
		"AUTH_TOKEN":    true,
		"APP_NAME":      false,
		"PRIVATE_KEY":   true,
		"PORT":          false,
	}
	for key, want := range cases {
		if got := isSecret(key); got != want {
			t.Errorf("isSecret(%q) = %v, want %v", key, got, want)
		}
	}
}
