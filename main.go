package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/pivotal-cf-experimental/mkman/commands"
)

const (
	Version = "0.0.1"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("mkman %s\n", Version)
		os.Exit(0)
	}
	parser := flags.NewParser(&commands.Mkman, flags.HelpFlag|flags.PassDoubleDash)
	parser.NamespaceDelimiter = "-"

	_, err := parser.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
