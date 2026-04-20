package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"envoy-cli/internal/env"
)

var reorderCmd = &cobra.Command{
	Use:   "reorder [file]",
	Short: "Reorder keys in a .env file, pinning specified keys to the top",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		src := args[0]
		keys, _ := cmd.Flags().GetString("keys")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		output, _ := cmd.Flags().GetString("output")

		var keyList []string
		if keys != "" {
			for _, k := range strings.Split(keys, ",") {
				k = strings.TrimSpace(k)
				if k != "" {
					keyList = append(keyList, k)
				}
			}
		}

		dst := src
		if output != "" {
			dst = output
		}

		res, err := env.ReorderFile(src, dst, env.ReorderOptions{
			Keys:   keyList,
			DryRun: dryRun,
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			return err
		}

		if dryRun {
			for _, e := range res.Entries {
				fmt.Printf("%s=%s\n", e.Key, e.Value)
			}
			return nil
		}
		fmt.Println(res.Format())
		return nil
	},
}

func init() {
	reorderCmd.Flags().String("keys", "", "comma-separated list of keys to pin to the top")
	reorderCmd.Flags().Bool("dry-run", false, "preview reordered output without writing")
	reorderCmd.Flags().String("output", "", "output file (defaults to input file)")
	rootCmd.AddCommand(reorderCmd)
}
