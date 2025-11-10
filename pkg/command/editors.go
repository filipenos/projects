package command

import (
	"fmt"

	"github.com/filipenos/projects/pkg/config"
	"github.com/filipenos/projects/pkg/editor"
	"github.com/filipenos/projects/pkg/log"
	"github.com/spf13/cobra"
)

func init() {

	editorsCmd := &cobra.Command{
		Use:   "editors",
		Short: "Manage editors for your projects",
		Long:  `Manage and configure editors for your projects.`,
	}

	// Comando para inicializar arquivo de configuração de editores
	initEditorsCmd := &cobra.Command{
		Use:   "init",
		Short: "Create editors configuration file with examples",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := config.InitEditors(); err != nil {
				return err
			}
			log.Infof("Editors configuration file created at: %s", config.GetEditorsConfigPath())
			log.Infof("You can now add custom editors to this file.")
			return nil
		},
	}
	editorsCmd.AddCommand(initEditorsCmd)

	// Comando para listar editores disponíveis
	listEditorsCmd := &cobra.Command{
		Use:   "list",
		Short: "List all available editors",
		RunE: func(cmd *cobra.Command, args []string) error {
			if editorService == nil {
				return fmt.Errorf("editor service not initialized")
			}

			avaliables, notAvailable := editorService.GetEditors()
			log.Infof("Available editors:")
			for _, editor := range avaliables {
				log.Infof("  - %s", editor)
			}
			if len(notAvailable) > 0 {
				log.Infof("Not available editors: (not found in PATH or not executable): ")
				for _, editor := range notAvailable {
					log.Infof("  - %s", editor)
				}
			}
			return nil
		},
	}
	editorsCmd.AddCommand(listEditorsCmd)

	// Comando para recarregar configuração de editores
	reloadEditorsCmd := &cobra.Command{
		Use:   "reload",
		Short: "Reload editors configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.Load()

			var err error
			editorService, err = editor.NewService(cfg)
			if err != nil {
				return fmt.Errorf("failed to reload editor service: %w", err)
			}

			log.Infof("Editors configuration reloaded successfully!")

			// Atualiza aliases do comando code
			for _, c := range rootCmd.Commands() {
				if c.Use == "code" {
					c.Aliases = editorService.Aliases()
					break
				}
			}

			return nil
		},
	}
	editorsCmd.AddCommand(reloadEditorsCmd)

	// Adiciona o comando de editores ao comando raiz
	rootCmd.AddCommand(editorsCmd)
}
