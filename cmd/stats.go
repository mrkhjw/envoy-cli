package cmd

import (
	"fmt"
	"os"

	"github.com/envoy-cli/envoy/internal/env"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats [file]",
	Short: "Show statistics about a .env file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := env.ParseFile(args[0])
		if err != nil {
			return fmt.Errorf("failed to parse file: %w", err)
		}
		result := env.Stats(entries)
		fmt.Fprintln(os.Stdout, result.Format())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
