package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/env"
)

var revertCmd = &cobra.Command{
	Use:   "revert <current-file> <baseline-file>",
	Short: "Revert env entries to their baseline values",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		currentFile := args[0]
		baselineFile := args[1]

		current, err := env.ParseFile(currentFile)
		if err != nil {
			return fmt.Errorf("reading current file: %w", err)
		}

		baseline, err := env.ParseFile(baselineFile)
		if err != nil {
			return fmt.Errorf("reading baseline file: %w", err)
		}

		keys, _ := cmd.Flags().GetStringSlice("keys")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		overwrite, _ := cmd.Flags().GetBool("overwrite")

		opts := env.RevertOptions{
			Keys:      keys,
			DryRun:    dryRun,
			Overwrite: overwrite,
		}

		updated, result := env.Revert(current, baseline, opts)
		fmt.Print(result.Format())

		if dryRun {
			return nil
		}

		if len(result.Reverted) == 0 {
			return nil
		}

		f, err := os.Create(currentFile)
		if err != nil {
			return fmt.Errorf("writing file: %w", err)
		}
		defer f.Close()

		for _, e := range updated {
			if e.IsComment {
				fmt.Fprintln(f, e.Raw)
			} else if e.Key != "" {
				fmt.Fprintf(f, "%s=%s\n", e.Key, e.Value)
			}
		}
		return nil
	},
}

func init() {
	revertCmd.Flags().StringSlice("keys", nil, "Specific keys to revert (default: all)")
	revertCmd.Flags().Bool("dry-run", false, "Preview changes without writing")
	revertCmd.Flags().Bool("overwrite", false, "Revert even if value is already identical")
	rootCmd.AddCommand(revertCmd)
}
