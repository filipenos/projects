package cmd

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:                "exec",
		Short:              "Exec command inside your project",
		DisableFlagParsing: true,
		RunE:               execCmd,
	})
}

func execCmd(cmdParam *cobra.Command, params []string) error {
	projects, err := Load(LoadSettings())
	if err != nil {
		return err
	}

	p, _ := projects.Find(safeName(params...))
	if p == nil {
		return errorf("project not found")

	}
	if p.Path == "" {
		return errorf("project '%s' dont have path", p.Name)
	}
	if !isExist(p.Path) {
		return errorf("path '%s' of project '%s' not exists", p.Path, p.Name)
	}
	if len(params) < 2 {
		return errorf("required at least one argument")
	}

	toExec := params[1]
	args := []string{}
	if len(params) > 2 {
		args = params[2:]
	}

	cmd := exec.Command(toExec, args...)
	cmd.Dir = p.Path
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()

}
