package testhelpers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/gomega"
)

type TestSetup struct {
	fixturesDir   string
	cfReleasePath string
	stemcellPath  string
	etcdPath      string

	TempDirPath         string
	ExampleManifestPath string
	ConfigPath          string
	ConfigContents      string
}

func SetupTest() *TestSetup {
	fixturesDir := getFixturesDir()
	cfReleasePath := getCfReleasePath()
	stemcellPath := filepath.Join(fixturesDir, "no-image-stemcell.tgz")
	etcdPath := filepath.Join(fixturesDir, "etcd-release.tgz")
	tempDir := setupTempDir()

	By("Writing config paths")
	configPath := filepath.Join(tempDir, "config.yml")
	stubPath := filepath.Join(fixturesDir, "stub.yml")
	configContents := fmt.Sprintf(`
{
  "cf": "%s",
  "stemcell": "%s",
	"etcd": "%s",
  "stubs": ["%s"]
}
`,
		cfReleasePath,
		stemcellPath,
		etcdPath,
		stubPath,
	)

	setup := &TestSetup{
		fixturesDir:   fixturesDir,
		cfReleasePath: cfReleasePath,
		stemcellPath:  stemcellPath,
		etcdPath:      etcdPath,

		TempDirPath:    tempDir,
		ConfigPath:     configPath,
		ConfigContents: configContents,
	}
	setup.createExampleManifest()
	return setup
}

func (setup *TestSetup) WriteConfig() {
	err := ioutil.WriteFile(setup.ConfigPath, []byte(setup.ConfigContents), os.ModePerm)
	Expect(err).ShouldNot(HaveOccurred())
}

func (setup *TestSetup) createExampleManifest() {
	By("Creating manifest template")
	manifestTemplatePath := filepath.Join(setup.fixturesDir, "manifest.yml.template")
	templateBytes, err := ioutil.ReadFile(manifestTemplatePath)
	templateContents := string(templateBytes)
	Expect(err).NotTo(HaveOccurred())

	templateContents = strings.Replace(string(templateContents), "$CF_RELEASE_DIR", setup.cfReleasePath, -1)
	templateContents = strings.Replace(string(templateContents), "$STEMCELL_PATH", setup.stemcellPath, -1)
	templateContents = strings.Replace(string(templateContents), "$ETCD_RELEASE_PATH", setup.etcdPath, -1)

	exampleManifestPath := filepath.Join(setup.TempDirPath, "manifest.yml")
	err = ioutil.WriteFile(exampleManifestPath, []byte(templateContents), os.ModePerm)
	Expect(err).NotTo(HaveOccurred())

	setup.ExampleManifestPath = exampleManifestPath
}

func getCfReleasePath() string {
	By("Ensuring $CF_RELEASE_DIR is set")
	cfReleasePath := os.Getenv("CF_RELEASE_DIR")
	Expect(cfReleasePath).NotTo(BeEmpty(), "$CF_RELEASE_DIR must be provided")
	return cfReleasePath
}

func setupTempDir() string {
	tempDirPath, err := ioutil.TempDir("", "")
	Expect(err).NotTo(HaveOccurred())

	return tempDirPath
}

func getFixturesDir() string {
	By("Locating fixtures dir")
	testDir := GetDirOfCurrentFile()
	fixturesDir := filepath.Join(testDir, "../fixtures")
	return fixturesDir
}
