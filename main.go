package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli"
)

var (
	filepath = fmt.Sprintf("%s/.projects.json", os.Getenv("HOME"))
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:      "add",
			Aliases:   []string{"a"},
			Usage:     "add new project",
			UsageText: "project add <name> <path>",
			ArgsUsage: "name path",
			Action:    add,
		},
		{
			Name:   "current",
			Usage:  "add the current path as project",
			Action: addCurrent,
		},
		{
			Name:      "remove",
			Usage:     "remove project",
			UsageText: "project remove <name>",
			ArgsUsage: "name",
			Action:    remove,
		},
		{
			Name:   "list",
			Usage:  "list projects",
			Action: list,
		},
		{
			Name:      "open",
			Usage:     "open project",
			UsageText: "project open <name>",
			ArgsUsage: "name",
			Action:    open,
		},
	}
	app.Name = "Projects"
	app.Usage = "Simple manager for your projects"
	app.Description = "Manage local projects"
	app.HideVersion = true
	app.ExitErrHandler = func(c *cli.Context, err error) {
		if err != nil {
			log("%v", err)
		}
		return
	}
	app.Run(os.Args)
}

func add(c *cli.Context) error {
	c.Args()
	name := strings.TrimSpace(c.Args().Get(0))
	if name == "" {
		return fmt.Errorf("name is required")
	}
	path := strings.TrimSpace(c.Args().Get(1))
	if path == "" {
		return fmt.Errorf("path is required")
	}
	if !checkPath(path) {
		return fmt.Errorf("path is no exists")
	}

	projects, err := Load(filepath)
	if err != nil {
		return err
	}
	projects.Add(name, path)
	return projects.Save()
}

func addCurrent(c *cli.Context) error {
	projects, err := Load(filepath)
	if err != nil {
		return err
	}
	pwd := os.Getenv("PWD")
	paths := strings.Split(pwd, "/")
	projects.Add(paths[len(paths)-1], pwd)
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

func log(msg string, args ...interface{}) {
	fmt.Printf("[projects] %s\n", fmt.Sprintf(msg, args...))
}
