package cmd

import (
	"github.com/spf13/cobra"
)

// openCmd represents the open command
var openCmd = &cobra.Command{
	Use:     "open project_name",
	Aliases: []string{"op", "o"},
	Short:   "Open your project",
	RunE:    open,
}

func init() {
	rootCmd.AddCommand(openCmd)
}
