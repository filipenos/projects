package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/urfave/cli"
)

var (
	filepath = fmt.Sprintf("%s/.projects.json", os.Getenv("HOME"))
)

func add(c *cli.Context) error {
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
	} else if c.Bool("editor") {
		p, err = editProject(p)
		if err != nil {
			return err
		}
	} else {
		p.Name = strings.TrimSpace(c.Args().Get(0))
		p.Path = strings.TrimSpace(c.Args().Get(1))
	}

	if p.Name == "" {
		return fmt.Errorf("name is required")
	}
	if p.Path == "" {
		return fmt.Errorf("path is required")
	}

	if !isExist(p.Path) {
		return fmt.Errorf("path is no exists")
	}

	projects, err := Load(filepath)
	if err != nil {
		return err
	}
	projects.AddProject(*p)
	return projects.Save()
}

func remove(c *cli.Context) error {
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
		return fmt.Errorf("Project %s not found", name)
	}

	log("Project %s removed successfully!", name)
	projects.Projects = aux
	return projects.Save()
}

func list(c *cli.Context) error {
	projects, err := Load(filepath)
	if err != nil {
		return err
	}
	for _, p := range projects.Projects {
		fmt.Printf("%s\n  %s\n", p.Name, p.Path)
	}
	return nil
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
		return fmt.Errorf("Project %s not found", name)
	}

	log("open path %s", path)

	cmd := exec.Command("tmux", "new", "-s", name, "-c", path)
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

func edit(c *cli.Context) error {
	name := strings.TrimSpace(c.Args().First())
	if name == "" {
		return fmt.Errorf("name is required")
	}

	f, err := Load(filepath)
	if err != nil {
		return err
	}

	var index int
	for i := range f.Projects {
		if f.Projects[i].Name == name {
			index = i
			break
		}
	}
	p := &f.Projects[index]
	if p == nil {
		return fmt.Errorf("project %s not found", name)
	}

	edited, err := editProject(p)
	if err != nil {
		return err
	}

	f.Projects[index] = *edited
	return f.Save()
}

func editProject(p *Project) (*Project, error) {
	tmp, err := NewTempFile()
	if err != nil {
		return nil, err
	}
	defer tmp.Remove()

	d := `name={{.Name}}
path={{.Path}}`

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
		}
	}
	return p
}
