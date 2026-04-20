package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeTempReorderCmdEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "reorder-cmd-*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func runReorderCmd(args ...string) (string, error) {
	buf := &bytes.Buffer{}
	reorderCmd.SetOut(buf)
	reorderCmd.SetErr(buf)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(append([]string{"reorder"}, args...))
	err := rootCmd.Execute()
	return buf.String(), err
}

func resetReorderCmd() {
	reorderCmd.ResetFlags()
	reorderCmd.Flags().String("keys", "", "")
	reorderCmd.Flags().Bool("dry-run", false, "")
	reorderCmd.Flags().String("output", "", "")
}

func TestReorderCmd_DryRun(t *testing.T) {
	f := writeTempReorderCmdEnv(t, "ZEBRA=z\nALPHA=a\n")
	t.Cleanup(func() { os.Remove(f) })

	out, err := runReorderCmd(f, "--keys", "ALPHA", "--dry-run")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "dry-run") {
		t.Errorf("expected dry-run in output, got %q", out)
	}
}

func TestReorderCmd_MissingFile(t *testing.T) {
	_, err := runReorderCmd("/nonexistent/.env", "--keys", "A")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

var _ = reorderCmd.Flags().Lookup // keep cobra import used
var _ *cobra.Command = reorderCmd
