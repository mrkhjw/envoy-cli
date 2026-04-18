package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/env"
)

var generateCmd = &cobra.Command{
	Use:   "generate [file]",
	Short: "Generate random values for secret keys in a .env file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file := args[0]
		length, _ := cmd.Flags().GetInt("length")
		format, _ := cmd.Flags().GetString("format")
		keys, _ := cmd.Flags().GetStringSlice("keys")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		entries, err := env.ParseFile(file)
		if err != nil {
			return fmt.Errorf("parse file: %w", err)
		}

		upper := make([]string, len(keys))
		for i, k := range keys {
			upper[i] = strings.ToUpper(k)
		}

		updated, result, err := env.Generate(entries, env.GenerateOptions{
			Length: length,
			Format: format,
			Keys:   upper,
			DryRun: dryRun,
		})
		if err != nil {
			return err
		}

		if dryRun {
			fmt.Println("[dry-run] would generate:")
			fmt.Print(result.Format())
			return nil
		}

		if err := env.ExportFile(updated, file, env.ExportOptions{}); err != nil {
			return fmt.Errorf("write file: %w", err)
		}

		fmt.Printf("Updated %s\n", file)
		fmt.Print(result.Format())
		return nil
	},
}

func init() {
	generateCmd.Flags().Int("length", 32, "Length of generated value")
	generateCmd.Flags().String("format", "hex", "Format: hex, alphanumeric, full")
	generateCmd.Flags().StringSlice("keys", nil, "Specific keys to regenerate (default: all secrets)")
	generateCmd.Flags().Bool("dry-run", false, "Preview without writing")
	rootCmd.AddCommand(generateCmd)
}
