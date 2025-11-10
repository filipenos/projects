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
	tmuxCmd := &cobra.Command{
		Use:                "tmux <project> [tmux args...]",
		Short:              "Open or attach a tmux session inside the project directory",
		DisableFlagParsing: true,
		RunE:               runTmux,
	}
	rootCmd.AddCommand(tmuxCmd)
}

func runTmux(cmdParam *cobra.Command, params []string) error {
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

	switch p.ProjectType {
	case project.ProjectTypeLocal, project.ProjectTypeWSL:
	default:
		return fmt.Errorf("project type %s not supported for tmux", p.ProjectType)
	}

	sessionName := sanitizeSessionName(p.Name)
	workingDir := projectWorkingDir(p)

	var tmuxArgs []string
	if len(params) > 1 {
		tmuxArgs = append(tmuxArgs, params[1:]...)
	}

	if err := path.EnsureExecutable("tmux"); err != nil {
		return err
	}

	sessionExists, err := tmuxSessionExists(sessionName)
	if err != nil {
		return err
	}

	if sessionExists && len(tmuxArgs) > 0 {
		return fmt.Errorf("tmux session '%s' already exists; close it before executing a new command", sessionName)
	}

	var args []string
	if sessionExists {
		// Attach detaching other clients to honor "fechando as outras"
		args = []string{"attach-session", "-d", "-t", sessionName}
	} else {
		args = []string{"new-session", "-s", sessionName, "-c", workingDir}
	}

	args = append(args, tmuxArgs...)

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

func projectWorkingDir(p *project.Project) string {
	pathValue := p.Path
	if pathValue == "" {
		pathValue = p.RootPath
	}
	if p.IsWorkspace {
		parts := strings.Split(pathValue, "/")
		if len(parts) > 1 {
			pathValue = strings.Join(parts[:len(parts)-1], "/")
		}
	}
	return pathValue
}

func sanitizeSessionName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "project"
	}

	var b strings.Builder
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '-' || r == '_' || r == '.':
			b.WriteRune(r)
		default:
			b.WriteRune('-')
		}
	}
	return b.String()
}
