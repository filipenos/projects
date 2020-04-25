package cmd

import (
	"github.com/spf13/cobra"
)

// closeCmd represents the close command
var closeCmd = &cobra.Command{
	Use:   "close",
	Short: "close your project",
	RunE:  close,
}

func init() {
	rootCmd.AddCommand(closeCmd)
}
