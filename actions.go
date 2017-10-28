package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/urfave/cli"
)

func add(c *cli.Context) error {
	c.Args()

	var name, path string
	if c.Bool("current") {
		pwd := os.Getenv("PWD")
		paths := strings.Split(pwd, "/")

		name = paths[len(paths)-1]
		path = pwd
	} else {
		name = strings.TrimSpace(c.Args().Get(0))
		if name == "" {
			return fmt.Errorf("name is required")
		}
		path = strings.TrimSpace(c.Args().Get(1))
		if path == "" {
			return fmt.Errorf("path is required")
		}
	}

	if !isExist(path) {
		return fmt.Errorf("path is no exists")
	}

	projects, err := Load(filepath)
	if err != nil {
		return err
	}
	projects.Add(name, path)
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

	var p *Project
	var index int
	for i := range f.Projects {
		if f.Projects[i].Name == name {
			p = &f.Projects[i]
			index = i
			break
		}
	}
	if p == nil {
		return fmt.Errorf("project %s not found", name)
	}

	tmp, err := NewTempFile()
	if err != nil {
		return err
	}
	defer tmp.Remove()

	d := `name={{.Name}}
path={{.Path}}`

	tmpl := template.Must(template.New("test").Parse(d))
	if err := tmpl.Execute(tmp, p); err != nil {
		return err
	}

	tmp.ReadFromUser()

	if err := tmp.Close(); err != nil {
		return err
	}

	content, err := tmp.GetContent()
	if err != nil {
		return err
	}
	editProject := parseContent(content)
	f.Projects[index] = editProject
	return f.Save()
}

func parseContent(data []byte) (p Project) {
	lines := strings.Split(string(data), "\n")
	for i := range lines {
		line := strings.TrimSpace(lines[i])
		values := strings.Split(line, "=")
		switch values[0] {
		case "name":
			p.Name = values[1]
		case "path":
			p.Path = values[1]
		}
	}
	return
}
