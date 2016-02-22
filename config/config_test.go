package config_test

import (
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/gomega"
	"github.com/cloudfoundry/mkman/config"
)

var _ = Describe("Config", func() {
	var (
		c config.Config
	)

	BeforeEach(func() {
		c = config.Config{
			CFPath:       "/path/to/cf",
			StemcellPath: "/path/to/stemcell",
			StubPaths:    []string{"/path/to/stub"},
		}
	})

	Context("All the fields available", func() {
		It("should not return any error", func() {
			err := c.Validate()
			Expect(err.HasAny()).Should(BeFalse())
		})
	})

	Describe("Handling errors", func() {
		Describe("on the CFPath", func() {
			Context("when it is an empty string", func() {
				BeforeEach(func() {
					c.CFPath = ""
				})

				It("should return an error", func() {
					err := c.Validate()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("value for cf is required"))
				})
			})

			Context("when it is not an absolute path", func() {
				BeforeEach(func() {
					c.CFPath = "./path/to/cf"
				})

				It("should return an error", func() {
					err := c.Validate()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("value for cf must be absolute path"))
				})
			})
		})

		Describe("on the StemcellPath", func() {
			Context("when it is an empty string", func() {
				BeforeEach(func() {
					c.StemcellPath = ""
				})

				It("should return an error", func() {
					err := c.Validate()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("value for stemcell is required"))
				})
			})
			Context("when it is not an absolute path", func() {
				BeforeEach(func() {
					c.StemcellPath = "./path/to/stemcell"
				})

				It("should return an error", func() {
					err := c.Validate()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("value for stemcell must be absolute path"))
				})
			})
		})

		Describe("on the StubPaths", func() {
			Context("when there are no stub paths", func() {
				BeforeEach(func() {
					c.StubPaths = []string{}
				})

				It("should return an error", func() {
					err := c.Validate()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("value for stub path is required"))
				})
			})

			Context("when there is an empty stub path", func() {
				BeforeEach(func() {
					c.StubPaths = []string{""}
				})

				It("should return an error", func() {
					err := c.Validate()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("value for stub path is required"))
				})
			})

			Context("when it is not an absolute path", func() {
				BeforeEach(func() {
					c.StubPaths = []string{"./path/to/stub"}
				})

				It("should return an error", func() {
					err := c.Validate()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("value for stub path must be absolute path"))
				})
			})
		})

		Describe("multiple errors", func() {
			BeforeEach(func() {
				c.CFPath = ""
				c.StubPaths = []string{}
			})

			Context("when there are multiple errors", func() {
				It("should return the errors", func() {
					err := c.Validate()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("value for stub path is required"))
					Expect(err.Error()).To(ContainSubstring("value for cf is required"))
				})
			})
		})
	})
})
