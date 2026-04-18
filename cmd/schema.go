package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/env"
)

var schemaFile string

func init() {
	schemaCmd := &cobra.Command{
		Use:   "schema",
		Short: "Validate a .env file against a schema definition",
		RunE: func(cmd *cobra.Command, args []string) error {
			envPath, _ := cmd.Flags().GetString("file")
			if envPath == "" {
				envPath = ".env"
			}

			entries, err := env.ParseFile(envPath)
			if err != nil {
				return fmt.Errorf("could not read env file: %w", err)
			}

			schema, err := env.LoadSchema(schemaFile)
			if err != nil {
				return fmt.Errorf("could not load schema: %w", err)
			}

			result := env.ValidateSchema(entries, schema)
			fmt.Print(result.Format())

			if len(result.Missing) > 0 {
				os.Exit(1)
			}
			return nil
		},
	}

	schemaCmd.Flags().StringP("file", "f", ".env", "Path to the .env file")
	schemaCmd.Flags().StringVarP(&schemaFile, "schema", "s", ".env.schema", "Path to schema definition file")
	rootCmd.AddCommand(schemaCmd)
}
