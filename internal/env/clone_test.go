package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempCloneEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestClone_CopiesAllKeys(t *testing.T) {
	src := writeTempCloneEnv(t, "FOO=bar\nBAZ=qux\n")
	dst := filepath.Join(t.TempDir(), ".env")

	res, err := Clone(src, dst, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.KeysCopied != 2 {
		t.Errorf("expected 2 keys copied, got %d", res.KeysCopied)
	}
	if res.Skipped != 0 {
		t.Errorf("expected 0 skipped, got %d", res.Skipped)
	}
}

func TestClone_NoOverwriteSkipsExisting(t *testing.T) {
	src := writeTempCloneEnv(t, "FOO=new\nBAR=baz\n")
	dst := writeTempCloneEnv(t, "FOO=old\n")

	res, err := Clone(src, dst, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.KeysCopied != 1 {
		t.Errorf("expected 1 key copied, got %d", res.KeysCopied)
	}
	if res.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", res.Skipped)
	}

	parsed, _ := ParseFile(dst)
	if parsed["FOO"] != "old" {
		t.Errorf("expected FOO=old, got %s", parsed["FOO"])
	}
}

func TestClone_OverwriteUpdatesExisting(t *testing.T) {
	src := writeTempCloneEnv(t, "FOO=new\n")
	dst := writeTempCloneEnv(t, "FOO=old\n")

	_, err := Clone(src, dst, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	parsed, _ := ParseFile(dst)
	if parsed["FOO"] != "new" {
		t.Errorf("expected FOO=new, got %s", parsed["FOO"])
	}
}

func TestClone_SummaryFormat(t *testing.T) {
	r := CloneResult{Source: "a.env", Destination: "b.env", KeysCopied: 3, Skipped: 1}
	got := r.Summary()
	expected := "Cloned 3 keys from a.env to b.env (1 skipped)"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
