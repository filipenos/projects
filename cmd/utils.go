package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

func isExist(path string) bool {
	if path == "" {
		return false
	}
	info, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return info != nil
}

func log(msg string, args ...interface{}) {
	fmt.Printf("[projects] %s\n", fmt.Sprintf(msg, args...))
}

func logDebug(debug bool, msg string, args ...interface{}) {
	if debug {
		fmt.Printf("[projects debug] %s\n", fmt.Sprintf(msg, args...))
	}
}

func errorf(msg string, args ...interface{}) error {
	return fmt.Errorf("%s", fmt.Sprintf(msg, args...))
}

// safeName return the name or find from pwd
func safeName(args ...string) (string, string) {
	if len(args) == 0 {
		return currentPwd()
	}
	name := strings.TrimSpace(args[0])
	if name != "" {
		return name, ""
	}
	return currentPwd()

}

// currentPwd return the current path and last path of dir
func currentPwd() (string, string) {
	pwd := os.Getenv("PWD")
	paths := strings.Split(pwd, "/")
	return strings.TrimSpace(paths[len(paths)-1]), strings.TrimSpace(pwd)
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
