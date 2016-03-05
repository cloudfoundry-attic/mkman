package config_test

import (
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/gomega"
	"github.com/cloudfoundry/mkman/config"
)

var _ = Describe("Validation", func() {

	var (
		validation config.Validation
	)

	BeforeEach(func() {
		validation = config.Validation{AllowedType: config.DirType}
	})

	Describe("Validate", func() {
		It("should not return any error", func() {
			validator := config.NewValidator("/tmp/", "test_name")
			err := validator.Validate(validation)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when testing for aliases", func() {
			Context("when there is an alias present", func() {
				BeforeEach(func() {
					validation.VersionAliases = &[]string{"foo-alias"}
				})

				Context("when the alias is valid", func() {
					It("should not return any error", func() {
						validator := config.NewValidator("foo-alias", "foo-name")
						err := validator.Validate(validation)
						Expect(err).NotTo(HaveOccurred())
					})
				})

				Context("when the alias is invalid", func() {
					It("should not return any error", func() {
						validator := config.NewValidator("bar-alias", "foo-name")
						err := validator.Validate(validation)
						Expect(err).To(HaveOccurred())
					})
				})
			})

			Context("handling errors", func() {
				It("should return error when object is empty", func() {
					validator := config.NewValidator("", "foo-name")
					err := validator.Validate(validation)
					Expect(err).To(HaveOccurred())
				})

				It("should return error when object expects a directory", func() {
					validator := config.NewValidator("foo.tgz", "foo-name")
					err := validator.Validate(validation)
					Expect(err).To(HaveOccurred())
				})

				It("should return error when object expects a file", func() {
					validation = config.Validation{AllowedType: config.FileType}
					validator := config.NewValidator("/foo", "foo-name")
					err := validator.Validate(validation)
					Expect(err).To(HaveOccurred())
				})
			})
		})
	})
})
