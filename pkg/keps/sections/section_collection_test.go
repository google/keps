package sections_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-multierror"

	"github.com/calebamiles/keps/pkg/keps/sections"
)

var _ = Describe("A collection of sections", func() {
	Describe("Open", func() {
		It("returns a collection representing sections on disk", func() {
			tmpDir, err := ioutil.TempDir("", "kep-collection")
			Expect(err).ToNot(HaveOccurred())
			defer os.RemoveAll(tmpDir)

			summaryFilename := "summary.md"
			motivationFilename := "motivation.md"
			testContent := []byte("test section content")

			summaryLocation := filepath.Join(tmpDir, summaryFilename)
			motivationLocation := filepath.Join(tmpDir, motivationFilename)

			err = ioutil.WriteFile(summaryLocation, testContent, os.ModePerm)
			Expect(err).ToNot(HaveOccurred())

			err = ioutil.WriteFile(motivationLocation, testContent, os.ModePerm)
			Expect(err).ToNot(HaveOccurred())

			locations := newMockLocationProvider()
			locations.SectionsOutput.Ret0 <- []string{summaryFilename, motivationFilename}
			locations.ContentDirOutput.Ret0 <- tmpDir

			col, err := sections.OpenCollection(locations)
			Expect(err).ToNot(HaveOccurred())

			Expect(col.Sections()).To(ConsistOf([]string{summaryFilename, motivationFilename}))
		})

		Context("when there is an error reading one of the sections", func() {
			It("returns the error after trying to read the remaining sections", func() {
				tmpDir, err := ioutil.TempDir("", "kep-collection")
				Expect(err).ToNot(HaveOccurred())
				defer os.RemoveAll(tmpDir)

				summaryFilename := "summary.md"
				motivationFilename := "motivation.md"

				locations := newMockLocationProvider()
				locations.SectionsOutput.Ret0 <- []string{summaryFilename, motivationFilename}
				locations.ContentDirOutput.Ret0 <- tmpDir

				_, err = sections.OpenCollection(locations)
				merr, ok := err.(*multierror.Error)
				Expect(ok).To(BeTrue())
				Expect(merr.Errors).To(HaveLen(2))
			})
		})

		Describe("#Persist", func() {
			It("persists the sections to disk", func() {
				tmpDir, err := ioutil.TempDir("", "kep-collection-test")
				Expect(err).ToNot(HaveOccurred())
				defer os.RemoveAll(tmpDir)

				sectionOne := newMockSection()
				sectionOne.PersistOutput.Ret0 <- nil // no error when calling Persist()

				sectionTwo := newMockSection()
				sectionTwo.PersistOutput.Ret0 <- nil // no error when calling Persist()

				col := sections.NewCollection(tmpDir, sectionOne, sectionTwo)
				By("persisting each section to disk")
				err = col.Persist()
				Expect(err).ToNot(HaveOccurred())

				Expect(<-sectionOne.PersistCalled).To(BeTrue())
				Expect(<-sectionTwo.PersistCalled).To(BeTrue())
			})

			Context("when a section cannot be written", func() {
				It("writes none of the sections", func() {
					sectionOne := newMockSection()
					sectionOne.PersistOutput.Ret0 <- errors.New("test error during persist")
					sectionOne.EraseOutput.Ret0 <- nil

					sectionTwo := newMockSection()
					sectionTwo.PersistOutput.Ret0 <- nil // no error when calling Persist()
					sectionTwo.EraseOutput.Ret0 <- nil

					col := sections.NewCollection("nonExistantContentRoot", sectionOne, sectionTwo)
					err := col.Persist()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("test error during persist"))

					By("calling erase on each section")
					Expect(<-sectionOne.PersistCalled).To(BeTrue())
					Expect(<-sectionTwo.PersistCalled).To(BeTrue())
					Expect(<-sectionOne.EraseCalled).To(BeTrue())
					Expect(<-sectionTwo.EraseCalled).To(BeTrue())
				})
			})
		})

		Describe("#Erase", func() {
			It("attempts to remove all sections", func() {
				sectionOne := newMockSection()
				sectionOne.EraseOutput.Ret0 <- nil

				sectionTwo := newMockSection()
				sectionTwo.EraseOutput.Ret0 <- nil

				col := sections.NewCollection("nonExistantContentRoot", sectionOne, sectionTwo)
				err := col.Erase()
				Expect(err).ToNot(HaveOccurred())

				Expect(<-sectionOne.EraseCalled).To(BeTrue())
				Expect(<-sectionTwo.EraseCalled).To(BeTrue())
			})
		})

		Describe("#Sections", func() {
			It("returns the filenames of the sections in the collection", func() {
				sectionOneFilename := "section_one.md"
				sectionTwoFilename := "section_two.md"

				sectionOne := newMockSection()
				sectionOne.FilenameOutput.Ret0 <- sectionOneFilename

				sectionTwo := newMockSection()
				sectionTwo.FilenameOutput.Ret0 <- sectionTwoFilename

				col := sections.NewCollection("nonExistantContentRoot", sectionOne, sectionTwo)
				secs := col.Sections()
				Expect(secs).To(ConsistOf(sectionOneFilename, sectionTwoFilename))
			})
		})
	})
})
