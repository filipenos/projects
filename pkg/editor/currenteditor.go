package editor

import (
	"github.com/filipenos/projects/pkg/project"
	"github.com/filipenos/projects/pkg/workspace"
)

// VSCodeEditor implementa suporte ao VS Code
type VSCodeEditor struct{}

func (e *VSCodeEditor) Name() string {
	return "vscode"
}

func (e *VSCodeEditor) Aliases() []string {
	return []string{"code"}
}

func (e *VSCodeEditor) SupportsProjectType(projectType project.ProjectType) bool {
	return true
}

func (e *VSCodeEditor) GetExecutable() string {
	return "code"
}

func (e *VSCodeEditor) BuildArgs(p *project.Project, window WindowType) ([]string, error) {
	args := make([]string, 0)

	switch window {
	case WindowTypeReuse:
		args = append(args, "--reuse-window")
	case WindowTypeAdd:
		args = append(args, "--add")
	default:
		args = append(args, "--new-window")
	}

	if p.IsWorkspace {
		args = append(args, "--file-uri")
	} else {
		args = append(args, "--folder-uri")
	}
	args = append(args, p.RootPath)

	return args, nil
}

// ZedEditor implementa suporte ao Zed
type ZedEditor struct{}

func (e *ZedEditor) Name() string {
	return "zed"
}

func (e *ZedEditor) Aliases() []string {
	return []string{}
}

func (e *ZedEditor) SupportsProjectType(projectType project.ProjectType) bool {
	return projectType == project.ProjectTypeLocal
}

func (e *ZedEditor) GetExecutable() string {
	return "zed"
}

func (e *ZedEditor) BuildArgs(p *project.Project, window WindowType) ([]string, error) {
	args := make([]string, 0)

	switch window {
	case WindowTypeAdd:
		args = append(args, "--add")
	}

	if p.IsWorkspace {
		w, err := workspace.Load(p.RootPath)
		if err != nil {
			return nil, err
		}
		args = append(args, w.FoldersPath()...)
	} else {
		args = append(args, p.RootPath)
	}

	return args, nil
}

// SublimeEditor implementa suporte ao Sublime Text
type SublimeEditor struct{}

func (e *SublimeEditor) Name() string {
	return "sublime"
}

func (e *SublimeEditor) Aliases() []string {
	return []string{"subl"}
}

func (e *SublimeEditor) SupportsProjectType(projectType project.ProjectType) bool {
	return projectType == project.ProjectTypeLocal
}

func (e *SublimeEditor) GetExecutable() string {
	return "subl"
}

func (e *SublimeEditor) BuildArgs(p *project.Project, window WindowType) ([]string, error) {
	args := make([]string, 0)

	switch window {
	case WindowTypeNew:
		args = append(args, "--new-window")
	case WindowTypeAdd:
		args = append(args, "--add")
	}

	if p.IsWorkspace {
		w, err := workspace.Load(p.RootPath)
		if err != nil {
			return nil, err
		}
		args = append(args, w.FoldersPath()...)
	} else {
		args = append(args, p.RootPath)
	}

	return args, nil
}

// NvimEditor implementa suporte ao Neovim
type NvimEditor struct{}

func (e *NvimEditor) Name() string {
	return "nvim"
}

func (e *NvimEditor) Aliases() []string {
	return []string{}
}

func (e *NvimEditor) SupportsProjectType(projectType project.ProjectType) bool {
	return projectType == project.ProjectTypeLocal || projectType == project.ProjectTypeWSL
}

func (e *NvimEditor) GetExecutable() string {
	return "nvim"
}

func (e *NvimEditor) BuildArgs(p *project.Project, window WindowType) ([]string, error) {
	return []string{p.Path}, nil
}
