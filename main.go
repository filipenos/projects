package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
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
			Usage:     "add new project",
			ArgsUsage: "name path",
			Action:    add,
		},
		{
			Name:   "current",
			Usage:  "add the current path as project",
			Action: addCurrent,
		},
		{
			Name:   "list",
			Usage:  "list the managed projects",
			Action: list,
		},
		{
			Name:   "open",
			Usage:  "open the path of project",
			Action: open,
		},
	}
	app.Name = "Projects"
	app.Description = "Manage projects"
	app.HideVersion = true
	app.ExitErrHandler = func(c *cli.Context, err error) {
		log.Println(err)
	}
	app.Run(os.Args)
}

func add(c *cli.Context) error {
	name := strings.TrimSpace(c.Args().Get(0))
	if name == "" {
		return fmt.Errorf("name is required")
	}
	path := strings.TrimSpace(c.Args().Get(1))
	if path == "" {
		return fmt.Errorf("path is required")
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

func list(c *cli.Context) error {
	projects, err := Load(filepath)
	if err != nil {
		return err
	}
	for _, p := range projects.Projects {
		fmt.Printf("  Project= '%s' Path='%s'\n", p.Name, p.Path)
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

	cmd := exec.Command("xdg-open", path)
	// cmd.Dir = path
	out, err := cmd.Output()
	log.Println(string(out))
	return err
}

//Project represent then project
type Project struct {
	Name string
	Path string
}

//Projects represet all projects managed by
type Projects struct {
	Path     string
	Projects []Project
}

func (p Projects) Len() int           { return len(p.Projects) }
func (p Projects) Swap(i, j int)      { p.Projects[i], p.Projects[j] = p.Projects[j], p.Projects[i] }
func (p Projects) Less(i, j int) bool { return p.Projects[i].Name < p.Projects[j].Name }

//Add new project to manage
func (p *Projects) Add(name, path string) {
	p.Projects = append(p.Projects, Project{Name: name, Path: path})
}

//Save save the current projects on conf file
func (p *Projects) Save() error {
	file, err := os.Create(p.Path)
	if err != nil {
		return err
	}
	defer file.Close()
	sort.Sort(p)

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	return enc.Encode(p.Projects)
}

//Load retrieve projects from config file
func Load(path string) (*Projects, error) {
	p := &Projects{Path: path}

	file, err := os.Open(p.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return p, nil
		}
		return nil, err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&p.Projects); err != nil {
		return nil, err
	}
	return p, nil
}
