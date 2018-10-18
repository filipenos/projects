package main

import (
	"fmt"
	"os"
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
