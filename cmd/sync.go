package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/env"
)

func init() {
	var overwrite bool
	var dryRun bool

	syncCmd := &cobra.Command{
		Use:   "sync <source> <destination>",
		Short: "Sync variables from source .env into destination .env",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			srcPath, dstPath := args[0], args[1]

			src, err := env.ParseFile(srcPath)
			if err != nil {
				return fmt.Errorf("reading source: %w", err)
			}

			dst, err := env.ParseFile(dstPath)
			if err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("reading destination: %w", err)
			}
			if dst == nil {
				dst = map[string]string{}
			}

			opts := env.SyncOptions{Overwrite: overwrite, DryRun: dryRun}
			result, err := env.Sync(dst, src, dstPath, opts)
			if err != nil {
				return err
			}

			if dryRun {
				fmt.Println("[dry-run] No changes written.")
			}
			fmt.Printf("Applied: %d  Skipped: %d\n", len(result.Applied), len(result.Skipped))
			for _, k := range result.Applied {
				fmt.Printf("  + %s\n", k)
			}
			for _, k := range result.Skipped {
				fmt.Printf("  ~ %s (skipped)\n", k)
			}
			return nil
		},
	}

	syncCmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "Overwrite existing keys in destination")
	syncCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "Preview changes without writing")

	rootCmd.AddCommand(syncCmd)
}
