package command

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/filipenos/projects/pkg/config"
	"github.com/filipenos/projects/pkg/log"
	"github.com/filipenos/projects/pkg/path"
	"github.com/filipenos/projects/pkg/project"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:     "create",
		Aliases: []string{"add"},
		Short:   "Create new project",
		RunE:    create,
	})
}

func create(cmdParam *cobra.Command, params []string) error {
	var (
		p   = &project.Project{}
		err error
	)

	switch len(params) {
	case 0:
		p.Name, p.Path = path.CurrentPwd()
	case 1:
		p.Name = strings.TrimSpace(params[0])
		_, p.Path = path.CurrentPwd()
	case 2:
		p.Name = strings.TrimSpace(params[0])
		p.Path = strings.TrimSpace(params[1])
	default:
		return fmt.Errorf("invalid size of arguments")
	}

	if p.Path != "" && path.Exist(fmt.Sprintf("%s/.git", p.Path)) {
		cmd := exec.Command("git", "-C", p.Path, "remote", "get-url", "origin")
		out, _ := cmd.CombinedOutput()
		if scm := strings.TrimSpace(string(out)); scm != "" {
			p.SCM = scm
		}

	}

	if SafeBoolFlag(cmdParam, "editor") {
		p, err = project.EditProject(p)
		if err != nil {
			return err
		}
	}

	if p.Name == "" {
		return project.ErrNameRequired
	}

	if !SafeBoolFlag(cmdParam, "no-validate") {
		if p.Path == "" {
			return project.ErrPathRequired
		}
		if !path.Exist(p.Path) {
			return project.ErrPathNoExist
		}
	}

	projects, err := project.Load(config.Load())
	if err != nil {
		return err
	}

	if p, _ := projects.Get(p.Name); p != nil {
		return fmt.Errorf("project '%s' already add to projects", p.Name)
	}

	projects = append(projects, *p)
	if err := projects.Save(config.Load()); err != nil {
		return err
	}
	log.Infof("Add project: '%s' path: '%s'", p.Name, p.Path)

	return nil
}
