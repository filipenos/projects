package command

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/filipenos/projects/pkg/config"
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
	projects, err := project.Load(config.Load())
	if err != nil {
		return err
	}

	p, _ := projects.Find(path.SafeName(params...))
	if p == nil {
		return fmt.Errorf("project not found")

	}
	if p.Path == "" {
		return fmt.Errorf("project '%s' dont have path", p.Name)
	}
	if !path.Exist(p.Path) {
		return fmt.Errorf("path '%s' of project '%s' not exists", p.Path, p.Name)
	}
	if len(params) < 2 {
		return fmt.Errorf("required at least one argument")
	}

	var (
		toExec = params[1]
		args   []string
	)
	if len(params) > 2 {
		args = params[2:]
	}

	cmd := exec.Command(toExec, args...)
	cmd.Dir = p.Path
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	return nil
}
