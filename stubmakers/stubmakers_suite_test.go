package stubmakers_test

import (
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/gomega"

	"testing"
)

func TestStubmakers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Stubmakers Suite")
}
