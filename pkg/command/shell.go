package command

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/filipenos/projects/pkg/log"
	"github.com/filipenos/projects/pkg/path"
	"github.com/filipenos/projects/pkg/project"
	"github.com/filipenos/projects/pkg/workspace"
	"github.com/spf13/cobra"
)

func init() {
	shellCmd := &cobra.Command{
		Use:     "shell",
		Short:   fmt.Sprintf("Open project using Shell (%s current)", CurrentShell()),
		Aliases: []string{"sh", "nu", "bash", "zsh"},
		RunE:    shell,
	}
	shellCmd.Flags().BoolVar(&shellNoWorkspaceDir, "no-workspace-dir", false, "Don't create workspace shell directory, use workspace parent")
	shellCmd.Flags().BoolVar(&shellNoRebuildWorkspaceDir, "no-rebuild-workspace-dir", false, "Don't rebuild workspace shell directory, create a new temp dir if needed")
	rootCmd.AddCommand(shellCmd)
}

var (
	shellNoWorkspaceDir        bool
	shellNoRebuildWorkspaceDir bool
)

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

	if shellNoWorkspaceDir && shellNoRebuildWorkspaceDir {
		return fmt.Errorf("flags --no-workspace-dir and --no-rebuild-workspace-dir are mutually exclusive")
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
			if shellNoWorkspaceDir {
				execDir = workspaceBaseDir(p)
			} else {
				workspaceDir, cleanup, err := buildWorkspaceShellDir(p, !shellNoRebuildWorkspaceDir)
				if err != nil {
					return err
				}
				defer cleanup()
				execDir = workspaceDir
			}
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
		sshHost, sshPath, err := p.SSHInfo()
		if err != nil {
			return err
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
	if err := cmd.Run(); err != nil {
		if isUserShellExit(err) {
			return nil
		}
		return err
	}
	return nil
}

func buildWorkspaceShellDir(p *project.Project, rebuild bool) (string, func(), error) {
	workspacePath := p.Path
	if workspacePath == "" {
		workspacePath = p.RootPath
	}

	ws, err := workspace.Load(workspacePath)
	if err != nil {
		return "", nil, err
	}

	workspaceName := strings.TrimSuffix(filepath.Base(workspacePath), ".code-workspace")
	if workspaceName == "" || workspaceName == ".code-workspace" {
		workspaceName = "workspace"
	}

	baseDir := filepath.Dir(workspacePath)
	workspaceDir := filepath.Join(baseDir, workspaceName)
	if rebuild {
		if err := os.MkdirAll(workspaceDir, 0o755); err != nil {
			return "", nil, err
		}
		if err := clearWorkspaceSymlinks(workspaceDir); err != nil {
			return "", nil, err
		}
	} else {
		if _, err := os.Stat(workspaceDir); err == nil {
			workspaceDir, err = os.MkdirTemp(baseDir, workspaceName+"-")
			if err != nil {
				return "", nil, err
			}
		} else {
			if !os.IsNotExist(err) {
				return "", nil, err
			}
			if err := os.MkdirAll(workspaceDir, 0o755); err != nil {
				return "", nil, err
			}
		}
	}

	cleanup := func() {
		if rebuild {
			return
		}
		_ = os.RemoveAll(workspaceDir)
	}

	seen := map[string]int{}
	for _, folderPath := range ws.FoldersPath() {
		baseName := filepath.Base(folderPath)
		if baseName == "" || baseName == "." || baseName == string(filepath.Separator) {
			baseName = "folder"
		}
		count := seen[baseName]
		seen[baseName] = count + 1
		if count > 0 {
			baseName = fmt.Sprintf("%s-%d", baseName, count+1)
		}

		linkPath := filepath.Join(workspaceDir, baseName)
		if err := os.Symlink(folderPath, linkPath); err != nil {
			cleanup()
			return "", nil, err
		}
	}

	return workspaceDir, cleanup, nil
}

func clearWorkspaceSymlinks(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink == 0 {
			continue
		}
		if err := os.Remove(filepath.Join(dir, entry.Name())); err != nil {
			return err
		}
	}
	return nil
}

func workspaceBaseDir(p *project.Project) string {
	workspacePath := p.Path
	if workspacePath == "" {
		workspacePath = p.RootPath
	}
	return filepath.Dir(workspacePath)
}

func CurrentShell() (s string) {
	s = os.Getenv("SHELL")
	if s == "" {
		s = "bash"
	}
	return
}

func shellQuote(s string) string {
	if s == "" {
		return "''"
	}
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

func commandSeparator(shell string) string {
	// nushell doesn't support && operator
	if strings.Contains(shell, "nu") {
		return ";"
	}
	return "&&"
}

func isUserShellExit(err error) bool {
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		return false
	}
	if exitErr.ExitCode() == 130 {
		return true
	}
	if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
		if status.Signaled() && status.Signal() == os.Interrupt {
			return true
		}
	}
	return false
}
