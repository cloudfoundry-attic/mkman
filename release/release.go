package release

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

func StubFromReleasePath(path string) (string, error) {
	fmt.Printf("@@@ DEBUG release looking for release.MF\n")

	releaseStub := releaseStub{
		Releases: []release{
			{
				Name:    "cf",
				URL:     "file://" + path,
				Version: "create",
			},
		},
	}

	b, err := yaml.Marshal(releaseStub)
	return string(b), err
}

type releaseStub struct {
	Releases []release `json:"release,omitempty"`
}

type release struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	URL     string `json:"url"`
}
