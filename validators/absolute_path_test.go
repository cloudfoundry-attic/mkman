package validators_test

import (
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/gomega"

	"github.com/cloudfoundry/mkman/validators"
)

var _ = Describe("AbsolutePath Validator", func() {
	var (
		v validators.Validator
	)

	BeforeEach(func() {
		v = validators.AbsolutePath()
	})

	It("should return the correct composable name", func() {
		Expect(v.ComposableName()).To(Equal("absolute path"))
	})

	Describe("Validate", func(){
		Context("when the target is valid", func(){
			var validationTarget validators.ValidationTarget

			BeforeEach(func(){
				validationTarget = validators.NewValidationTarget("/test/absolute/path", "path")
			})

			It("returns without error",func(){
				err := validationTarget.ValidateWith(v)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the target is invalid", func(){
			var validationTarget validators.ValidationTarget

			BeforeEach(func(){
				validationTarget = validators.NewValidationTarget("some/relative/path", "path")
			})

			It("returns without error",func(){
				err := validationTarget.ValidateWith(v)
				Expect(err).To(HaveOccurred())
				Expect(err.Length()).To(Equal(1))
				Expect(err.Error()).To(MatchRegexp("value must be absolute path.*some/relative/path"))
			})
		})
	})
})
