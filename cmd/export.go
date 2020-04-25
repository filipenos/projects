package cmd

import (
	"github.com/spf13/cobra"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export the list of projects to avaliables formats",
	RunE:  export,
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
