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

//TODO (filipenos) o project manager possui a variavel $home, podemos utilizala para gravar o path, ai nao importa qual distro esteja usando

var (
	ErrNameRequired = errorf("name is required")
	ErrPathRequired = errorf("path is required")
	ErrPathNoExist  = errorf("path is no exists")
)

func create(c *cli.Context) error {
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

	if p.Path != "" && isExist(fmt.Sprintf("%s/.git", p.Path)) {
		cmd := exec.Command("git", "-C", p.Path, "remote", "get-url", "origin")
		out, _ := cmd.CombinedOutput()
		if scm := strings.TrimSpace(string(out)); scm != "" {
			p.SCM = scm
		}

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

	projects = append(projects, *p)
	if err := projects.Save(s); err != nil {
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
	aux := make([]Project, 0, len(projects))
	for i := range projects {
		if projects[i].Name == name && !excluded {
			excluded = true
		} else {
			aux = append(aux, projects[i])
		}
	}

	if !excluded {
		return errorf("Project '%s' not found", name)
	}

	log("Project '%s' removed successfully!", name)
	projects = aux
	return projects.Save(s)
}

func list(c *cli.Context) error {
	projects, err := Load(LoadSettings())
	if err != nil {
		return errorf("error on load file: %v", err)
	}
	sort.Sort(projects)

	t := `{{range .Projects}}{{.Name}}{{if $.ExtraInfo}}{{if .Opened}} (opened){{end}}{{if .Attached}} (attached){{end}}{{if not .ValidPath}} (invalid-path){{end}}{{end}}{{if $.Path}}
  Path: {{.Path}}{{end}}
{{else}}No projects yeat!
{{end}}`
	tmpl := template.Must(template.New("editor").Parse(t))
	ctx := map[string]interface{}{
		"Projects":  projects,
		"Path":      c.Bool("path"),
		"ExtraInfo": !c.Bool("simple"),
	}
	err = tmpl.Execute(os.Stdout, ctx)
	if err != nil {
		return errorf("error on execute template: %v", err)
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
	p := &projects[index]
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

	projects[index] = *edited
	return projects.Save(s)
}

func editProject(p *Project) (*Project, error) {
	tmp, err := NewTempFile()
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
		return fmt.Errorf("Expected format (nerdtree|vimcommand|vim-project|aliases)")
	}

	out := bytes.NewBufferString("")
	sort.Sort(projects)
	for _, p := range projects {
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
		case "vim-project":
			fmt.Fprintf(out, `Project '%s' , '%s'
`, p.Path, p.Name)
		case "alias", "aliases":
			fmt.Fprintf(out, `alias %s="cd %s"
`, p.Name, p.Path)
		}
	}

	if c.Bool("override") {
		filename := os.Getenv("HOME")
		switch format {
		case "vimcommand":
			filename += "/.vimrc.projects"
		case "nerdtree":
			filename += "/.NERDTreeBookmarks"
		case "vim-project":
			filename += "/.vim-project.projects"
		}
		if err = ioutil.WriteFile(filename, out.Bytes(), 0666); err != nil {
			return err
		}
	} else {
		os.Stdout.Write(out.Bytes())
	}

	return nil
}

func open(c *cli.Context) error {
	var (
		name        string
		path        string
		withoutName bool
	)

	name = strings.TrimSpace(c.Args().First())
	if name == "" {
		withoutName = true
		name, path = current_pwd()
	}

	projects, err := Load(LoadSettings())
	if err != nil {
		return err
	}

	var p *Project
	if c.Bool("r") && withoutName {
		log("search any project on path %s", path)
		paths := strings.Split(path, "/")
		for i := len(paths) - 1; i >= 0; i-- {
			namePath := strings.TrimSpace(paths[i])
			if namePath == "" {
				continue
			}
			p, _ = projects.Get(namePath)
			if p != nil {
				log("found project %s on path", namePath)
				name = namePath
				break
			}
		}
	} else {
		p, _ = projects.Get(name)
	}

	if p == nil {
		p, _ = projects.GetByPath(path)
		if p == nil {
			return errorf("project '%s' not found", name)
		}
	}
	if p.Path == "" {
		return errorf("project '%s' dont have path", p.Name)
	}
	if !isExist(p.Path) {
		return errorf("path '%s' of project '%s' not exists", p.Path, p.Name)
	}

	log("open path '%s'", p.Path)

	//TODO (filipenos) validar se esta dentro de um tmux, acontece algo estranhoooo

	var hasSession bool
	sessions, err := getSessions()
	if err != nil {
		return err
	}
	_, hasSession = sessions[p.Name]

	if !hasSession {
		cmd := exec.Command("tmux", "new", "-s", p.Name, "-n", p.Name, "-c", p.Path, "-d")
		//option -d run tmux with daemon
		// tmux new-session -d -s mySession -n myWindow
		// tmux send-keys -t mySession:myWindow "cd /my/directory" Enter
		// tmux send-keys -t mySession:myWindow "vim" Enter
		// tmux attach -t mySession:myWindow
		out, err := cmd.CombinedOutput()
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
	args = append(args, []string{"-t", p.Name}...)
	cmd := exec.Command("tmux", args...)
	cmd.Stdin = os.Stdin
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errorf("error on attach: %v", err)
	}
	logDebug(c.Bool("debug"), "attch return: %v", string(out))
	return nil
}

func code(c *cli.Context) error {
	var (
		name        string
		path        string
		withoutName bool
		settings    = LoadSettings()
	)

	name = strings.TrimSpace(c.Args().First())
	if name == "" {
		name, path = current_pwd()
	}

	projects, err := Load(settings)
	if err != nil {
		return err
	}

	var p *Project
	if c.Bool("r") && withoutName {
		log("search any project on path %s", path)
		paths := strings.Split(path, "/")
		for i := len(paths) - 1; i >= 0; i-- {
			namePath := strings.TrimSpace(paths[i])
			if namePath == "" {
				continue
			}
			p, _ = projects.Get(namePath)
			if p != nil {
				log("found project %s on path", namePath)
				name = namePath
				break
			}
		}
	} else {
		p, _ = projects.Get(name)
	}

	if p == nil {
		p, _ = projects.GetByPath(path)
		if p == nil {
			return errorf("project '%s' not found", name)
		}
	}
	if p.Path == "" {
		return errorf("project '%s' dont have path", p.Name)
	}
	if !isExist(p.Path) {
		return errorf("path '%s' of project '%s' not exists", p.Path, p.Name)
	}

	editor := "code"
	args := make([]string, 0)

	if edit := c.String("e"); edit != "" {
		editor = edit
	} else {
		pos := "--new-window"
		if c.Bool("r") {
			pos = "--reuse-window"
		}
		args = append(args, pos)
	}
	args = append(args, p.Path)

	log("open path '%s' on ", p.Path, editor)

	cmd := exec.Command(editor, args...)
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func close(c *cli.Context) error {
	projects, err := Load(LoadSettings())
	if err != nil {
		return errorf("error on load file: %v", err)
	}

	toClose := make([]string, 0, 0)

	if c.Bool("all") {
		for _, p := range projects {
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

func scm(c *cli.Context) error {
	var (
		name string
		path string
	)

	name = strings.TrimSpace(c.Args().Get(0))
	if name == "" {
		name, path = current_pwd()
	}
	log("using '%s' name", name)

	s := LoadSettings()
	projects, err := Load(s)
	if err != nil {
		return err
	}

	p, index := projects.Get(name)
	if p == nil {
		p, _ = projects.GetByPath(path)
		if p == nil {
			return errorf("project '%s' not found", name)
		}
	}
	if p.Path == "" {
		return errorf("project '%s' dont have path", p.Name)
	}
	if !isExist(p.Path) {
		return errorf("path '%s' of project '%s' not exists", p.Path, p.Name)
	}

	log("using path: %s", p.Path)

	if c.Bool("set") {
		url := c.Args().Get(1)
		if url == "" {
			cmd := exec.Command("git", "-C", path, "remote", "get-url", "origin")
			out, err := cmd.CombinedOutput()
			if err != nil {
				return err
			}
			url = strings.TrimSpace(string(out))
		}
		log("setting scm url %s", url)
		p.SCM = url
		projects[index] = *p
		return projects.Save(s)
	}

	return nil
}
