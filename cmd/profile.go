package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/env"
)

var (
	profileFile    string
	profileName    string
	profileMask    bool
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Load and display a named profile from a multi-profile env file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if profileFile == "" {
			return fmt.Errorf("--file is required")
		}
		if profileName == "" {
			return fmt.Errorf("--name is required")
		}
		profiles, err := env.LoadProfiles(profileFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		res, err := env.Profile(profiles, profileName)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		fmt.Println(res.Format(profileMask))
		return nil
	},
}

func init() {
	profileCmd.Flags().StringVarP(&profileFile, "file", "f", "", "Path to multi-profile env file")
	profileCmd.Flags().StringVarP(&profileName, "name", "n", "", "Profile name to load")
	profileCmd.Flags().BoolVar(&profileMask, "mask", false, "Mask secret values in output")
	rootCmd.AddCommand(profileCmd)
}
