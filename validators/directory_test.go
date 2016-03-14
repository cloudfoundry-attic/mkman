package validators_test

import (
	"io/ioutil"
	"os"

	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/gomega"

	"github.com/cloudfoundry/mkman/validators"
)

var _ = Describe("Directory Validator", func() {
	var (
		v validators.Validator
	)

	BeforeEach(func() {
		v = validators.Directory()
	})

	It("should return the correct composable name", func() {
		Expect(v.ComposableName()).To(Equal("a directory"))
	})

	Describe("Validate", func() {
		Context("when the target is valid", func() {
			var validationTarget validators.ValidationTarget
			var tempDirPath string

			BeforeEach(func() {
				var err error
				tempDirPath, err = ioutil.TempDir("", "")
				Expect(err).NotTo(HaveOccurred())

				validationTarget = validators.NewValidationTarget(tempDirPath, "path")
			})

			AfterEach(func() {
				err := os.RemoveAll(tempDirPath)
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns without error", func() {
				err := v.Validate(validationTarget)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the target is invalid", func() {
			var validationTarget validators.ValidationTarget
			var tempFile *os.File

			BeforeEach(func() {
				var err error
				tempFile, err = ioutil.TempFile("", "")
				Expect(err).NotTo(HaveOccurred())

				validationTarget = validators.NewValidationTarget(tempFile.Name(), "path")
			})

			AfterEach(func() {
				err := os.Remove(tempFile.Name())
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns error", func() {
				err := v.Validate(validationTarget)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(v.ComposableName()))
			})
		})
	})
})
