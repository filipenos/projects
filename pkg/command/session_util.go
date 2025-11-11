package command

import (
	"strings"

	"github.com/filipenos/projects/pkg/project"
)

func projectWorkingDir(p *project.Project) string {
	pathValue := p.Path
	if pathValue == "" {
		pathValue = p.RootPath
	}
	if p.IsWorkspace {
		parts := strings.Split(pathValue, "/")
		if len(parts) > 1 {
			pathValue = strings.Join(parts[:len(parts)-1], "/")
		}
	}
	return pathValue
}

func sanitizeSessionName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "project"
	}

	var b strings.Builder
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '-' || r == '_' || r == '.':
			b.WriteRune(r)
		default:
			b.WriteRune('-')
		}
	}
	return b.String()
}
