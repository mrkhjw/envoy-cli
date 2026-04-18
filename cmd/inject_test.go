package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func writeTempInjectEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func runInjectCmd(args ...string) (string, error) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(append([]string{"inject"}, args...))
	err := rootCmd.Execute()
	return buf.String(), err
}

func resetInjectCmd() {
	injectCmd.ResetFlags()
	injectCmd.Flags().Bool("overwrite", false, "Overwrite existing environment variables")
	_ = cobra.Command{}
}

func TestInjectCmd_Basic(t *testing.T) {
	os.Unsetenv("CMD_INJECT_KEY")
	f := writeTempInjectEnv(t, "CMD_INJECT_KEY=val\n")
	out, err := runInjectCmd(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if os.Getenv("CMD_INJECT_KEY") != "val" {
		t.Error("expected CMD_INJECT_KEY=val")
	}
	_ = out
}

func TestInjectCmd_MissingFile(t *testing.T) {
	_, err := runInjectCmd("/no/such/file.env")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
