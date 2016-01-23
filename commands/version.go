package commands

import (
	"fmt"
	"os"
)

const (
	Version = "0.0.1"
)

var VersionFunc = func() {
	fmt.Printf("mkman %s\n", Version)
	os.Exit(0)
}
