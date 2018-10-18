package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var (
	configPath = fmt.Sprintf("%s/.projects-settings.json", os.Getenv("HOME"))
)

//Settings save configuration
type Settings struct {
	ProjectLocation string `json:"project.location,omitempty"`
}

//Project represent then project
type Project struct {
	Name string `json:"name,omitempty"`
	Path string `json:"rootPath,omitempty"`

	Opened   bool `json:"-"`
	Attached bool `json:"-"`
}

//File represet all projects managed by
type File struct {
	Path     string
	Projects []Project
}

func (f File) Len() int           { return len(f.Projects) }
func (f File) Swap(i, j int)      { f.Projects[i], f.Projects[j] = f.Projects[j], f.Projects[i] }
func (f File) Less(i, j int) bool { return f.Projects[i].Name < f.Projects[j].Name }

//Add new project to manage
func (f *File) Add(name, path string) {
	f.Projects = append(f.Projects, Project{Name: name, Path: path})
}

func (f *File) AddProject(p Project) {
	f.Projects = append(f.Projects, p)
}

func (f *File) Get(name string) (*Project, int) {
	name = strings.TrimSpace(name)
	for i := range f.Projects {
		if f.Projects[i].Name == name {
			return &f.Projects[i], i
		}
	}
	return nil, -1
}

//Save save the current projects on conf file
func Save(s Settings, f *File) error {
	b, err := json.MarshalIndent(f.Projects, " ", "  ")
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(f.Path, b, 0644); err != nil {
		return err
	}
	return nil
}

//Load retrieve projects from config file
func Load(s Settings) (*File, error) {
	f := &File{Path: s.ProjectLocation}

	file, err := os.Open(f.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return f, nil
		}
		return nil, err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&f.Projects); err != nil {
		return nil, err
	}

	sessions, err := getSessions()
	if err != nil {
		return nil, err
	}
	for i, p := range f.Projects {
		attached, ok := sessions[p.Name]
		if ok {
			f.Projects[i].Opened = true
			f.Projects[i].Attached = attached
		}
	}

	return f, nil
}

func getSessions() (map[string]bool, error) {
	m := make(map[string]bool, 0)

	cmd := exec.Command("tmux", "list-sessions")
	out, err := cmd.CombinedOutput()
	if err != nil && !strings.Contains(string(out), "no server running") {
		return m, err
	}
	for _, l := range strings.Split(string(out), "\n") {
		l = strings.TrimSpace(l)
		if len(l) == 0 {
			continue
		}
		p := strings.Split(l, ":")
		if len(p) == 0 {
			log("%v", p)
			continue
		}
		m[p[0]] = strings.Contains(strings.Join(p[1:], ""), "attached")
	}

	return m, nil
}

func LoadSettings() Settings {
	settings := Settings{
		ProjectLocation: fmt.Sprintf("%s/.projects.json", os.Getenv("HOME")),
	}
	file, err := os.Open(configPath)
	if err != nil {
		return settings
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&settings); err != nil {
		return settings
	}
	return settings
}
