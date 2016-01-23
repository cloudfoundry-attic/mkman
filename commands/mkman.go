package commands

type MkmanCommand struct {
	Version   func()           `long:"version" description:"Print version"`
	PrintAmit PrintAmitCommand `command:"print-amit" description:"Prints the man behind 'mkman'"`
}

var Mkman MkmanCommand = MkmanCommand{
	Version: VersionFunc,
}
