package project

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/filipenos/projects/pkg/config"
)

func TestParseProjectType(t *testing.T) {
	cases := map[string]ProjectType{
		"local":             ProjectTypeLocal,
		"ssh+user@host":     ProjectTypeSSH,
		"wsl+Ubuntu":        ProjectTypeWSL,
		"tunnel+example":    ProjectTypeTunnel,
		"ssh-remote+server": ProjectTypeSSH,
		"custom":            ProjectType("custom"),
		"":                  ProjectType(""),
	}

	for input, expected := range cases {
		if got := ParseProjectType(input); got != expected {
			t.Fatalf("ParseProjectType(%q) = %s, expected %s", input, got, expected)
		}
	}
}

func TestParseURL(t *testing.T) {
	scheme, domain, path := parseURL("ssh://user@host/path")
	if scheme != "ssh" || domain != "user@host" || path != "/path" {
		t.Fatalf("unexpected parse result: %s %s %s", scheme, domain, path)
	}

	_, domain, path = parseURL("host")
	if domain != "host" || path != "" {
		t.Fatalf("expected host with empty path, got domain=%s path=%s", domain, path)
	}
}

func TestValidateLocalProject(t *testing.T) {
	dir := t.TempDir()
	p := &Project{
		Name:        "proj",
		RootPath:    dir,
		ProjectType: ProjectTypeLocal,
	}
	if err := p.Validate(); err != nil {
		t.Fatalf("unexpected error validating local project: %v", err)
	}

	p.RootPath = filepath.Join(dir, "missing")
	if err := p.Validate(); err == nil {
		t.Fatalf("expected validation error when path missing")
	}

	p.ProjectType = "unsupported"
	if err := p.Validate(); err == nil {
		t.Fatalf("expected error for unsupported project type")
	}
}

func TestProjectsFind(t *testing.T) {
	projects := Projects{
		{Name: "alpha", RootPath: "/tmp/alpha"},
		{Name: "beta", RootPath: "/tmp/beta"},
	}

	if p, _ := projects.Find("alpha", ""); p == nil || p.Name != "alpha" {
		t.Fatalf("expected to find project by name")
	}

	if p, _ := projects.Find("", "/tmp/beta"); p == nil || p.Name != "beta" {
		t.Fatalf("expected to find project by path")
	}

	if p, _ := projects.Find("", "/tmp/beta/subdir"); p == nil || p.Name != "beta" {
		t.Fatalf("expected to resolve by parent path")
	}
}

func TestProjectsSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	cfg := config.Config{
		ProjectLocation: filepath.Join(dir, "projects.json"),
	}

	root := filepath.Join(dir, "proj")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatalf("failed to create project dir: %v", err)
	}

	projects := Projects{
		{
			Name:     "proj",
			RootPath: root,
		},
	}

	if err := projects.Save(cfg); err != nil {
		t.Fatalf("failed to save projects: %v", err)
	}

	loaded, err := Load(cfg)
	if err != nil {
		t.Fatalf("failed to load projects: %v", err)
	}
	if len(loaded) != 1 || loaded[0].Name != "proj" {
		t.Fatalf("unexpected loaded content: %+v", loaded)
	}
	if !loaded[0].ValidPath || loaded[0].ProjectType != ProjectTypeLocal {
		t.Fatalf("expected metadata to be populated: %+v", loaded[0])
	}
}

func TestParseContent(t *testing.T) {
	data := []byte("name=proj\npath=/tmp\ngroup=dev\nenabled=true\n")
	p := ParseContent(data)
	if p.Name != "proj" || p.RootPath != "/tmp" || p.Group != "dev" || !p.Enabled {
		t.Fatalf("unexpected parsed result: %+v", p)
	}
}
