package env

import (
	"strings"
	"testing"
)

func baseScope() []Entry {
	return []Entry{
		{Key: "PROD_DB_HOST", Value: "prod.db"},
		{Key: "PROD_API_SECRET", Value: "s3cr3t"},
		{Key: "DEV_DB_HOST", Value: "localhost"},
		{Key: "APP_NAME", Value: "envoy"},
	}
}

func TestScope_FiltersByPrefix(t *testing.T) {
	res := Scope(baseScope(), ScopeOptions{Prefix: "PROD"})
	if len(res.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(res.Entries))
	}
	if res.Total != 4 {
		t.Errorf("expected total 4, got %d", res.Total)
	}
}

func TestScope_StripPrefix(t *testing.T) {
	res := Scope(baseScope(), ScopeOptions{Prefix: "PROD", StripPrefix: true})
	for _, e := range res.Entries {
		if strings.HasPrefix(e.Key, "PROD_") {
			t.Errorf("key %q should have prefix stripped", e.Key)
		}
	}
	keys := map[string]bool{}
	for _, e := range res.Entries {
		keys[e.Key] = true
	}
	if !keys["DB_HOST"] || !keys["API_SECRET"] {
		t.Errorf("unexpected keys: %v", keys)
	}
}

func TestScope_EmptyPrefix(t *testing.T) {
	res := Scope(baseScope(), ScopeOptions{Prefix: ""})
	if len(res.Entries) != 4 {
		t.Errorf("expected all 4 entries, got %d", len(res.Entries))
	}
}

func TestScope_NoMatches(t *testing.T) {
	res := Scope(baseScope(), ScopeOptions{Prefix: "STAGING"})
	if len(res.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(res.Entries))
	}
}

func TestScopeResult_Format_MasksSecrets(t *testing.T) {
	res := Scope(baseScope(), ScopeOptions{Prefix: "PROD", StripPrefix: true})
	out := res.Format(true)
	if strings.Contains(out, "s3cr3t") {
		t.Errorf("expected secret to be masked")
	}
	if !strings.Contains(out, "****") {
		t.Errorf("expected masked placeholder")
	}
}

func TestScopeResult_Format_Empty(t *testing.T) {
	res := Scope(baseScope(), ScopeOptions{Prefix: "STAGING"})
	out := res.Format(false)
	if !strings.Contains(out, "no entries found") {
		t.Errorf("expected no entries message, got: %s", out)
	}
}
