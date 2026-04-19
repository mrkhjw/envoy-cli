package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/smatthewenglish/envoy-cli/internal/env"
	"github.com/spf13/cobra"
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt [file]",
	Short: "Encrypt secret values in a .env file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		passKey, _ := cmd.Flags().GetString("key")
		keys, _ := cmd.Flags().GetStringSlice("keys")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		if passKey == "" {
			passKey = os.Getenv("ENVOY_ENCRYPT_KEY")
		}
		if passKey == "" {
			return fmt.Errorf("encryption key required: use --key or ENVOY_ENCRYPT_KEY")
		}

		entries, err := env.ParseFile(filePath)
		if err != nil {
			return fmt.Errorf("parse error: %w", err)
		}

		result, err := env.Encrypt(entries, passKey, keys)
		if err != nil {
			return fmt.Errorf("encrypt error: %w", err)
		}

		if dryRun {
			fmt.Println("[dry-run] encrypted keys:")
			for k, v := range result.Encrypted {
				fmt.Printf("  %s=%s\n", k, v[:min(16, len(v))]+"...")
			}
			return nil
		}

		var sb strings.Builder
		for _, e := range entries {
			if enc, ok := result.Encrypted[e.Key]; ok {
				sb.WriteString(fmt.Sprintf("%s=%s\n", e.Key, enc))
			} else {
				sb.WriteString(fmt.Sprintf("%s=%s\n", e.Key, e.Value))
			}
		}
		if err := os.WriteFile(filePath, []byte(sb.String()), 0644); err != nil {
			return fmt.Errorf("write error: %w", err)
		}
		fmt.Print(result.Format())
		return nil
	},
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func init() {
	encryptCmd.Flags().String("key", "", "Encryption passphrase (or set ENVOY_ENCRYPT_KEY)")
	encryptCmd.Flags().StringSlice("keys", nil, "Specific keys to encrypt (default: all secrets)")
	encryptCmd.Flags().Bool("dry-run", false, "Preview encrypted output without writing")
	rootCmd.AddCommand(encryptCmd)
}
