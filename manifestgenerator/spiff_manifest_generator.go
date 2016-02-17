package manifestgenerator

import (
	"fmt"
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

func (g *SpiffManifestGenerator) GenerateManifest() (string, error) {
	stemcellStubPath, err := g.stemcellStubMaker.MakeStub()
	if err != nil {
		return "", err
	}

	releaseStubPath, err := g.releaseStubMaker.MakeStub()
	if err != nil {
		return "", err
	}

	stubPaths := append(g.stubPaths, stemcellStubPath, releaseStubPath)
	cmdArgs := append([]string{"aws"}, stubPaths...)

	generateManifestScriptPath := filepath.Join(g.cfPath, "scripts/generate_deployment_manifest")
	cmd := exec.Command(generateManifestScriptPath, cmdArgs...)

	outBytes, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", string(outBytes))
		return "", err
	}

	return string(outBytes), nil
}
