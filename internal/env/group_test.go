package env

import (
	"strings"
	"testing"
)

func TestGroup_ByPrefix(t *testing.T) {
	entries := []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "APP_NAME", Value: "envoy"},
	}
	result := Group(entries, "_")
	if len(result.Groups["DB"]) != 2 {
		t.Errorf("expected 2 DB keys, got %d", len(result.Groups["DB"]))
	}
	if len(result.Groups["APP"]) != 1 {
		t.Errorf("expected 1 APP key, got %d", len(result.Groups["APP"]))
	}
}

func TestGroup_Ungrouped(t *testing.T) {
	entries := []Entry{
		{Key: "PORT", Value: "8080"},
	}
	result := Group(entries, "_")
	if _, ok := result.Groups["(ungrouped)"]; !ok {
		t.Error("expected ungrouped key")
	}
}

func TestGroup_EmptySeparatorDefaultsToUnderscore(t *testing.T) {
	entries := []Entry{
		{Key: "REDIS_URL", Value: "redis://localhost"},
	}
	result := Group(entries, "")
	if len(result.Groups["REDIS"]) != 1 {
		t.Errorf("expected REDIS group, got %+v", result.Groups)
	}
}

func TestGroupResult_Format_MasksSecrets(t *testing.T) {
	entries := []Entry{
		{Key: "DB_PASSWORD", Value: "supersecret"},
		{Key: "DB_HOST", Value: "localhost"},
	}
	result := Group(entries, "_")
	out := result.Format(entries, true)
	if strings.Contains(out, "supersecret") {
		t.Error("expected secret to be masked")
	}
	if !strings.Contains(out, "localhost") {
		t.Error("expected non-secret value to be visible")
	}
}

func TestGroupResult_Format_ShowsValues(t *testing.T) {
	entries := []Entry{
		{Key: "APP_NAME", Value: "envoy"},
	}
	result := Group(entries, "_")
	out := result.Format(entries, false)
	if !strings.Contains(out, "envoy") {
		t.Error("expected value in output")
	}
}
