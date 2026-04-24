package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/envoy-cli/envoy/internal/env"
	"github.com/spf13/cobra"
)

var tokenizeCmd = &cobra.Command{
	Use:   "tokenize [file]",
	Short: "Split env values into tokens by a delimiter",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		delimiter, _ := cmd.Flags().GetString("delimiter")
		keys, _ := cmd.Flags().GetStringSlice("keys")
		mask, _ := cmd.Flags().GetBool("mask")

		entries, err := env.ParseFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return err
		}

		opts := env.TokenizeOptions{
			Delimiter:   delimiter,
			Keys:        keys,
			MaskSecrets: mask,
		}

		result := env.Tokenize(entries, opts)

		for key, tokens := range result.Tokens {
			if mask && isSecretKey(key) {
				fmt.Printf("%s = [REDACTED]\n", key)
				continue
			}
			fmt.Printf("%s = [%s]\n", key, strings.Join(tokens, ", "))
		}

		fmt.Printf("\ntotal=%d skipped=%d\n", result.Total, result.Skipped)
		return nil
	},
}

func isSecretKey(key string) bool {
	upper := strings.ToUpper(key)
	for _, kw := range []string{"SECRET", "PASSWORD", "PASS", "TOKEN", "KEY", "AUTH", "PRIVATE"} {
		if strings.Contains(upper, kw) {
			return true
		}
	}
	return false
}

func init() {
	tokenizeCmd.Flags().StringP("delimiter", "d", ",", "delimiter to split values on")
	tokenizeCmd.Flags().StringSliceP("keys", "k", []string{}, "specific keys to tokenize")
	tokenizeCmd.Flags().Bool("mask", false, "mask secret values in output")
	rootCmd.AddCommand(tokenizeCmd)
}
