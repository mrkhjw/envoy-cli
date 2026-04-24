package env

import (
	"strings"
	"testing"
)

var baseFormatEntries = []Entry{
	{Key: "app_name", Value: "envoy"},
	{Key: "db_password", Value: "secret123"},
	{Key: "port", Value: "8080"},
}

func TestFormat_NoOptions(t *testing.T) {
	res := Format(baseFormatEntries, FormatOptions{})
	if res.Total != 3 {
		t.Fatalf("expected 3 lines, got %d", res.Total)
	}
	if res.Lines[0] != "app_name=envoy" {
		t.Errorf("unexpected line: %s", res.Lines[0])
	}
}

func TestFormat_UppercaseKeys(t *testing.T) {
	res := Format(baseFormatEntries, FormatOptions{UppercaseKeys: true})
	for _, l := range res.Lines {
		key := strings.Split(l, "=")[0]
		if key != strings.ToUpper(key) {
			t.Errorf("expected uppercase key, got %s", key)
		}
	}
}

func TestFormat_ExportPrefix(t *testing.T) {
	res := Format(baseFormatEntries, FormatOptions{ExportPrefix: true})
	for _, l := range res.Lines {
		if !strings.HasPrefix(l, "export ") {
			t.Errorf("expected export prefix, got %s", l)
		}
	}
}

func TestFormat_QuoteValues(t *testing.T) {
	res := Format(baseFormatEntries, FormatOptions{QuoteValues: true})
	for _, l := range res.Lines {
		val := strings.SplitN(l, "=", 2)[1]
		if !strings.HasPrefix(val, "\"") {
			t.Errorf("expected quoted value, got %s", val)
		}
	}
}

func TestFormat_MaskSecrets(t *testing.T) {
	res := Format(baseFormatEntries, FormatOptions{MaskSecrets: true})
	for _, l := range res.Lines {
		parts := strings.SplitN(l, "=", 2)
		if isSecret(parts[0]) && parts[1] != "****" {
			t.Errorf("expected masked value for %s, got %s", parts[0], parts[1])
		}
	}
}

func TestFormat_String(t *testing.T) {
	res := Format(baseFormatEntries, FormatOptions{})
	out := res.String()
	if !strings.Contains(out, "\n") {
		t.Error("expected newline-separated output")
	}
}

func TestFormat_EmptyEntries(t *testing.T) {
	res := Format([]Entry{}, FormatOptions{})
	if res.Total != 0 {
		t.Fatalf("expected 0 lines, got %d", res.Total)
	}
	if len(res.Lines) != 0 {
		t.Errorf("expected empty lines slice, got %v", res.Lines)
	}
	out := res.String()
	if out != "" {
		t.Errorf("expected empty string output, got %q", out)
	}
}
