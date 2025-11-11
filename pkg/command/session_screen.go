package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/filipenos/projects/pkg/log"
	"github.com/filipenos/projects/pkg/path"
	"github.com/filipenos/projects/pkg/project"
)

type screenBackend struct{}

func newScreenBackend() sessionBackend {
	return &screenBackend{}
}

func (b *screenBackend) Name() string {
	return "screen"
}

func (b *screenBackend) Aliases() []string {
	return nil
}

func (b *screenBackend) Run(p *project.Project, backendArgs []string) error {
	switch p.ProjectType {
	case project.ProjectTypeLocal, project.ProjectTypeWSL:
	default:
		return fmt.Errorf("project type %s not supported for screen", p.ProjectType)
	}

	if err := path.EnsureExecutable("screen"); err != nil {
		return err
	}

	sessionName := sanitizeSessionName(p.Name)
	workingDir := projectWorkingDir(p)

	sessionExists, err := screenSessionExists(sessionName)
	if err != nil {
		return err
	}

	if sessionExists && len(backendArgs) > 0 {
		return fmt.Errorf("screen session '%s' already exists; close it before executing a new command", sessionName)
	}

	args := []string{"-S", sessionName, "-d", "-RR"}
	args = append(args, backendArgs...)

	log.Infof("screen %s", strings.Join(args, " "))

	cmd := exec.Command("screen", args...)
	cmd.Dir = workingDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func screenSessionExists(session string) (bool, error) {
	cmd := exec.Command("screen", "-S", session, "-Q", "select", ".")
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 1 {
				return false, nil
			}
		}
		return false, err
	}
	return true, nil
}
