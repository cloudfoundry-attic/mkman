package releasemakers_test

import (
	"github.com/cloudfoundry/mkman/releasemakers"

	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/gomega"
)

var _ = Describe("ReleaseMaker", func() {
	var releaseMaker releasemakers.ReleaseMaker
	var releasePath string

	BeforeEach(func() {
		releasePath = "/some/path/to/release"
		releaseMaker = releasemakers.NewCfReleaseMaker(releasePath)
	})

	It("returns a release stub", func() {
		cfRelease, err := releaseMaker.MakeRelease()

		Expect(err).NotTo(HaveOccurred())
		Expect(cfRelease.Name).To(Equal("cf"))
		Expect(cfRelease.URL).To(Equal("file:///some/path/to/release"))
		Expect(cfRelease.Version).To(Equal("create"))
	})
})
