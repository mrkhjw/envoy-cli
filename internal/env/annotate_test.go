package env

import (
	"os"
	"strings"
	"testing"
)

func writeTempAnnotateEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "annotate-*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestAnnotate_AllKeys(t *testing.T) {
	entries := []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASSWORD", Value: "secret"},
	}
	result := Annotate(entries, "managed", nil, false)
	if result.Modified != 2 {
		t.Errorf("expected 2 modified, got %d", result.Modified)
	}
	for _, ae := range result.Entries {
		if ae.Annotation != "managed" {
			t.Errorf("expected annotation 'managed' on %s, got %q", ae.Key, ae.Annotation)
		}
	}
}

func TestAnnotate_SpecificKeys(t *testing.T) {
	entries := []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_HOST", Value: "localhost"},
	}
	result := Annotate(entries, "reviewed", []string{"APP_NAME"}, false)
	if result.Modified != 1 {
		t.Errorf("expected 1 modified, got %d", result.Modified)
	}
	if result.Entries[0].Annotation != "reviewed" {
		t.Errorf("expected annotation on APP_NAME")
	}
	if result.Entries[1].Annotation != "" {
		t.Errorf("expected no annotation on DB_HOST")
	}
}

func TestAnnotate_DryRunDoesNotCount(t *testing.T) {
	entries := []Entry{
		{Key: "API_KEY", Value: "xyz"},
	}
	result := Annotate(entries, "note", nil, true)
	if result.Modified != 0 {
		t.Errorf("dry run should not increment Modified, got %d", result.Modified)
	}
}

func TestAnnotate_Format(t *testing.T) {
	r := AnnotateResult{Modified: 3}
	if !strings.Contains(r.Format(), "3") {
		t.Errorf("expected format to mention count")
	}
	r2 := AnnotateResult{Modified: 0}
	if !strings.Contains(r2.Format(), "No") {
		t.Errorf("expected 'No' in empty result format")
	}
}

func TestAnnotateFile_WritesAnnotations(t *testing.T) {
	path := writeTempAnnotateEnv(t, "APP_ENV=production\nSECRET_KEY=abc123\n")
	result, err := AnnotateFile(path, "auto-generated", nil, false)
	if err != nil {
		t.Fatal(err)
	}
	if result.Modified != 2 {
		t.Errorf("expected 2 modified, got %d", result.Modified)
	}
	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), "# auto-generated") {
		t.Errorf("expected annotation in written file")
	}
}
