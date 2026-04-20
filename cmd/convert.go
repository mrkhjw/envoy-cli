package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/env"
)

var (
	convertFormat string
	convertOutput string
)

var convertCmd = &cobra.Command{
	Use:   "convert [file]",
	Short: "Convert a .env file to another format",
	Long: `Convert a .env file to another format such as shell export, YAML, TOML, or JSON.

Supported formats: env, export, yaml, toml, json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]

		entries, err := env.ParseFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to parse file: %w", err)
		}

		format := strings.ToLower(convertFormat)
		result, err := env.Convert(entries, format)
		if err != nil {
			return fmt.Errorf("conversion failed: %w", err)
		}

		if convertOutput != "" {
			if err := os.WriteFile(convertOutput, []byte(result.Output), 0644); err != nil {
				return fmt.Errorf("failed to write output file: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Converted %d entries to %s format → %s\n",
				result.Count, format, convertOutput)
		} else {
			fmt.Fprint(cmd.OutOrStdout(), result.Output)
		}

		return nil
	},
}

func init() {
	convertCmd.Flags().StringVarP(&convertFormat, "format", "f", "json",
		"Output format: env, export, yaml, toml, json")
	convertCmd.Flags().StringVarP(&convertOutput, "output", "o", "",
		"Write output to file instead of stdout")
	rootCmd.AddCommand(convertCmd)
}
