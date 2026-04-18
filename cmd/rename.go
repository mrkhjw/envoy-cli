package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/env"
)

func init() {
	var overwrite bool

	renameCmd := &cobra.Command{
		Use:   "rename <file> <old-key> <new-key>",
		Short: "Rename a key in a .env file",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			oldKey := args[1]
			newKey := args[2]

			result, err := env.RenameFile(path, oldKey, newKey, overwrite)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
				return err
			}

			fmt.Println(result.Format())
			return nil
		},
	}

	renameCmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite new key if it already exists")
	rootCmd.AddCommand(renameCmd)
}
