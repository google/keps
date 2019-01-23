package sections_test

import (
	"github.com/calebamiles/keps/pkg/keps/sections"
	"github.com/calebamiles/keps/pkg/keps/sections/internal/rendering"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("working with KEP section names and filenames", func() {
	Describe("Filename()", func() {
		Context("given a top level KEP section name", func() {
			It("returns the corresponding top level section filename", func() {
				Expect(sections.Filename(sections.Summary)).To(Equal(rendering.SummaryFilename), "Summary -> summary.md")
				Expect(sections.Filename(sections.Motivation)).To(Equal(rendering.MotivationFilename), "Motivation -> motivation.md")
				Expect(sections.Filename(sections.DeveloperGuide)).To(Equal(rendering.DeveloperGuideFilename), "Developer Guide -> guides/developer.md")
				Expect(sections.Filename(sections.OperatorGuide)).To(Equal(rendering.OperatorGuideFilename), "Operator Guide -> guides/operator.md")
				Expect(sections.Filename(sections.TeacherGuide)).To(Equal(rendering.TeacherGuideFilename), "Teacher Guide -> guides/teacher.md")
				Expect(sections.Filename(sections.GraduationCriteria)).To(Equal(rendering.GraduationCriteriaFilename), "Graduation Criteria -> graduation_criteria.md")
				Expect(sections.Filename(sections.Readme)).To(Equal(rendering.ReadmeFilename), "README -> readme.md")
			})
		})

		Context("given an arbitrary section name", func() {
			It("returns a downcased filename with spaces converted to underscore", func() {
				givenSectionName := "A Good Section Title"
				expectedFilename := "a_good_section_title.md"

				Expect(sections.Filename(givenSectionName)).To(Equal(expectedFilename), "downcasing, removing spaces, adding .md")
			})
		})
	})

	Describe("NameForFilename()", func() {
		Context("given a top level KEP section filename", func() {
			It("returns the corresponding top level section name", func() {
				Expect(sections.Filename(sections.Summary)).To(Equal(rendering.SummaryFilename), "Summary -> summary.md")
				Expect(sections.Filename(sections.Motivation)).To(Equal(rendering.MotivationFilename), "Motivation -> motivation.md")
				Expect(sections.Filename(sections.DeveloperGuide)).To(Equal(rendering.DeveloperGuideFilename), "Developer Guide -> guides/developer.md")
				Expect(sections.Filename(sections.OperatorGuide)).To(Equal(rendering.OperatorGuideFilename), "Operator Guide -> guides/operator.md")
				Expect(sections.Filename(sections.TeacherGuide)).To(Equal(rendering.TeacherGuideFilename), "Teacher Guide -> guides/teacher.md")
				Expect(sections.Filename(sections.GraduationCriteria)).To(Equal(rendering.GraduationCriteriaFilename), "Graduation Criteria -> graduation_criteria.md")
				Expect(sections.Filename(sections.Readme)).To(Equal(rendering.ReadmeFilename), "README -> readme.md")
			})
		})

		Context("given an arbitrary section filename", func() {
			It("returns a name where undercases are converted to spaces, strips out a .md extension, and attempts some reasonable capitalization", func() {
				givenFilename := "a_good_section_title.md"
				expectedSectionName := "A Good Section Title"

				Expect(sections.NameForFilename(givenFilename)).To(Equal(expectedSectionName), "title text, removing underscores, removing .md")
			})
		})
	})
})
