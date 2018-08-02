package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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
