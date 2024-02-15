package command

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/filipenos/projects/pkg/config"
	"github.com/filipenos/projects/pkg/log"
	"github.com/filipenos/projects/pkg/path"
	"github.com/filipenos/projects/pkg/project"
	"github.com/spf13/cobra"
)

// codeCmd represents the code command
var codeCmd = &cobra.Command{
	Use:   "code",
	Short: "Edit your project using the editor (code as default)",
	RunE:  code,
}

func init() {
	codeCmd.Flags().StringP("editor", "e", "code", "Code using vscode editor")
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

	if edit := SafeStringFlag(cmdParam, "e"); edit != "" {
		editor = edit
	} else {
		pos := "--new-window"
		if SafeBoolFlag(cmdParam, "r") {
			pos = "--reuse-window"
		}
		args = append(args, pos)
	}
	openType := "folder"
	if p.IsWorkspace && p.ProjectType == project.ProjectTypeSSH {
		args = append(args, "--file-uri")
		openType = "file"
	} else {
		args = append(args, "--folder-uri")
	}
	args = append(args, p.Path)

	log.Infof("open %s '%s' on '%s'", openType, p.Path, editor)

	cmd := exec.Command(editor, args...)
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
