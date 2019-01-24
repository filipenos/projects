package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strings"
	"text/template"

	"github.com/urfave/cli"
)

var (
	ErrNameRequired = errorf("name is required")
	ErrPathRequired = errorf("path is required")
	ErrPathNoExist  = errorf("path is no exists")
)

func create(c *cli.Context) error {
	c.Args()

	var (
		p   = &Project{}
		err error
	)

	switch len(c.Args()) {
	case 0:
		p.Name, p.Path = current_pwd()
	case 1:
		p.Name = strings.TrimSpace(c.Args().Get(0))
		_, p.Path = current_pwd()
	case 2:
		p.Name = strings.TrimSpace(c.Args().Get(0))
		p.Path = strings.TrimSpace(c.Args().Get(1))
	default:
		return errorf("invalid size of arguments")
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

	if !c.Bool("no-validate") {
		if p.Path == "" {
			return ErrPathRequired
		}
		if !isExist(p.Path) {
			return ErrPathNoExist
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
	projects, err := Load(LoadSettings())
	if err != nil {
		return errorf("error on load file: %v", err)
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
	err = tmpl.Execute(os.Stdout, ctx)
	if err != nil {
		return errorf("error on execute template: %v", err)
	}
	return nil
}

func open(c *cli.Context) error {
	name := strings.TrimSpace(c.Args().First())
	if name == "" {
		name, _ = current_pwd()
	}

	projects, err := Load(LoadSettings())
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

	if c.Bool("code") {
		cmd := exec.Command("code", p.Path)
		cmd.Stdin = os.Stdin
		return cmd.Run()
	}

	log("open path '%s'", p.Path)

	if isRunning := os.Getenv("TMUX"); isRunning != "" {
		return errorf("can not open tmux inside another running: %s", isRunning)
		//TODO: não funcionou
		/*		cmd := exec.Command("tmux", "display-message", "-p", "#S")
				out, err := cmd.CombinedOutput()
				if err != nil {
					return errorf("error on get attach session name: %v", err)
				}
				currentAttached := strings.TrimSpace(string(out))
				log("detach running session %s", currentAttached)
				if currentAttached != "" {
					cmd = exec.Command("tmux", "detach-client", "-s", currentAttached)
					if err := cmd.Run(); err != nil {
						return errorf("error on detach-client: %v", err)
					}
				}
				if err := os.Setenv("TMUX", ""); err != nil {
					return errorf("error on clean tmux env: %v", err)
				}*/
	}

	var hasSession bool
	cmd := exec.Command("tmux", "has-session", "-t", name)
	out, err := cmd.CombinedOutput()
	logDebug(c.Bool("debug"), "has-session return: %v", string(out))
	if err == nil {
		hasSession = true
	}

	if !hasSession {
		cmd := exec.Command("tmux", "new", "-s", name, "-n", name, "-c", p.Path, "-d")
		//option -d run tmux with daemon
		// tmux new-session -d -s mySession -n myWindow
		// tmux send-keys -t mySession:myWindow "cd /my/directory" Enter
		// tmux send-keys -t mySession:myWindow "vim" Enter
		// tmux attach -t mySession:myWindow
		out, err = cmd.CombinedOutput()
		logDebug(c.Bool("debug"), "new-session return: %v", string(out))
		if err != nil {
			return errorf("error on new-session: %v", err)
		}
		if c.Bool("vim") {
			//args  = append(args, []string{"\\;", "new-window", "-n", "vim"}...)
			cmd := exec.Command("tmux", "new-window", "-n", "vim", "vim")
			out, err := cmd.CombinedOutput()
			logDebug(c.Bool("debug"), "new-window return: %v", string(out))
			if err != nil {
				return errorf("error on new-window: %v", err)
			}
		}
	}

	args := []string{"attach"}
	if !c.Bool("d") {
		args = append(args, "-d")
	}
	args = append(args, []string{"-t", name}...)
	cmd = exec.Command("tmux", args...)
	cmd.Stdin = os.Stdin
	out, err = cmd.CombinedOutput()
	if err != nil {
		return errorf("error on attach: %v", err)
	}
	logDebug(c.Bool("debug"), "attch return: %v", string(out))
	return nil
}

func close(c *cli.Context) error {
	projects, err := Load(LoadSettings())
	if err != nil {
		return errorf("error on load file: %v", err)
	}

	toClose := make([]string, 0, 0)

	if c.Bool("all") {
		for _, p := range projects.Projects {
			if p.Opened {
				toClose = append(toClose, p.Name)
			}
		}
		if len(toClose) == 0 {
			return errorf("no projects to close")
		}
	} else {
		name := strings.TrimSpace(c.Args().First())
		if name == "" && os.Getenv("TMUX") != "" {
			cmd := exec.Command("tmux", "display-message", "-p", "#S")
			out, err := cmd.CombinedOutput()
			if err != nil {
				return errorf("error on get attach session name: %v", err)
			}
			name = strings.TrimSpace(string(out))
			log("close running session %s", name)
			toClose = append(toClose, name)
		}

		if name == "" {
			return errorf("name of project is required")
		}
	}

	for _, name := range toClose {
		args := []string{"detach-client", "-s", name}
		if c.Bool("kill") {
			log("kill opened project %s", name)
			args = []string{"kill-session", "-t", name}
		}

		cmd := exec.Command("tmux", args...)
		if err := cmd.Run(); err != nil {
			return errorf("error on detach-client: %v", err)
		}
	}
	return nil
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
	if !c.Bool("no-validate") {
		if edited.Path == "" {
			return ErrPathRequired
		}
		if !isExist(edited.Path) {
			return ErrPathNoExist
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

	projects, err := Load(LoadSettings())
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

func export(c *cli.Context) error {
	projects, err := Load(LoadSettings())
	if err != nil {
		return err
	}

	format := strings.TrimSpace(c.String("format"))
	if len(format) == 0 {
		return fmt.Errorf("Expected format (nerdtree|vimcommand)")
	}

	out := bytes.NewBufferString("")
	sort.Sort(projects)
	for _, p := range projects.Projects {
		switch format {
		case "vimcommand":
			title := strings.ToUpper(p.Name[:1]) + p.Name[1:]
			title = strings.Replace(title, "-", "", -1)
			fmt.Fprintf(out, `
function! %s()
  cd %s
endfunction
command! %s call %s()`, title, p.Path, title, title)
		case "nerdtree":
			fmt.Fprintf(out, `%s %s
`, p.Name, p.Path)
		}
	}

	if c.Bool("override") {
		filename := os.Getenv("HOME")
		if format == "nerdtree" {
			filename += "/.NERDTreeBookmarks"
		} else {
			filename += "/.vimrc.projects"
		}
		if err = ioutil.WriteFile(filename, out.Bytes(), 0666); err != nil {
			return err
		}
	} else {
		os.Stdout.Write(out.Bytes())
	}

	return nil
}
