package project

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/filipenos/projects/pkg/config"
	"github.com/filipenos/projects/pkg/file"
	"github.com/filipenos/projects/pkg/path"
)

var (
	ErrNameRequired = fmt.Errorf("name is required")
	ErrPathRequired = fmt.Errorf("path is required")
	ErrPathNoExist  = fmt.Errorf("path is no exists")
)

// Project represent then project
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

// Save save the current projects on conf file
func (projects Projects) Save(s config.Config) error {
	b, err := json.MarshalIndent(projects, " ", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.ProjectLocation, b, 0644)
}

// Load retrieve projects from config file
func Load(s config.Config) (Projects, error) {
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

	for i, p := range projects {
		projects[i].ValidPath = path.Exist(p.Path)
	}

	return projects, nil
}

func EditProject(p *Project) (*Project, error) {
	tmp, err := file.NewTempFile()
	if err != nil {
		return nil, err
	}
	defer tmp.Remove()

	d := `name={{.Name}}
path={{.Path}}
group={{.Group}}
enabled={{.Enabled}}`

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
	return ParseContent(content), nil
}

func ParseContent(data []byte) *Project {
	lines := strings.Split(string(data), "\n")
	p := &Project{}
	for i := range lines {
		line := strings.TrimSpace(lines[i])
		values := strings.Split(line, "=")
		if len(values) != 2 {
			continue
		}
		v := strings.TrimSpace(values[1])
		switch values[0] {
		case "name":
			p.Name = v
		case "path":
			p.Path = v
		case "group":
			p.Group = v
		case "enabled":
			p.Enabled = v == "true"
		}
	}
	return p
}
