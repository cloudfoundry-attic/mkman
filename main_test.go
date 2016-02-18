package main_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	. "github.com/pivotal-cf-experimental/mkman/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/pivotal-cf-experimental/mkman/Godeps/_workspace/src/github.com/onsi/gomega"
	"github.com/pivotal-cf-experimental/mkman/Godeps/_workspace/src/github.com/onsi/gomega/gbytes"
	"github.com/pivotal-cf-experimental/mkman/Godeps/_workspace/src/github.com/onsi/gomega/gexec"
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
			tempDirPath         string
			configPath          string
			fixturesDir         string
			exampleManifestPath string
		)

		BeforeEach(func() {
			By("Locating fixtures dir")
			testDir := getDirOfCurrentFile()
			fixturesDir = filepath.Join(testDir, "fixtures")

			By("Ensuring $CF_RELEASE_DIR is set")
			cfReleasePath := os.Getenv("CF_RELEASE_DIR")
			Expect(cfReleasePath).NotTo(BeEmpty(), "$CF_RELEASE_DIR must be provided")

			var err error
			tempDirPath, err = ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())

			By("Creating manifest template")
			manifestTemplatePath := filepath.Join(fixturesDir, "manifest.yml.template")
			templateContents, err := ioutil.ReadFile(manifestTemplatePath)
			Expect(err).NotTo(HaveOccurred())
			templateContents2 := strings.Replace(string(templateContents), "$CF_RELEASE_DIR", cfReleasePath, -1)

			stemcellPath := filepath.Join(fixturesDir, "no-image-stemcell.tgz")
			templateContents3 := strings.Replace(string(templateContents2), "$STEMCELL_PATH", stemcellPath, -1)

			exampleManifestPath = filepath.Join(tempDirPath, "manifest.yml")
			err = ioutil.WriteFile(exampleManifestPath, []byte(templateContents3), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			By("Writing config paths")
			args = []string{"create-manifests"}

			configPath = filepath.Join(tempDirPath, "config.yml")

			stubPath := filepath.Join(fixturesDir, "stub.yml")

			configPathContents := fmt.Sprintf(`
{
  "cf": "%s",
  "stemcell": "%s",
  "stubs": ["%s"]
}
`,
				cfReleasePath,
				stemcellPath,
				stubPath,
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
			manifest := session.Out.Contents()

			manifestPath := filepath.Join(tempDirPath, "output_manifest.yml")
			err = ioutil.WriteFile(manifestPath, manifest, os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			diffCommand := exec.Command("diff", "-C3", manifestPath, exampleManifestPath)
			diffSession, err := gexec.Start(diffCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(diffSession).Should(gexec.Exit())
			Expect(diffSession.Out.Contents()).To(BeEmpty())
			Expect(diffSession.Err.Contents()).To(BeEmpty())
		})
	})
})

func getDirOfCurrentFile() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}
