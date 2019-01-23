package sections_test

import (
	"github.com/calebamiles/keps/pkg/keps/sections/sectionsfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/calebamiles/keps/pkg/keps/sections"
)

var _ = Describe("determining section locations relative to contentDir", func() {
	Describe("Locations()", func() {
		It("returns the filename for each section", func() {

			fakeSectionOneFilename := "section_one.md"
			fakeSectionOne := &sectionsfakes.FakeEntry{}
			fakeSectionOne.FilenameReturns(fakeSectionOneFilename)

			fakeSectionTwoFilename := "section_two.md"
			fakeSectionTwo := &sectionsfakes.FakeEntry{}
			fakeSectionTwo.FilenameReturns(fakeSectionTwoFilename)

			entries := []sections.Entry{fakeSectionOne, fakeSectionTwo}

			locs := sections.Locations(entries)
			Expect(locs).To(ConsistOf(fakeSectionOneFilename, fakeSectionTwoFilename), "section filenames should be extracted from sections.Entry")
		})
	})
})
