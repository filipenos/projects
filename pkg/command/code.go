package command

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/filipenos/projects/pkg/config"
	"github.com/filipenos/projects/pkg/log"
	"github.com/filipenos/projects/pkg/path"
	"github.com/filipenos/projects/pkg/project"
	"github.com/filipenos/projects/pkg/workspace"
	"github.com/spf13/cobra"
)

// codeCmd represents the code command
var codeCmd = &cobra.Command{
	Use:   "code",
	Short: "Edit your project using the editor (code as default)",
	RunE:  code,
}

func init() {
	codeCmd.Flags().StringP("editor", "e", "code", "Code using (code|zed|subl) editor")
	codeCmd.Flags().BoolP("reuse", "r", false, "Reuse same window")

	rootCmd.AddCommand(codeCmd)
}

func code(cmdParam *cobra.Command, params []string) error {
	projects, err := project.Load(config.Load())
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

	editor := "code"
	args := make([]string, 0)

	editor = SafeStringFlag(cmdParam, "editor")
	reuse := SafeBoolFlag(cmdParam, "reuse")

	openType := "folder"

	switch editor {
	case "code":
		pos := "--new-window"
		if reuse {
			pos = "--reuse-window"
		}
		args = append(args, pos)

		if p.IsWorkspace && p.ProjectType == project.ProjectTypeSSH {
			args = append(args, "--file-uri")
			openType = "file"
		} else {
			args = append(args, "--folder-uri")
		}
		args = append(args, p.Path)

	case "subl", "sublime":
		if p.ProjectType != project.ProjectTypeLocal {
			return fmt.Errorf("sublime not support remote project")
		}
		if !reuse {
			args = append(args, "--new-window")
		}
		if p.IsWorkspace {
			w, err := workspace.Load(p.Path)
			if err != nil {
				return err
			}
			args = append(args, w.FoldersPath()...)
		} else {
			args = append(args, p.Path)
		}

	case "zed":
		if p.ProjectType != project.ProjectTypeLocal {
			return fmt.Errorf("zed not support remote project")
		}
		if p.IsWorkspace {
			w, err := workspace.Load(p.Path)
			if err != nil {
				return err
			}
			args = append(args, w.FoldersPath()...)
		} else {
			args = append(args, p.Path)
		}
	default:
		return fmt.Errorf("editor not supported")
	}

	log.Infof("open %s '%s' on '%s'", openType, p.Path, editor)

	cmd := exec.Command(editor, args...)
	cmd.Stdin = os.Stdin
	return cmd.Run()

}
