package command

import (
	"testing"

	"github.com/filipenos/projects/pkg/project"
)

type fakeSessionBackend struct {
	name    string
	aliases []string
}

func (f *fakeSessionBackend) Name() string { return f.name }
func (f *fakeSessionBackend) Aliases() []string {
	return f.aliases
}
func (f *fakeSessionBackend) Run(_ *project.Project, _ []string) error {
	return nil
}

func TestParseSessionParamsDefaults(t *testing.T) {
	backend, project, args, err := parseSessionParams([]string{"proj"}, "tmux")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if backend != "tmux" || project != "proj" || len(args) != 0 {
		t.Fatalf("unexpected results: backend=%s project=%s args=%v", backend, project, args)
	}
}

func TestParseSessionParamsHandlesFlags(t *testing.T) {
	backend, project, args, err := parseSessionParams([]string{"--backend", "screen", "proj", "split-window"}, "tmux")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if backend != "screen" || project != "proj" || len(args) != 1 || args[0] != "split-window" {
		t.Fatalf("unexpected parsed values backend=%s project=%s args=%v", backend, project, args)
	}

	backend, project, _, err = parseSessionParams([]string{"-b", "screen", "proj"}, "tmux")
	if err != nil || backend != "screen" || project != "proj" {
		t.Fatalf("failed to parse shorthand flag: backend=%s project=%s err=%v", backend, project, err)
	}
}

func TestParseSessionParamsErrors(t *testing.T) {
	if _, _, _, err := parseSessionParams([]string{}, "tmux"); err == nil {
		t.Fatalf("expected error when no project provided")
	}
	if _, _, _, err := parseSessionParams([]string{"--backend"}, "tmux"); err == nil {
		t.Fatalf("expected error when backend flag missing value")
	}
	if _, _, _, err := parseSessionParams([]string{"-b"}, "tmux"); err == nil {
		t.Fatalf("expected error when shorthand flag missing value")
	}
}

func TestGetSessionBackend(t *testing.T) {
	original := availableSessionBackends
	defer func() { availableSessionBackends = original }()

	availableSessionBackends = []sessionBackend{
		&fakeSessionBackend{name: "demo", aliases: []string{"d"}},
	}

	if _, err := getSessionBackend("demo"); err != nil {
		t.Fatalf("expected to find backend by name: %v", err)
	}
	if _, err := getSessionBackend("d"); err != nil {
		t.Fatalf("expected to find backend by alias: %v", err)
	}
	if _, err := getSessionBackend("missing"); err == nil {
		t.Fatalf("expected error for unknown backend")
	}
}
