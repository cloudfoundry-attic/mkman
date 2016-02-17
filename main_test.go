package main_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
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
			tempDirPath string
			configPath  string
		)

		BeforeEach(func() {
			args = []string{"create-manifests"}

			var err error
			tempDirPath, err = ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())

			configPath = filepath.Join(tempDirPath, "config.json")

			stemcellPath := filepath.Join(fixturesDir, "no-image-stemcell.tgz")
			stubsPath := filepath.Join(fixturesDir, "stub.yml")

			configPathContents := fmt.Sprintf(`
{
  "cf": "%s",
  "stemcell": "%s",
  "stubs": ["%s"]
}
`,
				cfReleasePath,
				stemcellPath,
				stubsPath,
			)

			err = ioutil.WriteFile(configPath, []byte(configPathContents), os.ModePerm)
			Expect(err).ShouldNot(HaveOccurred())

			args = append(args, configPath)
		})

		AfterEach(func() {
			err := os.RemoveAll(tempDirPath)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("creates manifest without error", func() {
			command := exec.Command(binPath, args...)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, executableTimeout).Should(gexec.Exit(0))

			expectedManifestPath := filepath.Join(fixturesDir, "manifest.yml")

			cwd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			outputManifestPath := filepath.Join(
				cwd,
				"outputs",
				"manifests",
				"cf.yml",
			)

			diffCommand := exec.Command("diff", "-C3", outputManifestPath, expectedManifestPath)
			diffSession, err := gexec.Start(diffCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(diffSession).Should(gexec.Exit())
			Expect(diffSession.Out.Contents()).To(BeEmpty())
			Expect(diffSession.Err.Contents()).To(BeEmpty())
		})

		Context("when path is not provided", func() {
			BeforeEach(func() {
				args = []string{args[0]}
			})

			It("exits with error", func() {
				command := exec.Command(binPath, args...)
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session, executableTimeout).Should(gexec.Exit(1))
				Expect(session.Err).To(gbytes.Say("error: create-manifests requires PATH_TO_CONFIG"))
			})
		})

		Context("when path is not valid", func() {
			BeforeEach(func() {
				args[1] = "/bad/path"
			})

			It("exits with error", func() {
				command := exec.Command(binPath, args...)
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session, executableTimeout).Should(gexec.Exit(1))
				Expect(session.Err).To(gbytes.Say("error: open /bad/path: no such file or directory"))
			})
		})
	})
})
