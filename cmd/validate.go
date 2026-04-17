package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/env"
)

var validateCmd = &cobra.Command{
	Use:   "validate [file]",
	Short: "Validate a .env file for formatting and duplicate keys",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		result, err := env.ValidateFile(path)
		if err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}
		fmt.Println(result.Summary())
		if !result.Valid() {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
