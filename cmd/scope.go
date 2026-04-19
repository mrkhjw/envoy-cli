package cmd

import (
	"fmt"
	"os"

	"github.com/eze-kiel/envoy-cli/internal/env"
	"github.com/spf13/cobra"
)

var (
	scopePrefix      string
	scopeStripPrefix bool
	scopeMask        bool
)

var scopeCmd = &cobra.Command{
	Use:   "scope [file]",
	Short: "Filter entries by environment prefix (e.g. PROD_, DEV_)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := env.ParseFile(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return err
		}

		result := env.Scope(entries, env.ScopeOptions{
			Prefix:      scopePrefix,
			StripPrefix: scopeStripPrefix,
		})

		fmt.Print(result.Format(scopeMask))
		return nil
	},
}

func init() {
	scopeCmd.Flags().StringVarP(&scopePrefix, "prefix", "p", "", "Environment prefix to filter by (e.g. PROD)")
	scopeCmd.Flags().BoolVar(&scopeStripPrefix, "strip", false, "Strip prefix from key names in output")
	scopeCmd.Flags().BoolVar(&scopeMask, "mask", true, "Mask secret values in output")
	rootCmd.AddCommand(scopeCmd)
}
