package validators_test

import (
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/gomega"
	. "github.com/cloudfoundry/mkman/config"
)

var _ = Describe("Validation", func() {

	var (
		validation Validation
	)

	Describe("Validate", func() {
		Context("When the type is a directory", func() {

			BeforeEach(func() {
				validation = Validation{AllowedType: DirType}
			})

			It("should not return any error", func() {
				validator := NewValidator("/tmp/", "test_name") //Replace with a tmpdir
				err := validator.Validate(validation)
				Expect(err).NotTo(HaveOccurred())
			})
			Context("when there are errors", func() {
				It("should return error when object expects a directory", func() {
					validator := NewValidator("foo.tgz", "foo-name")
					err := validator.Validate(validation)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("value must be absolute path to directory"))
				})
			})
		})

		Context("when the type is a file", func() {
			BeforeEach(func() {
				validation = Validation{AllowedType: FileType}
			})

			It("should not return any error", func() {
				validator := NewValidator("/tmp/foo", "test_name") //Replace with a tmpdir
				err := validator.Validate(validation)
				Expect(err).NotTo(HaveOccurred())
			})

			Context("when there are errors", func() {
				It("should return error when object expects a file", func() {
					validator := NewValidator("/tmp/", "foo-name")
					err := validator.Validate(validation)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("value must be absolute path to file"))
				})
			})
		})

		Context("when the type is an alias", func() {
			BeforeEach(func() {
				validation.VersionAliases = &[]string{"foo-alias"}
			})

			Context("when the alias is valid", func() {
				It("should not return any error", func() {
					validator := NewValidator("foo-alias", "foo-name")
					err := validator.Validate(validation)
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("when the alias is invalid", func() {
				It("should return an error", func() {
					validator := NewValidator("bar-alias", "foo-name")
					err := validator.Validate(validation)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("a valid version alias or an absolute path"))
				})
			})
		})

		Context("when there are errors", func() {
			It("should return error when object is empty", func() {
				validator := NewValidator("", "foo-name")
				err := validator.Validate(validation)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("value is required"))
			})

			It("should return error when object does not exist", func() {
				validator := NewValidator("/foo", "foo-name")
				err := validator.Validate(validation)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("does not exist"))
			})
		})
	})
})
