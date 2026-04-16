package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/envoy-cli/envoy/internal/env"
	"github.com/spf13/cobra"
)

var (
	flagReveal  bool
	flagEnvFile string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all variables in a .env file",
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := env.ParseFile(flagEnvFile)
		if err != nil {
			return fmt.Errorf("parse env file: %w", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "KEY\tVALUE\tSECRET")
		fmt.Fprintln(w, "---\t-----\t------")
		for _, e := range entries {
			value := e.Value
			if e.Secret && !flagReveal {
				value = "********"
			}
			secretMark := ""
			if e.Secret {
				secretMark = "yes"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\n", e.Key, value, secretMark)
		}
		return w.Flush()
	},
}

func init() {
	listCmd.Flags().BoolVar(&flagReveal, "reveal", false, "Show secret values in plain text")
	listCmd.Flags().StringVarP(&flagEnvFile, "file", "f", ".env", "Path to the .env file")
	RootCmd.AddCommand(listCmd)
}
