package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

//Version, Commit This variables is filled on build time
var (
	Version string
	Commit  string
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version of projects",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version %s, Build %s\n", Version, Commit)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
