package section_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"time"

	"github.com/calebamiles/keps/pkg/keps/section"
	"github.com/calebamiles/keps/pkg/keps/section/internal/rendering"
	"github.com/calebamiles/keps/pkg/keps/states"
)

var _ = Describe("A Collection of Sections for a KEP state", func() {
	Describe("ForProvisionalState()", func() {
		It("returns a Motivation and Summary", func() {
			title := "The Kubernetes Enhancement Proposal Process"
			authors := []string{"jbeda", "calebamiles"}
			owningSIG := "sig-architecture"
			kepState := states.Implementable
			now := time.Now().UTC()

			info := newMockRenderingInfoProvider()

			for i := 0; i < 3; i++ {
				info.TitleOutput.Ret0 <- title
				info.AuthorsOutput.Ret0 <- authors
				info.OwningSIGOutput.Ret0 <- owningSIG
				info.ContentDirOutput.Ret0 <- ""
				info.StateOutput.Ret0 <- kepState
				info.LastUpdatedOutput.Ret0 <- now
			}

			col, err := section.ForProvisionalState(info)
			Expect(err).ToNot(HaveOccurred())

			secs := col.Sections()
			Expect(secs).ToNot(BeEmpty())

			sectionNames := []string{}
			for _, s := range secs {
				sectionNames = append(sectionNames, s.Name())
			}

			Expect(sectionNames).To(HaveLen(3))
			Expect(sectionNames[0]).To(Equal(rendering.SummaryName))
			Expect(sectionNames[1]).To(Equal(rendering.MotivationName))
			Expect(sectionNames[2]).To(Equal(rendering.ReadmeName))
		})
	})

	Describe("ForImplementableState()", func() {
		It("returns a Developer Guide, Operator Guide, Teacher Guide, and Graduation Criteria", func() {
			title := "The Kubernetes Enhancement Proposal Process"
			authors := []string{"jbeda", "calebamiles"}
			owningSIG := "sig-architecture"
			kepState := states.Implementable
			now := time.Now().UTC()

			info := newMockRenderingInfoProvider()
			for i := 0; i < 7; i++ {
				info.TitleOutput.Ret0 <- title
				info.AuthorsOutput.Ret0 <- authors
				info.OwningSIGOutput.Ret0 <- owningSIG
				info.ContentDirOutput.Ret0 <- ""
				info.StateOutput.Ret0 <- kepState
				info.LastUpdatedOutput.Ret0 <- now
			}

			col, err := section.ForImplementableState(info)
			Expect(err).ToNot(HaveOccurred())

			secs := col.Sections()
			Expect(secs).ToNot(BeEmpty())

			sectionNames := []string{}
			for _, s := range secs {
				sectionNames = append(sectionNames, s.Name())
			}

			Expect(sectionNames).To(HaveLen(7))
			Expect(sectionNames[0]).To(Equal(rendering.SummaryName))
			Expect(sectionNames[1]).To(Equal(rendering.MotivationName))
			Expect(sectionNames[2]).To(Equal(rendering.DeveloperGuideName))
			Expect(sectionNames[3]).To(Equal(rendering.OperatorGuideName))
			Expect(sectionNames[4]).To(Equal(rendering.TeacherGuideName))
			Expect(sectionNames[5]).To(Equal(rendering.GraduationCriteriaName))
			Expect(sectionNames[6]).To(Equal(rendering.ReadmeName))
		})
	})
})
