package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	projectsConf    = fmt.Sprintf("%s/.projects.conf.json", os.Getenv("HOME"))
	projectsPath    = fmt.Sprintf("%s/.projects.json", os.Getenv("HOME"))
	defaultSettings = Config{ProjectLocation: projectsPath, Editor: "code"}
)

// Config save configuration
type Config struct {
	ProjectLocation string `json:"projects_location"`
	Editor          string `json:"editor"`
}

// Load load configuration used on projects
func Load() Config {
	file, err := loadFile(projectsConf)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if file == nil {
		return defaultSettings
	}

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		panic(err)
	}
	return config
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
