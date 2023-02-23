package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a existent project",
	RunE:  delete,
}

func delete(cmdParam *cobra.Command, params []string) error {
	name, _ := safeName(params...)
	if name == "" {
		return ErrNameRequired
	}

	projects, err := Load(LoadSettings())
	if err != nil {
		return err
	}

	excluded := false
	aux := make([]Project, 0, len(projects))
	for i := range projects {
		if projects[i].Name == name && !excluded {
			excluded = true
		} else {
			aux = append(aux, projects[i])
		}
	}

	if !excluded {
		return errorf("Project '%s' not found", name)
	}

	log("Project '%s' removed successfully!", name)
	projects = aux
	return projects.Save(LoadSettings())
}
