package editor

import (
	"testing"

	"github.com/filipenos/projects/pkg/project"
)

func TestNewServiceRegistersBuiltinEditors(t *testing.T) {
	service := NewService()

	aliases := service.Aliases()
	if len(aliases) == 0 {
		t.Fatalf("expected builtin editor aliases")
	}

	for _, name := range []string{"cursor", "windsurf", "antigravity", "sublime"} {
		if _, ok := service.byName[name]; !ok {
			t.Errorf("expected editor %q to be registered", name)
		}
	}
}

func TestAliasesExcludesCodeName(t *testing.T) {
	service := NewService()
	for _, a := range service.Aliases() {
		if a == "code" {
			t.Fatal("aliases should not include 'code' (it's the main command name)")
		}
	}
}

func TestVSCodeBuildArgs(t *testing.T) {
	e := vscodeEditor("code", "code", nil)

	p := &project.Project{RootPath: "/tmp/test"}
	args, err := e.BuildArgs(p, WindowTypeNew)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(args) != 3 || args[0] != "--new-window" || args[1] != "--folder-uri" || args[2] != "/tmp/test" {
		t.Fatalf("unexpected args: %v", args)
	}

	p.IsWorkspace = true
	args, _ = e.BuildArgs(p, WindowTypeReuse)
	if len(args) != 3 || args[0] != "--reuse-window" || args[1] != "--file-uri" {
		t.Fatalf("unexpected workspace args: %v", args)
	}
}
