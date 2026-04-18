package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTemplate_BasicSubstitution(t *testing.T) {
	env := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080"}
	res := Template("http://${APP_HOST}:${APP_PORT}", env)
	if res.Rendered != "http://localhost:8080" {
		t.Errorf("unexpected rendered: %s", res.Rendered)
	}
	if res.Replaced != 2 {
		t.Errorf("expected 2 replaced, got %d", res.Replaced)
	}
	if len(res.Missing) != 0 {
		t.Errorf("expected no missing keys")
	}
}

func TestTemplate_MissingKey(t *testing.T) {
	env := map[string]string{"APP_HOST": "localhost"}
	res := Template("http://${APP_HOST}:${APP_PORT}", env)
	if res.Replaced != 1 {
		t.Errorf("expected 1 replaced, got %d", res.Replaced)
	}
	if len(res.Missing) != 1 || res.Missing[0] != "APP_PORT" {
		t.Errorf("expected APP_PORT in missing, got %v", res.Missing)
	}
	if !strings.Contains(res.Rendered, "${APP_PORT}") {
		t.Errorf("expected placeholder preserved in output")
	}
}

func TestTemplate_DuplicateMissingKey(t *testing.T) {
	env := map[string]string{}
	res := Template("${FOO} and ${FOO}", env)
	if len(res.Missing) != 1 {
		t.Errorf("expected deduplicated missing keys, got %v", res.Missing)
	}
}

func TestTemplate_NoPlaceholders(t *testing.T) {
	res := Template("no vars here", map[string]string{})
	if res.Rendered != "no vars here" {
		t.Errorf("unexpected rendered: %s", res.Rendered)
	}
	if res.Replaced != 0 {
		t.Errorf("expected 0 replaced")
	}
}

func TestTemplateFile_WritesOutput(t *testing.T) {
	dir := t.TempDir()
	tmplPath := filepath.Join(dir, "config.tmpl")
	outPath := filepath.Join(dir, "config.out")

	os.WriteFile(tmplPath, []byte("host=${DB_HOST}\nport=${DB_PORT}\n"), 0644)

	env := map[string]string{"DB_HOST": "db.local", "DB_PORT": "5432"}
	res, err := TemplateFile(tmplPath, env, outPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Replaced != 2 {
		t.Errorf("expected 2 replaced, got %d", res.Replaced)
	}
	data, _ := os.ReadFile(outPath)
	if !strings.Contains(string(data), "db.local") {
		t.Errorf("output missing substituted value")
	}
}

func TestTemplateFile_NotFound(t *testing.T) {
	_, err := TemplateFile("/nonexistent/file.tmpl", map[string]string{}, "/tmp/out")
	if err == nil {
		t.Error("expected error for missing template file")
	}
}

func TestTemplateResult_Format(t *testing.T) {
	res := TemplateResult{Replaced: 3, Missing: []string{"SECRET_KEY"}}
	out := res.Format()
	if !strings.Contains(out, "3") || !strings.Contains(out, "SECRET_KEY") {
		t.Errorf("unexpected format output: %s", out)
	}
}
