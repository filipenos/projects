package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update data of existing project",
	RunE:  update,
}

func update(cmdParam *cobra.Command, params []string) error {
	name, _ := safeName(params...)
	if name == "" {
		return ErrNameRequired
	}

	projects, err := Load(LoadSettings())
	if err != nil {
		return err
	}

	_, index := projects.Get(name)
	if index == -1 {
		return errorf("project '%s' not found", name)
	}
	p := &projects[index]
	if p == nil {
		return errorf("project '%s' not found", name)
	}

	edited, err := editProject(p)
	if err != nil {
		return err
	}
	if edited.Name == "" {
		return ErrNameRequired
	}
	if !SafeBoolFlag(cmdParam, "no-validate") {
		if edited.Path == "" {
			return ErrPathRequired
		}
		if !isExist(edited.Path) {
			return ErrPathNoExist
		}
	}

	projects[index] = *edited
	return projects.Save(LoadSettings())
}
