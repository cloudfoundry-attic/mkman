package validators_test

import (
	"os"
	"io/ioutil"
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/gomega"

	"github.com/cloudfoundry/mkman/validators"
)

var _ = Describe("AbsolutePath Validator", func() {
	var (
		v validators.Validator
	)

	BeforeEach(func() {
		v = validators.Directory()
	})

	It("should return the correct composable name", func(){
		Expect(v.ComposableName()).To(Equal("path to directory"))
	})

	Describe("Validate", func(){
		Context("when the target is valid", func(){
			var validationTarget validators.ValidationTarget
			var tempDirPath string
			var err error

			BeforeEach(func(){
				tempDirPath, err = ioutil.TempDir("", "")
				Expect(err).NotTo(HaveOccurred())

				validationTarget = validators.NewValidationTarget(tempDirPath, "path")
			})

			AfterEach(func(){
				err = os.RemoveAll(tempDirPath)
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns without error",func(){
				err := validationTarget.ValidateWith(v)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the target in invalid", func(){
			var validationTarget validators.ValidationTarget
			var tempFile *os.File
			var err error

			BeforeEach(func(){
				tempFile, err = ioutil.TempFile("", "")
				Expect(err).NotTo(HaveOccurred())

				validationTarget = validators.NewValidationTarget(tempFile.Name(), "path")
			})

			AfterEach(func(){
				err = os.Remove(tempFile.Name())
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns without error",func(){
				errs := validationTarget.ValidateWith(v)
				Expect(errs).To(HaveOccurred())
				Expect(errs.Length()).To(Equal(1))
				Expect(errs.Error()).To(ContainSubstring(validators.Directory().ComposableName()))
			})
		})
	})
})
