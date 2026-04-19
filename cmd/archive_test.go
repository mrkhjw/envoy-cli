package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func writeTempArchiveCmdEnv(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	_ = os.WriteFile(p, []byte("APP=myapp\nSECRET_KEY=abc123\nPORT=9000\n"), 0644)
	return p
}

func runArchiveCmd(args []string) (string, error) {
	buf := &bytes.Buffer{}
	archiveCmd.SetOut(buf)
	archiveCmd.SetErr(buf)
	archiveCmd.SetArgs(args)
	_, err := archiveCmd.ExecuteC()
	return buf.String(), err
}

func resetArchiveCmd() {
	archiveCmd.ResetFlags()
	archiveCmd.Flags().StringVarP(&archiveSource, "file", "f", ".env", "source .env file")
	archiveCmd.Flags().StringVarP(&archiveDest, "out", "o", "env.archive.json", "destination archive file")
	archiveCmd.Flags().StringVarP(&archiveLabel, "label", "l", "default", "label for this archive")
}

func TestArchiveCmd_Basic(t *testing.T) {
	src := writeTempArchiveCmdEnv(t)
	out := filepath.Join(t.TempDir(), "out.json")

	archiveCmd.RunE = archiveCmd.RunE
	archiveCmd.SetArgs([]string{"--file", src, "--out", out, "--label", "ci"})

	buf := &bytes.Buffer{}
	archiveCmd.SetOut(buf)
	archiveCmd.SetErr(buf)

	if err := archiveCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(out); os.IsNotExist(err) {
		t.Error("expected output archive file to be created")
	}
	resetArchiveCmd()
}

func TestArchiveCmd_MissingFile(t *testing.T) {
	c := &cobra.Command{Use: "archive"}
	c.SetArgs([]string{})
	archiveCmd.SetArgs([]string{"--file", "/nonexistent/.env", "--out", "/tmp/x.json"})
	buf := &bytes.Buffer{}
	archiveCmd.SetErr(buf)
	_ = archiveCmd.Execute()
	resetArchiveCmd()
}
