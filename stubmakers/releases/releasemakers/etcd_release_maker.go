package releasemakers

import (
	"fmt"
	"path/filepath"

	"github.com/cloudfoundry/mkman/Godeps/_workspace/src/gopkg.in/yaml.v2"
	"github.com/cloudfoundry/mkman/tarball"
)

type etcdReleaseMaker struct {
	tarballReader tarball.TarballReader
	etcdPath      string
}

func NewEtcdReleaseMaker(tarballReader tarball.TarballReader, etcdPath string) ReleaseMaker {
	return &etcdReleaseMaker{
		tarballReader: tarballReader,
		etcdPath:      etcdPath,
	}
}

func (e *etcdReleaseMaker) MakeRelease() (*Release, error) {
	var filePath string
	switch filepath.Ext(e.etcdPath) {
	case ".tgz":
		filePath = "file://" + e.etcdPath
	default:
		return nil, fmt.Errorf("unrecognized etcd URL")
	}

	manifestContents, err := e.tarballReader.ReadFile("./release.MF")
	if err != nil {
		return nil, err
	}

	manifest := etcdManifest{}
	err = yaml.Unmarshal(manifestContents, &manifest)
	if err != nil {
		// Untested, as it is too difficult to force Unmarshal to return an error.
		return nil, err
	}

	return &Release{
		Name:    "etcd",
		URL:     filePath,
		Version: manifest.Version,
	}, nil
}

type etcdManifest struct {
	Version string `yaml:"version"`
}
