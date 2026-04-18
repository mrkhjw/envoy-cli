package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeTempSortEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "sort-*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func runSortCmd(args ...string) (string, error) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(append([]string{"sort"}, args...))
	_, err := rootCmd.ExecuteC()
	return buf.String(), err
}

func resetSortCmd() {
	sortReverse = false
	sortByValue = false
	sortSecretsLast = false
	sortOutput = ""
	sortMask = false
	// reset cobra state
	sortCmd.ResetFlags()
	_ = cobra.Command{}
}

func TestSortCmd_BasicOutput(t *testing.T) {
	f := writeTempSortEnv(t, "ZEBRA=z\nAPPLE=a\nMANGO=m\n")
	out, err := runSortCmd(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	idx1 := strings.Index(out, "APPLE")
	idx2 := strings.Index(out, "ZEBRA")
	if idx1 > idx2 {
		t.Error("expected APPLE before ZEBRA in output")
	}
}

func TestSortCmd_MissingFile(t *testing.T) {
	_, err := runSortCmd("/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
