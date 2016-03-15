package validators_test

import (
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/gomega"

	"github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/cloudfoundry/multierror"
	"github.com/cloudfoundry/mkman/validators"
)

var _ = Describe("AllOf Validator", func() {
	var v validators.Validator

	BeforeEach(func() {
		v = validators.AllOf(
			validators.AbsolutePath(),
			validators.ExistsOnFilesystem(),
		)
	})

	It("should return empty composable name", func() {
		Expect(v.ComposableName()).To(Equal(""))
	})

	Describe("Validate", func() {
		Context("when the target fails to satisfy the first validation", func() {
			var validationTarget validators.ValidationTarget
			var err *multierror.MultiError

			BeforeEach(func() {
				validationTarget = validators.NewValidationTarget("some/relative/path", "path")
				err = validationTarget.ValidateWith(v)
			})

			It("returns the first error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Length()).To(Equal(1))
				Expect(err.Error()).To(ContainSubstring(validators.AbsolutePath().ComposableName()))
			})

			It("should have a composable name matching the failed validation", func() {
				Expect(v.ComposableName()).To(ContainSubstring(validators.AbsolutePath().ComposableName()))
			})
		})

		Context("when the target fails to satisfy the second validation", func() {
			var validationTarget validators.ValidationTarget
			var err *multierror.MultiError

			BeforeEach(func() {
				validationTarget = validators.NewValidationTarget("/some/nonexisting/absolute/path", "path")
				err = validationTarget.ValidateWith(v)
			})

			It("returns the first error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Length()).To(Equal(1))
				Expect(err.Error()).To(ContainSubstring(validators.ExistsOnFilesystem().ComposableName()))
			})

			It("should have a composable name matching the failed validation", func() {
				Expect(v.ComposableName()).To(ContainSubstring(validators.ExistsOnFilesystem().ComposableName()))
			})
		})
	})
})
