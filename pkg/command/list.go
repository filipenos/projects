package command

import (
	"fmt"
	"sort"

	"github.com/filipenos/projects/pkg/log"
	"github.com/filipenos/projects/pkg/project"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "l"},
	Short:   "List projects",
	RunE:    list,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func list(cmdParam *cobra.Command, params []string) error {
	projects, err := project.Load(cfg)
	if err != nil {
		return fmt.Errorf("error on load file: %v", err)
	}
	sort.Sort(projects)

	for _, p := range projects {
		print := fmt.Sprintf("%s %s", p.Name, string(p.ProjectType))
		if p.IsWorkspace {
			print += " (w)"
		}
		if p.Opened {
			print += " (opened)"
		}
		if p.Attached {
			print += " (attached)"
		}
		if !p.ValidPath {
			print += " (invalid-path)"
		}
		log.Infof("%s", print)
	}

	return nil
}
