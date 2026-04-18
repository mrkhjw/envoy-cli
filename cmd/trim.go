package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/your-org/envoy-cli/internal/env"
	"github.com/spf13/cobra"
)

func init() {
	var dryRun bool
	var keys string
	var trimLeft bool
	var trimRight bool

	cmd := &cobra.Command{
		Use:   "trim [file]",
		Short: "Trim leading/trailing whitespace from env values",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			file := args[0]
			entries, err := env.ParseFile(file)
			if err != nil {
				return fmt.Errorf("failed to parse file: %w", err)
			}

			var keyList []string
			if keys != "" {
				for _, k := range strings.Split(keys, ",") {
					if k = strings.TrimSpace(k); k != "" {
						keyList = append(keyList, k)
					}
				}
			}

			opts := env.TrimOptions{
				Keys:      keyList,
				TrimLeft:  trimLeft,
				TrimRight: trimRight,
				DryRun:    dryRun,
			}

			out, result := env.Trim(entries, opts)
			fmt.Print(result.Format())

			if dryRun {
				return nil
			}

			f, err := os.Create(file)
			if err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}
			defer f.Close()
			for _, e := range out {
				fmt.Fprintf(f, "%s=%s\n", e.Key, e.Value)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview changes without writing")
	cmd.Flags().StringVar(&keys, "keys", "", "Comma-separated keys to trim (default: all)")
	cmd.Flags().BoolVar(&trimLeft, "left", false, "Trim left whitespace only")
	cmd.Flags().BoolVar(&trimRight, "right", false, "Trim right whitespace only")

	rootCmd.AddCommand(cmd)
}
