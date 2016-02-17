package stubmakers

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type releaseStubMaker struct {
	releasePath string
}

func NewReleaseStubMaker(releasePath string) *releaseStubMaker {
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

	releaseStubContents, err := yaml.Marshal(releaseStub)
	if err != nil {
		return "", nil
	}

	intermediateDir, err := ioutil.TempDir("", "")
	if err != nil {
		// We cannot test this because it is too hard to get TempDir to return error
		return "", err
	}

	releaseStubPath := filepath.Join(intermediateDir, "releases.yml")
	err = ioutil.WriteFile(releaseStubPath, releaseStubContents, os.ModePerm)
	if err != nil {
		panic(err)
	}

	return releaseStubPath, nil
}

type releaseStub struct {
	Releases []release `json:"release,omitempty"`
}

type release struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	URL     string `json:"url"`
}
