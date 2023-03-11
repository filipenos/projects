package file

import (
	"testing"

	"github.com/filipenos/projects/pkg/path"
)

func TestNewTempFile(t *testing.T) {

	s := "initial data to test"
	f, err := NewTempFile()
	if err != nil {
		t.Fatalf("Unexpected error on create new temp file: %v", err)
	}
	if f == nil {
		t.Fatalf("Unexpected <nil> on create new temp file")
	}

	n, err := f.Write([]byte(s))
	if err != nil {
		t.Fatalf("Unexpected error on write data: %v", err)
	}
	if len(s) != n {
		t.Fatalf("Unexpected data write %d, expect %d", n, len(s))
	}

	if err := f.Close(); err != nil {
		t.Fatalf("Unexpected error on close temp file: %v", err)
	}

	name := f.Name()
	t.Logf("Created file %s", name)

	if !path.Exist(name) {
		t.Errorf("Temp file not exists")
	}

	content, err := f.GetContent()
	if err != nil {
		t.Fatalf("Unexpected error on get content: %v", err)
	}
	if string(content) != s {
		t.Fatalf("Unexpected content '%s', expect '%s", string(content), s)
	}

	if err := f.Remove(); err != nil {
		t.Fatalf("Unexpected error on remove temp file: %v", err)
	}

	if path.Exist(name) {
		t.Fatalf("Temp file exist, no removed")
	}
}
