// command/code.go
package command

import (
	"fmt"

	"github.com/filipenos/projects/pkg/config"
	"github.com/filipenos/projects/pkg/editor"
	"github.com/filipenos/projects/pkg/path"
	"github.com/filipenos/projects/pkg/project"
	"github.com/spf13/cobra"
)

var editorService *editor.Service

func init() {
	// Carrega a configuração
	cfg := config.Load()

	// Inicializa o serviço de editores
	var err error
	editorService, err = editor.NewService(cfg)
	if err != nil {
		fmt.Printf("Warning: failed to initialize editor service: %v\n", err)
	}

	codeCmd := &cobra.Command{
		Use:     "code",
		Short:   fmt.Sprintf("Edit your project using the editor (%s as default)", cfg.Editor),
		Aliases: editorService.GetAvailableEditors(),
		RunE:    code,
	}
	codeCmd.Flags().StringP("window", "w", "new", "Window type (new|reuse|add)")
	rootCmd.AddCommand(codeCmd)
}

func code(cmdParam *cobra.Command, params []string) error {
	projects, err := project.Load(cfg)
	if err != nil {
		return err
	}

	p, _ := projects.Find(path.SafeName(params...))
	if p == nil {
		return fmt.Errorf("project not found")
	}

	if err := p.Validate(); err != nil {
		return err
	}

	// Converte a flag window para o tipo apropriado
	windowFlag := SafeStringFlag(cmdParam, "window")
	var window editor.WindowType

	switch windowFlag {
	case "reuse":
		window = editor.WindowTypeReuse
	case "add":
		window = editor.WindowTypeAdd
	default:
		window = editor.WindowTypeNew
	}

	// Usa o serviço de editores para abrir o projeto
	return editorService.OpenProject(cmdParam.CalledAs(), p, window)
}

// Comandos auxiliares para gerenciar editores
