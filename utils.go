package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func isExist(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return info != nil
}

func log(msg string, args ...interface{}) {
	fmt.Printf("[projects] %s\n", fmt.Sprintf(msg, args...))
}

func logDebug(debug bool, msg string, args ...interface{}) {
	if debug {
		fmt.Printf("[projects debug] %s\n", fmt.Sprintf(msg, args...))
	}
}

func errorf(msg string, args ...interface{}) error {
	return fmt.Errorf("%s", fmt.Sprintf(msg, args...))
}

func tmux(args ...string) (string, error) {
	cmd := exec.Command("tmux", args...)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func startServer() error {
	_, err := tmux("start-server")
	if err != nil {
		return errorf("error on start server: %v", err)
	}
	return nil
}

func getSessions() (map[string]bool, error) {
	m := make(map[string]bool, 0)

	if err := startServer(); err != nil {
		return m, errorf("error on start server: %v", err)
	}

	out, err := tmux("list-sessions")
	if err != nil && !strings.Contains(out, "no server running") {
		return m, err
	}
	for _, l := range strings.Split(out, "\n") {
		l = strings.TrimSpace(l)
		if len(l) == 0 {
			continue
		}
		p := strings.Split(l, ":")
		if len(p) == 0 {
			log("%v", p)
			continue
		}
		m[p[0]] = strings.Contains(strings.Join(p[1:], ""), "attached")
	}

	return m, nil
}

//pwd return the current path name and location
func current_pwd() (string, string) {
	pwd := os.Getenv("PWD")
	paths := strings.Split(pwd, "/")
	return strings.TrimSpace(paths[len(paths)-1]), strings.TrimSpace(pwd)
}
