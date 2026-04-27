package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/env"
)

var envdiffCmd = &cobra.Command{
	Use:   "envdiff <v1-label> <file1> <v2-label> <file2>",
	Short: "Compare two versioned .env files and show a labeled diff",
	Args:  cobra.ExactArgs(4),
	RunE: func(cmd *cobra.Command, args []string) error {
		v1Label := args[0]
		file1 := args[1]
		v2Label := args[2]
		file2 := args[3]

		maskSecrets, _ := cmd.Flags().GetBool("mask-secrets")

		result, err := env.VersionDiffFile(v1Label, file1, v2Label, file2)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			return err
		}

		total := len(result.Added) + len(result.Removed) + len(result.Changed)
		if total == 0 {
			fmt.Printf("No differences between %s and %s\n", v1Label, v2Label)
			return nil
		}

		fmt.Print(result.Format(maskSecrets))
		fmt.Printf("\nSummary: +%d added, -%d removed, ~%d changed\n",
			len(result.Added), len(result.Removed), len(result.Changed))
		return nil
	},
}

func init() {
	envdiffCmd.Flags().Bool("mask-secrets", true, "Mask secret values in output")
	rootCmd.AddCommand(envdiffCmd)
}
