package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/envoy-cli/envoy/internal/env"
	"github.com/spf13/cobra"
)

var placeholderCmd = &cobra.Command{
	Use:   "placeholder [file]",
	Short: "Fill placeholder values in a .env file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file := args[0]
		token, _ := cmd.Flags().GetString("token")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		setFlags, _ := cmd.Flags().GetStringArray("set")

		entries, err := env.ParseFile(file)
		if err != nil {
			return fmt.Errorf("failed to parse file: %w", err)
		}

		replacements := map[string]string{}
		for _, s := range setFlags {
			parts := strings.SplitN(s, "=", 2)
			if len(parts) == 2 {
				replacements[parts[0]] = parts[1]
			}
		}

		result := env.FillPlaceholders(entries, replacements, env.PlaceholderOptions{
			Token:  token,
			DryRun: dryRun,
		})

		fmt.Println(result.Format())

		if !dryRun && result.Filled > 0 {
			f, err := os.Create(file)
			if err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}
			defer f.Close()
			for _, e := range result.Entries {
				fmt.Fprintf(f, "%s=%s\n", e.Key, e.Value)
			}
		}
		return nil
	},
}

func init() {
	placeholderCmd.Flags().String("token", "CHANGEME", "Placeholder token to match")
	placeholderCmd.Flags().Bool("dry-run", false, "Preview changes without writing")
	placeholderCmd.Flags().StringArray("set", []string{}, "KEY=VALUE replacements")
	rootCmd.AddCommand(placeholderCmd)
}
