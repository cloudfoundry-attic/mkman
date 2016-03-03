package releasemakers_test

import (
	"fmt"

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
	)

	BeforeEach(func() {
		version = "some-version"

		etcdURL = "/path/to/tarball.tgz"
		fakeTarballReader = &fakes.FakeTarballReader{}

		tarballErr = nil

		tarballFileContents = []byte(fmt.Sprintf(`---
name: some-name
version: %s
`, version))
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
	})

	Context("when the path extension is not .tgz", func() {
		BeforeEach(func() {
			etcdURL = "/path/to/tarball"
		})

		It("returns an error", func() {
			release, err := releaseMaker.MakeRelease()
			Expect(release).To(BeNil())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unrecognized etcd URL"))
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
