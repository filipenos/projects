package editor

import "github.com/filipenos/projects/pkg/project"

// ConfigurableEditor representa um editor configurável via JSON
type ConfigurableEditor struct {
	name           string
	aliases        []string
	executable     string
	supportedTypes []project.ProjectType
	windowArgs     map[WindowType][]string
	workspaceArgs  []string
	folderArgs     []string
	pathPosition   int
}

func (e *ConfigurableEditor) Name() string {
	return e.name
}

func (e *ConfigurableEditor) Aliases() []string {
	return e.aliases
}

func (e *ConfigurableEditor) SupportsProjectType(projectType project.ProjectType) bool {
	if len(e.supportedTypes) == 0 {
		return true
	}

	for _, supported := range e.supportedTypes {
		if supported == projectType {
			return true
		}
	}
	return false
}

func (e *ConfigurableEditor) GetExecutable() string {
	return e.executable
}

func (e *ConfigurableEditor) BuildArgs(p *project.Project, window WindowType) ([]string, error) {
	args := make([]string, 0)

	// Adiciona argumentos da janela
	if windowArgs, exists := e.windowArgs[window]; exists {
		args = append(args, windowArgs...)
	}

	// Determina se é workspace ou pasta
	var projectArgs []string
	if p.IsWorkspace && len(e.workspaceArgs) > 0 {
		projectArgs = e.workspaceArgs
	} else if len(e.folderArgs) > 0 {
		projectArgs = e.folderArgs
	}

	// Adiciona argumentos do projeto
	args = append(args, projectArgs...)

	// Adiciona o path na posição correta
	if e.pathPosition >= 0 && e.pathPosition < len(args) {
		args = append(args[:e.pathPosition], append([]string{p.RootPath}, args[e.pathPosition:]...)...)
	} else {
		args = append(args, p.RootPath)
	}

	return args, nil
}
