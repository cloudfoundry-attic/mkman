package leafnodes_test

import (
	. "github.com/pivotal-cf-experimental/mkman/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/pivotal-cf-experimental/mkman/Godeps/_workspace/src/github.com/onsi/gomega"

	"testing"
)

func TestLeafNode(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LeafNode Suite")
}
