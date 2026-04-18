package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"envoy-cli/internal/env"
)

func writeTempRotateEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "rotate-*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func runRotateCmd(args ...string) (string, error) {
	root := &cobra.Command{Use: "envoy"}
	var dryRun bool
	var keys string
	rotateCmd := &cobra.Command{
		Use:  "rotate <file>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var targetKeys []string
			if keys != "" {
				for _, k := range strings.Split(keys, ",") {
					if t := strings.TrimSpace(k); t != "" {
						targetKeys = append(targetKeys, t)
					}
				}
			}
			_, err := env.RotateFile(args[0], env.RotateOptions{Keys: targetKeys, DryRun: dryRun})
			return err
		},
	}
	rotateCmd.Flags().BoolVar(&dryRun, "dry-run", false, "")
	rotateCmd.Flags().StringVar(&keys, "keys", "", "")
	root.AddCommand(rotateCmd)
	root.SetArgs(args)
	buf := new(strings.Builder)
	root.SetOut(buf)
	err := root.Execute()
	return buf.String(), err
}

func TestRotateCmd_DryRun(t *testing.T) {
	path := writeTempRotateEnv(t, "API_KEY=secret\nAPP=myapp\n")
	_, err := runRotateCmd("rotate", "--dry-run", path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// File should be unchanged
	envMap, _ := env.ParseFile(path)
	if envMap["API_KEY"] != "secret" {
		t.Error("dry-run should not modify file")
	}
}

func TestRotateCmd_MissingFile(t *testing.T) {
	_, err := runRotateCmd("rotate", "/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
