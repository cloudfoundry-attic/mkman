package testhelpers

import (
	"path"
	"runtime"
)

func GetDirOfCurrentFile() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}
