package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/env"
)

func init() {
	var dryRun bool
	var keys string
	var timestamp bool

	rotateCmd := &cobra.Command{
		Use:   "rotate <file>",
		Short: "Rotate secret values in an env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]

			var targetKeys []string
			if keys != "" {
				for _, k := range strings.Split(keys, ",") {
					trimmed := strings.TrimSpace(k)
					if trimmed != "" {
						targetKeys = append(targetKeys, trimmed)
					}
				}
			}

			opts := env.RotateOptions{
				Keys:      targetKeys,
				DryRun:    dryRun,
				Timestamp: timestamp,
			}

			result, err := env.RotateFile(path, opts)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				return err
			}
			fmt.Print(result.Format())
			return nil
		},
	}

	rotateCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview rotation without writing changes")
	rotateCmd.Flags().StringVar(&keys, "keys", "", "Comma-separated list of keys to rotate (default: all secrets)")
	rotateCmd.Flags().BoolVar(&timestamp, "timestamp", false, "Append timestamp to rotated placeholder value")

	rootCmd.AddCommand(rotateCmd)
}
