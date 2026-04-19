package cmd

import (
	"fmt"
	"os"

	"github.com/seanomeara96/envoy-cli/internal/env"
	"github.com/spf13/cobra"
)

var pinCmd = &cobra.Command{
	Use:   "pin [file]",
	Short: "Pin specific keys to their current values",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		keys, _ := cmd.Flags().GetStringSlice("keys")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		entries, err := env.ParseFile(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
			return err
		}

		result := env.Pin(entries, keys, dryRun)
		fmt.Print(result.Format())
		return nil
	},
}

func init() {
	pinCmd.Flags().StringSlice("keys", nil, "Keys to pin (default: all)")
	pinCmd.Flags().Bool("dry-run", false, "Preview without writing")
	rootCmd.AddCommand(pinCmd)
}
