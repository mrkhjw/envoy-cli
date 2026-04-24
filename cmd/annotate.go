package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/env"
)

var annotateCmd = &cobra.Command{
	Use:   "annotate [file]",
	Short: "Add inline comments to keys in a .env file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		note, _ := cmd.Flags().GetString("note")
		keysStr, _ := cmd.Flags().GetString("keys")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		if note == "" {
			return fmt.Errorf("--note is required")
		}

		var keys []string
		if keysStr != "" {
			for _, k := range strings.Split(keysStr, ",") {
				k = strings.TrimSpace(k)
				if k != "" {
					keys = append(keys, k)
				}
			}
		}

		result, err := env.AnnotateFile(path, note, keys, dryRun)
		if err != nil {
			return err
		}

		fmt.Println(result.Format())
		if dryRun {
			for _, ae := range result.Entries {
				if ae.Key == "" {
					continue
				}
				val := ae.Value
				if ae.IsSecret {
					val = "***"
				}
				line := fmt.Sprintf("%s=%s", ae.Key, val)
				if ae.Annotation != "" {
					line += fmt.Sprintf(" # %s", ae.Annotation)
				}
				fmt.Println(line)
			}
		}
		return nil
	},
}

func init() {
	annotateCmd.Flags().String("note", "", "Annotation text to add as inline comment (required)")
	annotateCmd.Flags().String("keys", "", "Comma-separated list of keys to annotate (default: all)")
	annotateCmd.Flags().Bool("dry-run", false, "Preview changes without writing to file")
	rootCmd.AddCommand(annotateCmd)
}
