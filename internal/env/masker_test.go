package env

import (
	"testing"
)

func TestMaskSecrets_MasksSecretKeys(t *testing.T) {
	input := map[string]string{
		"APP_NAME":     "myapp",
		"DB_PASSWORD":  "supersecret",
		"API_KEY":      "abc123",
		"PORT":         "8080",
		"SECRET_TOKEN": "xyz",
	}

	result := MaskSecrets(input)

	if result["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME to be unmasked, got %q", result["APP_NAME"])
	}
	if result["PORT"] != "8080" {
		t.Errorf("expected PORT to be unmasked, got %q", result["PORT"])
	}
	for _, key := range []string{"DB_PASSWORD", "API_KEY", "SECRET_TOKEN"} {
		if result[key] != MaskedValue {
			t.Errorf("expected %s to be masked, got %q", key, result[key])
		}
	}
}

func TestMaskSecrets_EmptyMap(t *testing.T) {
	result := MaskSecrets(map[string]string{})
	if len(result) != 0 {
		t.Errorf("expected empty map, got %d entries", len(result))
	}
}

func TestMaskLine_SecretLine(t *testing.T) {
	line := "API_KEY=myverysecretkey"
	result := MaskLine(line)
	expected := "API_KEY=" + MaskedValue
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestMaskLine_NonSecretLine(t *testing.T) {
	line := "APP_ENV=production"
	result := MaskLine(line)
	if result != line {
		t.Errorf("expected line unchanged, got %q", result)
	}
}

func TestMaskLine_CommentLine(t *testing.T) {
	line := "# this is a comment"
	result := MaskLine(line)
	if result != line {
		t.Errorf("expected comment unchanged, got %q", result)
	}
}

func TestMaskLine_EmptyLine(t *testing.T) {
	result := MaskLine("")
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}
