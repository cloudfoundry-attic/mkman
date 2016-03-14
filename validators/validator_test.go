package validators_test

import (
	"fmt"

	"github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/cloudfoundry/multierror"
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/gomega"
	"github.com/cloudfoundry/mkman/validators"
	"github.com/cloudfoundry/mkman/validators/fakes"
)

var _ = Describe("Validator", func() {

	var (
		fakeValidator    fakes.FakeValidator
		validationTarget validators.ValidationTarget
	)
	BeforeEach(func() {
		fakeValidator = fakes.FakeValidator{}
		fakeValidator.ComposableNameReturns("fake validator")
		fakeValidator.ValidateReturns(nil)
		validationTarget = validators.NewValidationTarget("test", "test")
	})

	Describe("ValidateWith", func() {
		It("should return nil", func() {
			err := validationTarget.ValidateWith(&fakeValidator)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when validation fails", func() {
			BeforeEach(func() {
				fakeValidator.ValidateReturns(fmt.Errorf(fakeValidator.ComposableName()))
			})

			It("should return the multierror", func() {
				err := validationTarget.ValidateWith(&fakeValidator)
				errors := err.(*multierror.MultiError)
				Expect(errors).To(HaveOccurred())
				Expect(errors.Length()).To(Equal(1))
				Expect(errors.Error()).To(ContainSubstring("fake validator"))
			})
		})
	})
})
