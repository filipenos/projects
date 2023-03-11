package command

import (
	"fmt"

	"github.com/filipenos/projects/pkg/config"
	"github.com/filipenos/projects/pkg/log"
	"github.com/filipenos/projects/pkg/path"
	"github.com/filipenos/projects/pkg/project"
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
	name, _ := path.SafeName(params...)
	if name == "" {
		return project.ErrNameRequired
	}

	projects, err := project.Load(config.Load())
	if err != nil {
		return err
	}

	excluded := false
	aux := make([]project.Project, 0, len(projects))
	for i := range projects {
		if projects[i].Name == name && !excluded {
			excluded = true
		} else {
			aux = append(aux, projects[i])
		}
	}

	if !excluded {
		return fmt.Errorf("project '%s' not found", name)
	}

	log.Infof("Project '%s' removed successfully!", name)

	projects = aux
	return projects.Save(config.Load())
}
