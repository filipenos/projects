package command

import (
	"fmt"
	"sort"

	"github.com/filipenos/projects/pkg/log"
	"github.com/filipenos/projects/pkg/project"
	"github.com/spf13/cobra"
)

var (
	listSSH       bool
	listLocal     bool
	listWorkspace bool
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "l"},
	Short:   "List projects",
	RunE:    list,
}

func init() {
	listCmd.Flags().BoolVar(&listSSH, "ssh", false, "List only SSH projects")
	listCmd.Flags().BoolVar(&listLocal, "local", false, "List only local projects")
	listCmd.Flags().BoolVar(&listWorkspace, "workspace", false, "List only workspace projects")
	rootCmd.AddCommand(listCmd)
}

func list(cmdParam *cobra.Command, params []string) error {
	projects, err := project.Load(cfg)
	if err != nil {
		return fmt.Errorf("error on load file: %v", err)
	}
	sort.Sort(projects)

	// If no filters are specified, show all
	showAll := !listSSH && !listLocal && !listWorkspace

	for _, p := range projects {
		// Apply filters with AND logic
		if !showAll {
			// Check SSH filter
			if listSSH && !(p.ProjectType == project.ProjectTypeSSH) {
				continue
			}
			// Check Local filter
			if listLocal && !(p.ProjectType == project.ProjectTypeLocal || p.ProjectType == project.ProjectTypeWSL) {
				continue
			}
			// Check Workspace filter
			if listWorkspace && !p.IsWorkspace {
				continue
			}
		}

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
		log.Println(print)
	}

	return nil
}
