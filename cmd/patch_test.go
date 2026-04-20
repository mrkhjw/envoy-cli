package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeTempPatchEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "patch-*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func runPatchCmd(args ...string) (string, error) {
	buf := new(strings.Builder)
	patchCmd.ResetFlags()
	patchCmd.Flags().BoolVar(&patchDryRun, "dry-run", false, "")
	patchCmd.SetOut(buf)
	patchCmd.SetErr(buf)
	patchCmd.SetArgs(args)
	_, err := patchCmd.ExecuteC()
	return buf.String(), err
}

func resetPatchCmd() {
	patchDryRun = false
	patchCmd.ResetFlags()
	patchCmd.Flags().BoolVar(&patchDryRun, "dry-run", false, "Preview changes without writing")
}

func TestPatchCmd_DryRun(t *testing.T) {
	file := writeTempPatchEnv(t, "PORT=8080\nAPP=myapp\n")
	defer resetPatchCmd()
	_, err := runPatchCmd(file, "--dry-run", "set:PORT=9090")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(file)
	if strings.Contains(string(data), "9090") {
		t.Error("dry-run should not modify file")
	}
}

func TestPatchCmd_MissingFile(t *testing.T) {
	defer resetPatchCmd()
	_, err := runPatchCmd("/nonexistent/.env", "set:X=1")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestPatchCmd_InvalidOp(t *testing.T) {
	file := writeTempPatchEnv(t, "PORT=8080\n")
	defer resetPatchCmd()
	cmd := &cobra.Command{}
	_ = cmd
	_, err := runPatchCmd(file, "badop:KEY")
	if err == nil {
		t.Error("expected error for unknown op")
	}
}
