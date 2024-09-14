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

type ProjectType string

const (
	ProjectTypeLocal  ProjectType = "local"
	ProjectTypeSSH    ProjectType = "ssh"
	ProjectTypeWSL    ProjectType = "wsl"
	ProjectTypeTunnel ProjectType = "tunnel"
)

func ParseProjectType(typ string) ProjectType {
	if i := strings.Index(typ, "+"); i > -1 {
		switch typ[:i] {
		case "wsl":
			return ProjectTypeWSL
		case "ssh":
			return ProjectTypeSSH
		case "ssh-remote":
			return ProjectTypeSSH
		case "tunnel":
			return ProjectTypeTunnel
		default:
			return ProjectType(typ[:i])
		}
	}
	return ProjectType(typ)
}

func parseURL(input string) (string, string, string) {
	var scheme, domain, path string

	// Separa o schema e o resto
	if strings.Contains(input, "://") {
		parts := strings.SplitN(input, "://", 2)
		scheme = parts[0]
		input = parts[1]
	}

	// Separa o domain e o path
	if strings.Contains(input, "/") {
		parts := strings.SplitN(input, "/", 2)
		domain = parts[0]
		path = "/" + parts[1]
	} else {
		domain = input
		path = ""
	}

	return scheme, domain, path
}

var (
	ErrNameRequired = fmt.Errorf("name is required")
	ErrPathRequired = fmt.Errorf("path is required")
	ErrPathNoExist  = fmt.Errorf("path is no exists")
)

// Project represent then project
type Project struct {
	Name     string   `json:"name,omitempty"`
	Alias    string   `json:"alias,omitempty"`
	RootPath string   `json:"rootPath,omitempty"`
	Group    string   `json:"group,omitempty"`
	Enabled  bool     `json:"enabled,omitempty"`
	SCM      string   `json:"scm,omitempty"`
	Tags     []string `json:"tags,omitempty"`

	Scheme string `json:"-"`
	Domain string `json:"-"`
	Path   string `json:"-"`

	ProjectType ProjectType `json:"-"`
	Opened      bool        `json:"-"`
	Attached    bool        `json:"-"`
	ValidPath   bool        `json:"-"`
	IsWorkspace bool        `json:"-"`
}

func (p *Project) Validate() error {
	if p.Name == "" {
		return ErrNameRequired
	}
	if p.RootPath == "" {
		return fmt.Errorf("project '%s' dont have path", p.Name)
	}
	switch p.ProjectType {
	case ProjectTypeLocal:
		if !path.Exist(p.RootPath) {
			return fmt.Errorf("path '%s' of project '%s' not exists", p.RootPath, p.Name)
		}
	case ProjectTypeSSH, ProjectTypeWSL, ProjectTypeTunnel:
	default:
		return fmt.Errorf("invalid project type: %s", p.ProjectType)
	}
	return nil
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
		if projects[i].RootPath == path {
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
		if strings.Contains(p.RootPath, "~") {
			projects[i].RootPath = strings.Replace(p.RootPath, "~", os.Getenv("HOME"), 1)
		}

		projects[i].Scheme, projects[i].Domain, projects[i].Path = parseURL(p.RootPath)
		if projects[i].Scheme != "" {
			projects[i].ProjectType = ParseProjectType(projects[i].Domain)
			projects[i].ValidPath = true
		} else {
			projects[i].ProjectType = ProjectTypeLocal
			projects[i].ValidPath = path.Exist(p.RootPath)
		}

		if strings.HasSuffix(p.RootPath, ".code-workspace") {
			projects[i].IsWorkspace = true
		}

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
			p.RootPath = v
		case "group":
			p.Group = v
		case "enabled":
			p.Enabled = v == "true"
		}
	}
	return p
}
