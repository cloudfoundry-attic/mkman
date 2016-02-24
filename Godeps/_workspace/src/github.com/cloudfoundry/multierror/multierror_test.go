package multierror_test

import (
	"fmt"

	"github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/cloudfoundry/multierror"

	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/onsi/gomega"
)

var _ = Describe("Multierror", func() {
	var (
		m multierror.MultiError
	)

	BeforeEach(func() {
		m = multierror.MultiError{}
	})

	Describe("Add", func() {
		It("adds an error", func() {
			m.Add(fmt.Errorf("Sample Error"))
			Expect(m.HasAny()).To(BeTrue())
			Expect(m.Error()).To(ContainSubstring("Sample Error"))
		})

		Describe("adding a MultiError", func() {
			It("adds all the errors", func() {
				m.Add(fmt.Errorf("Error 1"))
				m.Add(fmt.Errorf("Error 2"))
				m2 := multierror.MultiError{}
				m2.Add(fmt.Errorf("Error 3"))
				m2.Add(fmt.Errorf("Error 4"))

				m.Add(m2)

				Expect(m.Error()).To(ContainSubstring("Error 1"))
				Expect(m.Error()).To(ContainSubstring("Error 2"))
				Expect(m.Error()).To(ContainSubstring("Error 3"))
				Expect(m.Error()).To(ContainSubstring("Error 4"))
			})

			Context("when the multierror is empty", func() {
				It("retains existing errors", func() {
					m.Add(fmt.Errorf("Error 1"))
					m.Add(fmt.Errorf("Error 2"))
					m2 := multierror.MultiError{}

					m.Add(m2)

					Expect(m.Error()).To(ContainSubstring("Error 1"))
					Expect(m.Error()).To(ContainSubstring("Error 2"))
				})
			})
		})
	})

	Describe("AddWithPrefix", func() {
		It("adds an error with prefix", func() {
			m.AddWithPrefix(fmt.Errorf("Error 1"), "Prefix:")

			Expect(m.Error()).To(ContainSubstring("Prefix:Error 1"))
		})

		Context("adding a MultiError", func() {
			It("adds all the errors with prefix", func() {
				m.AddWithPrefix(fmt.Errorf("Error 1"), "Prefix:")
				m.AddWithPrefix(fmt.Errorf("Error 2"), "Prefix:")
				m2 := multierror.MultiError{}
				m2.Add(fmt.Errorf("Error 3"))
				m2.Add(fmt.Errorf("Error 4"))

				m.AddWithPrefix(m2, "Prefix:")

				Expect(m.Error()).To(ContainSubstring("Prefix:Error 1"))
				Expect(m.Error()).To(ContainSubstring("Prefix:Error 2"))
				Expect(m.Error()).To(ContainSubstring("Prefix:Error 3"))
				Expect(m.Error()).To(ContainSubstring("Prefix:Error 4"))
			})
		})
	})

	Describe("Error", func() {
		Context("when there are errors", func() {
			BeforeEach(func() {
				m.Add(fmt.Errorf("Error 1"))
				m.Add(fmt.Errorf("Error 2"))
			})

			It("prints all the errors", func() {
				Expect(m.Error()).To(ContainSubstring("encountered 2 errors during validation"))
				Expect(m.Error()).To(ContainSubstring("Error 1"))
				Expect(m.Error()).To(ContainSubstring("Error 2"))
			})
		})

		It("should say encountered 0 errors", func() {
			Expect(m.Error()).To(ContainSubstring("encountered 0 errors during validation"))
		})
	})

	Describe("HasAny", func() {
		Context("when there are errors", func() {
			BeforeEach(func() {
				m.Add(fmt.Errorf("Error 1"))
			})

			It("returns true", func() {
				Expect(m.HasAny()).To(BeTrue())
			})
		})

		Context("when there are no errors", func() {
			It("returns false", func() {
				Expect(m.HasAny()).To(BeFalse())
			})
		})
	})
})
