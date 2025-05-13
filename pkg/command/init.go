package command

import (
	"github.com/filipenos/projects/pkg/config"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{"i"},
	Short:   "Initialize new config file",
	Example: `projects init`,
	Args:    cobra.NoArgs,
	RunE:    initConfig,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func initConfig(cmdParam *cobra.Command, params []string) error {
	return config.Init()
}
