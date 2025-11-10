package log

import (
	"fmt"
	"io"
	"os"
)

var (
	stdout     io.Writer = os.Stdout
	stderr     io.Writer = os.Stderr
	infoPrefix           = "[projects]"
	warnPrefix           = "[projects] WARNING"
)

func SetOutput(w io.Writer) {
	if w != nil {
		stdout = w
	}
}

func SetErrorOutput(w io.Writer) {
	if w != nil {
		stderr = w
	}
}

func Infof(msg string, args ...any) {
	fmt.Fprintf(stdout, "%s %s\n", infoPrefix, fmt.Sprintf(msg, args...))
}

func Warnf(msg string, args ...any) {
	fmt.Fprintf(stderr, "%s %s\n", warnPrefix, fmt.Sprintf(msg, args...))
}

func Println(args ...any) {
	fmt.Fprintln(stdout, args...)
}

func Printf(format string, args ...any) {
	fmt.Fprintf(stdout, format, args...)
}
