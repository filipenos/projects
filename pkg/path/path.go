package path

import (
	"os"
	"strings"
)

func Exist(path string) bool {
	if path == "" {
		return false
	}
	info, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return info != nil
}

// SafeName return the name or find from pwd
func SafeName(args ...string) (string, string) {
	if len(args) == 0 {
		return CurrentPwd()
	}
	var name string
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			continue
		}
		name = a
		break
	}
	if name != "" {
		return name, ""
	}
	return CurrentPwd()
}

// CurrentPwd return the current path and last path of dir
func CurrentPwd() (string, string) {
	pwd := os.Getenv("PWD")
	paths := strings.Split(pwd, "/")
	return strings.TrimSpace(paths[len(paths)-1]), strings.TrimSpace(pwd)
}
