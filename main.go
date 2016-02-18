package main

import (
	"fmt"
	"os"

	"github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/jessevdk/go-flags"
	"github.com/cloudfoundry/mkman/commands"
)

func main() {
	parser := flags.NewParser(&commands.Mkman, flags.HelpFlag|flags.PassDoubleDash)

	_, err := parser.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
