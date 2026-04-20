package env

import "testing"

func baseLookupEntries() []Entry {
	return []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "API_SECRET", Value: "topsecret"},
		{Key: "PORT", Value: "3000"},
		{IsComment: true, Raw: "# ignored"},
	}
}

func TestLookup_Found(t *testing.T) {
	r := Lookup(baseLookupEntries(), "DB_HOST", LookupOptions{CaseSensitive: true})
	if !r.Found {
		t.Fatal("expected key to be found")
	}
	if r.Value != "localhost" {
		t.Errorf("expected localhost, got %q", r.Value)
	}
}

func TestLookup_NotFound(t *testing.T) {
	r := Lookup(baseLookupEntries(), "MISSING", LookupOptions{CaseSensitive: true})
	if r.Found {
		t.Fatal("expected key not to be found")
	}
}

func TestLookup_CaseInsensitive(t *testing.T) {
	r := Lookup(baseLookupEntries(), "db_host", LookupOptions{CaseSensitive: false})
	if !r.Found {
		t.Fatal("expected case-insensitive match")
	}
	if r.Value != "localhost" {
		t.Errorf("expected localhost, got %q", r.Value)
	}
}

func TestLookup_MasksSecret(t *testing.T) {
	r := Lookup(baseLookupEntries(), "API_SECRET", LookupOptions{CaseSensitive: true, MaskSecrets: true})
	if !r.Found {
		t.Fatal("expected key to be found")
	}
	if r.Value != "***" {
		t.Errorf("expected masked value, got %q", r.Value)
	}
	if !r.Masked {
		t.Error("expected Masked=true")
	}
}

func TestLookup_NonSecretNotMasked(t *testing.T) {
	r := Lookup(baseLookupEntries(), "PORT", LookupOptions{CaseSensitive: true, MaskSecrets: true})
	if r.Value != "3000" {
		t.Errorf("expected 3000, got %q", r.Value)
	}
	if r.Masked {
		t.Error("expected Masked=false for non-secret")
	}
}

func TestLookupResult_Format_Found(t *testing.T) {
	r := LookupResult{Key: "PORT", Value: "3000", Found: true}
	out := r.Format()
	if out != "PORT=3000" {
		t.Errorf("unexpected format: %q", out)
	}
}

func TestLookupResult_Format_NotFound(t *testing.T) {
	r := LookupResult{Key: "MISSING", Found: false}
	out := r.Format()
	if out != `key "MISSING" not found` {
		t.Errorf("unexpected format: %q", out)
	}
}

func TestLookupResult_Format_Masked(t *testing.T) {
	r := LookupResult{Key: "API_SECRET", Value: "***", Found: true, Masked: true}
	out := r.Format()
	if out != "API_SECRET=*** (masked)" {
		t.Errorf("unexpected format: %q", out)
	}
}
