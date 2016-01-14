package commands

type MkmanCommand struct {
	Version   string           `long:"version" description:"Print version"`
	PrintAmit PrintAmitCommand `command:"print-amit" description:"Prints the man behind 'mkman'"`
}

var Mkman MkmanCommand
