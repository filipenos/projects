package command

import (
	"fmt"
	"os"
	"sort"
	"text/template"

	"github.com/filipenos/projects/pkg/config"
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
	listCmd.Flags().BoolP("path", "p", false, "Reuse same window")
	listCmd.Flags().BoolP("simple", "s", false, "Reuse same window")

	rootCmd.AddCommand(listCmd)
}

func list(cmdParam *cobra.Command, params []string) error {
	projects, err := project.Load(config.Load())
	if err != nil {
		return fmt.Errorf("error on load file: %v", err)
	}
	sort.Sort(projects)

	t := `{{range .Projects}}{{.Name}}{{if $.ExtraInfo}} {{.ProjectType}} {{if .IsWorkspace}}(w){{end}} {{if .Opened}} (opened){{end}}{{if .Attached}} (attached){{end}}{{if not .ValidPath}} (invalid-path){{end}}{{end}}{{if $.Path}}
  Path: {{.Path}}{{end}}
{{else}}No projects yeat!
{{end}}`
	tmpl := template.Must(template.New("editor").Parse(t))
	ctx := map[string]interface{}{
		"Projects":  projects,
		"Path":      SafeBoolFlag(cmdParam, "path"),
		"ExtraInfo": !SafeBoolFlag(cmdParam, "simple"),
	}
	err = tmpl.Execute(os.Stdout, ctx)
	if err != nil {
		return fmt.Errorf("error on execute template: %v", err)
	}
	return nil
}
