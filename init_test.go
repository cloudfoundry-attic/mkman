package main_test

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var (
	binPath       string
	fixturesDir   string
	cfReleasePath string
)

var _ = BeforeSuite(func() {
	By("Locating fixtures dir")
	testDir := getDirOfCurrentFile()
	fixturesDir = filepath.Join(testDir, "fixtures")

	By("Ensuring $CF_RELEASE_DIR is set")
	cfReleasePath = os.Getenv("CF_RELEASE_DIR")
	Expect(cfReleasePath).NotTo(BeEmpty(), "$CF_RELEASE_DIR must be provided")

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

func getDirOfCurrentFile() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}
