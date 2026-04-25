package env

import (
	"strings"
	"testing"
)

var baseObfuscateEntries = []Entry{
	{Key: "APP_NAME", Value: "myapp"},
	{Key: "API_SECRET", Value: "supersecret"},
	{Key: "DB_PASSWORD", Value: "hunter2"},
	{Key: "PORT", Value: "8080"},
	{Key: "", Value: "# a comment", Comment: true},
}

func TestObfuscate_StarStyle(t *testing.T) {
	res := Obfuscate(baseObfuscateEntries, ObfuscateOptions{Style: "star"})
	for _, e := range res.Entries {
		if e.Comment {
			continue
		}
		if isSecret(e.Key) {
			if !strings.ContainsAny(e.Result, "*") {
				t.Errorf("expected %s to be obfuscated, got %s", e.Key, e.Result)
			}
		}
	}
}

func TestObfuscate_HashStyle(t *testing.T) {
	res := Obfuscate(baseObfuscateEntries, ObfuscateOptions{Style: "hash"})
	for _, e := range res.Entries {
		if e.Comment || !isSecret(e.Key) {
			continue
		}
		if !strings.ContainsAny(e.Result, "#") {
			t.Errorf("expected hash obfuscation for %s, got %s", e.Key, e.Result)
		}
		if len(e.Result) != len(e.Original) {
			t.Errorf("hash length mismatch for %s: want %d got %d", e.Key, len(e.Original), len(e.Result))
		}
	}
}

func TestObfuscate_PartialStyle(t *testing.T) {
	entries := []Entry{
		{Key: "API_SECRET", Value: "abcdefgh"},
	}
	res := Obfuscate(entries, ObfuscateOptions{Style: "partial", RevealChars: 3})
	if len(res.Entries) == 0 {
		t.Fatal("expected entries")
	}
	got := res.Entries[0].Result
	if !strings.HasSuffix(got, "fgh") {
		t.Errorf("expected suffix 'fgh', got %s", got)
	}
	if !strings.HasPrefix(got, "*****") {
		t.Errorf("expected leading stars, got %s", got)
	}
}

func TestObfuscate_SpecificKeys(t *testing.T) {
	res := Obfuscate(baseObfuscateEntries, ObfuscateOptions{
		Keys:  []string{"APP_NAME"},
		Style: "star",
	})
	for _, e := range res.Entries {
		if e.Key == "APP_NAME" && !e.Changed {
			t.Errorf("expected APP_NAME to be obfuscated")
		}
		if e.Key == "PORT" && e.Changed {
			t.Errorf("expected PORT to be unchanged")
		}
	}
}

func TestObfuscate_NonSecretsUnchanged(t *testing.T) {
	res := Obfuscate(baseObfuscateEntries, ObfuscateOptions{})
	for _, e := range res.Entries {
		if e.Key == "PORT" && e.Changed {
			t.Errorf("PORT should not be obfuscated")
		}
		if e.Key == "APP_NAME" && e.Changed {
			t.Errorf("APP_NAME should not be obfuscated")
		}
	}
}

func TestObfuscateResult_Format(t *testing.T) {
	res := Obfuscate(baseObfuscateEntries, ObfuscateOptions{Style: "star"})
	out := res.Format()
	if !strings.Contains(out, "obfuscated") {
		t.Errorf("expected summary line, got: %s", out)
	}
	if !strings.Contains(out, "~") {
		t.Errorf("expected changed marker '~' in output, got: %s", out)
	}
}
