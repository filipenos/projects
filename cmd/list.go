package cmd

import (
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "l"},
	Short:   "List projects",
	RunE:    list,
}

func init() {
	rootCmd.AddCommand(listCmd)
	codeCmd.Flags().BoolP("path", "p", false, "Reuse same window")
	codeCmd.Flags().BoolP("simple", "s", false, "Reuse same window")
}
