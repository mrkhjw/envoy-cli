package env

import (
	"strings"
	"testing"
)

func TestLint_Clean(t *testing.T) {
	lines := []string{"APP_NAME=myapp", "PORT=8080"}
	result := Lint(lines)
	if len(result.Issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(result.Issues))
	}
	if result.HasErrors() {
		t.Error("expected no errors")
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	result := Lint([]string{"app_name=test"})
	if len(result.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(result.Issues))
	}
	if result.Issues[0].Severity != "warning" {
		t.Errorf("expected warning, got %s", result.Issues[0].Severity)
	}
}

func TestLint_EmptyValue(t *testing.T) {
	result := Lint([]string{"API_KEY="})
	if len(result.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(result.Issues))
	}
	if result.Issues[0].Severity != "warning" {
		t.Errorf("expected warning, got %s", result.Issues[0].Severity)
	}
}

func TestLint_InvalidFormat(t *testing.T) {
	result := Lint([]string{"NOTAVALIDLINE"})
	if len(result.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(result.Issues))
	}
	if result.Issues[0].Severity != "error" {
		t.Errorf("expected error, got %s", result.Issues[0].Severity)
	}
	if !result.HasErrors() {
		t.Error("expected HasErrors to be true")
	}
}

func TestLint_DuplicateKey(t *testing.T) {
	result := Lint([]string{"FOO=bar", "FOO=baz"})
	if len(result.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(result.Issues))
	}
	if result.Issues[0].Severity != "error" {
		t.Errorf("expected error, got %s", result.Issues[0].Severity)
	}
}

func TestLint_SkipsComments(t *testing.T) {
	lines := []string{"# this is a comment", "APP=ok"}
	result := Lint(lines)
	if len(result.Issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(result.Issues))
	}
}

func TestLint_SkipsBlankLines(t *testing.T) {
	lines := []string{"", "   ", "APP=ok"}
	result := Lint(lines)
	if len(result.Issues) != 0 {
		t.Fatalf("expected no issues for blank lines, got %d", len(result.Issues))
	}
}

func TestLintResult_Format_NoIssues(t *testing.T) {
	r := &LintResult{}
	if r.Format() != "No lint issues found." {
		t.Errorf("unexpected format output: %s", r.Format())
	}
}

func TestLintResult_Format_WithIssues(t *testing.T) {
	r := Lint([]string{"NOTVALID"})
	out := r.Format()
	if !strings.Contains(out, "[ERROR]") {
		t.Errorf("expected ERROR in output, got: %s", out)
	}
}
