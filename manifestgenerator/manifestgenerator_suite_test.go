package manifestgenerator_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestManifestgenerator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Manifestgenerator Suite")
}
