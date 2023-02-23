package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"add"},
	Short:   "Create new project",
	RunE:    create,
}

func create(cmdParam *cobra.Command, params []string) error {
	var (
		p   = &Project{}
		err error
	)

	switch len(params) {
	case 0:
		p.Name, p.Path = currentPwd()
	case 1:
		p.Name = strings.TrimSpace(params[0])
		_, p.Path = currentPwd()
	case 2:
		p.Name = strings.TrimSpace(params[0])
		p.Path = strings.TrimSpace(params[1])
	default:
		return errorf("invalid size of arguments")
	}

	if p.Path != "" && isExist(fmt.Sprintf("%s/.git", p.Path)) {
		cmd := exec.Command("git", "-C", p.Path, "remote", "get-url", "origin")
		out, _ := cmd.CombinedOutput()
		if scm := strings.TrimSpace(string(out)); scm != "" {
			p.SCM = scm
		}

	}

	if SafeBoolFlag(cmdParam, "editor") {
		p, err = editProject(p)
		if err != nil {
			return err
		}
	}

	if p.Name == "" {
		return ErrNameRequired
	}

	if !SafeBoolFlag(cmdParam, "no-validate") {
		if p.Path == "" {
			return ErrPathRequired
		}
		if !isExist(p.Path) {
			return ErrPathNoExist
		}
	}

	projects, err := Load(LoadSettings())
	if err != nil {
		return err
	}

	if p, _ := projects.Get(p.Name); p != nil {
		return errorf("project '%s' already add to projects", p.Name)
	}

	projects = append(projects, *p)
	if err := projects.Save(LoadSettings()); err != nil {
		return err
	}
	log("Add project: '%s' path: '%s'", p.Name, p.Path)
	return nil
}
