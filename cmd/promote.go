package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/env"
)

var (
	promoteOverwrite bool
	promoteKeys      []string
)

var promoteCmd = &cobra.Command{
	Use:   "promote <src> <dst>",
	Short: "Promote env vars from one file to another",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		srcPath := args[0]
		dstPath := args[1]

		result, err := env.PromoteFile(srcPath, dstPath, promoteKeys, promoteOverwrite)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return err
		}

		fmt.Printf("Promoted %s → %s\n", srcPath, dstPath)
		fmt.Print(result.Format())

		if len(result.Added)+len(result.Updated) == 0 {
			fmt.Println("Nothing to promote.")
		}
		return nil
	},
}

func init() {
	promoteCmd.Flags().BoolVarP(&promoteOverwrite, "overwrite", "o", false, "Overwrite existing keys in destination")
	promoteCmd.Flags().StringSliceVarP(&promoteKeys, "keys", "k", nil, "Specific keys to promote (default: all)")
	rootCmd.AddCommand(promoteCmd)
}
