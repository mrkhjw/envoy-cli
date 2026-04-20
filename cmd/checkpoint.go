package cmd

import (
	"fmt"
	"os"

	"github.com/saurabh/envoy-cli/internal/env"
	"github.com/spf13/cobra"
)

var (
	checkpointLabel  string
	checkpointOutput string
	checkpointDryRun bool
)

var checkpointCmd = &cobra.Command{
	Use:   "checkpoint [file]",
	Short: "Save a labeled checkpoint of an env file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := env.ParseFile(args[0])
		if err != nil {
			return fmt.Errorf("checkpoint: %w", err)
		}

		if checkpointOutput == "" {
			checkpointOutput = checkpointLabel + ".checkpoint.json"
		}

		result, err := env.Checkpoint(entries, checkpointLabel, checkpointOutput, checkpointDryRun)
		if err != nil {
			return err
		}

		fmt.Fprintln(os.Stdout, result.Format())
		return nil
	},
}

var checkpointLoadCmd = &cobra.Command{
	Use:   "load [checkpoint-file]",
	Short: "Load and display a saved checkpoint",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cp, err := env.LoadCheckpoint(args[0])
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "checkpoint: %q  saved: %s  entries: %d\n",
			cp.Label, cp.Timestamp.Format("2006-01-02T15:04:05Z"), len(cp.Entries))
		for _, e := range cp.Entries {
			if env.IsSecret(e.Key) {
				fmt.Fprintf(os.Stdout, "  %s=***\n", e.Key)
			} else {
				fmt.Fprintf(os.Stdout, "  %s=%s\n", e.Key, e.Value)
			}
		}
		return nil
	},
}

func init() {
	checkpointCmd.Flags().StringVarP(&checkpointLabel, "label", "l", "checkpoint", "Label for the checkpoint")
	checkpointCmd.Flags().StringVarP(&checkpointOutput, "output", "o", "", "Output file path (default: <label>.checkpoint.json)")
	checkpointCmd.Flags().BoolVar(&checkpointDryRun, "dry-run", false, "Preview without writing")
	checkpointCmd.AddCommand(checkpointLoadCmd)
	rootCmd.AddCommand(checkpointCmd)
}
