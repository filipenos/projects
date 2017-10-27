package main

import (
	"io/ioutil"
	"os"
	"os/exec"
)

//TempFile represent the Temporary File
type TempFile struct {
	Filename string
	osFile   *os.File
}

//NewTempFile create new TempFile with initial content
//initial content can be empty
func NewTempFile(data []byte) (*TempFile, error) {
	//TODO usar o tempdir, acho que já e
	tmpFile, err := ioutil.TempFile("", "temp")
	if err != nil {
		return nil, err
	}

	tmp := &TempFile{osFile: tmpFile, Filename: tmpFile.Name()}
	if len(data) > 0 {
		tmp.Write(data)
	}

	return tmp, nil
}

//GetContent retrieve content from file
func (f *TempFile) GetContent() []byte {
	b, _ := ioutil.ReadFile(f.Filename)
	return b
}

func (f *TempFile) Write(data []byte) error {
	_, err := f.osFile.Write(data)
	return err
}

//Remove temp file clean up
func (f *TempFile) Remove() {
	os.Remove(f.Filename)
}

//Close close temp file
func (f *TempFile) Close() error {
	return f.osFile.Close()
}

//ReadFromUser show editor to user
func (f *TempFile) ReadFromUser() error {
	cmd := exec.Command("editor", f.Filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
