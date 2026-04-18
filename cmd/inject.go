package cmd

import (
	"fmt"
	"os"

	"github.com/envoy-cli/envoy/internal/env"
	"github.com/spf13/cobra"
)

var injectCmd = &cobra.Command{
	Use:   "inject [file]",
	Short: "Inject variables from a .env file into the current process environment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		overwrite, _ := cmd.Flags().GetBool("overwrite")

		result, err := env.InjectFile(path, overwrite)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return err
		}

		fmt.Print(result.Format())
		return nil
	},
}

func init() {
	injectCmd.Flags().Bool("overwrite", false, "Overwrite existing environment variables")
	rootCmd.AddCommand(injectCmd)
}
