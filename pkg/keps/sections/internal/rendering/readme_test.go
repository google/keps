package rendering_test

import (
	"fmt"
	"time"

	"github.com/calebamiles/keps/pkg/keps/metadata/metadatafakes"
	"github.com/calebamiles/keps/pkg/keps/states"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/calebamiles/keps/pkg/keps/sections/internal/rendering"
)

var _ = Describe("The README section", func() {
	Describe("NewReadme()", func() {
		It("renders a new README from the provided sections", func() {
			title := "Kubernetes Enhancement Proposal Process"
			authors := []string{"jbeda", "calebamiles"}
			sectionLocations := []string{rendering.SummaryFilename, rendering.MotivationFilename, rendering.ReadmeFilename}
			owningSIG := "sig-architecture"
			contentDir := ""

			now := time.Now().UTC()

			info := &metadatafakes.FakeKEP{}
			info.TitleReturns(title)
			info.AuthorsReturns(authors)
			info.OwningSIGReturns(owningSIG)
			info.ContentDirReturns(contentDir)
			info.StateReturns(states.Provisional)
			info.LastUpdatedReturns(now)
			info.SectionLocationsReturns(sectionLocations)

			readmeBytes, err := rendering.NewReadme(info)
			Expect(err).ToNot(HaveOccurred())

			Expect(string(readmeBytes)).To(ContainSubstring("**Authors: jbeda, calebamiles**"), "expected jbeda and calebamiles to be listed as authors")
			Expect(string(readmeBytes)).To(ContainSubstring("**Sponsoring SIG: [Architecture](https://github.com/kubernetes/community/tree/master/sig-architecture/README.md)**"), "expected SIG Architecture to be listed as owning SIG")
			Expect(string(readmeBytes)).To(ContainSubstring("**Status: provisional**"), "expected KEP to have `provisional` state")
			Expect(string(readmeBytes)).To(ContainSubstring("## Table of Contents"), "expected to find `Table of Contents` heading")
			Expect(string(readmeBytes)).To(ContainSubstring("[Summary](summary.md)"), "expected to find `Summary` listed in table of contents")
			Expect(string(readmeBytes)).To(ContainSubstring("[Motivation](motivation.md)"), "expected to find `Motivation` listed in table of contents")
			Expect(string(readmeBytes)).To(ContainSubstring(fmt.Sprintf("Last Updated: %s", now.UTC().String())), "expected last updated time to appear in README")
			Expect(string(readmeBytes)).ToNot(ContainSubstring("[README](README.md)"), "expected to remove README.md reference within README.md")
		})
	})
})
