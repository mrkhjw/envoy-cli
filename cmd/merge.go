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

	mergeCmd := &cobra.Command{
		Use:   "merge <base> <source>",
		Short: "Merge source .env into base .env",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			baseFile := args[0]
			srcFile := args[1]

			base, err := env.ParseFile(baseFile)
			if err != nil {
				return fmt.Errorf("reading base: %w", err)
			}
			src, err := env.ParseFile(srcFile)
			if err != nil {
				return fmt.Errorf("reading source: %w", err)
			}

			merged, result := env.Merge(base, src, overwrite)
			fmt.Print(result.Format())

			if dryRun {
				fmt.Fprintln(os.Stderr, "dry-run: no changes written")
				return nil
			}

			f, err := os.Create(baseFile)
			if err != nil {
				return fmt.Errorf("writing base: %w", err)
			}
			defer f.Close()
			for k, v := range merged {
				fmt.Fprintf(f, "%s=%s\n", k, v)
			}
			return nil
		},
	}

	mergeCmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite existing keys")
	mergeCmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview changes without writing")
	rootCmd.AddCommand(mergeCmd)
}
