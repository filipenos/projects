package log

import "fmt"

func Infof(msg string, args ...any) {
	fmt.Printf("[projects] %s\n", fmt.Sprintf(msg, args...))
}

func Warnf(msg string, args ...any) {
	fmt.Printf("[projects] WARNING: %s\n", fmt.Sprintf(msg, args...))
}
