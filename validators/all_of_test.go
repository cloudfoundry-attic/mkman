package validators_test

import (
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/gomega"

	"github.com/cloudfoundry/mkman/validators"
)

var _ = Describe("AllOf Validator", func() {
	var v validators.Validator

	BeforeEach(func() {
		v = validators.AllOf(
			validators.AbsolutePath(),
			validators.ExistsOnFilesystem(validators.File(), validators.Directory()),
		)
	})

	It("should return empty composable name", func() {
		Expect(v.ComposableName()).To(Equal(""))
	})

	Describe("Validate", func() {
		Context("when the target fails to satisfy the first validation", func() {
			var validationTarget validators.ValidationTarget

			BeforeEach(func() {
				validationTarget = validators.NewValidationTarget("some/relative/path", "path")
			})

			It("returns the first validation error", func() {
				err := v.Validate(validationTarget)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(
					validators.AbsolutePath().Validate(validationTarget).Error(),
				))
			})

			It("should have a composable name matching the failed validation", func() {
				v.Validate(validationTarget) // Set the composable name to some value
				Expect(v.ComposableName()).To(ContainSubstring(validators.AbsolutePath().ComposableName()))
			})
		})

		Context("when the target fails to satisfy the second validation", func() {
			var validationTarget validators.ValidationTarget

			BeforeEach(func() {
				validationTarget = validators.NewValidationTarget("/some/nonexistant/absolute/path", "path")
			})

			It("returns the the error from the validator that failed", func() {
				err := v.Validate(validationTarget)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(
					validators.ExistsOnFilesystem(
						validators.File(),
						validators.Directory(),
					).Validate(validationTarget).Error(),
				))
			})

			It("should have a composable name matching the failed validation", func() {
				v.Validate(validationTarget) // Set the composable name to some value
				Expect(v.ComposableName()).To(ContainSubstring(
					validators.ExistsOnFilesystem(
						validators.File(),
						validators.Directory(),
					).ComposableName(),
				))
			})
		})
	})
})
