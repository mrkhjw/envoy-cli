package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/envoy-cli/envoy/internal/env"
	"github.com/spf13/cobra"
)

var protectCmd = &cobra.Command{
	Use:   "protect [file]",
	Short: "Mark keys as protected to prevent accidental overwrite",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		keysFlag, _ := cmd.Flags().GetString("keys")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		entries, err := env.ParseFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to parse file: %w", err)
		}

		var keys []string
		if keysFlag != "" {
			for _, k := range strings.Split(keysFlag, ",") {
				if k = strings.TrimSpace(k); k != "" {
					keys = append(keys, k)
				}
			}
		}

		opts := env.ProtectOptions{
			Keys:   keys,
			DryRun: dryRun,
		}

		_, result := env.Protect(entries, opts)
		fmt.Fprintln(os.Stdout, result.Format())
		return nil
	},
}

func init() {
	protectCmd.Flags().String("keys", "", "Comma-separated list of keys to protect (default: all secrets)")
	protectCmd.Flags().Bool("dry-run", false, "Preview which keys would be protected without modifying")
	rootCmd.AddCommand(protectCmd)
}
