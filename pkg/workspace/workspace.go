package workspace

import (
	"bytes"
	"encoding/json"
	"net/url"
	"os"
	"strings"
)

type Workspace struct {
	Path     string
	Folders  []Folders `json:"folders"`
	Settings Settings  `json:"settings"`
}
type Folders struct {
	Path string `json:"path"`
}
type RemoteSSHDefaultForwardedPorts struct {
	Name       string `json:"name"`
	LocalPort  int    `json:"localPort"`
	RemotePort int    `json:"remotePort"`
}
type Settings struct {
	RemoteSSHDefaultForwardedPorts []RemoteSSHDefaultForwardedPorts `json:"remote.SSH.defaultForwardedPorts"`
}

func Load(path string) (*Workspace, error) {
	w := new(Workspace)
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	file = stripJSONC(file)
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
	if strings.HasPrefix(path, "/") {
		return path
	}
	return w.Basepath() + "/" + path
}

func (w *Workspace) FoldersPath() []string {
	var paths []string
	for _, f := range w.Folders {
		if f.Path == "" {
			continue
		}
		paths = append(paths, w.FolderPath(f.Path))
	}
	return paths
}

func (f *Folders) UnmarshalJSON(data []byte) error {
	if len(bytes.TrimSpace(data)) == 0 {
		return nil
	}
	if data[0] == '"' {
		var p string
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		f.Path = p
		return nil
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if v, ok := raw["path"]; ok {
		var p string
		if err := json.Unmarshal(v, &p); err != nil {
			return err
		}
		f.Path = p
		return nil
	}
	if v, ok := raw["uri"]; ok {
		var u string
		if err := json.Unmarshal(v, &u); err != nil {
			return err
		}
		parsed, err := url.Parse(u)
		if err != nil {
			return err
		}
		if parsed.Scheme == "file" {
			f.Path = parsed.Path
		}
	}
	return nil
}

func stripJSONC(input []byte) []byte {
	// Remove // and /* */ comments, and trailing commas, preserving strings.
	const (
		stateNormal = iota
		stateString
		stateStringEscape
		stateLineComment
		stateBlockComment
		stateBlockCommentStar
	)

	var out []byte
	state := stateNormal
	for i := 0; i < len(input); i++ {
		c := input[i]
		switch state {
		case stateNormal:
			if c == '"' {
				out = append(out, c)
				state = stateString
				continue
			}
			if c == '/' && i+1 < len(input) {
				n := input[i+1]
				if n == '/' {
					state = stateLineComment
					i++
					continue
				}
				if n == '*' {
					state = stateBlockComment
					i++
					continue
				}
			}
			out = append(out, c)
		case stateString:
			out = append(out, c)
			if c == '\\' {
				state = stateStringEscape
				continue
			}
			if c == '"' {
				state = stateNormal
			}
		case stateStringEscape:
			out = append(out, c)
			state = stateString
		case stateLineComment:
			if c == '\n' {
				out = append(out, c)
				state = stateNormal
			}
		case stateBlockComment:
			if c == '*' {
				state = stateBlockCommentStar
			}
		case stateBlockCommentStar:
			if c == '/' {
				state = stateNormal
			} else if c != '*' {
				state = stateBlockComment
			}
		}
	}

	out = fixMissingArrayCommas(out)
	return stripTrailingCommas(out)
}

func stripTrailingCommas(input []byte) []byte {
	const (
		stateNormal = iota
		stateString
		stateStringEscape
	)

	out := make([]byte, 0, len(input))
	state := stateNormal
	for i := 0; i < len(input); i++ {
		c := input[i]
		if state == stateString {
			out = append(out, c)
			if c == '\\' {
				state = stateStringEscape
			} else if c == '"' {
				state = stateNormal
			}
			continue
		}
		if state == stateStringEscape {
			out = append(out, c)
			state = stateString
			continue
		}
		if c == '"' {
			out = append(out, c)
			state = stateString
			continue
		}
		if c == ',' {
			j := i + 1
			for j < len(input) {
				if input[j] == ' ' || input[j] == '\n' || input[j] == '\r' || input[j] == '\t' {
					j++
					continue
				}
				break
			}
			if j < len(input) && (input[j] == ']' || input[j] == '}') {
				continue
			}
		}
		out = append(out, c)
	}
	return out
}

func fixMissingArrayCommas(input []byte) []byte {
	const (
		stateNormal = iota
		stateString
		stateStringEscape
	)

	out := make([]byte, 0, len(input))
	state := stateNormal
	arrayDepth := 0
	for i := 0; i < len(input); i++ {
		c := input[i]
		if state == stateString {
			out = append(out, c)
			if c == '\\' {
				state = stateStringEscape
			} else if c == '"' {
				state = stateNormal
			}
			continue
		}
		if state == stateStringEscape {
			out = append(out, c)
			state = stateString
			continue
		}
		if c == '"' {
			out = append(out, c)
			state = stateString
			continue
		}
		if c == '[' {
			arrayDepth++
			out = append(out, c)
			continue
		}
		if c == ']' {
			if arrayDepth > 0 {
				arrayDepth--
			}
			out = append(out, c)
			continue
		}
		if c == '}' && arrayDepth > 0 {
			out = append(out, c)
			j := i + 1
			for j < len(input) {
				if input[j] == ' ' || input[j] == '\n' || input[j] == '\r' || input[j] == '\t' {
					j++
					continue
				}
				break
			}
			if j < len(input) && input[j] == '{' {
				out = append(out, ',')
			}
			continue
		}
		out = append(out, c)
	}
	return out
}
