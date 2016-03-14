package commands_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/gomega"
	"github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/gomega/gexec"
	"github.com/cloudfoundry/mkman/commands"
	"github.com/cloudfoundry/mkman/testhelpers"
)

var _ = Describe("CreateManifestsCommand", func() {
	var (
		args           []string
		cmd            commands.CreateManifestsCommand
		outputManifest *bytes.Buffer
		setup          *testhelpers.TestSetup
	)

	BeforeEach(func() {
		setup = testhelpers.SetupTest()

		args = []string{}
		outputManifest = &bytes.Buffer{}

		cmd = commands.CreateManifestsCommand{
			OutputWriter: outputManifest,
			ConfigPath:   setup.ConfigPath,
		}
	})

	AfterEach(func() {
		err := os.RemoveAll(setup.TempDirPath)
		Expect(err).ShouldNot(HaveOccurred())
	})

	JustBeforeEach(func() {
		setup.WriteConfig()
	})

	It("creates manifest without error", func() {
		err := cmd.Execute(args)
		Expect(err).NotTo(HaveOccurred())

		manifestPath := filepath.Join(setup.TempDirPath, "output_manifest.yml")
		err = ioutil.WriteFile(manifestPath, outputManifest.Bytes(), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		diffCommand := exec.Command("diff", "-C3", manifestPath, setup.ExampleManifestPath)
		diffSession, err := gexec.Start(diffCommand, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(diffSession).Should(gexec.Exit())
		Expect(diffSession.Out.Contents()).To(BeEmpty())
		Expect(diffSession.Err.Contents()).To(BeEmpty())
	})

	Context("when path is not valid", func() {
		BeforeEach(func() {
			cmd.ConfigPath = "/bad/path"
		})

		It("returns an error", func() {
			err := cmd.Execute(args)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("open /bad/path: no such file or directory"))
		})
	})

	Context("when the path points to an invalid config", func() {
		BeforeEach(func() {
			setup.ConfigContents = "{{"
		})

		It("returns an error", func() {
			err := cmd.Execute(args)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("yaml: line 1: did not find expected node content"))
		})
	})

	Context("when the manifest generator returns an error", func() {
		BeforeEach(func() {
			// force an error by giving a bad cfReleasePath
			setup.ConfigContents = `
        cf: /not/a/valid/path
        stemcell:
        stubs:
        `
		})

		It("returns an error", func() {
			err := cmd.Execute(args)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("/not/a/valid/path"))
		})
	})

	Context("when the config file has empty values", func() {
		BeforeEach(func() {
			setup.ConfigContents = fmt.Sprintf(`
        cf: 
        stemcell: 
        stubs:
        - 
        `)
		})

		It("returns an error", func() {
			err := cmd.Execute(args)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("value must be non-empty"))
		})
	})

	Context("when writing the output fails", func() {
		BeforeEach(func() {
			cmd.OutputWriter = &testhelpers.AlwaysErrorWriter{}
		})

		It("forwards the error", func() {
			err := cmd.Execute(args)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("writer error"))
		})
	})

	Context("when there are any extra arguments", func() {
		BeforeEach(func() {
			args = []string{"extra-foo-arg1", "extra-foo-arg2"}
		})

		It("should return an error", func() {
			err := cmd.Execute(args)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid additional arguments"))
			Expect(err.Error()).To(ContainSubstring("extra-foo-arg1"))
			Expect(err.Error()).To(ContainSubstring("extra-foo-arg2"))
		})
	})
})
