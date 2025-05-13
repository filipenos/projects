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
		return fmt.Errorf("required at least one argument")
	}

	var (
		toExec = params[1]
		args   []string
		path   string
	)

	switch p.ProjectType {
	case project.ProjectTypeLocal:
		path = p.RootPath
	case project.ProjectTypeSSH:
		log.Infof("executing on ssh host, some features not working")
		i := strings.Index(p.RootPath, "+")
		if i == -1 {
			return fmt.Errorf("invalid path")
		}
		subPath := p.RootPath[i+1:]
		i = strings.Index(subPath, "/")
		if i == -1 {
			return fmt.Errorf("invalid path")
		}
		sshHost := subPath[:i]
		subPath = subPath[i:]

		args = append(args, sshHost)
		args = append(args, "-t") // Force pseudo-tty allocation
		args = append(args, "cd "+subPath+";"+toExec)
		toExec = "ssh"
	}

	if len(params) > 2 {
		args = append(args, params[2:]...)
	}

	cmd := exec.Command(toExec, args...)
	cmd.Dir = path
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
