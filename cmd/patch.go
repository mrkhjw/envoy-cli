package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/sierrasoftworks/envoy-cli/internal/env"
	"github.com/spf13/cobra"
)

var patchDryRun bool

var patchCmd = &cobra.Command{
	Use:   "patch <file> <op:key[=value|->newkey]...",
	Short: "Apply patch operations (set, delete, rename) to a .env file",
	Example: `  envoy patch .env set:PORT=9090 delete:OLD_KEY rename:DB_URL->DATABASE_URL`,
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		entries, err := env.ParseFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to parse file: %w", err)
		}

		var ops []env.PatchOp
		for _, raw := range args[1:] {
			op, err := parsePatchOp(raw)
			if err != nil {
				return err
			}
			ops = append(ops, op)
		}

		out, result := env.Patch(entries, ops, patchDryRun)
		fmt.Print(result.Format())

		if !patchDryRun {
			f, err := os.Create(filePath)
			if err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}
			defer f.Close()
			for _, e := range out {
				fmt.Fprintf(f, "%s=%s\n", e.Key, e.Value)
			}
		}
		return nil
	},
}

func parsePatchOp(raw string) (env.PatchOp, error) {
	if strings.HasPrefix(raw, "set:") {
		parts := strings.SplitN(strings.TrimPrefix(raw, "set:"), "=", 2)
		if len(parts) != 2 {
			return env.PatchOp{}, fmt.Errorf("invalid set op: %s", raw)
		}
		return env.PatchOp{Op: "set", Key: parts[0], Value: parts[1]}, nil
	}
	if strings.HasPrefix(raw, "delete:") {
		return env.PatchOp{Op: "delete", Key: strings.TrimPrefix(raw, "delete:")}, nil
	}
	if strings.HasPrefix(raw, "rename:") {
		parts := strings.SplitN(strings.TrimPrefix(raw, "rename:"), "->", 2)
		if len(parts) != 2 {
			return env.PatchOp{}, fmt.Errorf("invalid rename op: %s", raw)
		}
		return env.PatchOp{Op: "rename", Key: parts[0], NewKey: parts[1]}, nil
	}
	return env.PatchOp{}, fmt.Errorf("unknown op: %s", raw)
}

func init() {
	patchCmd.Flags().BoolVar(&patchDryRun, "dry-run", false, "Preview changes without writing")
	rootCmd.AddCommand(patchCmd)
}
