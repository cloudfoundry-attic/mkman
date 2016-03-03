package commands

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/cloudfoundry/mkman/config"
	"github.com/cloudfoundry/mkman/manifestgenerator"
	"github.com/cloudfoundry/mkman/stubmakers"
	"github.com/cloudfoundry/mkman/stubmakers/jobtemplates"
	releaseStubMakers "github.com/cloudfoundry/mkman/stubmakers/releases"
	"github.com/cloudfoundry/mkman/stubmakers/releases/releasemakers"
	stemcellStubMakers "github.com/cloudfoundry/mkman/stubmakers/stemcells"
	"github.com/cloudfoundry/mkman/tarball"

	"github.com/cloudfoundry/mkman/Godeps/_workspace/src/gopkg.in/yaml.v2"
)

type ManifestGenerator interface {
	GenerateManifest() (string, error)
}

type CreateManifestsCommand struct {
	OutputWriter io.Writer
	ConfigPath   string `long:"config" short:"c" required:"true" description:"Configuration file (required)"`
}

func (command *CreateManifestsCommand) Execute(args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("invalid additional arguments %v", args)
	}

	configFileContents, err := ioutil.ReadFile(command.ConfigPath)
	if err != nil {
		return err
	}

	config := config.Config{}
	err = yaml.Unmarshal(configFileContents, &config)
	if err != nil {
		return err
	}

	err = config.Validate()
	if err != nil {
		return err
	}

	stemcellTarballReader := tarball.NewTarballReader(config.StemcellPath)
	etcdTarballReader := tarball.NewTarballReader(config.EtcdPath)

	cfReleaseMaker := releasemakers.NewCfReleaseMaker(config.CFPath)
	etcdReleaseMaker := releasemakers.NewEtcdReleaseMaker(etcdTarballReader, config.EtcdPath)

	stemcellStubMaker := stemcellStubMakers.NewStemcellStubMaker(stemcellTarballReader, config.StemcellPath)
	releaseStubMaker := releaseStubMakers.NewReleaseStubMaker([]releasemakers.ReleaseMaker{
		cfReleaseMaker,
		etcdReleaseMaker,
	})
	jobTemplateStubMaker := jobtemplates.NewJobTemplateStubMaker()

	stubMakers := []stubmakers.StubMaker{
		stemcellStubMaker,
		releaseStubMaker,
		jobTemplateStubMaker,
	}
	manifestGenerator := manifestgenerator.NewSpiffManifestGenerator(stubMakers, config.StubPaths, config.CFPath)

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
