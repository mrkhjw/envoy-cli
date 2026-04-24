package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/env"
)

var verifyCmd = &cobra.Command{
	Use:   "verify <file>",
	Short: "Verify that env file values match expected KEY=VALUE pairs",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		expectStrs, _ := cmd.Flags().GetStringSlice("expect")

		entries, err := env.ParseFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to parse file: %w", err)
		}

		expected := make(map[string]string, len(expectStrs))
		for _, pair := range expectStrs {
			parts := strings.SplitN(pair, "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid expect format %q, use KEY=VALUE", pair)
			}
			expected[parts[0]] = parts[1]
		}

		result := env.Verify(entries, env.VerifyOptions{Expected: expected})
		fmt.Println(result.Format())

		if !result.OK() {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	verifyCmd.Flags().StringSlice("expect", nil, "Expected KEY=VALUE pairs to verify (repeatable)")
	rootCmd.AddCommand(verifyCmd)
}
