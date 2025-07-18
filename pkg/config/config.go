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
	editorsConf     = fmt.Sprintf("%s/.projects.editors.json", os.Getenv("HOME"))
	defaultSettings = Config{
		ProjectLocation: projectsPath,
		Editor:          "code",
		EditorsLocation: editorsConf,
	}
)

// Config save configuration
type Config struct {
	ProjectLocation string `json:"projects_location"`
	Editor          string `json:"editor"`
	EditorsLocation string `json:"editors_location,omitempty"`
}

// EditorConfig representa a configuração de um editor no arquivo JSON
type EditorConfig struct {
	Name           string              `json:"name"`
	Aliases        []string            `json:"aliases,omitempty"`
	Executable     string              `json:"executable"`
	SupportedTypes []string            `json:"supported_types,omitempty"`
	WindowArgs     map[string][]string `json:"window_args,omitempty"`
	WorkspaceArgs  []string            `json:"workspace_args,omitempty"`
	FolderArgs     []string            `json:"folder_args,omitempty"`
	PathPosition   int                 `json:"path_position,omitempty"`
}

// EditorsConfig representa o arquivo de configuração de editores
type EditorsConfig struct {
	Editors []EditorConfig `json:"editors"`
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

	// Se EditorsLocation não estiver definido, usa o padrão
	if config.EditorsLocation == "" {
		config.EditorsLocation = editorsConf
	}

	return config
}

// LoadEditors carrega a configuração de editores
func LoadEditors(config Config) (*EditorsConfig, error) {
	file, err := loadFile(config.EditorsLocation)
	if err != nil {
		return nil, err
	}

	if file == nil {
		return nil, nil // Arquivo não existe
	}
	defer file.Close()

	var editorsConfig EditorsConfig
	if err := json.NewDecoder(file).Decode(&editorsConfig); err != nil {
		return nil, fmt.Errorf("failed to decode editors config: %w", err)
	}

	return &editorsConfig, nil
}

// GetConfigDir retorna o diretório onde ficam os arquivos de configuração
func GetConfigDir() string {
	return filepath.Dir(projectsConf)
}

// GetEditorsConfigPath retorna o caminho do arquivo de configuração de editores
func GetEditorsConfigPath() string {
	return editorsConf
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

// InitEditors cria o arquivo de configuração de editores com exemplos
func InitEditors() error {
	if _, err := os.Stat(editorsConf); err == nil {
		return fmt.Errorf("editors config file already exists: %s", editorsConf)
	}

	defaultEditorsConfig := EditorsConfig{
		Editors: []EditorConfig{
			{
				Name:           "cursor",
				Aliases:        []string{"cursor"},
				Executable:     "cursor",
				SupportedTypes: []string{"local", "remote"},
				WindowArgs: map[string][]string{
					"new":   {"--new-window"},
					"reuse": {"--reuse-window"},
					"add":   {"--add"},
				},
				WorkspaceArgs: []string{"--file-uri"},
				FolderArgs:    []string{"--folder-uri"},
				PathPosition:  -1,
			},
			{
				Name:           "intellij",
				Aliases:        []string{"idea", "intellij"},
				Executable:     "idea",
				SupportedTypes: []string{"local"},
				WindowArgs: map[string][]string{
					"new": {"--new-window"},
				},
				PathPosition: -1,
			},
			{
				Name:           "emacs",
				Aliases:        []string{"emacs"},
				Executable:     "emacs",
				SupportedTypes: []string{"local", "wsl"},
				PathPosition:   -1,
			},
		},
	}

	b, err := json.MarshalIndent(defaultEditorsConfig, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(editorsConf, b, 0644)
}
