package section_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/calebamiles/keps/pkg/keps/section"
	"github.com/calebamiles/keps/pkg/keps/section/internal/rendering"
	"github.com/calebamiles/keps/pkg/keps/states"
)

var _ = Describe("A collection of sections", func() {
	Describe("common operations on collections", func() {
		Describe("#Persist", func() {
			It("persists the sections to disk", func() {
				title := "The Kubernetes Enhancement Proposal Process"
				authors := []string{"jbeda", "calebamiles"}
				owningSIG := "sig-architecture"
				kepState := states.Implementable
				now := time.Now().UTC()

				tmpDir, err := ioutil.TempDir("", "kep-collection-test")
				Expect(err).ToNot(HaveOccurred())
				defer os.RemoveAll(tmpDir)

				info := newMockRenderingInfoProvider()

				for i := 0; i < 3; i++ {
					info.TitleOutput.Ret0 <- title
					info.AuthorsOutput.Ret0 <- authors
					info.OwningSIGOutput.Ret0 <- owningSIG
					info.ContentDirOutput.Ret0 <- tmpDir
					info.StateOutput.Ret0 <- kepState
					info.LastUpdatedOutput.Ret0 <- now
				}

				col, err := section.ForProvisionalState(info)
				Expect(err).ToNot(HaveOccurred())

				summaryLoc := filepath.Join(tmpDir, rendering.SummaryFilename)
				motivationLoc := filepath.Join(tmpDir, rendering.MotivationFilename)

				Expect(summaryLoc).ToNot(BeAnExistingFile())
				Expect(motivationLoc).ToNot(BeAnExistingFile())

				By("persisting each section to disk")
				err = col.Persist()
				Expect(err).ToNot(HaveOccurred())

				Expect(summaryLoc).To(BeARegularFile())
				Expect(motivationLoc).To(BeARegularFile())

			})

			It("writes a README.md", func() {
				title := "The Kubernetes Enhancement Proposal Process"
				authors := []string{"jbeda", "calebamiles"}
				owningSIG := "sig-architecture"
				kepState := states.Implementable
				now := time.Now().UTC()

				tmpDir, err := ioutil.TempDir("", "kep-collection-test")
				Expect(err).ToNot(HaveOccurred())
				defer os.RemoveAll(tmpDir)

				info := newMockRenderingInfoProvider()

				for i := 0; i < 3; i++ {
					info.TitleOutput.Ret0 <- title
					info.AuthorsOutput.Ret0 <- authors
					info.OwningSIGOutput.Ret0 <- owningSIG
					info.ContentDirOutput.Ret0 <- tmpDir
					info.StateOutput.Ret0 <- kepState
					info.LastUpdatedOutput.Ret0 <- now
				}

				col, err := section.ForProvisionalState(info)
				Expect(err).ToNot(HaveOccurred())

				By("persisting a README.md to disk")
				loc := filepath.Join(tmpDir, rendering.ReadmeFilename)
				err = col.Persist()
				Expect(err).ToNot(HaveOccurred())

				Expect(loc).To(BeARegularFile())
			})

			Context("when a section cannot be written", func() {
				It("writes none of the sections", func() {
					title := "The Kubernetes Enhancement Proposal Process"
					authors := []string{"jbeda", "calebamiles"}
					owningSIG := "sig-architecture"
					kepState := states.Implementable
					now := time.Now().UTC()

					tmpDir, err := ioutil.TempDir("", "kep-collection-test")
					Expect(err).ToNot(HaveOccurred())
					defer os.RemoveAll(tmpDir)

					info := newMockRenderingInfoProvider()

					// first call: ok
					info.TitleOutput.Ret0 <- title
					info.AuthorsOutput.Ret0 <- authors
					info.OwningSIGOutput.Ret0 <- owningSIG
					info.ContentDirOutput.Ret0 <- tmpDir
					info.StateOutput.Ret0 <- kepState
					info.LastUpdatedOutput.Ret0 <- now

					// second call: bummer
					info.TitleOutput.Ret0 <- title
					info.AuthorsOutput.Ret0 <- authors
					info.OwningSIGOutput.Ret0 <- owningSIG
					info.ContentDirOutput.Ret0 <- filepath.Join(tmpDir, "not-a-dir")
					info.StateOutput.Ret0 <- kepState
					info.LastUpdatedOutput.Ret0 <- now

					// third call: ok
					info.TitleOutput.Ret0 <- title
					info.AuthorsOutput.Ret0 <- authors
					info.OwningSIGOutput.Ret0 <- owningSIG
					info.ContentDirOutput.Ret0 <- tmpDir
					info.StateOutput.Ret0 <- kepState
					info.LastUpdatedOutput.Ret0 <- now

					col, err := section.ForProvisionalState(info)
					Expect(err).ToNot(HaveOccurred())

					summaryLoc := filepath.Join(tmpDir, rendering.SummaryFilename)
					motivationLoc := filepath.Join(tmpDir, rendering.MotivationFilename)

					By("erasing all sections")
					err = col.Persist()
					Expect(err).To(HaveOccurred())

					Expect(summaryLoc).ToNot(BeAnExistingFile())
					Expect(motivationLoc).ToNot(BeAnExistingFile())
				})
			})
		})

		Describe("#Erase", func() {
			It("attempts to remove all sections", func() {
				title := "The Kubernetes Enhancement Proposal Process"
				authors := []string{"jbeda", "calebamiles"}
				owningSIG := "sig-architecture"
				kepState := states.Implementable
				now := time.Now().UTC()

				tmpDir, err := ioutil.TempDir("", "kep-collection-test")
				Expect(err).ToNot(HaveOccurred())
				defer os.RemoveAll(tmpDir)

				info := newMockRenderingInfoProvider()
				for i := 0; i < 3; i++ {
					info.TitleOutput.Ret0 <- title
					info.AuthorsOutput.Ret0 <- authors
					info.OwningSIGOutput.Ret0 <- owningSIG
					info.ContentDirOutput.Ret0 <- tmpDir
					info.StateOutput.Ret0 <- kepState
					info.LastUpdatedOutput.Ret0 <- now
				}

				col, err := section.ForProvisionalState(info)
				Expect(err).ToNot(HaveOccurred())

				summaryLoc := filepath.Join(tmpDir, rendering.SummaryFilename)
				motivationLoc := filepath.Join(tmpDir, rendering.MotivationFilename)

				err = col.Persist()
				Expect(err).ToNot(HaveOccurred())

				Expect(summaryLoc).To(BeARegularFile())
				Expect(motivationLoc).To(BeARegularFile())

				By("erasing the section content from disk")
				err = col.Erase()
				Expect(err).ToNot(HaveOccurred())

				Expect(summaryLoc).ToNot(BeAnExistingFile())
				Expect(motivationLoc).ToNot(BeAnExistingFile())
			})
		})

		Describe("#Sections", func() {
			It("returns basic info for the sections in the collection", func() {
				title := "The Kubernetes Enhancement Proposal Process"
				authors := []string{"jbeda", "calebamiles"}
				owningSIG := "sig-architecture"
				kepState := states.Implementable
				now := time.Now().UTC()

				tmpDir, err := ioutil.TempDir("", "kep-collection-test")
				Expect(err).ToNot(HaveOccurred())
				defer os.RemoveAll(tmpDir)

				info := newMockRenderingInfoProvider()
				for i := 0; i < 3; i++ {
					info.TitleOutput.Ret0 <- title
					info.AuthorsOutput.Ret0 <- authors
					info.OwningSIGOutput.Ret0 <- owningSIG
					info.ContentDirOutput.Ret0 <- tmpDir
					info.StateOutput.Ret0 <- kepState
					info.LastUpdatedOutput.Ret0 <- now
				}

				col, err := section.ForProvisionalState(info)
				Expect(err).ToNot(HaveOccurred())

				secs := col.Sections()
				sectionNames := []string{}

				for i := range secs {
					sectionNames = append(sectionNames, secs[i].Name())
				}

				Expect(sectionNames).To(ContainElement(rendering.SummaryName))
				Expect(sectionNames).To(ContainElement(rendering.MotivationName))
			})
		})
	})
})
