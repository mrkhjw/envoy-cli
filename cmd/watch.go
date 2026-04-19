package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sQVe/envoy-cli/internal/env"
	"github.com/spf13/cobra"
)

var watchInterval int

var watchCmd = &cobra.Command{
	Use:   "watch [file]",
	Short: "Watch a .env file for changes and report diffs",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		if _, err := os.Stat(path); err != nil {
			return fmt.Errorf("file not found: %s", path)
		}

		duration := time.Duration(watchInterval) * time.Millisecond
		fmt.Fprintf(cmd.OutOrStdout(), "Watching %s every %dms... (Ctrl+C to stop)\n", path, watchInterval)

		done := make(chan struct{})
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-sigs
			close(done)
		}()

		return env.Watch(path, duration, done, func(r env.WatchResult) {
			fmt.Fprintln(cmd.OutOrStdout(), r.Format())
		})
	},
}

func init() {
	watchCmd.Flags().IntVarP(&watchInterval, "interval", "i", 500, "Poll interval in milliseconds")
	rootCmd.AddCommand(watchCmd)
}
