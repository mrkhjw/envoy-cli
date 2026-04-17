package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"envoy-cli/internal/env"
)

func runExportCmd(t *testing.T, args []string) (string, error) {
	t.Helper()
	var buf bytes.Buffer
	cmd := &cobra.Command{
		Use:  "export [file]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			format, _ := cmd.Flags().GetString("format")
			mask, _ := cmd.Flags().GetBool("mask")
			content, err := env.ExportFile(args[0], env.ExportFormat(format), mask)
			if err != nil {
				return err
			}
			cmd.Print(content)
			return nil
		},
	}
	cmd.Flags().String("format", "shell", "")
	cmd.Flags().Bool("mask", false, "")
	cmd.SetOut(&buf)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return buf.String(), err
}

func writeTempExportEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestExportCmd_ShellOutput(t *testing.T) {
	f := writeTempExportEnv(t, "APP=myapp\n")
	out, err := runExportCmd(t, []string{"--format", "shell", f})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "export APP=") {
		t.Errorf("expected shell export in output, got: %s", out)
	}
}

func TestExportCmd_JSONOutput(t *testing.T) {
	f := writeTempExportEnv(t, "KEY=value\n")
	out, err := runExportCmd(t, []string{"--format", "json", f})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "KEY") {
		t.Errorf("expected KEY in JSON output, got: %s", out)
	}
}

func TestExportCmd_MissingFile(t *testing.T) {
	_, err := runExportCmd(t, []string{"--format", "shell", "/no/such/file.env"})
	if err == nil {
		t.Error("expected error for missing file")
	}
}
