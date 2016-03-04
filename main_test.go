package main_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/gomega"
	"github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/gomega/gbytes"
	"github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/gomega/gexec"
	"github.com/cloudfoundry/mkman/testhelpers"
)

const (
	executableTimeout = 30 * time.Second
)

var _ = Describe("Executing binary", func() {
	var (
		args []string
	)

	Context("with no arguments or flags", func() {
		BeforeEach(func() {
			args = []string{}
		})

		It("exits with error", func() {
			command := exec.Command(binPath, args...)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, executableTimeout).Should(gexec.Exit(1))
			Expect(session.Err).To(gbytes.Say("error"))
		})
	})

	Context("when --version is provided", func() {
		BeforeEach(func() {
			args = []string{"--version"}
		})

		It("prints version with --version", func() {
			command := exec.Command(binPath, args...)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, executableTimeout).Should(gexec.Exit(0))
			Expect(session.Out).To(gbytes.Say("mkman 0.0.1"))
		})
	})

	Context("when an invalid argument is provided", func() {
		BeforeEach(func() {
			args = []string{"--not-a-valid-arg"}
		})

		It("exits with error", func() {
			command := exec.Command(binPath, args...)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, executableTimeout).Should(gexec.Exit(1))
			Expect(session.Err).To(gbytes.Say("error"))
			Expect(session.Err).To(gbytes.Say("unknown flag"))
		})
	})

	Context("when an invalid command is provided", func() {
		BeforeEach(func() {
			args = []string{"not-a-valid-command"}
		})

		It("exits with error", func() {
			command := exec.Command(binPath, args...)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, executableTimeout).Should(gexec.Exit(1))
			Expect(session.Err).To(gbytes.Say("error"))
			Expect(session.Err).To(gbytes.Say("Unknown command"))
		})
	})

	Describe("create-manifests", func() {
		var (
			setup *testhelpers.TestSetup
		)

		BeforeEach(func() {
			setup = testhelpers.SetupTest()

			args = []string{"create-manifests"}
			args = append(args, "-c", setup.ConfigPath)

			setup.WriteConfig()
		})

		AfterEach(func() {
			err := os.RemoveAll(setup.TempDirPath)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("creates manifest without error", func() {
			command := exec.Command(binPath, args...)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, executableTimeout).Should(gexec.Exit(0))
			manifest := session.Out.Contents()

			manifestPath := filepath.Join(setup.TempDirPath, "output_manifest.yml")
			err = ioutil.WriteFile(manifestPath, manifest, os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			diffCommand := exec.Command("diff", "-C3", manifestPath, setup.ExampleManifestPath)
			diffSession, err := gexec.Start(diffCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(diffSession).Should(gexec.Exit())
			Expect(diffSession.Out.Contents()).To(BeEmpty())
			Expect(diffSession.Err.Contents()).To(BeEmpty())
		})

		Context("when the required config flag is not provided", func() {
			BeforeEach(func() {
				args = []string{"create-manifests"}
			})

			It("exits with error", func() {
				command := exec.Command(binPath, args...)
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session, executableTimeout).Should(gexec.Exit())
				Expect(session.Err).To(gbytes.Say("required flag"))
				Expect(session.Err).To(gbytes.Say("config"))
				Expect(session.Err).To(gbytes.Say("not specified"))
			})
		})

		Describe("--help", func() {
			BeforeEach(func() {
				args = append(args, "--help")
			})

			It("displays required flags", func() {
				command := exec.Command(binPath, args...)
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session, executableTimeout).Should(gexec.Exit())
				Expect(session.Err).To(gbytes.Say("config.*required"))
			})
		})
	})
})
