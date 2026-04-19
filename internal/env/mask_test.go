package env

import (
	"strings"
	"testing"
)

var baseMaskEntries = []Entry{
	{Key: "APP_NAME", Value: "myapp"},
	{Key: "DB_PASSWORD", Value: "supersecret"},
	{Key: "API_KEY", Value: "abc123xyz"},
	{Key: "PORT", Value: "8080"},
}

func TestMask_DefaultPlaceholder(t *testing.T) {
	result := Mask(baseMaskEntries, MaskOptions{})
	if result.Masked != 2 {
		t.Errorf("expected 2 masked, got %d", result.Masked)
	}
	for _, e := range result.Entries {
		if isSecret(e.Key) && e.Value != "***" {
			t.Errorf("expected *** for %s, got %s", e.Key, e.Value)
		}
	}
}

func TestMask_CustomPlaceholder(t *testing.T) {
	result := Mask(baseMaskEntries, MaskOptions{Placeholder: "[HIDDEN]"})
	for _, e := range result.Entries {
		if isSecret(e.Key) && e.Value != "[HIDDEN]" {
			t.Errorf("expected [HIDDEN] for %s, got %s", e.Key, e.Value)
		}
	}
}

func TestMask_RevealPrefix(t *testing.T) {
	result := Mask(baseMaskEntries, MaskOptions{Placeholder: "***", RevealPrefix: 3})
	for _, e := range result.Entries {
		if e.Key == "DB_PASSWORD" {
			if !strings.HasPrefix(e.Value, "sup") {
				t.Errorf("expected value to start with 'sup', got %s", e.Value)
			}
		}
	}
}

func TestMask_NonSecretsUnchanged(t *testing.T) {
	result := Mask(baseMaskEntries, MaskOptions{})
	for _, e := range result.Entries {
		if e.Key == "APP_NAME" && e.Value != "myapp" {
			t.Errorf("expected APP_NAME unchanged, got %s", e.Value)
		}
		if e.Key == "PORT" && e.Value != "8080" {
			t.Errorf("expected PORT unchanged, got %s", e.Value)
		}
	}
}

func TestMaskResult_Format(t *testing.T) {
	result := Mask(baseMaskEntries, MaskOptions{})
	out := result.Format()
	if !strings.Contains(out, "Masked 2/4") {
		t.Errorf("expected summary line in format output, got: %s", out)
	}
}
