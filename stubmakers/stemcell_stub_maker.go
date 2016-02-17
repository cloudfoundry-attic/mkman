package stubmakers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pivotal-cf-experimental/mkman/tar"
	"gopkg.in/yaml.v2"
)

type stemcellStubMaker struct {
	tarballPath string
}

func NewStemcellStubMaker(tarballPath string) *stemcellStubMaker {
	return &stemcellStubMaker{
		tarballPath: tarballPath,
	}
}

func (s *stemcellStubMaker) MakeStub() (string, error) {
	manifestContents, err := tar.ReadFileContentsFromTar(s.tarballPath, "stemcell.MF")
	if err != nil {
		panic(err)
	}

	fmt.Printf("@@@ DEBUG manifest contents: %s\n", string(manifestContents))
	manifest := stemcellManifest{}
	err = yaml.Unmarshal(manifestContents, &manifest)
	if err != nil {
		panic(err)
	}

	stemcellStubContents, err := stub(
		manifest.Name,
		manifest.Version,
		s.tarballPath,
	)
	if err != nil {
		panic(err)
	}

	intermediateDir, err := ioutil.TempDir("", "")
	if err != nil {
		// We cannot test this because it is too hard to get TempDir to return error
		return "", err
	}

	stemcellStubPath := filepath.Join(intermediateDir, "stemcell.yml")
	err = ioutil.WriteFile(stemcellStubPath, []byte(stemcellStubContents), os.ModePerm)
	if err != nil {
		panic(err)
	}

	return stemcellStubPath, nil
}

func stub(
	name,
	version string,
	stemcellURL string,
) (string, error) {
	stemcellStub := stemcellStub{
		Meta: meta{
			Stemcell: stemcell{
				Name:    name,
				Version: version,
			},
		},
	}

	if filepath.Ext(stemcellURL) == ".tgz" {
		stemcellStub.Meta.Stemcell.URL = "file://" + stemcellURL
	}

	b, err := yaml.Marshal(stemcellStub)
	return string(b), err
}

type stemcellStub struct {
	Meta meta `json:"meta"`
}

type meta struct {
	Stemcell stemcell `json:"stemcell,omitempty"`
}

type stemcell struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	URL     string `json:"url,omitempty"`
}

type stemcellManifest struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
