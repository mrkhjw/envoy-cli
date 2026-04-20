package cmd

import (
	"fmt"

	"github.com/envoy-cli/envoy/internal/env"
	"github.com/spf13/cobra"
)

var (
	transformUpperKeys   bool
	transformLowerKeys   bool
	transformUpperValues bool
	transformLowerValues bool
	transformTrim        bool
	transformKeys        []string
	transformMask        bool
)

var transformCmd = &cobra.Command{
	Use:   "transform [file]",
	Short: "Apply key/value transformations to a .env file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := env.ParseFile(args[0])
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		if transformUpperKeys && transformLowerKeys {
			return fmt.Errorf("--upper-keys and --lower-keys are mutually exclusive")
		}
		if transformUpperValues && transformLowerValues {
			return fmt.Errorf("--upper-values and --lower-values are mutually exclusive")
		}

		opts := env.TransformOpts{
			UppercaseKeys:   transformUpperKeys,
			LowercaseKeys:   transformLowerKeys,
			UppercaseValues: transformUpperValues,
			LowercaseValues: transformLowerValues,
			TrimValues:      transformTrim,
			Keys:            transformKeys,
		}

		result := env.Transform(entries, opts)
		fmt.Print(result.Format(transformMask))
		return nil
	},
}

func init() {
	transformCmd.Flags().BoolVar(&transformUpperKeys, "upper-keys", false, "Uppercase all keys")
	transformCmd.Flags().BoolVar(&transformLowerKeys, "lower-keys", false, "Lowercase all keys")
	transformCmd.Flags().BoolVar(&transformUpperValues, "upper-values", false, "Uppercase all values")
	transformCmd.Flags().BoolVar(&transformLowerValues, "lower-values", false, "Lowercase all values")
	transformCmd.Flags().BoolVar(&transformTrim, "trim", false, "Trim whitespace from values")
	transformCmd.Flags().StringSliceVar(&transformKeys, "keys", nil, "Limit transformations to specific keys")
	transformCmd.Flags().BoolVar(&transformMask, "mask", false, "Mask secret values in output")

	rootCmd.AddCommand(transformCmd)
}
