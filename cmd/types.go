package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	projectsPath    = fmt.Sprintf("%s/.projects.json", os.Getenv("HOME"))
	defaultSettings = Settings{ProjectLocation: projectsPath, Editor: "code"}
)

//Settings save configuration
type Settings struct {
	ProjectLocation string
	Editor          string
}

//Project represent then project
type Project struct {
	Name    string   `json:"name,omitempty"`
	Alias   string   `json:"alias,omitempty"`
	Path    string   `json:"rootPath,omitempty"`
	Group   string   `json:"group,omitempty"`
	Enabled bool     `json:"enabled,omitempty"`
	SCM     string   `json:"scm,omitempty"`
	Tags    []string `json:"tags,omitempty"`

	Opened    bool `json:"-"`
	Attached  bool `json:"-"`
	ValidPath bool `json:"-"`
}

type Projects []Project

func (projects Projects) Len() int           { return len(projects) }
func (projects Projects) Swap(i, j int)      { projects[i], projects[j] = projects[j], projects[i] }
func (projects Projects) Less(i, j int) bool { return projects[i].Name < projects[j].Name }

func (projects Projects) Get(name string) (*Project, int) {
	name = strings.TrimSpace(name)
	for i := range projects {
		if projects[i].Name == name || projects[i].Alias == name {
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

func (projects Projects) Find(name, path string) (*Project, int) {
	var (
		project *Project
		pos     int
	)
	if name != "" {
		project, pos = projects.Get(name)
		if project != nil {
			return project, pos
		}
	}
	if path != "" {
		project, pos = projects.GetByPath(path)
		if project != nil {
			return project, pos
		}

		paths := strings.Split(path, "/")
		for i := len(paths) - 1; i >= 0; i-- {
			namePath := strings.TrimSpace(paths[i])
			if namePath == "" {
				continue
			}
			project, pos = projects.Get(namePath)
			if project != nil {
				return project, pos
			}
		}
	}
	return nil, -1
}

//Save save the current projects on conf file
func (projects Projects) Save(s Settings) error {
	b, err := json.MarshalIndent(projects, " ", "  ")
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

	var projects Projects
	if err := json.NewDecoder(file).Decode(&projects); err != nil {
		return nil, err
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
	var settings Settings
	if err := viper.Unmarshal(&settings); err != nil {
		settings = defaultSettings
	}
	if settings == (Settings{}) {
		settings = defaultSettings
	}
	return settings
}

func SafeBoolFlag(cmd *cobra.Command, flagName string) bool {
	v, _ := cmd.Flags().GetBool(flagName)
	return v
}

func SafeStringFlag(cmd *cobra.Command, flagName string) string {
	v, _ := cmd.Flags().GetString(flagName)
	return v
}
