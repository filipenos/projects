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

func init() {
	shellCmd := &cobra.Command{
		Use:     "shell",
		Short:   "Open project using Shell",
		Aliases: []string{"sh", "nu", "bash", "zsh"},
		RunE:    shell,
	}
	rootCmd.AddCommand(shellCmd)
}

func shell(cmdParam *cobra.Command, params []string) error {
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

	switch p.ProjectType {
	case project.ProjectTypeLocal, project.ProjectTypeWSL:
	default:
		return fmt.Errorf("project type %s not supported", p.ProjectType)
	}

	var shell string
	args := make([]string, 0)
	switch cmdParam.CalledAs() {
	case "shell", "sh", "nu":
		shell = "nu"
		args = append(args, fmt.Sprintf("-e 'cd %s'", p.Path))
	case "zsh":
		shell = "zsh"
		args = append(args, "-c", fmt.Sprintf("cd %s; exec zsh", p.Path))
	case "bash":
		shell = "bash"
		args = append(args, "-c", fmt.Sprintf("cd %s; exec bash", p.Path))
	default:
		return fmt.Errorf("shell not supported")
	}

	log.Infof("shell %s on '%s'", shell, p.RootPath)

	cmd := exec.Command(shell, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
