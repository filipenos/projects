package cmd

import (
	"os"
	"sort"
	"text/template"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
	codeCmd.Flags().BoolP("path", "p", false, "Reuse same window")
	codeCmd.Flags().BoolP("simple", "s", false, "Reuse same window")
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "l"},
	Short:   "List projects",
	RunE:    list,
}

func list(cmdParam *cobra.Command, params []string) error {
	projects, err := Load(LoadSettings())
	if err != nil {
		return errorf("error on load file: %v", err)
	}
	sort.Sort(projects)

	t := `{{range .Projects}}{{.Name}}{{if $.ExtraInfo}}{{if .Opened}} (opened){{end}}{{if .Attached}} (attached){{end}}{{if not .ValidPath}} (invalid-path){{end}}{{end}}{{if $.Path}}
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
		return errorf("error on execute template: %v", err)
	}
	return nil
}
