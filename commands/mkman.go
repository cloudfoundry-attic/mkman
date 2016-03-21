package commands

type MkmanCommand struct {
	Version         func()                 `long:"version" description:"Print version"`
	CreateManifests CreateManifestsCommand `command:"create-manifests" description:"Make a manifest from a config file"`
}

var Mkman MkmanCommand = MkmanCommand{
	Version: VersionFunc,
}
