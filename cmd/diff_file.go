package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/env"
)

var diffFileCmd = &cobra.Command{
	Use:   "diff-file <fileA> <fileB>",
	Short: "Show differences between two .env files",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := env.DiffFile(args[0], args[1])
		if err != nil {
			return fmt.Errorf("diff-file: %w", err)
		}

		maskSecrets, _ := cmd.Flags().GetBool("mask-secrets")
		output := result.Format(maskSecrets)

		if output == "" {
			fmt.Fprintln(os.Stdout, "No differences found.")
			return nil
		}

		fmt.Fprint(os.Stdout, output)
		return nil
	},
}

func init() {
	diffFileCmd.Flags().Bool("mask-secrets", true, "Mask secret values in output")
	rootCmd.AddCommand(diffFileCmd)
}
