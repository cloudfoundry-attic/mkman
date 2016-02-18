package stubmakers_test

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/cloudfoundry/mkman/stubmakers"
)

var _ = Describe("ReleaseStubMaker", func() {
	var releaseStubMaker stubmakers.StubMaker
	var releasePath string

	BeforeEach(func() {
		releasePath = "/some/path/to/release"
		releaseStubMaker = stubmakers.NewReleaseStubMaker(releasePath)
	})

	It("writes a release stub and returns the path", func() {
		stubPath, err := releaseStubMaker.MakeStub()
		Expect(err).NotTo(HaveOccurred())

		fileBytes, err := ioutil.ReadFile(stubPath)
		Expect(err).NotTo(HaveOccurred())

		var releaseStub struct {
			Releases []struct {
				Name    string `yaml:"name"`
				Version string `yaml:"version"`
				URL     string `yaml:"url"`
			} `yaml:"releases"`
		}
		err = yaml.Unmarshal(fileBytes, &releaseStub)
		Expect(err).NotTo(HaveOccurred())

		Expect(releaseStub.Releases).To(HaveLen(1))
		Expect(releaseStub.Releases[0].Name).To(Equal("cf"))
		Expect(releaseStub.Releases[0].Version).To(Equal("create"))
		Expect(releaseStub.Releases[0].URL).To(Equal(fmt.Sprintf("file://%s", releasePath)))
	})
})
