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

// shellCmd represents the list command
var shellCmd = &cobra.Command{
	Use:     "shell",
	Aliases: []string{"sh", "nu", "n"},
	Short:   "Open project using Shell",
	RunE:    shell,
}

func init() {
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

	cmd := exec.Command("nu", fmt.Sprintf("-e 'cd %s'", p.Path))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
