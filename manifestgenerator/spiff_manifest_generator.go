package manifestgenerator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cloudfoundry/mkman/stubmakers"
)

type SpiffManifestGenerator struct {
	// stemcellStubMaker stubmakers.StubMaker
	// releaseStubMaker  stubmakers.StubMaker
	stubMakers []stubmakers.StubMaker
	stubPaths  []string
	cfPath     string
}

func NewSpiffManifestGenerator(stubMakers []stubmakers.StubMaker, stubPaths []string, cfPath string) *SpiffManifestGenerator {
	return &SpiffManifestGenerator{
		stubMakers: stubMakers,
		stubPaths:  stubPaths,
		cfPath:     cfPath,
	}
}

func (g *SpiffManifestGenerator) GenerateManifest() (string, error) {
	var stubPaths []string
	stubPaths = append(stubPaths, g.stubPaths...)

	for _, stubMaker := range g.stubMakers {
		stubPath, err := stubMaker.MakeStub()
		if err != nil {
			return "", err
		}
		stubPaths = append(stubPaths, stubPath)
	}

	cmdArgs := append([]string{"aws"}, stubPaths...)

	generateManifestScriptPath := filepath.Join(g.cfPath, "scripts/generate_deployment_manifest")
	cmd := exec.Command(generateManifestScriptPath, cmdArgs...)

	outBytes, err := cmd.CombinedOutput()
	if err != nil {
		// Don't print spurious empty output
		if len(outBytes) > 0 {
			fmt.Fprintf(os.Stderr, "%s\n", string(outBytes))
		}
		return "", err
	}

	return string(outBytes), nil
}
