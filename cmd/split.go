package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/envoy-cli/envoy/internal/env"
	"github.com/spf13/cobra"
)

var splitCmd = &cobra.Command{
	Use:   "split [file]",
	Short: "Split .env entries into matched and remainder groups by key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file := args[0]
		keys, _ := cmd.Flags().GetString("keys")
		invert, _ := cmd.Flags().GetBool("invert")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		entries, err := env.ParseFile(file)
		if err != nil {
			return fmt.Errorf("failed to parse file: %w", err)
		}

		var keyList []string
		if keys != "" {
			for _, k := range strings.Split(keys, ",") {
				keyList = append(keyList, strings.TrimSpace(k))
			}
		}

		result := env.Split(entries, env.SplitOptions{
			Keys:   keyList,
			Invert: invert,
			DryRun: dryRun,
		})

		fmt.Fprintln(os.Stdout, "=== Matched ===")
		for _, e := range result.Matched {
			fmt.Fprintf(os.Stdout, "%s=%s\n", e.Key, e.Value)
		}
		fmt.Fprintln(os.Stdout, "=== Remainder ===")
		for _, e := range result.Remainder {
			fmt.Fprintf(os.Stdout, "%s=%s\n", e.Key, e.Value)
		}
		fmt.Fprintln(os.Stdout, result.Format())
		return nil
	},
}

func init() {
	splitCmd.Flags().String("keys", "", "Comma-separated list of keys to match")
	splitCmd.Flags().Bool("invert", false, "Invert selection (match keys NOT in list)")
	splitCmd.Flags().Bool("dry-run", false, "Preview split without writing")
	rootCmd.AddCommand(splitCmd)
}
