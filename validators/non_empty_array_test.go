package validators_test

import (
	"fmt"

	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/gomega"

	"github.com/cloudfoundry/mkman/validators"
)

var _ = Describe("NonEmptyArray Validator", func() {
	var (
		v validators.Validator
	)

	BeforeEach(func() {
		v = validators.NonEmptyArray()
	})

	It("should return the correct composable name", func() {
		Expect(v.ComposableName()).To(Equal("non-empty array"))
	})

	Describe("Validate", func() {
		Context("when the target is empty array", func() {
			var validationTarget validators.ValidationTarget

			BeforeEach(func() {
				validationTarget = validators.NewValidationTarget([]string{}, "empty_array")
			})

			It("returns error", func() {
				err := v.Validate(validationTarget)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("value must be %s", v.ComposableName())))
			})
		})

		Context("when the target is not of type []string", func() {
			var validationTarget validators.ValidationTarget

			BeforeEach(func() {
				validationTarget = validators.NewValidationTarget("test", "empty_array")
			})

			It("returns error", func() {
				err := v.Validate(validationTarget)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("value must be of type string array"))
			})
		})

		Context("when the target is valid", func() {
			var validationTarget validators.ValidationTarget

			BeforeEach(func() {
				validationTarget = validators.NewValidationTarget([]string{"non_empty"}, "nonempty_object")
			})

			It("returns without error", func() {
				err := v.Validate(validationTarget)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
