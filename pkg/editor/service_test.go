package editor

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/filipenos/projects/pkg/config"
	"github.com/filipenos/projects/pkg/project"
)

func TestCreateConfigurableEditorValidatesInput(t *testing.T) {
	s := &Service{}

	_, err := s.createConfigurableEditor(config.EditorConfig{})
	if err == nil {
		t.Fatalf("expected error when name is missing")
	}

	_, err = s.createConfigurableEditor(config.EditorConfig{Name: "foo"})
	if err == nil {
		t.Fatalf("expected error when executable is missing")
	}
}

func TestCreateConfigurableEditorBuildsStruct(t *testing.T) {
	s := &Service{}
	cfg := config.EditorConfig{
		Name:       "custom",
		Aliases:    []string{"alias1"},
		Executable: "custom-bin",
		SupportedTypes: []string{
			string(project.ProjectTypeLocal),
		},
		WindowArgs: map[string][]string{
			string(WindowTypeNew): {"--new"},
		},
		WorkspaceArgs: []string{"--workspace"},
		FolderArgs:    []string{"--folder"},
		PathPosition:  1,
	}

	editor, err := s.createConfigurableEditor(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if editor.Name() != "custom" || editor.GetExecutable() != "custom-bin" {
		t.Fatalf("invalid editor created: %+v", editor)
	}
}

func TestNewServiceLoadsCustomEditors(t *testing.T) {
	tmp := t.TempDir()
	editorsPath := filepath.Join(tmp, "editors.json")

	custom := config.EditorsConfig{
		Editors: []config.EditorConfig{
			{
				Name:       "custom",
				Executable: "custom",
				Aliases:    []string{"custom"},
			},
		},
	}
	writeJSON(t, editorsPath, custom)

	cfg := config.Config{
		EditorsLocation: editorsPath,
	}

	service, err := NewService(cfg)
	if err != nil {
		t.Fatalf("NewService failed: %v", err)
	}

	foundAlias := false
	for _, alias := range service.Aliases() {
		if alias == "custom" {
			foundAlias = true
			break
		}
	}

	if !foundAlias {
		t.Fatalf("expected custom editor alias to be registered")
	}
}

func TestServiceGetEditorsCategorizesAvailability(t *testing.T) {
	service := &Service{
		registry: NewRegistry(),
	}

	availableExec := filepath.Join(t.TempDir(), "has-bin")
	if err := os.WriteFile(availableExec, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatalf("failed to write temp exec: %v", err)
	}

	service.RegisterEditor(&fakeEditor{
		name:       "available",
		executable: availableExec,
	})
	service.RegisterEditor(&fakeEditor{
		name:       "missing",
		executable: filepath.Join(t.TempDir(), "missing-bin"),
	})

	available, missing := service.GetEditors()
	if len(available) != 1 || available[0] != "available" {
		t.Fatalf("expected available editor, got %v", available)
	}
	if len(missing) != 1 || missing[0] != "missing" {
		t.Fatalf("expected missing editor, got %v", missing)
	}
}

type fakeEditor struct {
	name       string
	executable string
}

func (f *fakeEditor) Name() string                                 { return f.name }
func (f *fakeEditor) Aliases() []string                            { return []string{} }
func (f *fakeEditor) SupportsProjectType(project.ProjectType) bool { return true }
func (f *fakeEditor) BuildArgs(*project.Project, WindowType) ([]string, error) {
	return nil, nil
}
func (f *fakeEditor) GetExecutable() string { return f.executable }

func writeJSON(t *testing.T, path string, data any) {
	t.Helper()
	b, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("failed to marshal json: %v", err)
	}
	if err := os.WriteFile(path, b, 0o644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
}
