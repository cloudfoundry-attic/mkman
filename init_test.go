package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var (
	binPath string
)

var _ = BeforeSuite(func() {
	By("Compiling binary")
	var err error
	binPath, err = gexec.Build("github.com/pivotal-cf-experimental/mkman", "-race")
	Expect(err).ShouldNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

func TestMain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "mkman executable test suite")
}
