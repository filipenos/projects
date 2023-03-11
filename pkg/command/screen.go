package command

/*
import (
	"fmt"
	"os"
	"os/exec"

	"github.com/filipenos/projects/pkg/config"
	"github.com/filipenos/projects/pkg/file"
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
	cmd.Flags().BoolP("screen", "s", false, "Use screen to run command")

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
		toExec string
		args   []string
	)

	if SafeBoolFlag(cmdParam, "screen") || params[1] == "screen" {
		tmpFile, err := file.NewTempFile()
		if err != nil {
			return err
		}
		defer tmpFile.Remove()
		tmpFile.WriteString("startup_message off\n")
		tmpFile.WriteString("vbell off\n")
		tmpFile.WriteString("hardstatus alwayslastline\n")
		tmpFile.WriteString(`hardstatus string '%{= kG}[ %{G}%H %{g}][%= %{= kw}%?%-Lw%?%{r}(%{W}%n*%f%t%?(%u)%?%{r})%{w}%?%+Lw%?%?%= %{g}][%{B} %m-%d %{W}%c %{g}]'`)

		toExec = "screen"
		args = screenArgs(p.Name, tmpFile.Name(), params[1:])
	} else {
		toExec = params[1]
		if len(params) > 2 {
			args = params[2:]
		}
	}

	log.Infof("%s %v", toExec, args)

	cmd := exec.Command(toExec, args...)
	cmd.Dir = p.Path
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func screenArgs(projectName, configFile string, args []string) (r []string) {
	r = []string{"-t", projectName, "-S", projectName, "-c", configFile}
	if len(args) > 0 {
		r = append(r, args...)
	}
	return
}

*/
