package env

import (
	"strings"
	"testing"
)

var baseChainEntries = []Entry{
	{Key: "app_name", Value: "  myapp  "},
	{Key: "api_secret", Value: "supersecret"},
	{Key: "db_host", Value: " localhost "},
	{Key: "# comment", Value: "", Comment: true},
}

func TestChain_TrimValues(t *testing.T) {
	res := Chain(baseChainEntries, ChainOptions{TrimValues: true})
	for _, e := range res.Entries {
		if e.Comment {
			continue
		}
		if strings.HasPrefix(e.Value, " ") || strings.HasSuffix(e.Value, " ") {
			t.Errorf("expected trimmed value for %s, got %q", e.Key, e.Value)
		}
	}
	if len(res.Log) == 0 {
		t.Error("expected log entries for trim")
	}
}

func TestChain_UpperKeys(t *testing.T) {
	res := Chain(baseChainEntries, ChainOptions{UpperKeys: true})
	for _, e := range res.Entries {
		if e.Comment {
			continue
		}
		for _, c := range e.Key {
			if c >= 'a' && c <= 'z' {
				t.Errorf("expected uppercase key, got %q", e.Key)
				break
			}
		}
	}
}

func TestChain_MaskSecrets(t *testing.T) {
	res := Chain(baseChainEntries, ChainOptions{MaskSecrets: true})
	for _, e := range res.Entries {
		if isSecret(e.Key) && e.Value != "****" {
			t.Errorf("expected masked value for %s, got %q", e.Key, e.Value)
		}
	}
}

func TestChain_MultipleSteps(t *testing.T) {
	res := Chain(baseChainEntries, ChainOptions{TrimValues: true, UpperKeys: true, MaskSecrets: true})
	if len(res.Log) == 0 {
		t.Error("expected log entries for multiple steps")
	}
	steps := map[int]bool{}
	for _, l := range res.Log {
		steps[l.Step] = true
	}
	if len(steps) < 2 {
		t.Errorf("expected entries from multiple steps, got steps: %v", steps)
	}
}

func TestChain_NoOptions(t *testing.T) {
	res := Chain(baseChainEntries, ChainOptions{})
	if len(res.Log) != 0 {
		t.Errorf("expected empty log with no options, got %d entries", len(res.Log))
	}
}

func TestChainResult_Format_WithChanges(t *testing.T) {
	res := Chain(baseChainEntries, ChainOptions{TrimValues: true, MaskSecrets: true})
	out := res.Format()
	if !strings.Contains(out, "transformation") {
		t.Errorf("expected 'transformation' in format output, got: %s", out)
	}
	if !strings.Contains(out, "trim") {
		t.Errorf("expected 'trim' action in format output")
	}
}

func TestChainResult_Format_NoChanges(t *testing.T) {
	res := Chain(baseChainEntries, ChainOptions{})
	out := res.Format()
	if !strings.Contains(out, "no transformations") {
		t.Errorf("expected 'no transformations' message, got: %s", out)
	}
}
