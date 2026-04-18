package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/env"
)

func init() {
	var overwrite bool

	cloneCmd := &cobra.Command{
		Use:   "clone <source> <destination>",
		Short: "Clone an .env file to a new location",
		Long:  "Copy key-value pairs from a source .env file to a destination file, with optional overwrite support.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			srcPath := args[0]
			dstPath := args[1]

			result, err := env.Clone(srcPath, dstPath, overwrite)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				return err
			}

			fmt.Println(result.Summary())
			return nil
		},
	}

	cloneCmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "Overwrite existing keys in destination")
	rootCmd.AddCommand(cloneCmd)
}
