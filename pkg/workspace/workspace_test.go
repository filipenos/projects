package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadWorkspace(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "proj.code-workspace")
	content := `{
		"folders": [
			{"path": "../folder1"},
			{"path": "./folder2"}
		],
		"settings": {
			"remote.SSH.defaultForwardedPorts": [
				{"name":"db","localPort":5432,"remotePort":5432}
			]
		}
	}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write workspace file: %v", err)
	}

	ws, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if ws.Path != path {
		t.Fatalf("expected Path to be set, got %s", ws.Path)
	}
	if len(ws.Folders) != 2 {
		t.Fatalf("expected 2 folders, got %d", len(ws.Folders))
	}
	if len(ws.Settings.RemoteSSHDefaultForwardedPorts) != 1 {
		t.Fatalf("expected remote port settings to be parsed")
	}
}

func TestWorkspacePaths(t *testing.T) {
	ws := &Workspace{
		Path: "/projects/sample.code-workspace",
		Folders: []Folders{
			{Path: "folderA"},
			{Path: "../folderB"},
		},
	}

	base := ws.Basepath()
	if base != "/projects" {
		t.Fatalf("expected basepath '/projects', got %s", base)
	}

	if ws.FolderPath("folderA") != "/projects/folderA" {
		t.Fatalf("FolderPath returned unexpected path")
	}

	paths := ws.FoldersPath()
	if len(paths) != 2 || paths[0] != "/projects/folderA" || paths[1] != "/projects/../folderB" {
		t.Fatalf("FoldersPath returned unexpected result: %v", paths)
	}
}
