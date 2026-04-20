package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/env"
)

var (
	lookupFile          string
	lookupMaskSecrets   bool
	lookupCaseSensitive bool
)

var lookupCmd = &cobra.Command{
	Use:   "lookup <KEY>",
	Short: "Look up a single key in a .env file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]

		entries, err := env.ParseFile(lookupFile)
		if err != nil {
			return fmt.Errorf("failed to parse file: %w", err)
		}

		opts := env.LookupOptions{
			CaseSensitive: lookupCaseSensitive,
			MaskSecrets:   lookupMaskSecrets,
		}

		result := env.Lookup(entries, key, opts)
		fmt.Fprintln(os.Stdout, result.Format())

		if !result.Found {
			return fmt.Errorf("key not found: %s", key)
		}
		return nil
	},
}

func init() {
	lookupCmd.Flags().StringVarP(&lookupFile, "file", "f", ".env", "Path to the .env file")
	lookupCmd.Flags().BoolVar(&lookupMaskSecrets, "mask", false, "Mask secret values in output")
	lookupCmd.Flags().BoolVar(&lookupCaseSensitive, "case-sensitive", false, "Use case-sensitive key matching")
	rootCmd.AddCommand(lookupCmd)
}
