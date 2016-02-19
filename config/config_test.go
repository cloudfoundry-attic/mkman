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

	Describe("Handling errors", func() {
		Describe("on the CFPath", func() {
			BeforeEach(func() {
				c.CFPath = ""
			})

			Context("when it is an empty string", func() {
				It("should return an error", func() {
					err := c.Validate()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("path to cf is missing"))
				})
			})
		})

		Describe("on the StemcellPath", func() {
			BeforeEach(func() {
				c.StemcellPath = ""
			})

			Context("when it is an empty string", func() {
				It("should return an error", func() {
					err := c.Validate()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("path to stemcell is missing"))
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
					Expect(err.Error()).To(ContainSubstring("at least one stub path is required"))
				})
			})

			Context("when there is an empty stub path", func() {
				BeforeEach(func() {
					c.StubPaths = []string{""}
				})
				It("should return an error", func() {
					err := c.Validate()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("there is an empty stub path"))
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
					Expect(err.Error()).To(ContainSubstring("at least one stub path is required"))
					Expect(err.Error()).To(ContainSubstring("path to cf is missing"))
				})
			})
		})
	})
})
