package file

import (
	"os"
	"os/exec"
)

// TempFile represent the Temporary File
type TempFile struct {
	*os.File
}

// NewTempFile create new TempFile
func NewTempFile() (*TempFile, error) {
	tmp, err := os.CreateTemp("", "project_tmp_")
	if err != nil {
		return nil, err
	}
	return &TempFile{tmp}, nil
}

// GetContent retrieve content from file
func (f *TempFile) GetContent() ([]byte, error) {
	return os.ReadFile(f.Name())
}

// Remove temp file clean up
func (f *TempFile) Remove() error {
	return os.Remove(f.Name())
}

// ReadFromUser show editor to user
func (f *TempFile) ReadFromUser() error {
	cmd := exec.Command("vim", f.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
