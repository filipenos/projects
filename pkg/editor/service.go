// pkg/editor/service.go
package editor

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/filipenos/projects/pkg/config"
	"github.com/filipenos/projects/pkg/log"
	"github.com/filipenos/projects/pkg/project"
)

// WindowType representa os tipos de janela suportados
type WindowType string

const (
	WindowTypeNew   WindowType = "new"
	WindowTypeReuse WindowType = "reuse"
	WindowTypeAdd   WindowType = "add"
)

// Editor define a interface que todos os editores devem implementar
type Editor interface {
	Name() string
	Aliases() []string
	SupportsProjectType(projectType project.ProjectType) bool
	BuildArgs(p *project.Project, window WindowType) ([]string, error)
	GetExecutable() string
}

// Registry mantém o registro de editores disponíveis
type Registry struct {
	editors map[string]Editor
}

// NewRegistry cria uma nova instância do registro
func NewRegistry() *Registry {
	return &Registry{
		editors: make(map[string]Editor),
	}
}

// Register registra um editor no sistema
func (r *Registry) Register(editor Editor) {
	r.editors[editor.Name()] = editor
	for _, alias := range editor.Aliases() {
		r.editors[alias] = editor
	}
}

// Get retorna um editor pelo nome ou alias
func (r *Registry) Get(name string) (Editor, bool) {
	editor, exists := r.editors[name]
	return editor, exists
}

// GetAllNames retorna todos os nomes e aliases disponíveis
func (r *Registry) GetAllNames() []string {
	var names []string
	seen := make(map[string]bool)

	for name := range r.editors {
		if !seen[name] {
			names = append(names, name)
			seen[name] = true
		}
	}
	return names
}

// Service gerencia as operações com editores
type Service struct {
	registry *Registry
	config   config.Config
}

// NewService cria uma nova instância do serviço
func NewService(cfg config.Config) (*Service, error) {
	registry := NewRegistry()

	// Registra editores padrão
	registry.Register(&VSCodeEditor{})
	registry.Register(&SublimeEditor{})
	registry.Register(&ZedEditor{})
	registry.Register(&NvimEditor{})

	service := &Service{
		registry: registry,
		config:   cfg,
	}

	// Carrega editores customizados
	if err := service.loadCustomEditors(); err != nil {
		// Log do erro mas não falha na inicialização
		log.Warnf("Failed to load custom editors: %v", err)
	}

	return service, nil
}

// loadCustomEditors carrega editores do arquivo de configuração
func (s *Service) loadCustomEditors() error {
	editorsConfig, err := config.LoadEditors(s.config)
	if err != nil {
		return err
	}

	if editorsConfig == nil {
		return nil // Arquivo não existe
	}

	for _, editorConfig := range editorsConfig.Editors {
		editor, err := s.createConfigurableEditor(editorConfig)
		if err != nil {
			return fmt.Errorf("failed to create editor '%s': %w", editorConfig.Name, err)
		}
		s.registry.Register(editor)
	}

	return nil
}

func (s *Service) createConfigurableEditor(cfg config.EditorConfig) (*ConfigurableEditor, error) {
	if cfg.Name == "" {
		return nil, fmt.Errorf("editor name is required")
	}

	if cfg.Executable == "" {
		return nil, fmt.Errorf("editor executable is required")
	}

	// Converte tipos de projeto
	supportedTypes := make([]project.ProjectType, 0, len(cfg.SupportedTypes))
	for _, typeStr := range cfg.SupportedTypes {
		projectType := project.ProjectType(typeStr)
		supportedTypes = append(supportedTypes, projectType)
	}

	// Converte argumentos de janela
	windowArgs := make(map[WindowType][]string)
	for windowStr, args := range cfg.WindowArgs {
		windowType := WindowType(windowStr)
		windowArgs[windowType] = args
	}

	return &ConfigurableEditor{
		name:           cfg.Name,
		aliases:        cfg.Aliases,
		executable:     cfg.Executable,
		supportedTypes: supportedTypes,
		windowArgs:     windowArgs,
		workspaceArgs:  cfg.WorkspaceArgs,
		folderArgs:     cfg.FolderArgs,
		pathPosition:   cfg.PathPosition,
	}, nil
}

// RegisterEditor registra um novo editor
func (s *Service) RegisterEditor(editor Editor) {
	s.registry.Register(editor)
}

// GetAvailableEditors retorna todos os editores disponíveis
func (s *Service) GetAvailableEditors() []string {
	return s.registry.GetAllNames()
}

// OpenProject abre um projeto no editor especificado
func (s *Service) OpenProject(editorName string, p *project.Project, window WindowType) error {
	editor, exists := s.registry.Get(editorName)
	if !exists {
		return fmt.Errorf("editor '%s' not found", editorName)
	}

	if !editor.SupportsProjectType(p.ProjectType) {
		return fmt.Errorf("editor '%s' does not support project type '%s'",
			editor.Name(), p.ProjectType)
	}

	args, err := editor.BuildArgs(p, window)
	if err != nil {
		return fmt.Errorf("failed to build args for editor '%s': %w", editor.Name(), err)
	}

	executable := editor.GetExecutable()
	openType := "folder"
	if p.IsWorkspace {
		openType = "file"
	}

	log.Infof("open %s '%s' on '%s'", openType, p.RootPath, executable)

	cmd := exec.Command(executable, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
