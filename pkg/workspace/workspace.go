package workspace

import (
	"encoding/json"
	"os"
	"strings"
)

type Workspace struct {
	Path    string
	Folders []Folders `json:"folders"`
}
type Folders struct {
	Path string `json:"path"`
}

func Load(path string) (*Workspace, error) {
	w := new(Workspace)
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(file, w); err != nil {
		return nil, err
	}
	w.Path = path

	return w, nil
}

func (w *Workspace) Basepath() string {
	parts := strings.Split(w.Path, "/")
	return strings.Join(parts[:len(parts)-1], "/")
}

func (w *Workspace) FolderPath(path string) string {
	return w.Basepath() + "/" + path
}

func (w *Workspace) FoldersPath() []string {
	var paths []string
	for _, f := range w.Folders {
		paths = append(paths, w.FolderPath(f.Path))
	}
	return paths
}
