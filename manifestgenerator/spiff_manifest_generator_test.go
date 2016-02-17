package manifestgenerator_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/pivotal-cf-experimental/mkman/manifestgenerator"
	"github.com/pivotal-cf-experimental/mkman/manifestgenerator/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SpiffManifestGenerator", func() {
	var cfReleasePath string
	var stubPath string
	var stemcellStubMaker *fakes.FakeStubMaker
	var releaseStubMaker *fakes.FakeStubMaker
	var manifestGenerator *manifestgenerator.SpiffManifestGenerator
	var tempDirPath string

	BeforeEach(func() {
		cfReleasePath = "/Users/pivotal/workspace/cf-release"
		stubPath = filepath.Join("../fixtures", "stub.yml")

		var err error
		tempDirPath, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		stemcellStubMaker = &fakes.FakeStubMaker{}
		releaseStubMaker = &fakes.FakeStubMaker{}

		stemcellStubPath := filepath.Join(tempDirPath, "stemcell.yml")
		stemcellStubMaker.MakeStubReturns(stemcellStubPath, nil)

		releaseStubPath := filepath.Join(tempDirPath, "release.yml")
		releaseStubMaker.MakeStubReturns(releaseStubPath, nil)

		stemcellStubPathContents := fmt.Sprintf(`---
meta:
  stemcell:
    name: stemcell-name
    version: 123
    url: https://bosh.io/stemcell-name-123
`)

		releaseStubPathContents := fmt.Sprintf(`---
releases:
- name: release-name
  version: 123
  url: https://bosh.io/release-name-123
`)

		err = ioutil.WriteFile(stemcellStubPath, []byte(stemcellStubPathContents), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
		err = ioutil.WriteFile(releaseStubPath, []byte(releaseStubPathContents), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
	})

	JustBeforeEach(func() {
		manifestGenerator = manifestgenerator.NewSpiffManifestGenerator(
			stemcellStubMaker,
			releaseStubMaker,
			[]string{stubPath},
			cfReleasePath,
		)
	})

	AfterEach(func() {
		err := os.RemoveAll(tempDirPath)
		Expect(err).ShouldNot(HaveOccurred())
	})

	type outputManifest struct {
		Releases []struct {
			Name    string `yaml:"name"`
			Version string `yaml:"version"`
			URL     string `yaml:"url"`
		} `yaml:"releases"`
	}

	It("includes the information about releases", func() {
		manifestStr, err := manifestGenerator.GenerateManifest()
		Expect(err).NotTo(HaveOccurred())

		var manifest outputManifest
		err = yaml.Unmarshal([]byte(manifestStr), &manifest)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(manifest.Releases)).To(BeNumerically(">=", 1))
		Expect(manifest.Releases[0].Name).To(Equal("release-name"))
	})
})
