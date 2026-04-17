package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/env"
)

func init() {
	var format string
	var maskSecrets bool
	var output string

	exportCmd := &cobra.Command{
		Use:   "export [file]",
		Short: "Export .env file in various formats",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]

			content, err := env.ExportFile(filePath, env.ExportFormat(format), maskSecrets)
			if err != nil {
				return fmt.Errorf("export failed: %w", err)
			}

			if output != "" {
				if err := env.WriteExport(content, output); err != nil {
					return fmt.Errorf("failed to write output: %w", err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Exported to %s\n", output)
			} else {
				fmt.Fprint(cmd.OutOrStdout(), content)
			}
			return nil
		},
	}

	exportCmd.Flags().StringVarP(&format, "format", "f", "shell", "Output format: shell, docker, json")
	exportCmd.Flags().BoolVarP(&maskSecrets, "mask", "m", false, "Mask secret values in output")
	exportCmd.Flags().StringVarP(&output, "output", "o", "", "Write output to file instead of stdout")

	if root := rootCmd(); root != nil {
		root.AddCommand(exportCmd)
	} else {
		os.Exit(1)
	}
}
