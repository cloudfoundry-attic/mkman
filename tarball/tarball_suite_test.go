package tarball_test

import (
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/gomega"

	"testing"
)

func TestTarball(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tarball Suite")
}
