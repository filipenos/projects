package cmd

import (
	"github.com/spf13/cobra"
)

// pathCmd represents the path command
var pathCmd = &cobra.Command{
	Use:     "path",
	Aliases: []string{"pwd"},
	Short:   "Show path of project",
	RunE:    path,
}

func init() {
	rootCmd.AddCommand(pathCmd)
}
