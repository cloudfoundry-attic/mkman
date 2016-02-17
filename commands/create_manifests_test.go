package commands_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf-experimental/mkman/commands"
)

var _ = Describe("CreateManifestsCommand", func() {
	var args []string
	var cmd commands.CreateManifestsCommand
	var configPathContents string
	var cfReleasePath string
	var stemcellPath string
	var stubPath string
	var configPath string
	var tempDirPath string

	BeforeEach(func() {
		var err error
		tempDirPath, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		configPathContents = fmt.Sprintf(`
cf: %s
stemcell: %s
stubs:
- %s
`,
			cfReleasePath,
			stemcellPath,
			stubPath,
		)
		configPath = filepath.Join(tempDirPath, "config.yml")
		args = []string{configPath}

		cmd = commands.CreateManifestsCommand{}
	})

	AfterEach(func() {
		err := os.RemoveAll(tempDirPath)
		Expect(err).ShouldNot(HaveOccurred())
	})

	JustBeforeEach(func() {
		err := ioutil.WriteFile(configPath, []byte(configPathContents), os.ModePerm)
		Expect(err).ShouldNot(HaveOccurred())
	})

	Context("when path is not provided", func() {
		BeforeEach(func() {
			args = []string{}
		})

		It("returns an error", func() {
			err := cmd.Execute(args)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("create-manifests requires PATH_TO_CONFIG"))
		})
	})

	Context("when path is not valid", func() {
		BeforeEach(func() {
			args = []string{"/bad/path"}
		})

		It("returns an error", func() {
			err := cmd.Execute(args)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("open /bad/path: no such file or directory"))
		})
	})

	Context("when the path points to an invalid config", func() {
		BeforeEach(func() {
			configPathContents = "{{"
		})

		It("returns an error", func() {
			err := cmd.Execute(args)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("yaml: line 1: did not find expected node content"))
		})
	})

})
