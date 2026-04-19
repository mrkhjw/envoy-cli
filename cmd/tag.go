package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/env"
)

var (
	tagFile        string
	tagTags        []string
	tagKeys        []string
	tagMaskSecrets bool
)

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Annotate env entries with tags as inline comments",
	RunE: func(cmd *cobra.Command, args []string) error {
		if tagFile == "" {
			return fmt.Errorf("--file is required")
		}
		if len(tagTags) == 0 {
			return fmt.Errorf("at least one --tag is required")
		}

		entries, err := env.ParseFile(tagFile)
		if err != nil {
			return fmt.Errorf("failed to parse file: %w", err)
		}

		// Flatten comma-separated tags
		var flatTags []string
		for _, t := range tagTags {
			for _, part := range strings.Split(t, ",") {
				if s := strings.TrimSpace(part); s != "" {
					flatTags = append(flatTags, s)
				}
			}
		}

		result := env.Tag(entries, env.TagOptions{
			Tags:        flatTags,
			Keys:        tagKeys,
			MaskSecrets: tagMaskSecrets,
		})

		fmt.Print(result.Format())
		return nil
	},
}

func init() {
	tagCmd.Flags().StringVarP(&tagFile, "file", "f", "", "Path to .env file")
	tagCmd.Flags().StringArrayVarP(&tagTags, "tag", "t", []string{}, "Tag(s) to apply")
	tagCmd.Flags().StringArrayVarP(&tagKeys, "key", "k", []string{}, "Keys to tag (default: all)")
	tagCmd.Flags().BoolVar(&tagMaskSecrets, "mask-secrets", false, "Mask secret values in output")
	rootCmd.AddCommand(tagCmd)
}
