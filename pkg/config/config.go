// config/config.go
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var (
	projectsConf    = fmt.Sprintf("%s/.projects.conf.json", os.Getenv("HOME"))
	projectsPath    = fmt.Sprintf("%s/.projects.json", os.Getenv("HOME"))
	defaultSettings = Config{
		ProjectLocation: projectsPath,
		Editor:          "code",
		SessionBackend:  "tmux",
	}
)

// Config save configuration
type Config struct {
	ProjectLocation string `json:"projects_location"`
	Editor          string `json:"editor"`
	SessionBackend  string `json:"session_backend,omitempty"`
}

// Load load configuration used on projects
func Load() (Config, error) {
	file, err := loadFile(projectsConf)
	if err != nil {
		return defaultSettings, fmt.Errorf("failed to load config: %w", err)
	}

	if file == nil {
		return defaultSettings, nil
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return defaultSettings, fmt.Errorf("failed to decode config: %w", err)
	}

	return config, nil
}

// GetConfigDir retorna o diretório onde ficam os arquivos de configuração
func GetConfigDir() string {
	return filepath.Dir(projectsConf)
}

func loadFile(path string) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return file, nil
}

func Init() error {
	if _, err := os.Stat(projectsConf); err == nil {
		return fmt.Errorf("config file already exists: %s", projectsConf)
	}
	b, err := json.MarshalIndent(defaultSettings, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(projectsConf, b, 0644)
}
