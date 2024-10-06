package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/filipenos/projects/pkg/config"
	"github.com/filipenos/projects/pkg/log"
	"github.com/filipenos/projects/pkg/path"
	"github.com/filipenos/projects/pkg/project"
	"github.com/spf13/cobra"
)

func init() {
	shellCmd := &cobra.Command{
		Use:     "shell",
		Short:   fmt.Sprintf("Open project using Shell (%s current)", CurrentShell()),
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

	//TODO (filipenos) é possivel abrir o ssh, fiz no exec
	switch p.ProjectType {
	case project.ProjectTypeLocal, project.ProjectTypeWSL:
	default:
		return fmt.Errorf("project type %s not supported", p.ProjectType)
	}

	shell := cmdParam.CalledAs()
	switch cmdParam.CalledAs() {
	case "shell", "sh":
		shell = CurrentShell()
	case "zsh", "bash", "nu":
	default:
		return fmt.Errorf("shell not supported")
	}

	log.Infof("shell %s on '%s'", shell, p.RootPath)

	path := p.Path
	if p.IsWorkspace {
		parts := strings.Split(p.Path, "/")
		path = strings.Join(parts[:len(parts)-1], "/")
	}

	args := []string{"-c", fmt.Sprintf("cd %s; exec %s", path, shell)}

	cmd := exec.Command(shell, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func CurrentShell() (s string) {
	s = os.Getenv("SHELL")
	if s == "" {
		s = "bash"
	}
	return
}
