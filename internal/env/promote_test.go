package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempPromoteEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	os.WriteFile(p, []byte(content), 0644)
	return p
}

func TestPromote_AddsNewKeys(t *testing.T) {
	src := map[string]string{"FOO": "bar", "BAZ": "qux"}
	dst := map[string]string{}
	result := Promote(src, dst, nil, false)
	if len(result.Added) != 2 {
		t.Errorf("expected 2 added, got %d", len(result.Added))
	}
}

func TestPromote_SkipsExistingWithoutOverwrite(t *testing.T) {
	src := map[string]string{"FOO": "new"}
	dst := map[string]string{"FOO": "old"}
	result := Promote(src, dst, nil, false)
	if len(result.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(result.Skipped))
	}
	if dst["FOO"] != "old" {
		t.Error("expected FOO to remain old")
	}
}

func TestPromote_OverwriteUpdatesExisting(t *testing.T) {
	src := map[string]string{"FOO": "new"}
	dst := map[string]string{"FOO": "old"}
	result := Promote(src, dst, nil, true)
	if len(result.Updated) != 1 {
		t.Errorf("expected 1 updated, got %d", len(result.Updated))
	}
	if dst["FOO"] != "new" {
		t.Error("expected FOO to be updated")
	}
}

func TestPromote_SpecificKeys(t *testing.T) {
	src := map[string]string{"FOO": "1", "BAR": "2"}
	dst := map[string]string{}
	result := Promote(src, dst, []string{"FOO"}, false)
	if len(result.Added) != 1 || result.Added[0] != "FOO" {
		t.Error("expected only FOO to be added")
	}
	if _, ok := dst["BAR"]; ok {
		t.Error("BAR should not be promoted")
	}
}

func TestPromoteResult_Format(t *testing.T) {
	r := PromoteResult{
		Added:   []string{"NEW_KEY"},
		Updated: []string{"UPDATED_KEY"},
		Skipped: []string{"SKIP_KEY"},
	}
	out := r.Format()
	if !strings.Contains(out, "+ NEW_KEY") {
		t.Error("expected added key in format")
	}
	if !strings.Contains(out, "~ UPDATED_KEY") {
		t.Error("expected updated key in format")
	}
	if !strings.Contains(out, "= SKIP_KEY (skipped)") {
		t.Error("expected skipped key in format")
	}
}
