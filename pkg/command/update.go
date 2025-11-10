package command

import (
	"fmt"

	"github.com/filipenos/projects/pkg/path"
	"github.com/filipenos/projects/pkg/project"
	"github.com/spf13/cobra"
)

func init() {
	updateCmd.Flags().Bool("no-validate", false, "Skip path validation")
	rootCmd.AddCommand(updateCmd)
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update data of existing project",
	RunE:  update,
}

func update(cmdParam *cobra.Command, params []string) error {
	name, _ := path.SafeName(params...)
	if name == "" {
		return project.ErrNameRequired
	}

	projects, err := project.Load(cfg)
	if err != nil {
		return err
	}

	_, index := projects.Get(name)
	if index == -1 {
		return fmt.Errorf("project '%s' not found", name)
	}
	p := &projects[index]

	edited, err := project.EditProject(p)
	if err != nil {
		return err
	}
	if edited.Name == "" {
		return project.ErrNameRequired
	}
	if !SafeBoolFlag(cmdParam, "no-validate") {
		if edited.RootPath == "" {
			return project.ErrPathRequired
		}
		if !path.Exist(edited.RootPath) {
			return project.ErrPathNoExist
		}
	}

	projects[index] = *edited
	return projects.Save(cfg)
}
