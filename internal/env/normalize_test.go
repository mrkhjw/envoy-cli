package env

import (
	"strings"
	"testing"
)

func baseNormalizeEntries() []Entry {
	return []Entry{
		{Key: "db_host", Value: "  localhost  "},
		{Key: "api_key", Value: "secret123"},
		{Key: "export PORT", Value: "8080"},
	}
}

func TestNormalize_UppercaseKeys(t *testing.T) {
	entries := baseNormalizeEntries()
	res := Normalize(entries, NormalizeOptions{UppercaseKeys: true})
	if res.Entries[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", res.Entries[0].Key)
	}
	if len(res.Changed) != 3 {
		t.Errorf("expected 3 changed, got %d", len(res.Changed))
	}
}

func TestNormalize_TrimValues(t *testing.T) {
	entries := baseNormalizeEntries()
	res := Normalize(entries, NormalizeOptions{TrimValues: true})
	if res.Entries[0].Value != "localhost" {
		t.Errorf("expected trimmed value, got %q", res.Entries[0].Value)
	}
	if len(res.Changed) != 1 {
		t.Errorf("expected 1 changed, got %d", len(res.Changed))
	}
}

func TestNormalize_QuoteValues(t *testing.T) {
	entries := []Entry{{Key: "FOO", Value: "bar"}}
	res := Normalize(entries, NormalizeOptions{QuoteValues: true})
	if res.Entries[0].Value != `"bar"` {
		t.Errorf("expected quoted value, got %s", res.Entries[0].Value)
	}
}

func TestNormalize_StripExported(t *testing.T) {
	entries := baseNormalizeEntries()
	res := Normalize(entries, NormalizeOptions{StripExported: true})
	if res.Entries[2].Key != "PORT" {
		t.Errorf("expected PORT, got %s", res.Entries[2].Key)
	}
}

func TestNormalize_NoChanges(t *testing.T) {
	entries := []Entry{{Key: "FOO", Value: "bar"}}
	res := Normalize(entries, NormalizeOptions{})
	if len(res.Changed) != 0 {
		t.Errorf("expected no changes")
	}
	if res.Format() != "no changes" {
		t.Errorf("expected 'no changes'")
	}
}

func TestNormalizeResult_Format(t *testing.T) {
	res := NormalizeResult{Changed: []string{"db_host", "api_key"}}
	out := res.Format()
	if !strings.Contains(out, "db_host") || !strings.Contains(out, "api_key") {
		t.Errorf("expected keys in format output, got: %s", out)
	}
}
