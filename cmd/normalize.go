package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/env"
)

var normalizeCmd = &cobra.Command{
	Use:   "normalize [file]",
	Short: "Normalize key casing, value quoting, and whitespace in a .env file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file := args[0]
		uppercase, _ := cmd.Flags().GetBool("uppercase")
		trim, _ := cmd.Flags().GetBool("trim")
		quote, _ := cmd.Flags().GetBool("quote")
		strip, _ := cmd.Flags().GetBool("strip-export")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		entries, err := env.ParseFile(file)
		if err != nil {
			return fmt.Errorf("failed to parse file: %w", err)
		}

		opts := env.NormalizeOptions{
			UppercaseKeys: uppercase,
			TrimValues:    trim,
			QuoteValues:   quote,
			StripExported: strip,
		}

		result := env.Normalize(entries, opts)
		fmt.Println(result.Format())

		if dryRun {
			return nil
		}

		f, err := os.Create(file)
		if err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
		defer f.Close()
		for _, e := range result.Entries {
			fmt.Fprintf(f, "%s=%s\n", e.Key, e.Value)
		}
		return nil
	},
}

func init() {
	normalizeCmd.Flags().Bool("uppercase", false, "Convert all keys to uppercase")
	normalizeCmd.Flags().Bool("trim", false, "Trim whitespace from values")
	normalizeCmd.Flags().Bool("quote", false, "Wrap unquoted values in double quotes")
	normalizeCmd.Flags().Bool("strip-export", false, "Remove 'export ' prefix from keys")
	normalizeCmd.Flags().Bool("dry-run", false, "Preview changes without writing")
	rootCmd.AddCommand(normalizeCmd)
}
