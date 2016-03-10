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

var _ = Describe("EtcdReleaseMaker", func() {
	var (
		releaseMaker        releasemakers.ReleaseMaker
		fakeTarballReader   *fakes.FakeTarballReader
		tarballErr          error
		tarballFileContents []byte
		etcdURL             string
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
		etcdURL = tmpDir + "/tarball.tgz"
	})

	AfterEach(func() {
		err := os.RemoveAll(tmpDir)
		Expect(err).NotTo(HaveOccurred())
	})

	JustBeforeEach(func() {
		fakeTarballReader.ReadFileReturns(tarballFileContents, tarballErr)
		releaseMaker = releasemakers.NewEtcdReleaseMaker(fakeTarballReader, etcdURL)
	})

	It("returns a release stub", func() {
		etcdRelease, err := releaseMaker.MakeRelease()

		Expect(err).NotTo(HaveOccurred())
		Expect(etcdRelease.Name).To(Equal("etcd"))
		Expect(etcdRelease.URL).To(Equal("file://" + etcdURL))
		Expect(etcdRelease.Version).To(Equal(version))
	})

	It("gets the version from etcd's release.MF", func() {
		releaseMaker.MakeRelease()
		Expect(fakeTarballReader.ReadFileCallCount()).To(Equal(1))
		Expect(fakeTarballReader.ReadFileArgsForCall(0)).To(Equal("./release.MF"))
	})

	Context("when the file does not exist", func() {
		BeforeEach(func() {
			etcdURL = "/non/existant/file.tgz"
		})

		It("should return an error", func() {
			release, err := releaseMaker.MakeRelease()
			Expect(release).To(BeNil())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(MatchRegexp("no such file or directory"))
		})
	})

	Context("when the path extension is not .tgz", func() {
		Context("when the path is a directory", func() {
			BeforeEach(func() {
				var err error
				version = "create"
				etcdURL, err = ioutil.TempDir("", "")
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns a release stub", func() {
				etcdRelease, err := releaseMaker.MakeRelease()

				Expect(err).NotTo(HaveOccurred())
				Expect(etcdRelease.Name).To(Equal("etcd"))
				Expect(etcdRelease.URL).To(Equal(etcdURL))
				Expect(etcdRelease.Version).To(Equal(version))
			})

			Context("with errors", func() {
				Context("when the directory does not exist", func() {
					BeforeEach(func() {
						etcdURL = "/non/existant/directory"
					})

					It("should return an error", func() {
						release, err := releaseMaker.MakeRelease()
						Expect(release).To(BeNil())
						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(MatchRegexp("no such file or directory"))
					})
				})
			})
		})

		Context("when the path is a file without a tgz extension", func() {
			BeforeEach(func() {
				var tmpFile *os.File
				tmpDir, err := ioutil.TempDir("", "")
				Expect(err).NotTo(HaveOccurred())
				tmpFile, err = ioutil.TempFile(tmpDir, "bad_extension.png")
				Expect(err).NotTo(HaveOccurred())
				etcdURL = tmpFile.Name()
			})

			It("returns an error", func() {
				release, err := releaseMaker.MakeRelease()
				Expect(release).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("unrecognized etcd URL"))
			})
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

	Context("when the release.MF file has invalid yaml", func() {
		BeforeEach(func() {
			tarballFileContents = []byte(fmt.Sprintf(`---
name: &
version: %s
`, version))
		})
		It("should return an error", func() {
			release, err := releaseMaker.MakeRelease()
			Expect(release).To(BeNil())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("yaml: line"))
		})
	})
})
