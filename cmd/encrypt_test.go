package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeTempEncryptEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func runEncryptCmd(args []string) (string, error) {
	buf := new(strings.Builder)
	encryptCmd.SetOut(buf)
	encryptCmd.SetErr(buf)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(append([]string{"encrypt"}, args...))
	err := rootCmd.Execute()
	return buf.String(), err
}

func resetEncryptCmd() {
	encryptCmd.ResetFlags()
	encryptCmd.Flags().String("key", "", "Encryption passphrase")
	encryptCmd.Flags().StringSlice("keys", nil, "Specific keys")
	encryptCmd.Flags().Bool("dry-run", false, "Preview")
}

func TestEncryptCmd_DryRun(t *testing.T) {
	p := writeTempEncryptEnv(t, "API_SECRET=topsecret\nAPP_NAME=myapp\n")
	_ = p
	// dry-run should not error and should print preview
	rootCmd.SetArgs([]string{"encrypt", "--key", "passphrase12345678901234567890ab", "--dry-run", p})
	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestEncryptCmd_MissingFile(t *testing.T) {
	rootCmd.SetArgs([]string{"encrypt", "--key", "somekey", "/nonexistent/.env"})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestEncryptCmd_MissingKey(t *testing.T) {
	p := writeTempEncryptEnv(t, "API_SECRET=val\n")
	os.Unsetenv("ENVOY_ENCRYPT_KEY")
	_ = cobra.Command{}
	rootCmd.SetArgs([]string{"encrypt", p})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when no key provided")
	}
}
