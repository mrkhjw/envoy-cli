package env

import (
	"strings"
	"testing"
)

func TestRedact_MasksSecrets(t *testing.T) {
	vars := map[string]string{
		"APP_NAME":    "myapp",
		"DB_PASSWORD": "supersecret",
		"API_KEY":     "abc123",
	}
	res := Redact(vars, "***")
	lineMap := make(map[string]string)
	for _, l := range res.Lines {
		parts := strings.SplitN(l, "=", 2)
		if len(parts) == 2 {
			lineMap[parts[0]] = parts[1]
		}
	}
	if lineMap["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME=myapp, got %s", lineMap["APP_NAME"])
	}
	if lineMap["DB_PASSWORD"] != "***" {
		t.Errorf("expected DB_PASSWORD=***, got %s", lineMap["DB_PASSWORD"])
	}
	if lineMap["API_KEY"] != "***" {
		t.Errorf("expected API_KEY=***, got %s", lineMap["API_KEY"])
	}
	if res.Redacted != 2 {
		t.Errorf("expected 2 redacted, got %d", res.Redacted)
	}
}

func TestRedact_DefaultPlaceholder(t *testing.T) {
	vars := map[string]string{"SECRET_TOKEN": "tok"}
	res := Redact(vars, "")
	if len(res.Lines) != 1 || !strings.Contains(res.Lines[0], "***") {
		t.Errorf("expected default placeholder ***, got %v", res.Lines)
	}
}

func TestRedactString_MasksInlineSecrets(t *testing.T) {
	input := "APP_NAME=myapp\nDB_PASSWORD=secret\n# comment\nPORT=8080"
	out := RedactString(input, "REDACTED")
	if !strings.Contains(out, "DB_PASSWORD=REDACTED") {
		t.Errorf("expected DB_PASSWORD to be redacted, got:\n%s", out)
	}
	if !strings.Contains(out, "APP_NAME=myapp") {
		t.Errorf("expected APP_NAME to remain, got:\n%s", out)
	}
	if !strings.Contains(out, "# comment") {
		t.Errorf("expected comment to be preserved")
	}
}

func TestRedactString_EmptyInput(t *testing.T) {
	out := RedactString("", "***")
	if out != "" {
		t.Errorf("expected empty output, got %q", out)
	}
}
