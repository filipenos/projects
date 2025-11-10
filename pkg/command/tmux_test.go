package command

import (
	"testing"

	"github.com/filipenos/projects/pkg/project"
)

func TestSanitizeSessionName(t *testing.T) {
	t.Parallel()

	cases := map[string]string{
		"":             "project",
		"my proj":      "my-proj",
		"UPPER_lower":  "UPPER_lower",
		"weird!chars?": "weird-chars-",
	}

	for input, expected := range cases {
		if got := sanitizeSessionName(input); got != expected {
			t.Fatalf("sanitizeSessionName(%q) = %q, expected %q", input, got, expected)
		}
	}
}

func TestProjectWorkingDir(t *testing.T) {
	t.Parallel()

	p := &project.Project{
		Path:     "/tmp/project/.code-workspace",
		RootPath: "/tmp/project",
	}
	p.IsWorkspace = true

	got := projectWorkingDir(p)
	expected := "/tmp/project"

	if got != expected {
		t.Fatalf("projectWorkingDir returned %s, expected %s", got, expected)
	}

	p.IsWorkspace = false
	p.Path = ""

	if got = projectWorkingDir(p); got != p.RootPath {
		t.Fatalf("expected fallback to RootPath, got %s", got)
	}
}
