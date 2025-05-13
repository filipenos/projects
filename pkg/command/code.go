package command

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/filipenos/projects/pkg/log"
	"github.com/filipenos/projects/pkg/path"
	"github.com/filipenos/projects/pkg/project"
	"github.com/filipenos/projects/pkg/workspace"
	"github.com/spf13/cobra"
)

func init() {
	codeCmd := &cobra.Command{
		Use:     "code",
		Short:   fmt.Sprintf("Edit your project using the editor (%s as default)", cfg.Editor),
		Aliases: []string{"vscode", "subl", "zed", "nvim"},
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

	args := make([]string, 0)

	var editor string
	window := SafeStringFlag(cmdParam, "window")
	openType := "folder"

	switch cmdParam.CalledAs() {
	case "code", "vscode":
		editor = "code"

		switch window {
		case "reuse":
			args = append(args, "--reuse-window")
		case "add":
			args = append(args, "--add")
		default:
			args = append(args, "--new-window")
		}

		if p.IsWorkspace {
			args = append(args, "--file-uri")
			openType = "file"
		} else {
			args = append(args, "--folder-uri")
		}
		args = append(args, p.RootPath)

	case "subl", "sublime":
		editor = "subl"

		if p.ProjectType != project.ProjectTypeLocal {
			return fmt.Errorf("sublime not support remote project")
		}

		switch window {
		case "new":
			args = append(args, "--new-window")
		case "add":
			args = append(args, "--add")
			//default reuse window
		}

		if p.IsWorkspace {
			w, err := workspace.Load(p.RootPath)
			if err != nil {
				return err
			}
			args = append(args, w.FoldersPath()...)
		} else {
			args = append(args, p.RootPath)
		}

	case "zed":
		editor = "zed"

		if p.ProjectType != project.ProjectTypeLocal {
			return fmt.Errorf("zed not support remote project")
		}

		switch window {
		case "add":
			args = append(args, "--add")
			//default new window
		}

		if p.IsWorkspace {
			w, err := workspace.Load(p.RootPath)
			if err != nil {
				return err
			}
			args = append(args, w.FoldersPath()...)
		} else {
			args = append(args, p.RootPath)
		}

	case "nvim":
		editor = "nvim"

		switch p.ProjectType {
		case project.ProjectTypeLocal, project.ProjectTypeWSL:
		default:
			return fmt.Errorf("project type %s not supported", p.ProjectType)
		}

		args = append(args, p.Path)

	default:
		return fmt.Errorf("editor not supported")
	}

	log.Infof("open %s '%s' on '%s'", openType, p.RootPath, editor)

	cmd := exec.Command(editor, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()

}
