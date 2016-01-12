package main_test

import (
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

const (
	executableTimeout = 5 * time.Second
)

var _ = Describe("Executing binary", func() {
	It("runs and exits without error", func() {
		command := exec.Command(binPath)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session, executableTimeout).Should(gexec.Exit(0))
	})
})
