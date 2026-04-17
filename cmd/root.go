package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "envoy",
	Short: "envoy-cli — manage and sync .env files across environments",
	Long:  `A CLI tool for managing and syncing .env files with secret masking support.`,
}

// rootCmd returns the root cobra command
func rootCmd() *cobra.Command {
	return root
}

// Execute runs the root command
func Execute() {
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
