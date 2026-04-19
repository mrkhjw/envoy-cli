package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeTempTagEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "tag_test_*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func runTagCmd(args []string) (string, error) {
	buf := &bytes.Buffer{}
	tagCmd.ResetFlags()
	tagCmd.Flags().StringVarP(&tagFile, "file", "f", "", "")
	tagCmd.Flags().StringArrayVarP(&tagTags, "tag", "t", []string{}, "")
	tagCmd.Flags().StringArrayVarP(&tagKeys, "key", "k", []string{}, "")
	tagCmd.Flags().BoolVar(&tagMaskSecrets, "mask-secrets", false, "")
	tagCmd.SetOut(buf)
	tagCmd.SetErr(buf)
	root := &cobra.Command{Use: "root"}
	root.AddCommand(tagCmd)
	root.SetArgs(append([]string{"tag"}, args...))
	err := root.Execute()
	return buf.String(), err
}

func TestTagCmd_Basic(t *testing.T) {
	f := writeTempTagEnv(t, "APP_NAME=myapp\nPORT=8080\n")
	out, err := runTagCmd([]string{"--file", f, "--tag", "prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v (output: %s)", err, out)
	}
	if !strings.Contains(out, "APP_NAME") {
		t.Error("expected APP_NAME in output")
	}
	if !strings.Contains(out, "tags:prod") {
		t.Error("expected tags:prod in output")
	}
}

func TestTagCmd_MissingFile(t *testing.T) {
	_, err := runTagCmd([]string{"--file", "/nonexistent.env", "--tag", "prod"})
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestTagCmd_MissingTag(t *testing.T) {
	f := writeTempTagEnv(t, "APP_NAME=myapp\n")
	_, err := runTagCmd([]string{"--file", f})
	if err == nil {
		t.Error("expected error when no tags provided")
	}
}
