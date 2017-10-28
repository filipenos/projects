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
