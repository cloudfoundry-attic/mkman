package stemcell

import (
	"fmt"
	"path/filepath"

	"github.com/pivotal-cf-experimental/mkman/tar"

	"gopkg.in/yaml.v2"
)

func StubFromTar(stemcellPath string) (string, error) {
	fmt.Printf("@@@ DEBUG tar looking for stemcell.MF\n")
	manifestContents, err := tar.ReadFileContentsFromTar(stemcellPath, "stemcell.MF")
	if err != nil {
		panic(err)
	}

	fmt.Printf("@@@ DEBUG manifest contents: %s\n", string(manifestContents))
	manifest := stemcellManifest{}
	err = yaml.Unmarshal(manifestContents, &manifest)
	if err != nil {
		panic(err)
	}

	return stub(
		manifest.Name,
		manifest.Version,
		stemcellPath,
	)
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
