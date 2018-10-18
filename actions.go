package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/urfave/cli"
)

var (
	ErrNameRequired = errorf("name is required")
)

func create(c *cli.Context) error {
	c.Args()

	var (
		p   = &Project{}
		err error
	)

	if c.Bool("current") {
		pwd := os.Getenv("PWD")
		paths := strings.Split(pwd, "/")

		p.Name = strings.TrimSpace(paths[len(paths)-1])
		p.Path = strings.TrimSpace(pwd)
	} else {
		p.Name = strings.TrimSpace(c.Args().Get(0))
		p.Path = strings.TrimSpace(c.Args().Get(1))
	}

	if c.Bool("editor") {
		p, err = editProject(p)
		if err != nil {
			return err
		}
	}

	if p.Name == "" {
		return ErrNameRequired
	}

	if c.Bool("validate-path") {
		if p.Path == "" {
			return errorf("path is required")
		}
		if !isExist(p.Path) {
			return errorf("path is no exists")
		}
	}

	s := LoadSettings()
	projects, err := Load(s)
	if err != nil {
		return err
	}

	if p, _ := projects.Get(p.Name); p != nil {
		return errorf("project '%s' already add to projects", p.Name)
	}

	projects.AddProject(*p)
	if err := Save(s, projects); err != nil {
		return err
	}
	log("Add project: '%s' path: '%s'", p.Name, p.Path)
	return nil
}

func delete(c *cli.Context) error {
	name := strings.TrimSpace(c.Args().First())
	if name == "" {
		return ErrNameRequired
	}

	s := LoadSettings()
	projects, err := Load(s)
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
		return errorf("Project '%s' not found", name)
	}

	log("Project '%s' removed successfully!", name)
	projects.Projects = aux
	return Save(s, projects)
}

func list(c *cli.Context) error {
	s := LoadSettings()
	projects, err := Load(s)
	if err != nil {
		return err
	}

	t := `{{range .Projects.Projects}}{{.Name}}{{if .Opened}} (opened){{end}}{{if .Attached}} (attached){{end}}{{if $.Full}}
  Path: {{.Path}}{{end}}
{{else}}No projects yeat!
{{end}}`
	tmpl := template.Must(template.New("editor").Parse(t))
	ctx := map[string]interface{}{
		"Projects": projects,
		"Full":     c.Bool("full"),
	}
	return tmpl.Execute(os.Stdout, ctx)
}

func open(c *cli.Context) error {
	name := strings.TrimSpace(c.Args().First())
	if name == "" {
		return ErrNameRequired
	}

	s := LoadSettings()
	projects, err := Load(s)
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
		return errorf("Project '%s' not found", name)
	}

	log("open path '%s'", path)

	if isRunning := os.Getenv("TMUX"); isRunning != "" {
		return errorf("can not open tmux inside another running: %s", isRunning)
	}

	var hasSession bool
	cmd := exec.Command("tmux", "has-session", "-t", name)
	out, err := cmd.CombinedOutput()
	logDebug(c.Bool("debug"), "has-session return: %v", string(out))
	if err == nil {
		hasSession = true
	}

	if !hasSession {
		cmd := exec.Command("tmux", "new", "-s", name, "-n", name, "-c", path, "-d")
		//option -d run tmux with daemon
		// tmux new-session -d -s mySession -n myWindow
		// tmux send-keys -t mySession:myWindow "cd /my/directory" Enter
		// tmux send-keys -t mySession:myWindow "vim" Enter
		// tmux attach -t mySession:myWindow
		out, err = cmd.CombinedOutput()
		logDebug(c.Bool("debug"), "new-session return: %v", string(out))
		if err != nil {
			return err
		}
		if c.Bool("vim") {
			//args  = append(args, []string{"\\;", "new-window", "-n", "vim"}...)
			cmd := exec.Command("tmux", "new-window", "-n", "vim", "vim")
			out, err := cmd.CombinedOutput()
			logDebug(c.Bool("debug"), "new-window return: %v", string(out))
			if err != nil {
				return err
			}
		}
	}
	cmd = exec.Command("tmux", "attach", "-t", name)
	cmd.Stdin = os.Stdin
	out, err = cmd.CombinedOutput()
	logDebug(c.Bool("debug"), "attch return: %v", string(out))
	return err
}

func edit(c *cli.Context) error {
	name := strings.TrimSpace(c.Args().First())
	if name == "" {
		return ErrNameRequired
	}

	s := LoadSettings()
	projects, err := Load(s)
	if err != nil {
		return err
	}

	p, _ := projects.Get(name)
	if p == nil {
		return errorf("project '%s' not found", name)
	}
	if p.Path == "" {
		return errorf("project '%s' dont have path", p.Name)
	}
	if !isExist(p.Path) {
		return errorf("path '%s' of project '%s' not exists", p.Path, p.Name)
	}

	log("opening %s to edit", p.Name)

	cmd := exec.Command("code", p.Path)
	// cmd := exec.Command("vim", "--cmd", fmt.Sprintf(`"cd %s"`, p.Path))
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func update(c *cli.Context) error {
	name := strings.TrimSpace(c.Args().First())
	if name == "" {
		return ErrNameRequired
	}

	s := LoadSettings()
	projects, err := Load(s)
	if err != nil {
		return err
	}

	_, index := projects.Get(name)
	if index == -1 {
		return errorf("project '%s' not found", name)
	}
	p := &projects.Projects[index]
	if p == nil {
		return errorf("project '%s' not found", name)
	}

	edited, err := editProject(p)
	if err != nil {
		return err
	}
	if edited.Name == "" {
		return ErrNameRequired
	}
	if c.Bool("validate-path") {
		if edited.Path == "" {
			return errorf("path is required")
		}
		if !isExist(edited.Path) {
			return errorf("path is no exists")
		}
	}

	projects.Projects[index] = *edited
	return Save(s, projects)
}

func editProject(p *Project) (*Project, error) {
	tmp, err := NewTempFile()
	if err != nil {
		return nil, err
	}
	defer tmp.Remove()

	d := `name={{.Name}}
path={{.Path}}`

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
	return parseContent(content), nil
}

func parseContent(data []byte) *Project {
	lines := strings.Split(string(data), "\n")
	p := &Project{}
	for i := range lines {
		line := strings.TrimSpace(lines[i])
		values := strings.Split(line, "=")
		if len(values) != 2 {
			continue
		}
		switch values[0] {
		case "name":
			p.Name = strings.TrimSpace(values[1])
		case "path":
			p.Path = strings.TrimSpace(values[1])
		}
	}
	return p
}

func path(c *cli.Context) error {
	name := strings.TrimSpace(c.Args().First())
	if name == "" {
		return ErrNameRequired
	}

	s := LoadSettings()
	projects, err := Load(s)
	if err != nil {
		return err
	}

	p, _ := projects.Get(name)
	if p == nil {
		return errorf("project '%s' not found", name)
	}
	if p.Path == "" {
		return errorf("project '%s' dont have path", p.Name)
	}
	if !isExist(p.Path) {
		return errorf("path '%s' of project '%s' not exists", p.Path, p.Name)
	}

	fmt.Println(p.Path)

	return nil
}

func vimCommand(c *cli.Context) error {
	s := LoadSettings()
	projects, err := Load(s)
	if err != nil {
		return err
	}

	for _, p := range projects.Projects {
		title := strings.ToUpper(p.Name[:1]) + p.Name[1:]
		title = strings.Replace(title, "-", "", -1)
		fmt.Printf(`
function! %s()
  cd %s
endfunction
command! %s call %s()`, title, p.Path, title, title)
	}

	return nil
}
