package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pivotal-cf-experimental/mkman/config"
	"github.com/pivotal-cf-experimental/mkman/stemcell"

	"gopkg.in/yaml.v2"
)

type CreateManifestsCommand struct {
}

func (command *CreateManifestsCommand) Execute(args []string) error {
	mydir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	outputDirPath := filepath.Join(mydir, "outputs")
	manifestsDirPath := filepath.Join(outputDirPath, "manifests")
	fmt.Printf("creating manifests from: %s\n", args[0])

	configFileContents, err := ioutil.ReadFile(args[0])
	if err != nil {
		panic(err)
	}

	config := config.Config{}
	err = yaml.Unmarshal(configFileContents, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	intermediateDir, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}

	stemcellStubContents, err := stemcell.StubFromTar(config.StemcellPath)
	if err != nil {
		panic(err)
	}

	stemcellStubPath := filepath.Join(intermediateDir, "stemcell.yml")
	err = ioutil.WriteFile(stemcellStubPath, []byte(stemcellStubContents), os.ModePerm)
	if err != nil {
		panic(err)
	}

	fmt.Printf("@@@ Debug stemcellStubPath: %s\n", stemcellStubPath)

	var stubPaths []string
	var cmdArgs []string

	stubPaths = append(stubPaths, stemcellStubPath)
	stubPaths = append(stubPaths, config.StubPaths...)

	cmdArgs = append(cmdArgs, "aws")
	cmdArgs = append(cmdArgs, stubPaths...)

	generateManifestScriptPath := filepath.Join(config.CFPath, "scripts/generate_deployment_manifest")
	cmd := exec.Command(generateManifestScriptPath, cmdArgs...)

	fmt.Printf("@@@ debug cmd: %+v\n", cmd)

	outBytes, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("---\n%s\n", string(outBytes))
		panic(err)
	}

	err = os.MkdirAll(manifestsDirPath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	manifestFilePath := filepath.Join(manifestsDirPath, "cf.yml")

	fmt.Printf("writing manifest to: %s\n", manifestFilePath)
	err = ioutil.WriteFile(manifestFilePath, outBytes, os.ModePerm)
	if err != nil {
		panic(err)
	}

	return nil
}
