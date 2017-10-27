package main

import "testing"

func TestNewTempFile(t *testing.T) {

	s := "initial data to test"
	f, err := NewTempFile([]byte(s))
	if err != nil {
		t.Fatalf("Unexpected error on create new temp file: %v", err)
	}
	if f == nil {
		t.Fatalf("Unexpected <nil> on create new temp file")
	}

	if err := f.Close(); err != nil {
		t.Fatalf("Unexpected error on close temp file: %v", err)
	}

	content := f.GetContent()
	if len(content) == 0 {
		t.Fatalf("Unexpected 0 content")
	}
}
