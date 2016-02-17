package commands_test

import (
	"path"
	"path/filepath"
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var (
	fixturesDir string
)

var _ = BeforeSuite(func() {
	By("Locating fixtures dir")
	testDir := getDirOfCurrentFile()
	fixturesDir = filepath.Join(testDir, "..", "fixtures")
})

func TestCommands(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commands Suite")
}

func getDirOfCurrentFile() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}
