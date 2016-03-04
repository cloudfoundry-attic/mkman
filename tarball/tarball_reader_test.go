package tarball_test

import (
	"path/filepath"

	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/gomega"
	"github.com/cloudfoundry/mkman/tarball"
	"github.com/cloudfoundry/mkman/testhelpers"
)

var _ = Describe("TarballReader", func() {
	var tarballReader tarball.TarballReader
	var tarballPath string
	var fixturesDir string

	BeforeEach(func() {
		By("Locating fixtures dir")
		testDir := testhelpers.GetDirOfCurrentFile()
		fixturesDir = filepath.Join(testDir, "..", "fixtures")

		tarballPath = filepath.Join(fixturesDir, "no-image-stemcell.tgz")
	})

	JustBeforeEach(func() {
		tarballReader = tarball.NewTarballReader(tarballPath)
	})

	It("returns the contents of the specified file in the tarball", func() {
		fileBytes, err := tarballReader.ReadFile("stemcell.MF")
		Expect(err).NotTo(HaveOccurred())

		Expect(string(fileBytes)).To(ContainSubstring("name: bosh-warden-boshlite-ubuntu-trusty-go_agent"))
	})

	Context("when the provided tarball path cannot be opened", func() {
		BeforeEach(func() {
			tarballPath = "/some/bogus/path"
		})

		It("returns an error", func() {
			fileBytes, err := tarballReader.ReadFile("stemcell.MF")
			Expect(err).To(HaveOccurred())
			Expect(fileBytes).To(BeNil())
		})
	})

	Context("when the provided tarball path is not gzipped", func() {
		BeforeEach(func() {
			tarballPath = filepath.Join(fixturesDir, "stub.yml")
		})

		It("returns an error", func() {
			fileBytes, err := tarballReader.ReadFile("stemcell.MF")
			Expect(err).To(HaveOccurred())
			Expect(fileBytes).To(BeNil())
		})
	})

	Context("when the requested file is not present in the tarball", func() {
		It("returns an error", func() {
			fileBytes, err := tarballReader.ReadFile("non-existent-file")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("filename 'non-existent-file' not found in tarPath"))
			Expect(fileBytes).To(BeNil())
		})
	})
})
