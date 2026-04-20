package cmd

import (
	"fmt"
	"os"

	"github.com/envoy-cli/envoy/internal/env"
	"github.com/spf13/cobra"
)

var healthCmd = &cobra.Command{
	Use:   "health [file]",
	Short: "Run a health check on a .env file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		entries, err := env.ParseFile(path)
		if err != nil {
			return fmt.Errorf("failed to parse file: %w", err)
		}

		result := env.Health(entries)
		fmt.Print(result.Format())

		if result.Errors > 0 {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(healthCmd)
}
