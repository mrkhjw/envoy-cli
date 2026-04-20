package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/env"
)

var wrapCmd = &cobra.Command{
	Use:   "wrap [file]",
	Short: "Wrap or truncate long .env values",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file := args[0]
		maxLen, _ := cmd.Flags().GetInt("max-length")
		quote, _ := cmd.Flags().GetBool("quote")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		keysFlag, _ := cmd.Flags().GetString("keys")

		var keys []string
		if keysFlag != "" {
			for _, k := range strings.Split(keysFlag, ",") {
				k = strings.TrimSpace(k)
				if k != "" {
					keys = append(keys, k)
				}
			}
		}

		opts := env.WrapOptions{
			MaxLength: maxLen,
			Quote:     quote,
			DryRun:    dryRun,
			Keys:      keys,
		}

		result, err := env.WrapFile(file, opts)
		if err != nil {
			return err
		}

		fmt.Println(result.Format())
		if dryRun {
			for _, e := range result.Wrapped {
				if !e.Comment && e.Key != "" {
					fmt.Printf("  %s=%s\n", e.Key, e.Value)
				}
			}
		}
		return nil
	},
}

func init() {
	wrapCmd.Flags().Int("max-length", 80, "Maximum value length before truncation")
	wrapCmd.Flags().Bool("quote", false, "Quote values after wrapping")
	wrapCmd.Flags().Bool("dry-run", false, "Preview changes without writing")
	wrapCmd.Flags().String("keys", "", "Comma-separated list of keys to wrap (default: all)")
	rootCmd.AddCommand(wrapCmd)
}
