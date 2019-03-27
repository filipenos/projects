package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var (
	configPath      = fmt.Sprintf("%s/.projects.config.yaml", os.Getenv("HOME"))
	projectsPath    = fmt.Sprintf("%s/.projects.json", os.Getenv("HOME"))
	defaultSettings = Settings{ProjectLocation: projectsPath}
)

//Settings save configuration
type Settings struct {
	ProjectLocation string
}

type File struct {
	Groups   []Group   `json:"groups,omitempty"`
	Projects []Project `json:"projects,omitempty"`
}

type Group struct {
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Icon        string    `json:"icon,omitempty"`
	Projects    []Project `json:"projects,omitempty"`
	Groups      []Group   `json:"groups,omitempty"`
}

//Project represent then project
type Project struct {
	Name        string `json:"name,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Description string `json:"description,omitempty"`
	Path        string `json:"path,omitempty"`

	Group     string `json:"-"`
	Opened    bool   `json:"-"`
	Attached  bool   `json:"-"`
	ValidPath bool   `json:"-"`
}

type Projects []Project
type Groups map[string]Projects

func (projects Projects) Len() int           { return len(projects) }
func (projects Projects) Swap(i, j int)      { projects[i], projects[j] = projects[j], projects[i] }
func (projects Projects) Less(i, j int) bool { return projects[i].Name < projects[j].Name }

func (projects Projects) Get(name string) (*Project, int) {
	name = strings.TrimSpace(name)
	for i := range projects {
		if projects[i].Name == name {
			return &projects[i], i
		}
	}
	return nil, -1
}

func (projects Projects) GetByPath(path string) (*Project, int) {
	path = strings.TrimSpace(path)
	for i := range projects {
		if projects[i].Path == path {
			return &projects[i], i
		}
	}
	return nil, -1
}

//Save save the current projects on conf file
func (projects Projects) Save(s Settings) error {
	groups := make(Groups, 0)
	for i := range projects {
		p := projects[i]
		projs := groups[p.Group]
		projs = append(projs, p)
		groups[p.Group] = projs
	}

	file := File{}
	for k, v := range groups {
		if k == "" {
			file.Projects = v
		} else {
			file.Groups = append(file.Groups, Group{Name: k, Projects: v})
		}
	}

	b, err := json.MarshalIndent(file, " ", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(s.ProjectLocation, b, 0644)
}

//Load retrieve projects from config file
func Load(s Settings) (Projects, error) {
	file, err := os.Open(s.ProjectLocation)
	if err != nil {
		if os.IsNotExist(err) {
			return Projects{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var f File
	if err := json.NewDecoder(file).Decode(&f); err != nil {
		return nil, err
	}

	projects := make(Projects, 0)
	for _, p := range f.Projects {
		projects = append(projects, p)
	}
	for _, g := range f.Groups {
		for i := range g.Projects {
			p := g.Projects[i]
			p.Group = g.Name
			projects = append(projects, p)
		}
	}

	sessions, err := getSessions()
	if err != nil {
		return nil, errorf("error on get tmux sessions: %v", err)
	}
	for i, p := range projects {
		attached, ok := sessions[p.Name]
		if ok {
			projects[i].Opened = true
			projects[i].Attached = attached
		}
		projects[i].ValidPath = isExist(p.Path)
	}

	return projects, nil
}

//LoadSettings load configuration used on projects
func LoadSettings() Settings {
	var (
		settings Settings
		v        = viper.New()
	)

	v.SetConfigFile(configPath)
	if _, err := os.Stat(v.ConfigFileUsed()); os.IsNotExist(err) {
		v.SetDefault("projectlocation", defaultSettings.ProjectLocation)
	} else {
		if err := v.ReadInConfig(); err != nil {
			settings = defaultSettings
		}
		if err := v.Unmarshal(&settings); err != nil {
			settings = defaultSettings
		}
	}
	if err := v.WriteConfig(); err != nil {
		settings = defaultSettings
	}
	return settings
}
