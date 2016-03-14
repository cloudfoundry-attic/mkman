package validators_test

import (
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/gomega"

	"github.com/cloudfoundry/mkman/validators"
)

var _ = Describe("NonEmpty Validator", func() {
	var (
		v validators.Validator
	)

	BeforeEach(func() {
		v = validators.NonEmpty()
	})

	It("should return the correct composable name", func() {
		Expect(v.ComposableName()).To(Equal("non-empty"))
	})

	Describe("Validate", func() {
		Context("when the target is empty string", func() {
			var validationTarget validators.ValidationTarget

			BeforeEach(func() {
				validationTarget = validators.NewValidationTarget("", "empty_object")
			})

			It("returns error", func() {
				err := v.Validate(validationTarget)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(v.ComposableName()))
			})
		})

		Context("when the target is nonempty string", func() {
			var validationTarget validators.ValidationTarget

			BeforeEach(func() {
				validationTarget = validators.NewValidationTarget("non_empty", "nonempty_object")
			})

			It("returns without error", func() {
				err := v.Validate(validationTarget)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
