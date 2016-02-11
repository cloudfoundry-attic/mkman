package main_test

import (
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

const (
	executableTimeout = 5 * time.Second
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
		BeforeEach(func() {
			args = []string{"create-manifests"}
		})

		Context("with path to config", func() {
			var path string

			BeforeEach(func() {
				path = "/path/to/config"
				args = append(args, path)
			})

			It("exits successfully", func() {
				command := exec.Command(binPath, args...)
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session, executableTimeout).Should(gexec.Exit(0))
				Expect(session.Out).To(gbytes.Say("creating manifests from: %s", path))
			})
		})
	})
})
