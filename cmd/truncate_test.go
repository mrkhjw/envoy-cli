package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeTempTruncateEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func runTruncateCmd(args ...string) (string, error) {
	buf := &bytes.Buffer{}
	truncateCmd.SetOut(buf)
	truncateCmd.SetErr(buf)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	cmd := &cobra.Command{Use: "root"}
	cmd.AddCommand(truncateCmd)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(append([]string{"truncate"}, args...))
	err := cmd.Execute()
	return buf.String(), err
}

func TestTruncateCmd_DryRun(t *testing.T) {
	path := writeTempTruncateEnv(t,
		"API_KEY=averylongvaluethatiswaytoolong\nSHORT=ok\n")

	out, err := runTruncateCmd(path, "--max-len=5", "--dry-run")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected API_KEY in output, got: %q", out)
	}

	// File should be unchanged
	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), "averylongvaluethatiswaytoolong") {
		t.Error("dry-run should not modify the file")
	}
}

func TestTruncateCmd_MissingFile(t *testing.T) {
	_, err := runTruncateCmd("/nonexistent/.env", "--max-len=10")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
