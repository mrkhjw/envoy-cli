package cmd

import (
	"fmt"
	"os"

	"github.com/envoy-cli/envoy/internal/env"
	"github.com/spf13/cobra"
)

var (
	unsetFile   string
	unsetDryRun bool
	unsetWrite  bool
)

var unsetCmd = &cobra.Command{
	Use:   "unset [keys...]",
	Short: "Remove one or more keys from a .env file",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := env.ParseFile(unsetFile)
		if err != nil {
			return fmt.Errorf("failed to parse file: %w", err)
		}

		result := env.Unset(entries, args, unsetDryRun)
		fmt.Println(result.Format())

		if !unsetDryRun && unsetWrite {
			f, err := os.Create(unsetFile)
			if err != nil {
				return fmt.Errorf("failed to open file for writing: %w", err)
			}
			defer f.Close()
			for _, e := range result.Entries {
				if e.Key == "" {
					fmt.Fprintln(f, e.Raw)
				} else {
					fmt.Fprintf(f, "%s=%s\n", e.Key, e.Value)
				}
			}
		}

		return nil
	},
}

func init() {
	unsetCmd.Flags().StringVarP(&unsetFile, "file", "f", ".env", "Path to the .env file")
	unsetCmd.Flags().BoolVar(&unsetDryRun, "dry-run", false, "Preview changes without applying them")
	unsetCmd.Flags().BoolVar(&unsetWrite, "write", false, "Write changes back to the file")
	rootCmd.AddCommand(unsetCmd)
}
