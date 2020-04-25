package cmd

import (
	"github.com/spf13/cobra"
)

// codeCmd represents the code command
var codeCmd = &cobra.Command{
	Use:   "code",
	Short: "Edit your project using the editor (code as default)",
	RunE:  code,
}

func init() {
	rootCmd.AddCommand(codeCmd)
	codeCmd.Flags().StringP("editor", "e", "code", "Code using vscode editor")
	codeCmd.Flags().BoolP("reuse", "r", false, "Reuse same window")
}
