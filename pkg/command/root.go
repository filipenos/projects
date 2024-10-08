package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

//https://github.com/liamg/sunder

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:  "projects",
	Long: `Have all your work projects in one place. Open, edit in a much simpler way.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
