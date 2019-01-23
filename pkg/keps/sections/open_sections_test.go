package sections_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/calebamiles/keps/pkg/keps/metadata/metadatafakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/calebamiles/keps/pkg/keps/sections"
)

var _ = Describe("Opening existing section content from disk", func() {
	Describe("Open()", func() {
		It("returns sections entries", func() {
			tmpDir, err := ioutil.TempDir("", "sections-open-test")
			Expect(err).ToNot(HaveOccurred(), "creating temp dir for test section")
			defer os.RemoveAll(tmpDir)

			sectionOneFilename := "section_one.md"
			sectionTwoFilename := "section_two.md"

			sectionOneLoc := filepath.Join(tmpDir, sectionOneFilename)
			sectionTwoLoc := filepath.Join(tmpDir, sectionTwoFilename)

			err = ioutil.WriteFile(sectionOneLoc, []byte("example content"), os.ModePerm)
			Expect(err).ToNot(HaveOccurred(), "creating first example section")

			err = ioutil.WriteFile(sectionTwoLoc, []byte("example content"), os.ModePerm)
			Expect(err).ToNot(HaveOccurred(), "creating first example section")

			fakeMetadata := &metadatafakes.FakeKEP{}
			fakeMetadata.SectionLocationsReturns([]string{sectionOneFilename, sectionTwoFilename})
			fakeMetadata.ContentDirReturns(tmpDir)

			entries, err := sections.Open(fakeMetadata)
			Expect(err).ToNot(HaveOccurred(), "opening section content")

			Expect(entries).To(HaveLen(2))
		})
	})

	Context("when a section cannot be opened", func() {
		It("returns an error after attempting to open all sections", func() {
			tmpDir, err := ioutil.TempDir("", "sections-open-test")
			Expect(err).ToNot(HaveOccurred(), "creating temp dir for test section")

			sectionOneFilename := "section_one.md"
			sectionTwoFilename := "does_not_exist.md"

			sectionOneLoc := filepath.Join(tmpDir, sectionOneFilename)

			err = ioutil.WriteFile(sectionOneLoc, []byte("example content"), os.ModePerm)
			Expect(err).ToNot(HaveOccurred(), "creating first example section")

			fakeMetadata := &metadatafakes.FakeKEP{}
			fakeMetadata.SectionLocationsReturns([]string{sectionOneFilename, sectionTwoFilename})
			fakeMetadata.ContentDirReturns(tmpDir)

			entries, err := sections.Open(fakeMetadata)
			Expect(err.Error()).To(ContainSubstring("no such file or directory"), "expected file not found error for missing KEP section")
			Expect(entries).To(BeNil())
		})
	})

})
