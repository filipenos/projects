package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

var (
	configPath = fmt.Sprintf("%s/.projects-settings.json", os.Getenv("HOME"))
)

type Settings struct {
	ProjectLocation string `json:"project.location,omitempty"`
}

//Project represent then project
type Project struct {
	Name string `json:"name,omitempty"`
	Path string `json:"rootPath,omitempty"`
}

//File represet all projects managed by
//TODO adicionar map com id=index
//TODO adicionar id, pode ser o nome se único
//TODO validar path ao salvar, edit e add
//TODO flags para add -c current
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
func (f *File) Save(s Settings) error {
	file, err := os.Create(s.ProjectLocation)
	if err != nil {
		return err
	}
	defer file.Close()
	sort.Sort(f)

	enc := json.NewEncoder(file)
	//enc.SetIndent("", "  ") Only go 1.7
	return enc.Encode(f.Projects)
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
	return f, nil
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
