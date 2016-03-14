package validators_test

import (
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/gomega"

	"github.com/cloudfoundry/mkman/validators"
)

var _ = Describe("VersionAlias Validator", func() {
	var (
		v validators.Validator
	)

	BeforeEach(func() {
		v = validators.VersionAlias([]string{"director-latest"})
	})

	It("should return the correct composable name", func() {
		Expect(v.ComposableName()).To(Equal("valid version alias"))
	})

	Describe("Validate", func() {
		Context("when the target is invalid", func() {
			var validationTarget validators.ValidationTarget

			BeforeEach(func() {
				validationTarget = validators.NewValidationTarget("invalid alias", "invalid key")
			})

			It("returns error", func() {
				err := v.Validate(validationTarget)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(v.ComposableName()))
			})
		})

		Context("when the target is valid", func() {
			var validationTarget validators.ValidationTarget

			BeforeEach(func() {
				validationTarget = validators.NewValidationTarget("director-latest", "valid key")
			})

			It("returns without error", func() {
				err := v.Validate(validationTarget)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
