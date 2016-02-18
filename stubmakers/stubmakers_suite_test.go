package stubmakers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestStubmakers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Stubmakers Suite")
}
