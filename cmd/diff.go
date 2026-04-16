package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/env"
)

var (
	diffMaskSecrets bool
)

var diffCmd = &cobra.Command{
	Use:   "diff [base-file] [target-file]",
	Short: "Show differences between two .env files",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		baseFile := args[0]
		targetFile := args[1]

		base, err := env.ParseFile(baseFile)
		if err != nil {
			return fmt.Errorf("failed to parse base file %q: %w", baseFile, err)
		}

		target, err := env.ParseFile(targetFile)
		if err != nil {
			return fmt.Errorf("failed to parse target file %q: %w", targetFile, err)
		}

		result := env.Diff(base, target)

		if len(result.Added)+len(result.Removed)+len(result.Changed) == 0 {
			fmt.Println("No differences found.")
			return nil
		}

		output := result.Format(diffMaskSecrets)
		fmt.Fprint(os.Stdout, output)
		return nil
	},
}

func init() {
	diffCmd.Flags().BoolVarP(&diffMaskSecrets, "mask", "m", true, "Mask secret values in output")
	rootCmd.AddCommand(diffCmd)
}
