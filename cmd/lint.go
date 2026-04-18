package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/env"
)

func init() {
	var filePath string

	lintCmd := &cobra.Command{
		Use:   "lint",
		Short: "Lint a .env file for common issues",
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := os.Open(filePath)
			if err != nil {
				return fmt.Errorf("could not open file: %w", err)
			}
			defer f.Close()

			var lines []string
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("error reading file: %w", err)
			}

			result := env.Lint(lines)
			fmt.Println(result.Format())

			if result.HasErrors() {
				os.Exit(1)
			}
			return nil
		},
	}

	lintCmd.Flags().StringVarP(&filePath, "file", "f", ".env", "Path to the .env file")
	rootCmd.AddCommand(lintCmd)
}
