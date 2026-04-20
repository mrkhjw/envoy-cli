package cmd

import (
	"fmt"
	"os"

	"github.com/envoy-cli/envoy/internal/env"
	"github.com/spf13/cobra"
)

var compareMaskSecrets bool

var compareCmd = &cobra.Command{
	Use:   "compare <file1> <file2>",
	Short: "Compare two .env files and report differences",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := env.CompareFile(args[0], args[1])
		if err != nil {
			return fmt.Errorf("compare failed: %w", err)
		}

		output := result.Summary(compareMaskSecrets)
		fmt.Fprintln(os.Stdout, output)
		return nil
	},
}

func init() {
	compareCmd.Flags().BoolVar(&compareMaskSecrets, "mask-secrets", true, "Mask secret values in output")
	rootCmd.AddCommand(compareCmd)
}
