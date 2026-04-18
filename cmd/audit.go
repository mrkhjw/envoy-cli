package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/env"
)

var auditCmd = &cobra.Command{
	Use:   "audit [file]",
	Short: "Audit keys in an .env file and display a redacted access log",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]

		vars, err := env.ParseFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to parse file: %w", err)
		}

		action, _ := cmd.Flags().GetString("action")
		if action == "" {
			action = "read"
		}

		log := env.AuditMap(vars, action, filePath)
		fmt.Fprintln(os.Stdout, log.Format())
		return nil
	},
}

func init() {
	auditCmd.Flags().String("action", "read", "Action label to record in the audit log (e.g. read, export, sync)")
	rootCmd.AddCommand(auditCmd)
}
