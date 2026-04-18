package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/env"
)

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Save a snapshot of an .env file to a JSON file",
	RunE: func(cmd *cobra.Command, args []string) error {
		source, _ := cmd.Flags().GetString("file")
		dest, _ := cmd.Flags().GetString("out")

		if source == "" {
			return fmt.Errorf("--file is required")
		}
		if dest == "" {
			return fmt.Errorf("--out is required")
		}

		result, err := env.TakeSnapshot(source, dest)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return err
		}

		fmt.Println(result.Format())
		return nil
	},
}

var loadSnapshotCmd = &cobra.Command{
	Use:   "snapshot-load",
	Short: "Load and display a previously saved .env snapshot",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := cmd.Flags().GetString("file")
		if path == "" {
			return fmt.Errorf("--file is required")
		}

		snap, err := env.LoadSnapshot(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return err
		}

		fmt.Printf("Snapshot from %s (taken %s):\n", snap.Source, snap.Timestamp.Format("2006-01-02 15:04:05"))
		for k, v := range snap.Entries {
			if env.MaskLine(k+"="+v) != k+"="+v {
				v = "***"
			}
			fmt.Printf("  %s=%s\n", k, v)
		}
		return nil
	},
}

func init() {
	snapshotCmd.Flags().String("file", "", "Source .env file")
	snapshotCmd.Flags().String("out", "", "Destination JSON snapshot file")
	rootCmd.AddCommand(snapshotCmd)

	loadSnapshotCmd.Flags().String("file", "", "Snapshot JSON file to load")
	rootCmd.AddCommand(loadSnapshotCmd)
}
