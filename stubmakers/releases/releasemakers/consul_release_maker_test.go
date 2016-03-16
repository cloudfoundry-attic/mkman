package releasemakers_test

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cloudfoundry/mkman/stubmakers/releases/releasemakers"
	"github.com/cloudfoundry/mkman/tarball/fakes"

	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/gomega"
)

var _ = Describe("ConsulReleaseMaker", func() {
	var (
		releaseMaker        releasemakers.ReleaseMaker
		fakeTarballReader   *fakes.FakeTarballReader
		tarballErr          error
		tarballFileContents []byte
		consulURL           string
		version             string
		tmpDir              string
	)

	BeforeEach(func() {
		version = "some-version"

		var err error
		tmpDir, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		fakeTarballReader = &fakes.FakeTarballReader{}

		tarballErr = nil

		tarballFileContents = []byte(fmt.Sprintf(`---
name: some-name
version: %s
`, version))

		err = ioutil.WriteFile(tmpDir+"/tarball.tgz", []byte{}, os.ModeTemporary)
		Expect(err).NotTo(HaveOccurred())
		consulURL = tmpDir + "/tarball.tgz"
	})

	JustBeforeEach(func() {
		fakeTarballReader.ReadFileReturns(tarballFileContents, tarballErr)
		releaseMaker = releasemakers.NewConsulReleaseMaker(fakeTarballReader, consulURL)
	})

	AfterEach(func() {
		err := os.RemoveAll(tmpDir)
		Expect(err).NotTo(HaveOccurred())
	})

	It("returns a release stub", func() {
		consulRelease, err := releaseMaker.MakeRelease()

		Expect(err).NotTo(HaveOccurred())
		Expect(consulRelease.Name).To(Equal("consul"))
		Expect(consulRelease.URL).To(Equal("file://" + consulURL))
		Expect(consulRelease.Version).To(Equal(version))
	})

	It("gets the version from consul's release.MF", func() {
		releaseMaker.MakeRelease()
		Expect(fakeTarballReader.ReadFileCallCount()).To(Equal(1))
		Expect(fakeTarballReader.ReadFileArgsForCall(0)).To(Equal("./release.MF"))
	})

	Context("when the path extension is not .tgz", func() {
		BeforeEach(func() {
			var tmpFile *os.File
			tmpDir, err := ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())
			tmpFile, err = ioutil.TempFile(tmpDir, "bad_extension.png")
			Expect(err).NotTo(HaveOccurred())
			consulURL = tmpFile.Name()
		})

		It("returns an error", func() {
			release, err := releaseMaker.MakeRelease()
			Expect(release).To(BeNil())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unrecognized consul URL"))
		})
	})

	Context("when the tarball reader returns an error", func() {
		BeforeEach(func() {
			tarballFileContents = nil
			tarballErr = fmt.Errorf("reading tarball failed")
		})

		It("forwards the error", func() {
			release, err := releaseMaker.MakeRelease()
			Expect(release).To(BeNil())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("reading tarball failed"))
		})
	})
})
