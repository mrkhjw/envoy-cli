package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/env"
)

var (
	sanitizeTrimKeys          bool
	sanitizeTrimValues        bool
	sanitizeNormalizeLineEnd  bool
	sanitizeRemoveNullBytes   bool
	sanitizeStripControl      bool
	sanitizeDryRun            bool
)

var sanitizeCmd = &cobra.Command{
	Use:   "sanitize [file]",
	Short: "Clean up .env file values by removing unwanted characters",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]

		entries, err := env.ParseFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to parse file: %w", err)
		}

		opts := env.SanitizeOptions{
			TrimKeys:             sanitizeTrimKeys,
			TrimValues:           sanitizeTrimValues,
			NormalizeLineEndings: sanitizeNormalizeLineEnd,
			RemoveNullBytes:      sanitizeRemoveNullBytes,
			StripControlChars:    sanitizeStripControl,
		}

		result := env.Sanitize(entries, opts)
		fmt.Println(result.Format())

		if sanitizeDryRun {
			fmt.Fprintln(os.Stderr, "[dry-run] no changes written")
			return nil
		}

		f, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
		defer f.Close()

		for _, e := range result.Entries {
			_, err := fmt.Fprintf(f, "%s=%s\n", e.Key, e.Value)
			if err != nil {
				return fmt.Errorf("failed to write entry: %w", err)
			}
		}
		return nil
	},
}

func init() {
	sanitizeCmd.Flags().BoolVar(&sanitizeTrimKeys, "trim-keys", false, "Trim whitespace from keys")
	sanitizeCmd.Flags().BoolVar(&sanitizeTrimValues, "trim-values", false, "Trim whitespace from values")
	sanitizeCmd.Flags().BoolVar(&sanitizeNormalizeLineEnd, "normalize-line-endings", false, "Normalize line endings to LF")
	sanitizeCmd.Flags().BoolVar(&sanitizeRemoveNullBytes, "remove-null-bytes", false, "Remove null bytes from keys and values")
	sanitizeCmd.Flags().BoolVar(&sanitizeStripControl, "strip-control-chars", false, "Strip non-printable control characters from values")
	sanitizeCmd.Flags().BoolVar(&sanitizeDryRun, "dry-run", false, "Preview changes without writing to file")
	rootCmd.AddCommand(sanitizeCmd)
}
