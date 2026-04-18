package cmd

import (
	"fmt"
	"os"

	"github.com/your-org/envoy-cli/internal/env"
	"github.com/spf13/cobra"
)

var (
	sortReverse     bool
	sortByValue     bool
	sortSecretsLast bool
	sortOutput      string
	sortMask        bool
)

var sortCmd = &cobra.Command{
	Use:   "sort [file]",
	Short: "Sort .env file entries alphabetically",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		src := args[0]
		opts := env.SortOptions{
			Reverse:     sortReverse,
			ByValue:     sortByValue,
			SecretsLast: sortSecretsLast,
		}

		if sortOutput != "" {
			res, err := env.SortFile(src, sortOutput, opts)
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "Sorted %d entries → %s\n", res.Total, sortOutput)
			return nil
		}

		entries, err := env.ParseFile(src)
		if err != nil {
			return err
		}
		res := env.Sort(entries, opts)
		fmt.Println(res.Format(sortMask))
		return nil
	},
}

func init() {
	sortCmd.Flags().BoolVarP(&sortReverse, "reverse", "r", false, "Sort in descending order")
	sortCmd.Flags().BoolVar(&sortByValue, "by-value", false, "Sort by value instead of key")
	sortCmd.Flags().BoolVar(&sortSecretsLast, "secrets-last", false, "Place secret keys at the end")
	sortCmd.Flags().StringVarP(&sortOutput, "output", "o", "", "Write sorted output to file")
	sortCmd.Flags().BoolVar(&sortMask, "mask", false, "Mask secret values in output")
	rootCmd.AddCommand(sortCmd)
}
