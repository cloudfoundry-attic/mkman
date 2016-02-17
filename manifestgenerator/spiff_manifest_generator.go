package manifestgenerator

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

//go:generate counterfeiter . StubMaker
type StubMaker interface {
	MakeStub() (string, error)
}

type SpiffManifestGenerator struct {
	stemcellStubMaker StubMaker
	releaseStubMaker  StubMaker
	stubPaths         []string
	cfPath            string
}

func NewSpiffManifestGenerator(stemcellStubMaker, releaseStubMaker StubMaker, stubPaths []string, cfPath string) *SpiffManifestGenerator {
	return &SpiffManifestGenerator{
		stemcellStubMaker: stemcellStubMaker,
		releaseStubMaker:  releaseStubMaker,
		stubPaths:         stubPaths,
		cfPath:            cfPath,
	}
}

func (g *SpiffManifestGenerator) GenerateManifest() error {
	stemcellStubPath, err := g.stemcellStubMaker.MakeStub()
	if err != nil {
		panic(err)
	}

	releaseStubPath, err := g.releaseStubMaker.MakeStub()
	if err != nil {
		panic(err)
	}

	stubPaths := append(g.stubPaths, stemcellStubPath, releaseStubPath)

	var cmdArgs []string
	cmdArgs = append(cmdArgs, "aws")
	cmdArgs = append(cmdArgs, stubPaths...)

	generateManifestScriptPath := filepath.Join(g.cfPath, "scripts/generate_deployment_manifest")
	cmd := exec.Command(generateManifestScriptPath, cmdArgs...)

	outBytes, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("---\n%s\n", string(outBytes))
		panic(err)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		// We cannot test this because it is too hard to get Getwd to return error
		return err
	}

	outputDirPath := filepath.Join(currentDir, "outputs")
	manifestsDirPath := filepath.Join(outputDirPath, "manifests")

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
