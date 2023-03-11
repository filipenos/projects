package config

import (
	"fmt"
	"os"
)

var (
	projectsPath    = fmt.Sprintf("%s/.projects.json", os.Getenv("HOME"))
	defaultSettings = Config{ProjectLocation: projectsPath, Editor: "code"}
)

// Config save configuration
type Config struct {
	ProjectLocation string
	Editor          string
}

// Load load configuration used on projects
func Load() Config {
	return defaultSettings
}
