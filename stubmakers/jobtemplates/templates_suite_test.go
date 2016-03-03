package jobtemplates_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestJobTemplates(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Job Templates Suite")
}
