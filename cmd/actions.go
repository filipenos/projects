package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

//TODO (filipenos) o project manager possui a variavel $home, podemos utilizala para gravar o path, ai nao importa qual distro esteja usando

//Errors returned
var (
	ErrNameRequired = errorf("name is required")
	ErrPathRequired = errorf("path is required")
	ErrPathNoExist  = errorf("path is no exists")
)

func create(cmdParam *cobra.Command, params []string) error {
	var (
		p   = &Project{}
		err error
	)

	switch len(params) {
	case 0:
		p.Name, p.Path = currentPwd()
	case 1:
		p.Name = strings.TrimSpace(params[0])
		_, p.Path = currentPwd()
	case 2:
		p.Name = strings.TrimSpace(params[0])
		p.Path = strings.TrimSpace(params[1])
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

	if SafeBoolFlag(cmdParam, "editor") {
		p, err = editProject(p)
		if err != nil {
			return err
		}
	}

	if p.Name == "" {
		return ErrNameRequired
	}

	if !SafeBoolFlag(cmdParam, "no-validate") {
		if p.Path == "" {
			return ErrPathRequired
		}
		if !isExist(p.Path) {
			return ErrPathNoExist
		}
	}

	projects, err := Load(LoadSettings())
	if err != nil {
		return err
	}

	if p, _ := projects.Get(p.Name); p != nil {
		return errorf("project '%s' already add to projects", p.Name)
	}

	projects = append(projects, *p)
	if err := projects.Save(LoadSettings()); err != nil {
		return err
	}
	log("Add project: '%s' path: '%s'", p.Name, p.Path)
	return nil
}

func delete(cmdParam *cobra.Command, params []string) error {
	name, _ := safeName(params...)
	if name == "" {
		return ErrNameRequired
	}

	projects, err := Load(LoadSettings())
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
	return projects.Save(LoadSettings())
}

func list(cmdParam *cobra.Command, params []string) error {
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
		"Path":      SafeBoolFlag(cmdParam, "path"),
		"ExtraInfo": !SafeBoolFlag(cmdParam, "simple"),
	}
	err = tmpl.Execute(os.Stdout, ctx)
	if err != nil {
		return errorf("error on execute template: %v", err)
	}
	return nil
}

func update(cmdParam *cobra.Command, params []string) error {
	name, _ := safeName(params...)
	if name == "" {
		return ErrNameRequired
	}

	projects, err := Load(LoadSettings())
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
	if !SafeBoolFlag(cmdParam, "no-validate") {
		if edited.Path == "" {
			return ErrPathRequired
		}
		if !isExist(edited.Path) {
			return ErrPathNoExist
		}
	}

	projects[index] = *edited
	return projects.Save(LoadSettings())
}

func path(cmdParam *cobra.Command, params []string) error {
	name, _ := safeName(params...)
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

func export(cmdParam *cobra.Command, params []string) error {
	projects, err := Load(LoadSettings())
	if err != nil {
		return err
	}

	format := strings.TrimSpace(SafeStringFlag(cmdParam, "format"))
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

	if SafeBoolFlag(cmdParam, "override") {
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

func open(cmdParam *cobra.Command, params []string) error {
	projects, err := Load(LoadSettings())
	if err != nil {
		return err
	}

	p, _ := projects.Find(safeName(params...))
	if p == nil {
		return errorf("project not found")
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
		out, err := cmd.CombinedOutput()
		logDebug(SafeBoolFlag(cmdParam, "debug"), "new-session return: %v", string(out))
		if err != nil {
			return errorf("error on new-session: %v", err)
		}
		if SafeBoolFlag(cmdParam, "vim") {
			//args  = append(args, []string{"\\;", "new-window", "-n", "vim"}...)
			cmd := exec.Command("tmux", "new-window", "-n", "vim", "vim")
			out, err := cmd.CombinedOutput()
			logDebug(SafeBoolFlag(cmdParam, "debug"), "new-window return: %v", string(out))
			if err != nil {
				return errorf("error on new-window: %v", err)
			}
		}
	}

	args := []string{"attach"}
	if !SafeBoolFlag(cmdParam, "d") {
		args = append(args, "-d")
	}
	args = append(args, []string{"-t", p.Name}...)
	cmd := exec.Command("tmux", args...)
	cmd.Stdin = os.Stdin
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errorf("error on attach: %v", err)
	}
	logDebug(SafeBoolFlag(cmdParam, "debug"), "attch return: %v", string(out))
	return nil
}

func code(cmdParam *cobra.Command, params []string) error {
	projects, err := Load(LoadSettings())
	if err != nil {
		return err
	}

	p, _ := projects.Find(safeName(params...))
	if p == nil {
		return errorf("project not found")

	}
	if p.Path == "" {
		return errorf("project '%s' dont have path", p.Name)
	}
	if !isExist(p.Path) {
		return errorf("path '%s' of project '%s' not exists", p.Path, p.Name)
	}

	editor := "code"
	args := make([]string, 0)

	if edit := SafeStringFlag(cmdParam, "e"); edit != "" {
		editor = edit
	} else {
		pos := "--new-window"
		if SafeBoolFlag(cmdParam, "r") {
			pos = "--reuse-window"
		}
		args = append(args, pos)
	}
	args = append(args, p.Path)

	log("open path '%s' on '%s'", p.Path, editor)

	cmd := exec.Command(editor, args...)
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func close(cmdParam *cobra.Command, params []string) error {
	projects, err := Load(LoadSettings())
	if err != nil {
		return errorf("error on load file: %v", err)
	}

	toClose := make([]string, 0, 0)

	if SafeBoolFlag(cmdParam, "all") {
		for _, p := range projects {
			if p.Opened {
				toClose = append(toClose, p.Name)
			}
		}
		if len(toClose) == 0 {
			return errorf("no projects to close")
		}
	} else {
		name, _ := safeName(params...)
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
		if SafeBoolFlag(cmdParam, "kill") {
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

func scm(cmdParam *cobra.Command, params []string) error {
	projects, err := Load(LoadSettings())
	if err != nil {
		return err
	}

	p, index := projects.Find(safeName(params...))
	if p == nil {
		return errorf("project not found")
	}
	if p.Path == "" {
		return errorf("project '%s' dont have path", p.Name)
	}
	if !isExist(p.Path) {
		return errorf("path '%s' of project '%s' not exists", p.Path, p.Name)
	}

	log("using path: %s", p.Path)

	if SafeBoolFlag(cmdParam, "set") {
		var url string
		if len(params) <= 1 {
			cmd := exec.Command("git", "-C", p.Path, "remote", "get-url", "origin")
			out, err := cmd.CombinedOutput()
			if err != nil {
				return err
			}
			url = strings.TrimSpace(string(out))
		} else {
			url = params[1]
		}
		log("setting scm url %s", url)
		p.SCM = url
		projects[index] = *p
		return projects.Save(LoadSettings())
	}

	return nil
}
