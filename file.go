package main

import (
	"encoding/json"
	"os"
	"sort"
)

//Project represent then project
type Project struct {
	Name string
	Path string
}

//File represet all projects managed by
type File struct {
	Path     string
	Projects []Project
}

func (f File) Len() int           { return len(f.Projects) }
func (f File) Swap(i, j int)      { f.Projects[i], f.Projects[j] = f.Projects[j], f.Projects[i] }
func (f File) Less(i, j int) bool { return f.Projects[i].Name < f.Projects[j].Name }

//Add new project to manage
func (f *File) Add(name, path string) {
	f.Projects = append(f.Projects, Project{Name: name, Path: path})
}

//Save save the current projects on conf file
func (f *File) Save() error {
	file, err := os.Create(f.Path)
	if err != nil {
		return err
	}
	defer file.Close()
	sort.Sort(f)

	enc := json.NewEncoder(file)
	//enc.SetIndent("", "  ") Only go 1.7
	return enc.Encode(f.Projects)
}

//Load retrieve projects from config file
func Load(path string) (*File, error) {
	f := &File{Path: path}

	file, err := os.Open(f.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return f, nil
		}
		return nil, err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&f.Projects); err != nil {
		return nil, err
	}
	return f, nil
}

func checkPath(path string) bool {
	_, err := os.Stat(path)
	return os.IsExist(err)
}
