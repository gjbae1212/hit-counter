package util

import (
	"log"
	"path"
	"path/filepath"
	"runtime"
)

var (
	root string
)

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Panic("No caller information")
	}
	root = filepath.Join(path.Dir(filename), "../")
}

func GetRoot() string {
	return root
}
