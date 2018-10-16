package rendering_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"time"

	"github.com/calebamiles/keps/pkg/keps/section"
	"github.com/calebamiles/keps/pkg/keps/section/internal/rendering"
	"github.com/calebamiles/keps/pkg/keps/states"
)

var _ = Describe("The README section", func() {
	Describe("NewReadme()", func() {
		It("renders a new README from the provided sections", func() {
			title := "Kubernetes Enhancement Proposal Process"
			authors := []string{"jbeda", "calebamiles"}
			owningSIG := "sig-architecture"
			contentDir := ""

			now := time.Now().UTC()

			info := newMockInfoProvider()

			for i := 0; i < 3; i++ {
				info.TitleOutput.Ret0 <- title
				info.AuthorsOutput.Ret0 <- authors
				info.OwningSIGOutput.Ret0 <- owningSIG
				info.ContentDirOutput.Ret0 <- contentDir
				info.StateOutput.Ret0 <- states.Provisional
				info.LastUpdatedOutput.Ret0 <- now
			}

			summary, err := section.New(rendering.SummaryName, info)
			Expect(err).ToNot(HaveOccurred())

			motivation, err := section.New(rendering.MotivationName, info)
			Expect(err).ToNot(HaveOccurred())

			secs := []rendering.SectionProvider{summary, motivation}
			enhancedInfo := &infoWithSections{
				ss:           secs,
				InfoProvider: info,
			}

			readmeBytes, err := rendering.NewReadme(enhancedInfo)
			Expect(err).ToNot(HaveOccurred())

			Expect(readmeBytes).To(ContainSubstring("**Authors: jbeda, calebamiles**"))
			Expect(readmeBytes).To(ContainSubstring("**Sponsoring SIG: [Architecture](https://github.com/kubernetes/community/tree/master/sig-architecture/README.md)**"))
			Expect(readmeBytes).To(ContainSubstring("**Status: provisional**"))
			Expect(readmeBytes).To(ContainSubstring("## Table of Contents"))
			Expect(readmeBytes).To(ContainSubstring("[Summary](summary.md)"))
			Expect(readmeBytes).To(ContainSubstring("[Motivation](motivation.md)"))
		})
	})
})

type infoWithSections struct {
	ss []rendering.SectionProvider
	rendering.InfoProvider
}

func (i *infoWithSections) Sections() []rendering.SectionProvider { return i.ss }
