package path

import (
	"os"
	"os/exec"
	"path/filepath"
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

func ExistsInPathOrAsFile(name string) bool {
	// Se for caminho absoluto ou relativo (tem /), verifica se o arquivo existe e é executável
	if filepath.Base(name) != name {
		info, err := os.Stat(name)
		if err != nil {
			return false
		}
		mode := info.Mode()
		return !mode.IsDir() && mode&0111 != 0 // verificando se é executável (permissão de execução)
	}

	// Caso contrário, tenta procurar no PATH
	_, err := exec.LookPath(name)
	return err == nil
}
