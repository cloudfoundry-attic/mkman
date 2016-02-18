package commands

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/pivotal-cf-experimental/mkman/config"
	"github.com/pivotal-cf-experimental/mkman/manifestgenerator"
	"github.com/pivotal-cf-experimental/mkman/stubmakers"
	"github.com/pivotal-cf-experimental/mkman/tarball"

	"gopkg.in/yaml.v2"
)

type ManifestGenerator interface {
	GenerateManifest() (string, error)
}

type CreateManifestsCommand struct {
	OutputWriter io.Writer
}

func (command *CreateManifestsCommand) Execute(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("create-manifests requires PATH_TO_CONFIG")
	}

	configFileContents, err := ioutil.ReadFile(args[0])
	if err != nil {
		return err
	}

	config := config.Config{}
	err = yaml.Unmarshal(configFileContents, &config)
	if err != nil {
		return err
	}

	tarballReader := tarball.NewTarballReader(config.StemcellPath)
	stemcellStubMaker := stubmakers.NewStemcellStubMaker(tarballReader, config.StemcellPath)
	releaseStubMaker := stubmakers.NewReleaseStubMaker(config.CFPath)
	manifestGenerator := manifestgenerator.NewSpiffManifestGenerator(stemcellStubMaker, releaseStubMaker, config.StubPaths, config.CFPath)

	manifest, err := manifestGenerator.GenerateManifest()
	if err != nil {
		return err
	}

	if command.OutputWriter == nil {
		command.OutputWriter = os.Stdout
	}

	_, err = fmt.Fprintf(command.OutputWriter, manifest)
	return err
}
