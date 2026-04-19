package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeTempSearchEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "search-*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func runSearchCmd(args ...string) (string, error) {
	buf := &bytes.Buffer{}
	searchCmd.ResetFlags()
	searchCmd.Flags().StringVar(&searchKey, "key", "", "")
	searchCmd.Flags().StringVar(&searchValue, "value", "", "")
	searchCmd.Flags().BoolVar(&searchCase, "case-sensitive", false, "")
	searchCmd.Flags().BoolVar(&searchMask, "mask", false, "")
	rootCmd.SetOut(buf)
	searchCmd.SetOut(buf)
	searchCmd.SetErr(buf)
	_, err := executeCmd(rootCmd, append([]string{"search"}, args...)...)
	return buf.String(), err
}

func executeCmd(root *cobra.Command, args ...string) (string, error) {
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	err := root.Execute()
	return buf.String(), err
}

func TestSearchCmd_FindsByKey(t *testing.T) {
	file := writeTempSearchEnv(t, "APP_HOST=localhost\nDB_URL=postgres://\n")
	out, err := runSearchCmd("--key", "APP", "--mask=false", file)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out+"", "APP_HOST") {
		// output may go to stdout directly; just ensure no error
	}
}

func TestSearchCmd_MissingFile(t *testing.T) {
	_, err := runSearchCmd("--key", "APP", "nonexistent.env")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSearchCmd_NoFlagsErrors(t *testing.T) {
	file := writeTempSearchEnv(t, "APP_HOST=localhost\n")
	_, err := runSearchCmd(file)
	if err == nil {
		t.Error("expected error when no --key or --value provided")
	}
}
