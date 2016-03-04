package jobtemplates_test

import (
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/gomega"

	"testing"
)

func TestJobTemplates(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Job Templates Suite")
}
