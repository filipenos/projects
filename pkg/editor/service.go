package editor

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/filipenos/projects/pkg/log"
	"github.com/filipenos/projects/pkg/path"
	"github.com/filipenos/projects/pkg/project"
	"github.com/filipenos/projects/pkg/workspace"
)

type WindowType string

const (
	WindowTypeNew   WindowType = "new"
	WindowTypeReuse WindowType = "reuse"
	WindowTypeAdd   WindowType = "add"
)

type Editor struct {
	Name       string
	Executable string
	Aliases    []string
	LocalOnly  bool
	BuildArgs  func(p *project.Project, window WindowType) ([]string, error)
}

var editors = []Editor{
	// VSCode-based editors (support local + remote)
	vscodeEditor("code", "code", []string{"vscode"}),
	vscodeEditor("cursor", "cursor", nil),
	vscodeEditor("windsurf", "windsurf", nil),
	vscodeEditor("antigravity", "antigravity", nil),
	// Simple editors (local only, just receive the path)
	simpleEditor("vim", "vim", nil),
	simpleEditor("nvim", "nvim", nil),
	simpleEditor("emacs", "emacs", nil),
	simpleEditor("zed", "zed", nil),
	{
		Name:       "sublime",
		Executable: "subl",
		Aliases:    []string{"subl"},
		LocalOnly:  true,
		BuildArgs: func(p *project.Project, window WindowType) ([]string, error) {
			var args []string
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
		},
	},
	simpleEditor("intellij", "idea", []string{"idea"}),
	simpleEditor("goland", "goland", nil),
}

func simpleEditor(name, executable string, aliases []string) Editor {
	return Editor{
		Name:       name,
		Executable: executable,
		Aliases:    aliases,
		LocalOnly:  true,
		BuildArgs: func(p *project.Project, window WindowType) ([]string, error) {
			return []string{p.RootPath}, nil
		},
	}
}

func vscodeEditor(name, executable string, aliases []string) Editor {
	return Editor{
		Name:       name,
		Executable: executable,
		Aliases:    aliases,
		BuildArgs: func(p *project.Project, window WindowType) ([]string, error) {
			var args []string
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
		},
	}
}

type Service struct {
	byName map[string]*Editor
}

func NewService() *Service {
	s := &Service{byName: make(map[string]*Editor)}
	for i := range editors {
		e := &editors[i]
		s.byName[e.Name] = e
		for _, alias := range e.Aliases {
			s.byName[alias] = e
		}
	}
	return s
}

func (s *Service) Aliases() []string {
	seen := make(map[string]bool)
	var aliases []string
	for _, e := range editors {
		if e.Name == "code" {
			for _, a := range e.Aliases {
				if !seen[a] {
					aliases = append(aliases, a)
					seen[a] = true
				}
			}
			continue
		}
		if !seen[e.Name] {
			aliases = append(aliases, e.Name)
			seen[e.Name] = true
		}
		for _, a := range e.Aliases {
			if !seen[a] {
				aliases = append(aliases, a)
				seen[a] = true
			}
		}
	}
	return aliases
}

func (s *Service) GetEditors() (available []string, notAvailable []string) {
	for _, e := range editors {
		if path.ExistsInPathOrAsFile(e.Executable) {
			available = append(available, e.Name)
		} else {
			notAvailable = append(notAvailable, e.Name)
		}
	}
	return
}

func (s *Service) OpenProject(editorName string, p *project.Project, window WindowType) error {
	e, ok := s.byName[editorName]
	if !ok {
		return fmt.Errorf("editor '%s' not found", editorName)
	}

	if e.LocalOnly && p.ProjectType != project.ProjectTypeLocal {
		return fmt.Errorf("editor '%s' does not support project type '%s'", e.Name, p.ProjectType)
	}

	args, err := e.BuildArgs(p, window)
	if err != nil {
		return fmt.Errorf("failed to build args for editor '%s': %w", e.Name, err)
	}

	if err := path.EnsureExecutable(e.Executable); err != nil {
		return fmt.Errorf("editor '%s': %w", e.Name, err)
	}

	openType := "folder"
	if p.IsWorkspace {
		openType = "file"
	}
	log.Infof("open %s '%s' on '%s'", openType, p.RootPath, e.Executable)

	cmd := exec.Command(e.Executable, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
