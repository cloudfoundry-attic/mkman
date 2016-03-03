package releasemakers

//go:generate counterfeiter . ReleaseMaker
type ReleaseMaker interface {
	MakeRelease() (*Release, error)
}

type Release struct {
	Name    string `yaml:"name,omitempty"`
	Version string `yaml:"version,omitempty"`
	URL     string `yaml:"url,omitempty"`
}
