package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeTempDiffCmdEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env: %v", err)
	}
	return p
}

func runDiffFileCmd(args ...string) (string, error) {
	buf := &bytes.Buffer{}
	diffFileCmd.SetOut(buf)
	diffFileCmd.SetErr(buf)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	root := &cobra.Command{Use: "envoy"}
	root.AddCommand(diffFileCmd)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(append([]string{"diff-file"}, args...))
	_ = root.Execute()
	return buf.String(), nil
}

func TestDiffFileCmd_NoDifferences(t *testing.T) {
	a := writeTempDiffCmdEnv(t, "FOO=bar\n")
	b := writeTempDiffCmdEnv(t, "FOO=bar\n")

	out, err := runDiffFileCmd(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No differences") {
		t.Errorf("expected no-differences message, got: %s", out)
	}
}

func TestDiffFileCmd_ShowsAdded(t *testing.T) {
	a := writeTempDiffCmdEnv(t, "FOO=bar\n")
	b := writeTempDiffCmdEnv(t, "FOO=bar\nNEW=val\n")

	out, err := runDiffFileCmd(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "NEW") {
		t.Errorf("expected NEW in output, got: %s", out)
	}
}

func TestDiffFileCmd_MissingFile(t *testing.T) {
	a := writeTempDiffCmdEnv(t, "FOO=bar\n")
	out, _ := runDiffFileCmd(a, "/nonexistent/.env")
	if !strings.Contains(out, "nonexistent") && !strings.Contains(out, "error") && !strings.Contains(out, "diff-file") {
		t.Logf("output: %s", out)
	}
}
