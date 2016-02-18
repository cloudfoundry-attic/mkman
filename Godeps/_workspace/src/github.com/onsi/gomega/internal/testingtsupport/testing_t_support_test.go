package testingtsupport_test

import (
	. "github.com/pivotal-cf-experimental/mkman/Godeps/_workspace/src/github.com/onsi/gomega"

	"testing"
)

func TestTestingT(t *testing.T) {
	RegisterTestingT(t)
	Ω(true).Should(BeTrue())
}
