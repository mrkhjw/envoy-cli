package env

import (
	"testing"
)

func TestGenerate_SecretKeys(t *testing.T) {
	entries := []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "API_SECRET", Value: "old"},
		{Key: "DB_PASSWORD", Value: "old"},
	}
	out, result, err := Generate(entries, GenerateOptions{Length: 16})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Generated) != 2 {
		t.Errorf("expected 2 generated, got %d", len(result.Generated))
	}
	for _, e := range out {
		if isSecret(e.Key) && e.Value == "old" {
			t.Errorf("expected %s to be regenerated", e.Key)
		}
	}
}

func TestGenerate_SkipsNonSecrets(t *testing.T) {
	entries := []Entry{
		{Key: "APP_NAME", Value: "myapp"},
	}
	_, result, err := Generate(entries, GenerateOptions{Length: 16})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Skipped) != 1 || result.Skipped[0] != "APP_NAME" {
		t.Errorf("expected APP_NAME to be skipped")
	}
}

func TestGenerate_DryRun(t *testing.T) {
	entries := []Entry{
		{Key: "API_SECRET", Value: "original"},
	}
	out, result, err := Generate(entries, GenerateOptions{Length: 16, DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "original" {
		t.Errorf("dry run should not modify value")
	}
	if _, ok := result.Generated["API_SECRET"]; !ok {
		t.Errorf("expected API_SECRET in generated map")
	}
}

func TestGenerate_SpecificKeys(t *testing.T) {
	entries := []Entry{
		{Key: "API_SECRET", Value: "old"},
		{Key: "DB_PASSWORD", Value: "old"},
	}
	_, result, err := Generate(entries, GenerateOptions{Length: 16, Keys: []string{"API_SECRET"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Generated) != 1 {
		t.Errorf("expected 1 generated, got %d", len(result.Generated))
	}
	if _, ok := result.Generated["API_SECRET"]; !ok {
		t.Errorf("expected API_SECRET to be generated")
	}
}

func TestGenerate_HexLength(t *testing.T) {
	entries := []Entry{{Key: "APP_SECRET", Value: ""}}
	out, _, err := Generate(entries, GenerateOptions{Length: 20, Format: "hex"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out[0].Value) != 20 {
		t.Errorf("expected length 20, got %d", len(out[0].Value))
	}
}

func TestGenerate_UnknownFormat(t *testing.T) {
	entries := []Entry{{Key: "APP_SECRET", Value: ""}}
	_, _, err := Generate(entries, GenerateOptions{Length: 16, Format: "base58"})
	if err == nil {
		t.Error("expected error for unknown format")
	}
}
