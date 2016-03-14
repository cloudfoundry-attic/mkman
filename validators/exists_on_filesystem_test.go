package validators_test

import (
	"io/ioutil"
	"os"

	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/gomega"

	"github.com/cloudfoundry/mkman/validators"
)

var _ = Describe("ExistsOnFileSystem Validator", func() {
	var (
		v validators.Validator
	)

	Describe("ComposableName", func() {
		Context("when validating only for directory type", func() {
			BeforeEach(func() {
				v = validators.ExistsOnFilesystem(validators.Directory())
			})

			It("should return the correct composable name", func() {
				Expect(v.ComposableName()).To(Equal("a path to a directory that exists"))
			})
		})

		Context("when validating only for file type", func() {
			BeforeEach(func() {
				v = validators.ExistsOnFilesystem(validators.File())
			})

			It("should return the correct composable name", func() {
				Expect(v.ComposableName()).To(Equal("a path to a file that exists"))
			})
		})

		Context("when validating for multiple file type", func() {
			BeforeEach(func() {
				v = validators.ExistsOnFilesystem(validators.Directory(), validators.File())
			})

			It("should return the correct composable name", func() {
				Expect(v.ComposableName()).To(Equal("a path to a directory or a file that exists"))
			})
		})
	})

	Describe("Validate", func() {
		Context("when the target is invalid", func() {
			var validationTarget validators.ValidationTarget

			BeforeEach(func() {
				validationTarget = validators.NewValidationTarget("/not/a/path", "object")
				v = validators.ExistsOnFilesystem(validators.File(), validators.Directory())
			})

			It("returns error", func() {
				err := v.Validate(validationTarget)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(v.ComposableName()))
			})
		})

		Context("when the target is a directory that is present on the filesystem", func() {
			var validationTarget validators.ValidationTarget
			var tempDirPath string
			var err error

			BeforeEach(func() {
				tempDirPath, err = ioutil.TempDir("", "")
				Expect(err).NotTo(HaveOccurred())

				validationTarget = validators.NewValidationTarget(tempDirPath, "path")
			})

			AfterEach(func() {
				err = os.RemoveAll(tempDirPath)
				Expect(err).NotTo(HaveOccurred())
			})

			Context("when validating only for one file type", func() {
				BeforeEach(func() {
					v = validators.ExistsOnFilesystem(validators.Directory())
				})

				It("returns without error", func() {
					err := v.Validate(validationTarget)
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("when validating for multiple file type", func() {
				BeforeEach(func() {
					v = validators.ExistsOnFilesystem(validators.Directory(), validators.File())
				})

				It("returns without error", func() {
					err := v.Validate(validationTarget)
					Expect(err).NotTo(HaveOccurred())
				})
			})
		})

		Context("when the target is a file that is present on the filesystem", func() {
			var validationTarget validators.ValidationTarget
			var tempFile *os.File
			var err error

			BeforeEach(func() {
				tempFile, err = ioutil.TempFile("", "")
				Expect(err).NotTo(HaveOccurred())

				validationTarget = validators.NewValidationTarget(tempFile.Name(), "path")
			})

			AfterEach(func() {
				err = os.Remove(tempFile.Name())
				Expect(err).NotTo(HaveOccurred())
			})

			Context("when validating only for one file type", func() {
				BeforeEach(func() {
					v = validators.ExistsOnFilesystem(validators.File())
				})

				It("returns without error", func() {
					err := v.Validate(validationTarget)
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("when validating for multiple file type", func() {
				BeforeEach(func() {
					v = validators.ExistsOnFilesystem(validators.Directory(), validators.File())
				})

				It("returns without error", func() {
					err := v.Validate(validationTarget)
					Expect(err).NotTo(HaveOccurred())
				})
			})
		})
	})
})
