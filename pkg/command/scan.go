package command

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/filipenos/projects/pkg/log"
	"github.com/filipenos/projects/pkg/project"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "scan [directory]",
		Short: "Scan a directory and add all child directories as projects",
		RunE:  scan,
	}
	rootCmd.AddCommand(cmd)
}

func scan(cmdParam *cobra.Command, params []string) error {
	var dir string
	if len(params) > 0 {
		dir = strings.TrimSpace(params[0])
	} else {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
	}

	dir, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	projects, err := project.Load(cfg)
	if err != nil {
		return err
	}

	var added int
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		name := entry.Name()
		fullPath := filepath.Join(dir, name)

		if p, _ := projects.Get(name); p != nil {
			log.Infof("skip: '%s' already exists", name)
			continue
		}
		if p, _ := projects.GetByPath(fullPath); p != nil {
			log.Infof("skip: '%s' already exists (path match)", p.Name)
			continue
		}

		p := project.Project{
			Name:     name,
			RootPath: fullPath,
			Enabled:  true,
		}

		gitDir := filepath.Join(fullPath, ".git")
		if info, err := os.Stat(gitDir); err == nil && info.IsDir() {
			cmd := exec.Command("git", "-C", fullPath, "remote", "get-url", "origin")
			out, _ := cmd.CombinedOutput()
			if scm := strings.TrimSpace(string(out)); scm != "" {
				p.SCM = scm
			}
		}

		projects = append(projects, p)
		log.Infof("added: '%s' -> %s", name, fullPath)
		added++
	}

	if added == 0 {
		log.Infof("no new projects found in %s", dir)
		return nil
	}

	if err := projects.Save(cfg); err != nil {
		return err
	}
	log.Infof("added %d project(s)", added)

	return nil
}
