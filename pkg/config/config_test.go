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

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg != defaultSettings {
		t.Fatalf("expected default config %+v, got %+v", defaultSettings, cfg)
	}
}

func TestLoadReadsConfigFile(t *testing.T) {
	restore := setupTempConfig(t)
	defer restore()

	custom := Config{
		ProjectLocation: "/tmp/projects.json",
		Editor:          "vim",
	}
	writeJSON(t, projectsConf, custom)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ProjectLocation != custom.ProjectLocation {
		t.Fatalf("expected project location %s, got %s", custom.ProjectLocation, cfg.ProjectLocation)
	}
	if cfg.Editor != "vim" {
		t.Fatalf("expected editor vim, got %s", cfg.Editor)
	}
	if cfg.SessionBackend != defaultSettings.SessionBackend {
		t.Fatalf("expected default session backend %s, got %s", defaultSettings.SessionBackend, cfg.SessionBackend)
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

func setupTempConfig(t *testing.T) func() {
	t.Helper()
	tmp := t.TempDir()

	oldProjectsConf := projectsConf
	oldProjectsPath := projectsPath
	oldDefault := defaultSettings

	projectsConf = filepath.Join(tmp, "projects.conf.json")
	projectsPath = filepath.Join(tmp, "projects.json")
	defaultSettings = Config{
		ProjectLocation: projectsPath,
		Editor:          "code",
	}

	return func() {
		projectsConf = oldProjectsConf
		projectsPath = oldProjectsPath
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
