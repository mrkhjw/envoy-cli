package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/envoy-cli/envoy/internal/env"
	"github.com/spf13/cobra"
)

var freezeCmd = &cobra.Command{
	Use:   "freeze [file]",
	Short: "Mark env entries as frozen to prevent accidental changes",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]

		entries, err := env.ParseFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to parse file: %w", err)
		}

		keysRaw, _ := cmd.Flags().GetString("keys")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		maskSecrets, _ := cmd.Flags().GetBool("mask")

		var keys []string
		if keysRaw != "" {
			for _, k := range strings.Split(keysRaw, ",") {
				k = strings.TrimSpace(k)
				if k != "" {
					keys = append(keys, k)
				}
			}
		}

		result := env.Freeze(entries, env.FreezeOption{
			Keys:   keys,
			DryRun: dryRun,
		})

		if !dryRun {
			lines := make([]string, 0, len(result.Entries))
			for _, e := range result.Entries {
				lines = append(lines, e.RawLine)
			}
			content := strings.Join(lines, "\n") + "\n"
			if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}
		}

		fmt.Print(result.Format(maskSecrets))
		return nil
	},
}

func init() {
	freezeCmd.Flags().String("keys", "", "Comma-separated list of keys to freeze (default: all)")
	freezeCmd.Flags().Bool("dry-run", false, "Preview freeze without writing changes")
	freezeCmd.Flags().Bool("mask", false, "Mask secret values in output")
	rootCmd.AddCommand(freezeCmd)
}
