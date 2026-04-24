package env

import (
	"strings"
	"testing"
)

func baseTokenizeEntries() []Entry {
	return []Entry{
		{Key: "ALLOWED_HOSTS", Value: "localhost, 127.0.0.1, example.com"},
		{Key: "SECRET_KEY", Value: "val1,val2,val3"},
		{Key: "PLAIN", Value: "single"},
		{Key: "", Value: "", Comment: true},
	}
}

func TestTokenize_DefaultDelimiter(t *testing.T) {
	entries := baseTokenizeEntries()
	res := Tokenize(entries, TokenizeOptions{})

	if res.Total != 3 {
		t.Errorf("expected 3 total, got %d", res.Total)
	}
	toks, ok := res.Tokens["ALLOWED_HOSTS"]
	if !ok || len(toks) != 3 {
		t.Errorf("expected 3 tokens for ALLOWED_HOSTS, got %v", toks)
	}
}

func TestTokenize_CustomDelimiter(t *testing.T) {
	entries := []Entry{
		{Key: "PATHS", Value: "/usr/bin:/usr/local/bin:/opt/bin"},
	}
	res := Tokenize(entries, TokenizeOptions{Delimiter: ":"})

	toks := res.Tokens["PATHS"]
	if len(toks) != 3 {
		t.Errorf("expected 3 path tokens, got %v", toks)
	}
}

func TestTokenize_SpecificKeys(t *testing.T) {
	entries := baseTokenizeEntries()
	res := Tokenize(entries, TokenizeOptions{Keys: []string{"ALLOWED_HOSTS"}})

	if res.Total != 1 {
		t.Errorf("expected 1 total, got %d", res.Total)
	}
	if res.Skipped != 2 {
		t.Errorf("expected 2 skipped, got %d", res.Skipped)
	}
	if _, ok := res.Tokens["PLAIN"]; ok {
		t.Error("PLAIN should not be tokenized")
	}
}

func TestTokenize_SingleValue(t *testing.T) {
	entries := []Entry{{Key: "PLAIN", Value: "single"}}
	res := Tokenize(entries, TokenizeOptions{})

	toks := res.Tokens["PLAIN"]
	if len(toks) != 1 || toks[0] != "single" {
		t.Errorf("expected single token, got %v", toks)
	}
}

func TestTokenizeResult_Format_MasksSecrets(t *testing.T) {
	res := TokenizeResult{
		Tokens:  map[string][]string{"SECRET_KEY": {"a", "b"}, "HOST": {"localhost"}},
		Total:   2,
		Skipped: 0,
	}
	out := res.Format(true)
	if strings.Contains(out, "a") || strings.Contains(out, "b") {
		t.Error("expected secret tokens to be masked")
	}
	if !strings.Contains(out, "REDACTED") {
		t.Error("expected REDACTED in output")
	}
	if !strings.Contains(out, "localhost") {
		t.Error("expected HOST value to be visible")
	}
}

func TestTokenize_SkipsComments(t *testing.T) {
	entries := []Entry{
		{Key: "", Value: "# this is a comment", Comment: true},
		{Key: "A", Value: "x,y"},
	}
	res := Tokenize(entries, TokenizeOptions{})
	if res.Total != 1 {
		t.Errorf("expected 1, got %d", res.Total)
	}
}
