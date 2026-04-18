package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"envoy-cli/internal/env"
)

var redactCmd = &cobra.Command{
	Use:   "redact [file]",
	Short: "Print .env file with secret values redacted",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		placeholder, _ := cmd.Flags().GetString("placeholder")

		vars, err := env.ParseFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to parse file: %w", err)
		}

		result := env.Redact(vars, placeholder)

		sort.Strings(result.Lines)
		for _, line := range result.Lines {
			fmt.Fprintln(os.Stdout, line)
		}

		if result.Redacted > 0 {
			fmt.Fprintf(os.Stderr, "# %d secret(s) redacted\n", result.Redacted)
		}
		return nil
	},
}

func init() {
	redactCmd.Flags().String("placeholder", "***", "Placeholder text for redacted values")
	rootCmd.AddCommand(redactCmd)
}
