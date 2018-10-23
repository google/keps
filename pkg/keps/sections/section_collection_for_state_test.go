package sections_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"time"

	"github.com/calebamiles/keps/pkg/keps/sections"
	"github.com/calebamiles/keps/pkg/keps/sections/internal/rendering"
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

			col, err := sections.ForProvisionalState(info)
			Expect(err).ToNot(HaveOccurred())

			secs := col.Sections()

			Expect(secs).To(HaveLen(3))
			Expect(secs[0]).To(Equal(rendering.SummaryFilename))
			Expect(secs[1]).To(Equal(rendering.MotivationFilename))
			Expect(secs[2]).To(Equal(rendering.ReadmeFilename))
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

			col, err := sections.ForImplementableState(info)
			Expect(err).ToNot(HaveOccurred())

			secs := col.Sections()
			Expect(secs).ToNot(BeEmpty())

			Expect(secs).To(HaveLen(7))
			Expect(secs[0]).To(Equal(rendering.SummaryFilename))
			Expect(secs[1]).To(Equal(rendering.MotivationFilename))
			Expect(secs[2]).To(Equal(rendering.DeveloperGuideFilename))
			Expect(secs[3]).To(Equal(rendering.OperatorGuideFilename))
			Expect(secs[4]).To(Equal(rendering.TeacherGuideFilename))
			Expect(secs[5]).To(Equal(rendering.GraduationCriteriaFilename))
			Expect(secs[6]).To(Equal(rendering.ReadmeFilename))
		})
	})
})
