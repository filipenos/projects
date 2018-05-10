package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/urfave/cli"
)

var (
	filepath      = fmt.Sprintf("%s/.projects.json", os.Getenv("HOME"))
	ErrUnmodified = errors.New("Unmodified")
)

func create(c *cli.Context) error {
	c.Args()

	var (
		p   = &Project{}
		err error
	)

	if c.Bool("current") {
		pwd := os.Getenv("PWD")
		paths := strings.Split(pwd, "/")

		p.Name = strings.TrimSpace(paths[len(paths)-1])
		p.Path = strings.TrimSpace(pwd)
	} else {
		p.Name = strings.TrimSpace(c.Args().Get(0))
		p.Path = strings.TrimSpace(c.Args().Get(1))
		p.SCM = strings.TrimSpace(c.Args().Get(2))
	}

	if c.Bool("editor") {
		p, err = editProject(p)
		if err != nil {
			return err
		}
	}

	if p.Name == "" {
		return fmt.Errorf("name is required")
	}

	if c.Bool("validate-path") {
		if p.Path == "" {
			return fmt.Errorf("path is required")
		}
		if !isExist(p.Path) {
			return fmt.Errorf("path is no exists")
		}
	}

	projects, err := Load(filepath)
	if err != nil {
		return err
	}

	if p, _ := projects.Get(p.Name); p != nil {
		return fmt.Errorf("project '%s' already add to projects", p.Name)
	}

	projects.AddProject(*p)
	if err := projects.Save(); err != nil {
		return err
	}
	log("Add project: '%s' path: '%s'", p.Name, p.Path)
	return nil
}

func delete(c *cli.Context) error {
	name := strings.TrimSpace(c.Args().First())
	if name == "" {
		return fmt.Errorf("name is required")
	}

	projects, err := Load(filepath)
	if err != nil {
		return err
	}

	excluded := false
	aux := make([]Project, 0, len(projects.Projects))
	for i := range projects.Projects {
		if projects.Projects[i].Name == name && !excluded {
			excluded = true
		} else {
			aux = append(aux, projects.Projects[i])
		}
	}

	if !excluded {
		return fmt.Errorf("Project '%s' not found", name)
	}

	log("Project '%s' removed successfully!", name)
	projects.Projects = aux
	return projects.Save()
}

func list(c *cli.Context) error {
	projects, err := Load(filepath)
	if err != nil {
		return err
	}

	t := `{{range .Projects.Projects}}Name: {{.Name}}{{if $.Full}}
	Path: {{.Path}}
	SMC: {{.SCM}}{{end}}
{{else}}No projects yeat!
{{end}}`
	tmpl := template.Must(template.New("editor").Parse(t))
	ctx := map[string]interface{}{"Projects": projects, "Full": c.Bool("full")}
	return tmpl.Execute(os.Stdout, ctx)
}

func open(c *cli.Context) error {
	name := strings.TrimSpace(c.Args().First())
	if name == "" {
		return fmt.Errorf("name is required")
	}

	projects, err := Load(filepath)
	if err != nil {
		return err
	}

	var path string
	for _, p := range projects.Projects {
		if p.Name == name {
			path = p.Path
			break
		}
	}

	if path == "" {
		return fmt.Errorf("Project '%s' not found", name)
	}

	log("open path '%s'", path)

	cmd := exec.Command("tmux", "new", "-s", name, "-n", name, "-c", path)
	//option -d run tmux with daemon
	// tmux new-session -d -s mySession -n myWindow
	// tmux send-keys -t mySession:myWindow "cd /my/directory" Enter
	// tmux send-keys -t mySession:myWindow "vim" Enter
	// tmux attach -t mySession:myWindow
	cmd.Stdin = os.Stdin
	out, err := cmd.CombinedOutput()
	if strings.Contains(string(out), "duplicate session") {
		cmd = exec.Command("tmux", "attach", "-t", name)
		cmd.Stdin = os.Stdin
		_, err := cmd.CombinedOutput()
		if err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func update(c *cli.Context) error {
	name := strings.TrimSpace(c.Args().First())
	if name == "" {
		return fmt.Errorf("name is required")
	}

	projects, err := Load(filepath)
	if err != nil {
		return err
	}

	_, index := projects.Get(name)
	if index == -1 {
		return fmt.Errorf("project '%s' not found", name)
	}
	p := &projects.Projects[index]
	if p == nil {
		return fmt.Errorf("project '%s' not found", name)
	}

	edited, err := editProject(p)
	if err != nil {
		return err
	}
	if edited.Name == "" {
		return fmt.Errorf("name is required")
	}
	if c.Bool("validate-path") {
		if edited.Path == "" {
			return fmt.Errorf("path is required")
		}
		if !isExist(edited.Path) {
			return fmt.Errorf("path is no exists")
		}
	}

	projects.Projects[index] = *edited
	return projects.Save()
}

func editProject(p *Project) (*Project, error) {
	tmp, err := NewTempFile()
	if err != nil {
		return nil, err
	}
	defer tmp.Remove()

	d := `name={{.Name}}
path={{.Path}}
scm={{.SCM}}`

	tmpl := template.Must(template.New("editor").Parse(d))
	if err := tmpl.Execute(tmp, p); err != nil {
		return nil, err
	}

	tmp.ReadFromUser()

	if err := tmp.Close(); err != nil {
		return nil, err
	}

	content, err := tmp.GetContent()
	if err != nil {
		return nil, err
	}
	return parseContent(content), nil
}

func parseContent(data []byte) *Project {
	lines := strings.Split(string(data), "\n")
	p := &Project{}
	for i := range lines {
		line := strings.TrimSpace(lines[i])
		values := strings.Split(line, "=")
		if len(values) != 2 {
			continue
		}
		switch values[0] {
		case "name":
			p.Name = strings.TrimSpace(values[1])
		case "path":
			p.Path = strings.TrimSpace(values[1])
		case "scm":
			p.SCM = strings.TrimSpace(values[1])
		case "scmtype":
			p.SCMType = strings.TrimSpace(values[1])
		}
	}
	return p
}

func getProject(c *cli.Context) error {
	name := strings.TrimSpace(c.Args().First())
	if name == "" {
		return fmt.Errorf("name is required")
	}

	projects, err := Load(filepath)
	if err != nil {
		return err
	}

	_, index := projects.Get(name)
	if index == -1 {
		return fmt.Errorf("project '%s' not found", name)
	}
	p := &projects.Projects[index]
	if p == nil {
		return fmt.Errorf("project '%s' not found", name)
	}
	if p.SCM == "" {
		return errorf("project '%s' dont have scm configured", p.Name)
	}
	if isExist(p.Path) {
		return errorf("path '%s' of project '%s' already exists", p.Path, p.Name)
	}

	cmd := exec.Command("git", "clone", p.SCM, p.Path)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	log(string(out))
	return nil
}

func checkRequirements() {}
