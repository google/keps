package sections_test

import (
	"sort"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/calebamiles/keps/pkg/keps/sections"
)

var _ = Describe("Section Ordering", func() {
	Describe("ByOrder", func() {
		It("sorts the top level sections in a stable order", func() {
			secs := []string{
				sections.Summary,
				sections.Readme,
				sections.GraduationCriteria,
				sections.Motivation,
				sections.DeveloperGuide,
				sections.TeacherGuide,
				sections.OperatorGuide,
			}

			expectedOrder := []string{
				sections.Summary,
				sections.Motivation,
				sections.DeveloperGuide,
				sections.OperatorGuide,
				sections.TeacherGuide,
				sections.GraduationCriteria,
				sections.Readme,
			}

			sort.Sort(sections.ByOrder(secs))
			Expect(secs).To(Equal(expectedOrder), "top level KEP sections should have a stable order")
		})

		Context("when user defined sections are included", func() {
			It("sorts user defined sections after the summary", func() {
				userSection := "A Good Section"
				anotherUserSection := "Another Good Section"

				secs := []string{
					userSection,
					sections.Summary,
					anotherUserSection,
				}

				sort.Sort(sections.ByOrder(secs))

				Expect(secs[0]).To(Equal(sections.Summary), "summary should be the first section after ordering")
			})
		})
	})
})
