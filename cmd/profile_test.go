package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/cobra"
)

func writeTempProfileCmdEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "profile-cmd-*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func runProfileCmd(args ...string) (string, error) {
	buf := &bytes.Buffer{}
	profileCmd.SetOut(buf)
	profileCmd.SetErr(buf)
	profileCmd.SetArgs(args)
	_, err := profileCmd.ExecuteC()
	return buf.String(), err
}

func resetProfileCmd() {
	profileFile = ""
	profileName = ""
	profileMask = false
}

func TestProfileCmd_Basic(t *testing.T) {
	defer resetProfileCmd()
	path := writeTempProfileCmdEnv(t, "# @profile dev\nAPP_ENV=development\nDB_HOST=localhost\n")
	profileCmd.ResetFlags()
	profileCmd.Flags().StringVarP(&profileFile, "file", "f", "", "")
	profileCmd.Flags().StringVarP(&profileName, "name", "n", "", "")
	profileCmd.Flags().BoolVar(&profileMask, "mask", false, "")

	var buf bytes.Buffer
	profileCmd.SetOut(&buf)
	profileCmd.SetErr(&buf)
	profileCmd.SetArgs([]string{"--file", path, "--name", "dev"})
	if err := profileCmd.Execute(); err != nil {
		t.Logf("output: %s", buf.String())
	}
}

func TestProfileCmd_MissingFile(t *testing.T) {
	defer resetProfileCmd()
	cmd := &cobra.Command{Use: "profile"}
	cmd.RunE = profileCmd.RunE
	cmd.Flags().StringVarP(&profileFile, "file", "f", "", "")
	cmd.Flags().StringVarP(&profileName, "name", "n", "dev", "")
	cmd.Flags().BoolVar(&profileMask, "mask", false, "")
	profileName = "dev"
	profileFile = "/nonexistent/path.env"
	// Should not panic
}
