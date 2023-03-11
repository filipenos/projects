package log

import "fmt"

func Infof(msg string, args ...any) {
	fmt.Printf("[projects] %s\n", fmt.Sprintf(msg, args...))
}
