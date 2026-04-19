package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/env"
)

var (
	maskPlaceholder  string
	maskRevealPrefix int
)

var maskCmd = &cobra.Command{
	Use:   "mask [file]",
	Short: "Mask secret values in a .env file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := env.ParseFile(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return err
		}

		opts := env.MaskOptions{
			Placeholder:  maskPlaceholder,
			RevealPrefix: maskRevealPrefix,
		}

		result := env.Mask(entries, opts)
		fmt.Println(result.Format())
		return nil
	},
}

func init() {
	maskCmd.Flags().StringVar(&maskPlaceholder, "placeholder", "***", "Placeholder string for masked values")
	maskCmd.Flags().IntVar(&maskRevealPrefix, "reveal-prefix", 0, "Number of leading characters to reveal")
	rootCmd.AddCommand(maskCmd)
}
