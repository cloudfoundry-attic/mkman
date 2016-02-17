package commands

import (
	"fmt"
	"io/ioutil"

	"github.com/pivotal-cf-experimental/mkman/config"
	"github.com/pivotal-cf-experimental/mkman/manifestgenerator"
	"github.com/pivotal-cf-experimental/mkman/stubmakers"

	"gopkg.in/yaml.v2"
)

//go:generate counterfeiter . ManifestGenerator
type ManifestGenerator interface {
	GenerateManifest() error
}

type CreateManifestsCommand struct{}

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

	stemcellStubMaker := stubmakers.NewStemcellStubMaker(config.StemcellPath)
	releaseStubMaker := stubmakers.NewReleaseStubMaker(config.CFPath)
	manifestGenerator := manifestgenerator.NewSpiffManifestGenerator(stemcellStubMaker, releaseStubMaker, config.StubPaths, config.CFPath)
	err = manifestGenerator.GenerateManifest()
	return err
}
