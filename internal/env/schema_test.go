package env

import (
	"strings"
	"testing"
)

func TestValidateSchema_AllPresent(t *testing.T) {
	entries := []Entry{{Key: "APP_HOST", Value: "localhost"}, {Key: "APP_PORT", Value: "8080"}}
	schema := []SchemaEntry{
		{Key: "APP_HOST", Required: true},
		{Key: "APP_PORT", Required: true},
	}
	res := ValidateSchema(entries, schema)
	if len(res.Missing) != 0 {
		t.Errorf("expected no missing keys, got %v", res.Missing)
	}
	if len(res.Extra) != 0 {
		t.Errorf("expected no extra keys, got %v", res.Extra)
	}
}

func TestValidateSchema_MissingRequired(t *testing.T) {
	entries := []Entry{{Key: "APP_HOST", Value: "localhost"}}
	schema := []SchemaEntry{
		{Key: "APP_HOST", Required: true},
		{Key: "APP_PORT", Required: true},
	}
	res := ValidateSchema(entries, schema)
	if len(res.Missing) != 1 || res.Missing[0] != "APP_PORT" {
		t.Errorf("expected APP_PORT missing, got %v", res.Missing)
	}
}

func TestValidateSchema_ExtraKeys(t *testing.T) {
	entries := []Entry{{Key: "APP_HOST", Value: "localhost"}, {Key: "UNKNOWN", Value: "x"}}
	schema := []SchemaEntry{{Key: "APP_HOST", Required: true}}
	res := ValidateSchema(entries, schema)
	if len(res.Extra) != 1 || res.Extra[0] != "UNKNOWN" {
		t.Errorf("expected UNKNOWN extra, got %v", res.Extra)
	}
}

func TestValidateSchema_DefaultsApplied(t *testing.T) {
	entries := []Entry{}
	schema := []SchemaEntry{{Key: "LOG_LEVEL", Required: false, Default: "info"}}
	res := ValidateSchema(entries, schema)
	if res.Defaults["LOG_LEVEL"] != "info" {
		t.Errorf("expected default info, got %s", res.Defaults["LOG_LEVEL"])
	}
}

func TestSchemaResult_Format_AllPresent(t *testing.T) {
	res := SchemaResult{Defaults: make(map[string]string)}
	out := res.Format()
	if !strings.Contains(out, "all required keys present") {
		t.Errorf("unexpected format output: %s", out)
	}
}

func TestSchemaResult_Format_Missing(t *testing.T) {
	res := SchemaResult{Missing: []string{"SECRET_KEY"}, Defaults: make(map[string]string)}
	out := res.Format()
	if !strings.Contains(out, "missing: SECRET_KEY") {
		t.Errorf("unexpected format output: %s", out)
	}
}
