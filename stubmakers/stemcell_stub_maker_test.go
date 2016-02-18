package stubmakers_test

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf-experimental/mkman/stubmakers"
	"github.com/pivotal-cf-experimental/mkman/tarball/fakes"
)

var _ = Describe("StemcellStubMaker", func() {
	var (
		stemcellStubMaker   stubmakers.StubMaker
		stemcellURL         string
		fakeTarballReader   *fakes.FakeTarballReader
		tarballErr          error
		tarballFileContents []byte
	)

	BeforeEach(func() {
		stemcellURL = "/path/to/tarball.tgz"
		fakeTarballReader = &fakes.FakeTarballReader{}

		tarballErr = nil

		tarballFileContents = []byte(`---
name: some-name
version: some-version
`)
	})

	JustBeforeEach(func() {
		fakeTarballReader.ReadFileReturns(tarballFileContents, tarballErr)
		stemcellStubMaker = stubmakers.NewStemcellStubMaker(fakeTarballReader, stemcellURL)
	})

	It("returns a path to a stemcell stub", func() {
		stubPath, err := stemcellStubMaker.MakeStub()
		Expect(err).NotTo(HaveOccurred())

		stemcellContents, err := ioutil.ReadFile(stubPath)
		Expect(err).NotTo(HaveOccurred())

		var stemcellStub struct {
			Meta struct {
				Stemcell struct {
					Name    string `yaml:"name"`
					Version string `yaml:"version"`
					URL     string `yaml:"url"`
				} `yaml:"stemcell"`
			} `yaml:"meta"`
		}

		err = yaml.Unmarshal(stemcellContents, &stemcellStub)
		Expect(err).NotTo(HaveOccurred())

		Expect(stemcellStub.Meta.Stemcell.Name).To(Equal("some-name"))
		Expect(stemcellStub.Meta.Stemcell.Version).To(Equal("some-version"))
		Expect(stemcellStub.Meta.Stemcell.URL).To(Equal(fmt.Sprintf("file://%s", stemcellURL)))
	})

	Context("when tarballReader returns an error", func() {
		BeforeEach(func() {
			tarballErr = fmt.Errorf("tarball error")
		})

		It("forwards the error", func() {
			stubPath, err := stemcellStubMaker.MakeStub()
			Expect(stubPath).To(BeEmpty())
			Expect(err).To(Equal(tarballErr))
		})
	})

	Context("when the provided stemcellURL is in an unrecognized format", func() {
		BeforeEach(func() {
			stemcellURL = "/non/tgz/file"
		})

		It("returns an error", func() {
			stubPath, err := stemcellStubMaker.MakeStub()
			Expect(stubPath).To(BeEmpty())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unrecognized stemcell URL"))
		})
	})
})
