package path

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExistsInPathOrAsFileAbsolute(t *testing.T) {
	t.Parallel()

	execPath := createTempExecutable(t)
	if !ExistsInPathOrAsFile(execPath) {
		t.Fatalf("expected %s to be detected as executable", execPath)
	}
}

func TestExistsInPathOrAsFileRelativeViaPATH(t *testing.T) {
	dir := t.TempDir()
	execPath := filepath.Join(dir, "mytool")
	writeExecutable(t, execPath)
	t.Setenv("PATH", dir)

	if !ExistsInPathOrAsFile("mytool") {
		t.Fatalf("expected mytool to be resolved via PATH")
	}
}

func TestEnsureExecutableErrorsWhenMissing(t *testing.T) {
	t.Parallel()

	if err := EnsureExecutable("definitely-not-installed"); err == nil {
		t.Fatalf("expected error when executable is missing")
	}
}

func TestEnsureExecutableSuccess(t *testing.T) {
	t.Parallel()

	execPath := createTempExecutable(t)
	if err := EnsureExecutable(execPath); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSafeNameSkipsFlags(t *testing.T) {
	t.Parallel()

	name, _ := SafeName("-f", "--flag", "proj", "-x")
	if name != "proj" {
		t.Fatalf("expected proj, got %s", name)
	}
}

func TestSafeNameDefaultsToPWD(t *testing.T) {
	t.Setenv("PWD", "/tmp/myproj")

	name, path := SafeName()
	if name != "myproj" || path != "/tmp/myproj" {
		t.Fatalf("expected defaults from PWD, got name=%s path=%s", name, path)
	}
}

func TestExistFunction(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(file, []byte("data"), 0o644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	if !Exist(file) {
		t.Fatalf("expected file to exist")
	}

	if Exist(filepath.Join(dir, "missing")) {
		t.Fatalf("expected missing file to return false")
	}
}

func createTempExecutable(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "exec-file")
	writeExecutable(t, path)
	return path
}

func writeExecutable(t *testing.T, path string) {
	t.Helper()
	if err := os.WriteFile(path, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatalf("failed to write temp executable: %v", err)
	}
}
