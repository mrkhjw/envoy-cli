package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/env"
)

var truncateCmd = &cobra.Command{
	Use:   "truncate [file]",
	Short: "Truncate values that exceed a maximum length",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		maxLen, _ := cmd.Flags().GetInt("max-len")
		suffix, _ := cmd.Flags().GetString("suffix")
		keys, _ := cmd.Flags().GetStringSlice("keys")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		res, err := env.TruncateFile(path, env.TruncateOptions{
			MaxLen: maxLen,
			Suffix: suffix,
			Keys:   keys,
			DryRun: dryRun,
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			return err
		}

		fmt.Println(res.Format())
		return nil
	},
}

func init() {
	truncateCmd.Flags().Int("max-len", 64, "Maximum allowed value length")
	truncateCmd.Flags().String("suffix", "...", "Suffix to append to truncated values")
	truncateCmd.Flags().StringSlice("keys", nil, "Specific keys to truncate (default: all)")
	truncateCmd.Flags().Bool("dry-run", false, "Preview without writing changes")
	rootCmd.AddCommand(truncateCmd)
}
