package sections_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/calebamiles/keps/pkg/keps/sections"
	"github.com/calebamiles/keps/pkg/keps/sections/internal/rendering"
	"github.com/calebamiles/keps/pkg/keps/states"
)

var _ = Describe("Creating a KEP section", func() {
	Describe("New()", func() {
		Context("when the section name does not exist", func() {
			It("returns an error", func() {
				_, err := sections.New("this-is-not-a-section", nil)
				Expect(err.Error()).To(ContainSubstring("no top level KEP section: this-is-not-a-section exists"))
			})
		})

		Context("when the section is the README", func() {
			It("returns an error", func() {
				_, err := sections.New(rendering.ReadmeName, nil)
				Expect(err).To(MatchError("cannot render README section using section.New(), use section.NewReadme()"))
			})
		})

		Context("when section content exists on disk", func() {
			It("reads the content on disk", func() {
				title := "The Kubernetes Enhancement Proposal Process"
				authors := []string{"jbeda", "calebamiles"}
				owningSIG := "sig-architecture"
				existingContent := []byte("existing content")

				tmpDir, err := ioutil.TempDir("", "kep-section-test")
				Expect(err).ToNot(HaveOccurred())
				defer os.RemoveAll(tmpDir)

				loc := filepath.Join(tmpDir, rendering.SummaryFilename)
				err = ioutil.WriteFile(loc, existingContent, os.ModePerm)
				Expect(err).ToNot(HaveOccurred())

				info := newMockRenderingInfoProvider()
				info.TitleOutput.Ret0 <- title
				info.AuthorsOutput.Ret0 <- authors
				info.OwningSIGOutput.Ret0 <- owningSIG
				info.ContentDirOutput.Ret0 <- tmpDir

				sec, err := sections.New(rendering.SummaryName, info)
				Expect(err).ToNot(HaveOccurred())

				Expect(sec.Content()).To(Equal(existingContent))
			})

			It("does not allow the content to be erased", func() {
				title := "The Kubernetes Enhancement Proposal Process"
				authors := []string{"jbeda", "calebamiles"}
				owningSIG := "sig-architecture"
				existingContent := []byte("existing content")

				tmpDir, err := ioutil.TempDir("", "kep-section-test")
				Expect(err).ToNot(HaveOccurred())
				defer os.RemoveAll(tmpDir)

				loc := filepath.Join(tmpDir, rendering.SummaryFilename)
				err = ioutil.WriteFile(loc, existingContent, os.ModePerm)
				Expect(err).ToNot(HaveOccurred())

				infoBeforeErase, err := os.Stat(loc)
				Expect(err).ToNot(HaveOccurred())

				info := newMockRenderingInfoProvider()
				info.TitleOutput.Ret0 <- title
				info.AuthorsOutput.Ret0 <- authors
				info.OwningSIGOutput.Ret0 <- owningSIG
				info.ContentDirOutput.Ret0 <- tmpDir

				sec, err := sections.New(rendering.SummaryName, info)
				Expect(err).ToNot(HaveOccurred())

				err = sec.Erase()
				Expect(err).ToNot(HaveOccurred())

				infoAfterErase, err := os.Stat(loc)
				Expect(err).ToNot(HaveOccurred())

				Expect(infoBeforeErase.ModTime()).To(Equal(infoAfterErase.ModTime()))
			})

			It("does not allow the content to be persisted again", func() {
				title := "The Kubernetes Enhancement Proposal Process"
				authors := []string{"jbeda", "calebamiles"}
				owningSIG := "sig-architecture"
				existingContent := []byte("existing content")

				tmpDir, err := ioutil.TempDir("", "kep-section-test")
				Expect(err).ToNot(HaveOccurred())
				defer os.RemoveAll(tmpDir)

				loc := filepath.Join(tmpDir, rendering.SummaryFilename)
				err = ioutil.WriteFile(loc, existingContent, os.ModePerm)
				Expect(err).ToNot(HaveOccurred())

				infoBeforeErase, err := os.Stat(loc)
				Expect(err).ToNot(HaveOccurred())

				info := newMockRenderingInfoProvider()
				info.TitleOutput.Ret0 <- title
				info.AuthorsOutput.Ret0 <- authors
				info.OwningSIGOutput.Ret0 <- owningSIG
				info.ContentDirOutput.Ret0 <- tmpDir

				sec, err := sections.New(rendering.SummaryName, info)
				Expect(err).ToNot(HaveOccurred())

				err = sec.Persist()
				Expect(err).ToNot(HaveOccurred())

				infoAfterErase, err := os.Stat(loc)
				Expect(err).ToNot(HaveOccurred())

				Expect(infoBeforeErase.ModTime()).To(Equal(infoAfterErase.ModTime()))
			})
		})

		Context("when section content does not exist on disk", func() {
			It("renders a new section", func() {
				title := "The Kubernetes Enhancement Proposal Process"
				authors := []string{"jbeda", "calebamiles"}
				owningSIG := "sig-architecture"

				tmpDir, err := ioutil.TempDir("", "kep-section-test")
				Expect(err).ToNot(HaveOccurred())
				defer os.RemoveAll(tmpDir)

				info := newMockRenderingInfoProvider()
				info.TitleOutput.Ret0 <- title
				info.AuthorsOutput.Ret0 <- authors
				info.OwningSIGOutput.Ret0 <- owningSIG
				info.ContentDirOutput.Ret0 <- tmpDir

				sec, err := sections.New(rendering.SummaryName, info)
				Expect(err).ToNot(HaveOccurred())

				Expect(sec.Content()).To(ContainSubstring(title))
				Expect(sec.Content()).To(ContainSubstring(rendering.SummaryName))
			})

			It("allows the content to be erased", func() {
				title := "The Kubernetes Enhancement Proposal Process"
				authors := []string{"jbeda", "calebamiles"}
				owningSIG := "sig-architecture"

				tmpDir, err := ioutil.TempDir("", "kep-section-test")
				Expect(err).ToNot(HaveOccurred())
				defer os.RemoveAll(tmpDir)

				loc := filepath.Join(tmpDir, rendering.SummaryFilename)

				info := newMockRenderingInfoProvider()
				info.TitleOutput.Ret0 <- title
				info.AuthorsOutput.Ret0 <- authors
				info.OwningSIGOutput.Ret0 <- owningSIG
				info.ContentDirOutput.Ret0 <- tmpDir

				sec, err := sections.New(rendering.SummaryName, info)
				Expect(err).ToNot(HaveOccurred())

				err = sec.Persist()
				Expect(err).ToNot(HaveOccurred())

				readContent, err := ioutil.ReadFile(loc)
				Expect(err).ToNot(HaveOccurred())
				Expect(readContent).To(Equal(sec.Content()))

				err = sec.Erase()
				Expect(err).ToNot(HaveOccurred())

				Expect(loc).ToNot(BeAnExistingFile())
			})

			It("allows the content to be persisted", func() {
				title := "The Kubernetes Enhancement Proposal Process"
				authors := []string{"jbeda", "calebamiles"}
				owningSIG := "sig-architecture"

				tmpDir, err := ioutil.TempDir("", "kep-section-test")
				Expect(err).ToNot(HaveOccurred())
				defer os.RemoveAll(tmpDir)

				loc := filepath.Join(tmpDir, rendering.SummaryFilename)

				info := newMockRenderingInfoProvider()
				info.TitleOutput.Ret0 <- title
				info.AuthorsOutput.Ret0 <- authors
				info.OwningSIGOutput.Ret0 <- owningSIG
				info.ContentDirOutput.Ret0 <- tmpDir

				sec, err := sections.New(rendering.SummaryName, info)
				Expect(err).ToNot(HaveOccurred())

				err = sec.Persist()
				Expect(err).ToNot(HaveOccurred())

				readContent, err := ioutil.ReadFile(loc)
				Expect(err).ToNot(HaveOccurred())
				Expect(readContent).To(Equal(sec.Content()))
			})
		})

	})

	Describe("NewReadme()", func() {
		It("renders a new README", func() {
			title := "The Kubernetes Enhancement Proposal Process"
			authors := []string{"jbeda", "calebamiles"}
			owningSIG := "sig-architecture"

			tmpDir, err := ioutil.TempDir("", "kep-section-test")
			Expect(err).ToNot(HaveOccurred())
			defer os.RemoveAll(tmpDir)

			loc := filepath.Join(tmpDir, rendering.ReadmeFilename)
			now := time.Now().UTC()
			info := newMockRenderingInfoProvider()

			for i := 0; i < 4; i++ {
				info.TitleOutput.Ret0 <- title
				info.AuthorsOutput.Ret0 <- authors
				info.OwningSIGOutput.Ret0 <- owningSIG
				info.ContentDirOutput.Ret0 <- tmpDir
				info.StateOutput.Ret0 <- states.Provisional
				info.LastUpdatedOutput.Ret0 <- now
			}

			summary, err := sections.New(rendering.SummaryName, info)
			Expect(err).ToNot(HaveOccurred())

			ss := []rendering.SectionProvider{summary}

			infoWithSections := &enhancedInfo{
				InfoProvider: info,
				ss:           ss,
			}

			readme, err := sections.NewReadme(infoWithSections)
			Expect(err).ToNot(HaveOccurred())

			err = readme.Persist()
			Expect(err).ToNot(HaveOccurred())

			Expect(loc).To(BeARegularFile())

			readmeBytes, err := ioutil.ReadFile(loc)
			Expect(err).ToNot(HaveOccurred())
			Expect(readmeBytes).To(Equal(readme.Content()))
		})
	})
})

type enhancedInfo struct {
	ss []rendering.SectionProvider
	rendering.InfoProvider
}

func (i *enhancedInfo) Sections() []rendering.SectionProvider {
	return i.ss
}
