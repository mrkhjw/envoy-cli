package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/env"
)

var (
	searchKey   string
	searchValue string
	searchCase  bool
	searchMask  bool
)

var searchCmd = &cobra.Command{
	Use:   "search [file]",
	Short: "Search for keys or values in an env file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file := args[0]

		if searchKey == "" && searchValue == "" {
			return fmt.Errorf("at least one of --key or --value must be specified")
		}

		opts := env.SearchOptions{
			Key:           searchKey,
			Value:         searchValue,
			CaseSensitive: searchCase,
		}

		res, err := env.SearchFile(file, opts, searchMask)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			return err
		}

		fmt.Println(res.Format(searchMask))
		return nil
	},
}

func init() {
	searchCmd.Flags().StringVar(&searchKey, "key", "", "Substring to match against keys")
	searchCmd.Flags().StringVar(&searchValue, "value", "", "Substring to match against values")
	searchCmd.Flags().BoolVar(&searchCase, "case-sensitive", false, "Enable case-sensitive matching")
	searchCmd.Flags().BoolVar(&searchMask, "mask", true, "Mask secret values in output")
	rootCmd.AddCommand(searchCmd)
}
