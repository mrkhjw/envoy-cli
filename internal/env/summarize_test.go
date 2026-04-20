package env

import (
	"strings"
	"testing"
)

var baseSummarizeEntries = []Entry{
	{Key: "APP_NAME", Value: "myapp"},
	{Key: "APP_ENV", Value: "production"},
	{Key: "DB_HOST", Value: "localhost"},
	{Key: "DB_PASSWORD", Value: "s3cr3t"},
	{Key: "API_SECRET", Value: ""},
	{Key: "# a comment", Value: ""},
}

func TestSummarize_Total(t *testing.T) {
	r := Summarize(baseSummarizeEntries, "_")
	if r.Total != 5 {
		t.Errorf("expected Total=5, got %d", r.Total)
	}
}

func TestSummarize_Secrets(t *testing.T) {
	r := Summarize(baseSummarizeEntries, "_")
	// DB_PASSWORD and API_SECRET are secrets
	if r.Secrets != 2 {
		t.Errorf("expected Secrets=2, got %d", r.Secrets)
	}
}

func TestSummarize_Empty(t *testing.T) {
	r := Summarize(baseSummarizeEntries, "_")
	if r.Empty != 1 {
		t.Errorf("expected Empty=1, got %d", r.Empty)
	}
}

func TestSummarize_Comments(t *testing.T) {
	r := Summarize(baseSummarizeEntries, "_")
	if r.Comments != 1 {
		t.Errorf("expected Comments=1, got %d", r.Comments)
	}
}

func TestSummarize_Groups(t *testing.T) {
	r := Summarize(baseSummarizeEntries, "_")
	if r.Groups["APP"] != 2 {
		t.Errorf("expected APP group count=2, got %d", r.Groups["APP"])
	}
	if r.Groups["DB"] != 2 {
		t.Errorf("expected DB group count=2, got %d", r.Groups["DB"])
	}
}

func TestSummarize_DefaultSeparator(t *testing.T) {
	r := Summarize(baseSummarizeEntries, "")
	if r.Groups["APP"] != 2 {
		t.Errorf("expected APP group with default separator, got %d", r.Groups["APP"])
	}
}

func TestSummarizeResult_Format_MasksSecrets(t *testing.T) {
	r := Summarize(baseSummarizeEntries, "_")
	out := r.Format(true)
	if strings.Contains(out, "s3cr3t") {
		t.Error("expected secret value to be masked")
	}
	if !strings.Contains(out, "****") {
		t.Error("expected masked placeholder in output")
	}
}

func TestSummarizeResult_Format_ShowsValues(t *testing.T) {
	r := Summarize(baseSummarizeEntries, "_")
	out := r.Format(false)
	if !strings.Contains(out, "s3cr3t") {
		t.Error("expected plain value in unmasked output")
	}
	if !strings.Contains(out, "Total keys") {
		t.Error("expected header line in output")
	}
}
