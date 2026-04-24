package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeTempAnnotateCmdEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "annotate-cmd-*.env")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(content)
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func runAnnotateCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	buf := &bytes.Buffer{}
	annotateCmd.ResetFlags()
	annotateCmd.Flags().String("note", "", "")
	annotateCmd.Flags().String("keys", "", "")
	annotateCmd.Flags().Bool("dry-run", false, "")
	annotateCmd.SetOut(buf)
	annotateCmd.SetErr(buf)
	root := &cobra.Command{Use: "root"}
	root.AddCommand(annotateCmd)
	root.SetArgs(append([]string{"annotate"}, args...))
	err := root.Execute()
	return buf.String(), err
}

func TestAnnotateCmd_DryRun(t *testing.T) {
	path := writeTempAnnotateCmdEnv(t, "APP_ENV=staging\nAPI_KEY=secret123\n")
	out, err := runAnnotateCmd(t, path, "--note", "reviewed", "--dry-run")
	if err != nil {
		t.Fatalf("unexpected error: %v — output: %s", err, out)
	}
}

func TestAnnotateCmd_MissingNote(t *testing.T) {
	path := writeTempAnnotateCmdEnv(t, "APP_ENV=test\n")
	_, err := runAnnotateCmd(t, path)
	if err == nil {
		t.Error("expected error when --note is missing")
	}
}

func TestAnnotateCmd_MissingFile(t *testing.T) {
	_, err := runAnnotateCmd(t, "/nonexistent/.env", "--note", "test")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestAnnotateCmd_WritesFile(t *testing.T) {
	path := writeTempAnnotateCmdEnv(t, "APP_NAME=myapp\n")
	_, err := runAnnotateCmd(t, path, "--note", "auto")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), "# auto") {
		t.Errorf("expected annotation written to file, got: %s", string(data))
	}
}
