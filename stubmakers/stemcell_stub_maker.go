package stubmakers

import (
	"fmt"
	"path/filepath"

	"github.com/cloudfoundry/mkman/tarball"

	"gopkg.in/yaml.v2"
)

type stemcellStubMaker struct {
	tarballReader tarball.TarballReader
	stemcellURL   string
}

func NewStemcellStubMaker(tarballReader tarball.TarballReader, stemcellURL string) StubMaker {
	return &stemcellStubMaker{
		tarballReader: tarballReader,
		stemcellURL:   stemcellURL,
	}
}

func (s *stemcellStubMaker) MakeStub() (string, error) {
	stemcellStub := stemcellStub{
		Meta: meta{
			Stemcell: stemcell{},
		},
	}

	switch filepath.Ext(s.stemcellURL) {
	case ".tgz":
		stemcellStub.Meta.Stemcell.URL = "file://" + s.stemcellURL
	default:
		return "", fmt.Errorf("unrecognized stemcell URL")
	}

	manifestContents, err := s.tarballReader.ReadFile("stemcell.MF")
	if err != nil {
		return "", err
	}

	manifest := stemcellManifest{}
	err = yaml.Unmarshal(manifestContents, &manifest)
	if err != nil {
		// Untested, as it is too difficult to force Unmarshal to return an error.
		return "", err
	}

	stemcellStub.Meta.Stemcell.Name = manifest.Name
	stemcellStub.Meta.Stemcell.Version = manifest.Version

	return marshalTempStub(stemcellStub, "stemcell.yml")
}

type stemcellStub struct {
	Meta meta `yaml:"meta"`
}

type meta struct {
	Stemcell stemcell `yaml:"stemcell,omitempty"`
}

type stemcell struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	URL     string `yaml:"url,omitempty"`
}

type stemcellManifest struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}
