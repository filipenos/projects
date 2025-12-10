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
	cmd := &cobra.Command{
		Use:                "exec",
		Short:              "Exec command inside your project",
		DisableFlagParsing: true,
		RunE:               execCmd,
	}
	rootCmd.AddCommand(cmd)
}

func execCmd(cmdParam *cobra.Command, params []string) error {
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
	if len(params) < 2 {
		return fmt.Errorf("missing command to execute inside project")
	}

	var (
		command = params[1]
		args    []string
		workDir string
	)

	switch p.ProjectType {
	case project.ProjectTypeLocal:
		workDir = p.RootPath
		if len(params) > 2 {
			args = params[2:]
		}

	case project.ProjectTypeSSH:
		log.Infof("executing on ssh host")
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

		// Constrói o comando completo com argumentos
		var remoteCmd strings.Builder
		remoteCmd.WriteString(fmt.Sprintf("cd %s && %s", sshPath, command))
		if len(params) > 2 {
			for _, arg := range params[2:] {
				remoteCmd.WriteString(" ")
				remoteCmd.WriteString(arg)
			}
		}

		args = []string{sshHost, "-t", remoteCmd.String()}
		command = "ssh"
		workDir = "" // SSH doesn't use local workDir

	default:
		return fmt.Errorf("project type %s not supported for exec command", p.ProjectType)
	}

	cmd := exec.Command(command, args...)
	cmd.Dir = workDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
