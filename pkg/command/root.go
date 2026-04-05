package command

import (
	"fmt"
	"os"

	"github.com/filipenos/projects/pkg/config"
	"github.com/filipenos/projects/pkg/log"
	"github.com/spf13/cobra"
)

var cfg, cfgErr = config.Load()

//https://github.com/liamg/sunder

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:  "projects",
	Long: `Have all your work projects in one place. Open, edit in a much simpler way.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cfgErr != nil {
			return fmt.Errorf("failed to load configuration: %w", cfgErr)
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Warnf("command failed: %v", err)
		os.Exit(1)
	}
}
