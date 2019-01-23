package sections_test

import (
	"errors"

	"github.com/hashicorp/go-multierror"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/calebamiles/keps/pkg/keps/sections"
	"github.com/calebamiles/keps/pkg/keps/sections/sectionsfakes"
)

var _ = Describe("Persting Sections", func() {
	Describe("Persist()", func() {
		It("persists the sections", func() {
			sectionOne := &sectionsfakes.FakeEntry{}
			sectionTwo := &sectionsfakes.FakeEntry{}

			secs := []sections.Entry{sectionOne, sectionTwo}

			err := sections.Persist(secs)
			Expect(err).ToNot(HaveOccurred())

			Expect(sectionOne.PersistCallCount()).To(Equal(1))
			Expect(sectionTwo.PersistCallCount()).To(Equal(1))
		})

		Context("when a section cannot be persisted", func() {
			It("returns the error", func() {
				sectionOne := &sectionsfakes.FakeEntry{}
				sectionTwo := &sectionsfakes.FakeEntry{}

				expectedErr := errors.New("oh no, error")
				sectionTwo.PersistReturns(expectedErr)

				secs := []sections.Entry{sectionOne, sectionTwo}

				err := sections.Persist(secs)
				merr, ok := err.(*multierror.Error)
				Expect(ok).To(BeTrue(), "typecasting an error as a hashicorp/go-multierror")
				Expect(merr.Errors).To(HaveLen(1), "expected one error after persitsing")

				err = merr.Errors[0]
				Expect(err.Error()).To(ContainSubstring(expectedErr.Error()), "expected persist error to contain expected error string")
			})

			XIt("removes the sections", func() {
				Fail("not implemented")
			})
		})
	})
})
