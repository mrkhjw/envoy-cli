package cmd

import (
	"fmt"
	"os"

	"github.com/yourusername/envoy-cli/internal/env"
	"github.com/spf13/cobra"
)

var (
	groupFile      string
	groupSep       string
	groupMask      bool
)

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Group env keys by prefix",
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := env.ParseFile(groupFile)
		if err != nil {
			return fmt.Errorf("failed to parse file: %w", err)
		}
		result := env.Group(entries, groupSep)
		fmt.Print(result.Format(entries, groupMask))
		return nil
	},
}

func init() {
	groupCmd.Flags().StringVarP(&groupFile, "file", "f", ".env", "Path to .env file")
	groupCmd.Flags().StringVar(&groupSep, "sep", "_", "Key separator for grouping")
	groupCmd.Flags().BoolVar(&groupMask, "mask", false, "Mask secret values")
	rootCmd.AddCommand(groupCmd)
	_ = os.Stderr
}
