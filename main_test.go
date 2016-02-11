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

			cfReleasePath := "/Users/robdimsdale/workspace/cf-release"
			stemcellPath := "/Users/robdimsdale/Downloads/bosh-stemcell-2776-warden-boshlite-ubuntu-trusty-go_agent.tgz"
			stubsPath := "/Users/robdimsdale/workspace/cf-deployment/spec/assets/stub.yml"

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

		FIt("creates manifest without error", func() {
			command := exec.Command(binPath, args...)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, executableTimeout).Should(gexec.Exit(0))
			Expect(session.Out).To(gbytes.Say("creating manifests from: %s", configPath))
		})

		XContext("when path is not provided", func() {

		})

		XContext("when path is not valid", func() {

		})
	})
})
