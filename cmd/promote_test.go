package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeTempPromoteCmdEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	os.WriteFile(p, []byte(content), 0644)
	return p
}

func runPromoteCmd(args ...string) (string, error) {
	buf := &bytes.Buffer{}
	promoteCmd.SetOut(buf)
	promoteCmd.SetErr(buf)
	promoteOverwrite = false
	promoteKeys = nil
	rootCmd.SetArgs(append([]string{"promote"}, args...))
	_, err := rootCmd.ExecuteC()
	return buf.String(), err
}

func resetPromoteCmd() {
	promoteCmd.ResetFlags()
	promoteCmd.Flags().BoolVarP(&promoteOverwrite, "overwrite", "o", false, "Overwrite existing keys in destination")
	promoteCmd.Flags().StringSliceVarP(&promoteKeys, "keys", "k", nil, "Specific keys to promote (default: all)")
}

func TestPromoteCmd_Basic(t *testing.T) {
	_ = &cobra.Command{} // ensure import used
	src := writeTempPromoteCmdEnv(t, "FOO=bar\nBAZ=qux\n")
	dst := writeTempPromoteCmdEnv(t, "")
	out, err := runPromoteCmd(src, dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Promoted") && !strings.Contains(out, "FOO") {
		t.Errorf("unexpected output: %s", out)
	}
	resetPromoteCmd()
}

func TestPromoteCmd_MissingFile(t *testing.T) {
	_, err := runPromoteCmd("/nonexistent/.env", "/also/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing source file")
	}
	resetPromoteCmd()
}
