package releasemakers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/mkman/Godeps/_workspace/src/gopkg.in/yaml.v2"
	"github.com/cloudfoundry/mkman/tarball"
)

type consulReleaseMaker struct {
	tarballReader tarball.TarballReader
	consulPath    string
}

func NewConsulReleaseMaker(tarballReader tarball.TarballReader, consulPath string) ReleaseMaker {
	return &consulReleaseMaker{
		tarballReader: tarballReader,
		consulPath:    consulPath,
	}
}

func (c *consulReleaseMaker) MakeRelease() (*Release, error) {
	var filePath string
	fileInfo, err := os.Stat(c.consulPath)
	if err != nil {
		return nil, err
	}
	if fileInfo.IsDir() {
		return &Release{
			Name:    "consul",
			URL:     "file://" + c.consulPath,
			Version: "create",
		}, nil
	}

	switch filepath.Ext(c.consulPath) {
	case ".tgz":
		filePath = "file://" + c.consulPath
	default:
		return nil, fmt.Errorf("unrecognized consul URL")
	}

	manifestContents, err := c.tarballReader.ReadFile("./release.MF")
	if err != nil {
		return nil, err
	}

	manifest := consulManifest{}
	err = yaml.Unmarshal(manifestContents, &manifest)
	if err != nil {
		// Untested, as it is too difficult to force Unmarshal to return an error.
		return nil, err
	}

	return &Release{
		Name:    "consul",
		URL:     filePath,
		Version: manifest.Version,
	}, nil
}

type consulManifest struct {
	Version string `yaml:"version"`
}
