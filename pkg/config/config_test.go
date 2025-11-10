package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadReturnsDefaultWhenFileMissing(t *testing.T) {
	restore := setupTempConfig(t)
	defer restore()

	cfg := Load()
	if cfg != defaultSettings {
		t.Fatalf("expected default config %+v, got %+v", defaultSettings, cfg)
	}
}

func TestLoadReadsConfigFileAndBackfillsEditorsPath(t *testing.T) {
	restore := setupTempConfig(t)
	defer restore()

	custom := Config{
		ProjectLocation: "/tmp/projects.json",
		Editor:          "vim",
	}
	writeJSON(t, projectsConf, custom)

	cfg := Load()
	if cfg.ProjectLocation != custom.ProjectLocation {
		t.Fatalf("expected project location %s, got %s", custom.ProjectLocation, cfg.ProjectLocation)
	}
	if cfg.Editor != "vim" {
		t.Fatalf("expected editor vim, got %s", cfg.Editor)
	}
	if cfg.EditorsLocation != editorsConf {
		t.Fatalf("expected editors location fallback %s, got %s", editorsConf, cfg.EditorsLocation)
	}
}

func TestInitCreatesConfigFile(t *testing.T) {
	restore := setupTempConfig(t)
	defer restore()

	if err := Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	var cfg Config
	readJSON(t, projectsConf, &cfg)

	if cfg != defaultSettings {
		t.Fatalf("unexpected config content: %+v", cfg)
	}

	// Second call should error because file now exists
	if err := Init(); err == nil {
		t.Fatalf("expected error on second Init call")
	}
}

func TestInitEditorsCreatesFile(t *testing.T) {
	restore := setupTempConfig(t)
	defer restore()

	if err := InitEditors(); err != nil {
		t.Fatalf("InitEditors failed: %v", err)
	}

	var editors EditorsConfig
	readJSON(t, editorsConf, &editors)

	if len(editors.Editors) == 0 {
		t.Fatalf("expected default editors content")
	}

	if err := InitEditors(); err == nil {
		t.Fatalf("expected error when editors config already exists")
	}
}

func TestLoadEditors(t *testing.T) {
	restore := setupTempConfig(t)
	defer restore()

	cfg := defaultSettings

	loaded, err := LoadEditors(cfg)
	if err != nil {
		t.Fatalf("LoadEditors returned error on missing file: %v", err)
	}
	if loaded != nil {
		t.Fatalf("expected nil when file missing")
	}

	expected := EditorsConfig{
		Editors: []EditorConfig{
			{
				Name:       "custom",
				Executable: "custom",
			},
		},
	}
	writeJSON(t, editorsConf, expected)

	loaded, err = LoadEditors(cfg)
	if err != nil {
		t.Fatalf("LoadEditors failed: %v", err)
	}
	if len(loaded.Editors) != 1 || loaded.Editors[0].Name != "custom" {
		t.Fatalf("unexpected editors content: %+v", loaded)
	}
}

func setupTempConfig(t *testing.T) func() {
	t.Helper()
	tmp := t.TempDir()

	oldProjectsConf := projectsConf
	oldProjectsPath := projectsPath
	oldEditorsConf := editorsConf
	oldDefault := defaultSettings

	projectsConf = filepath.Join(tmp, "projects.conf.json")
	projectsPath = filepath.Join(tmp, "projects.json")
	editorsConf = filepath.Join(tmp, "editors.conf.json")
	defaultSettings = Config{
		ProjectLocation: projectsPath,
		Editor:          "code",
		EditorsLocation: editorsConf,
	}

	return func() {
		projectsConf = oldProjectsConf
		projectsPath = oldProjectsPath
		editorsConf = oldEditorsConf
		defaultSettings = oldDefault
	}
}

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

func readJSON(t *testing.T, path string, target any) {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if err := json.Unmarshal(b, target); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
}
