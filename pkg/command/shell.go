package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

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
	if len(params) == 0 {
		return fmt.Errorf("project name is required")
	}
	
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

	shell := cmdParam.CalledAs()
	switch cmdParam.CalledAs() {
	case "shell", "sh":
		shell = CurrentShell()
	case "zsh", "bash", "nu":
	default:
		return fmt.Errorf("shell not supported")
	}

	log.Infof("shell %s on '%s'", shell, p.RootPath)

	var (
		command string
		args    []string
		execDir string
	)

	switch p.ProjectType {
	case project.ProjectTypeLocal, project.ProjectTypeWSL:
		if err := path.EnsureExecutable(shell); err != nil {
			return err
		}

		execDir = p.Path
		if p.IsWorkspace {
			parts := strings.Split(p.Path, "/")
			execDir = strings.Join(parts[:len(parts)-1], "/")
		}

		command = shell
		sep := commandSeparator(shell)
		args = []string{"-c", fmt.Sprintf("cd %s %s exec %s", execDir, sep, shell)}

	case project.ProjectTypeSSH:
		// SSH connections use the remote default shell, not the local alias
		if cmdParam.CalledAs() != "shell" && cmdParam.CalledAs() != "sh" {
			log.Infof("warning: SSH connections use the remote server's default shell, ignoring '%s' alias", cmdParam.CalledAs())
			shell = CurrentShell()
		}

		log.Infof("opening shell on ssh host")
		i := strings.Index(p.RootPath, "+")
		if i == -1 {
			return fmt.Errorf("invalid path format")
		}
		subPath := p.RootPath[i+1:]
		i = strings.Index(subPath, "/")
		if i == -1 {
			return fmt.Errorf("invalid path format")
		}
		sshHost := subPath[:i]
		sshPath := subPath[i:]

		if p.IsWorkspace {
			parts := strings.Split(sshPath, "/")
			if len(parts) > 1 {
				sshPath = strings.Join(parts[:len(parts)-1], "/")
			}
		}

		command = "ssh"
		sep := commandSeparator(shell)
		args = []string{sshHost, "-t", fmt.Sprintf("cd %s %s exec %s", sshPath, sep, shell)}
		execDir = "" // SSH doesn't use local workDir

	default:
		return fmt.Errorf("project type %s not supported", p.ProjectType)
	}

	cmd := exec.Command(command, args...)
	cmd.Dir = execDir
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

func commandSeparator(shell string) string {
	// nushell doesn't support && operator
	if strings.Contains(shell, "nu") {
		return ";"
	}
	return "&&"
}
