package config

type Config struct {
	CFPath       string   `yaml:"cf"`
	StemcellPath string   `yaml:"stemcell"`
	StubPaths    []string `yaml:"stubs"`
}
