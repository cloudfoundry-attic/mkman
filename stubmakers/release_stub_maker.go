package stubmakers

type releaseStubMaker struct {
	releasePath string
}

func NewReleaseStubMaker(releasePath string) StubMaker {
	return &releaseStubMaker{
		releasePath: releasePath,
	}
}

func (r *releaseStubMaker) MakeStub() (string, error) {
	releaseStub := releaseStub{
		Releases: []release{
			{
				Name:    "cf",
				URL:     "file://" + r.releasePath,
				Version: "create",
			},
		},
	}

	return marshalTempStub(releaseStub, "release.yml")
}

type releaseStub struct {
	Releases []release `yaml:"releases,omitempty"`
}

type release struct {
	Name    string `yaml:"name,omitempty"`
	Version string `yaml:"version,omitempty"`
	URL     string `yaml:"url,omitempty"`
}
