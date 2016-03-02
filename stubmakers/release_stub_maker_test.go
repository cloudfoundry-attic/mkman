package stubmakers_test

import (
	"fmt"
	"io/ioutil"

	"github.com/cloudfoundry/mkman/Godeps/_workspace/src/gopkg.in/yaml.v2"
	"github.com/cloudfoundry/mkman/releasemakers"

	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/gomega"
	releasefakes "github.com/cloudfoundry/mkman/releasemakers/fakes"
	"github.com/cloudfoundry/mkman/stubmakers"
)

var _ = Describe("ReleaseStubMaker", func() {
	var releaseStubMaker stubmakers.StubMaker
	var releasePath string
	var fakeReleaseMaker *releasefakes.FakeReleaseMaker

	BeforeEach(func() {
		fakeReleaseMaker = &releasefakes.FakeReleaseMaker{}

		releasePath = "/some/path/to/release"
		dummyRelease := releasemakers.Release{
			Name:    "dummy name",
			Version: "some version",
			URL:     "file://" + releasePath,
		}
		fakeReleaseMaker.MakeReleaseReturns(&dummyRelease, nil)

		releaseStubMaker = stubmakers.NewReleaseStubMaker([]releasemakers.ReleaseMaker{
			fakeReleaseMaker,
		})
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
		Expect(releaseStub.Releases[0].Name).To(Equal("dummy name"))
		Expect(releaseStub.Releases[0].Version).To(Equal("some version"))
		Expect(releaseStub.Releases[0].URL).To(Equal(fmt.Sprintf("file://%s", releasePath)))
	})

	Context("when a release maker returns an error", func() {
		BeforeEach(func() {
			fakeReleaseMaker.MakeReleaseReturns(nil, fmt.Errorf("some error"))
		})

		It("returns the error", func() {
			stubPath, err := releaseStubMaker.MakeStub()
			Expect(stubPath).To(BeEmpty())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("some error"))
		})
	})
})
