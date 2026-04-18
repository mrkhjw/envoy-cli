package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeTempLintEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "lint-test-*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func runLintCmd(args []string) (string, error) {
	buf := new(strings.Builder)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(append([]string{"lint"}, args...))
	err := rootCmd.Execute()
	return buf.String(), err
}

func resetLintCmd() {
	rootCmd.SetArgs([]string{})
	_ = cobra.Command{}
}

func TestLintCmd_CleanFile(t *testing.T) {
	path := writeTempLintEnv(t, "APP_NAME=myapp\nPORT=8080\n")
	out, err := runLintCmd([]string{"--file", path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = out
	resetLintCmd()
}

func TestLintCmd_MissingFile(t *testing.T) {
	_, err := runLintCmd([]string{"--file", "/nonexistent/.env"})
	if err == nil {
		t.Error("expected error for missing file")
	}
	resetLintCmd()
}

func TestLintCmd_WithWarnings(t *testing.T) {
	path := writeTempLintEnv(t, "app_name=myapp\n")
	_, err := runLintCmd([]string{"--file", path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resetLintCmd()
}
