package cmd

import (
	"fmt"
	"os"

	"github.com/envoy-cli/envoy/internal/env"
	"github.com/spf13/cobra"
)

var (
	archiveLabel  string
	archiveDest   string
	archiveSource string
)

var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Archive a .env file to a timestamped JSON snapshot",
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := env.ParseFile(archiveSource)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return err
		}

		res, err := env.Archive(entries, archiveDest, archiveLabel)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return err
		}

		fmt.Println(res.Format())
		return nil
	},
}

func init() {
	archiveCmd.Flags().StringVarP(&archiveSource, "file", "f", ".env", "source .env file")
	archiveCmd.Flags().StringVarP(&archiveDest, "out", "o", "env.archive.json", "destination archive file")
	archiveCmd.Flags().StringVarP(&archiveLabel, "label", "l", "default", "label for this archive")
	rootCmd.AddCommand(archiveCmd)
}
