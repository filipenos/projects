package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

//TODO (filipenos) o project manager possui a variavel $home, podemos utilizala para gravar o path, ai nao importa qual distro esteja usando

// Errors returned
var (
	ErrNameRequired = errorf("name is required")
	ErrPathRequired = errorf("path is required")
	ErrPathNoExist  = errorf("path is no exists")
)

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

func code(cmdParam *cobra.Command, params []string) error {
	projects, err := Load(LoadSettings())
	if err != nil {
		return err
	}

	p, _ := projects.Find(safeName(params...))
	if p == nil {
		return errorf("project not found")

	}
	if p.Path == "" {
		return errorf("project '%s' dont have path", p.Name)
	}
	if !isExist(p.Path) {
		return errorf("path '%s' of project '%s' not exists", p.Path, p.Name)
	}

	editor := "code"
	args := make([]string, 0)

	if edit := SafeStringFlag(cmdParam, "e"); edit != "" {
		editor = edit
	} else {
		pos := "--new-window"
		if SafeBoolFlag(cmdParam, "r") {
			pos = "--reuse-window"
		}
		args = append(args, pos)
	}
	args = append(args, p.Path)

	log("open path '%s' on '%s'", p.Path, editor)

	cmd := exec.Command(editor, args...)
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
