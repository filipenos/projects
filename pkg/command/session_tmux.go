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

type tmuxBackend struct{}

func newTmuxBackend() sessionBackend {
	return &tmuxBackend{}
}

func (b *tmuxBackend) Name() string {
	return "tmux"
}

func (b *tmuxBackend) Aliases() []string {
	return nil
}

func (b *tmuxBackend) Run(p *project.Project, backendArgs []string) error {
	switch p.ProjectType {
	case project.ProjectTypeLocal, project.ProjectTypeWSL:
	default:
		return fmt.Errorf("project type %s not supported for tmux", p.ProjectType)
	}

	if err := path.EnsureExecutable("tmux"); err != nil {
		return err
	}

	sessionName := sanitizeSessionName(p.Name)
	workingDir := projectWorkingDir(p)

	sessionExists, err := tmuxSessionExists(sessionName)
	if err != nil {
		return err
	}

	if sessionExists && len(backendArgs) > 0 {
		return fmt.Errorf("tmux session '%s' already exists; close it before executing a new command", sessionName)
	}

	var args []string
	if sessionExists {
		args = []string{"attach-session", "-d", "-t", sessionName}
	} else {
		args = []string{"new-session", "-s", sessionName, "-c", workingDir}
	}

	args = append(args, backendArgs...)

	log.Infof("tmux %s", strings.Join(args, " "))

	cmd := exec.Command("tmux", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func tmuxSessionExists(session string) (bool, error) {
	cmd := exec.Command("tmux", "has-session", "-t", session)
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
