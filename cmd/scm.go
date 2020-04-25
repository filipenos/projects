package cmd

import (
	"github.com/spf13/cobra"
)

// scmCmd represents the scm command
var scmCmd = &cobra.Command{
	Use:   "scm",
	Short: "Handle scm project",
	RunE:  scm,
}

func init() {
	rootCmd.AddCommand(scmCmd)
}
