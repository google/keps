package metadata_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"
	"os"
	"time"

	"github.com/calebamiles/keps/pkg/keps/metadata"
	"github.com/calebamiles/keps/pkg/keps/sections"
	"github.com/calebamiles/keps/pkg/keps/states"
)

var _ = Describe("KEP Metdata", func() {
	Describe("New()", func() {
		It("creates metadata", func() {
			author := "Dawn Chen"
			title := "kubelet"
			owningSIG := "sig-node"
			subprojects := []string{"kubelet"}

			info := newMockRoutingInfoProvider()
			info.OwningSIGOutput.Ret0 <- owningSIG
			info.AffectedSubprojectsOutput.Ret0 <- subprojects
			info.SIGWideOutput.Ret0 <- true
			info.KubernetesWideOutput.Ret0 <- false
			info.ParticipatingSIGsOutput.Ret0 <- []string{}
			info.ContentDirOutput.Ret0 <- ""

			m, err := metadata.New([]string{author}, title, info)
			Expect(err).ToNot(HaveOccurred())

			Expect(m.Authors()).To(ContainElement(author))
			Expect(m.Title()).To(Equal(title))
			Expect(m.OwningSIG()).To(Equal(owningSIG))
			Expect(m.AffectedSubprojects()).To(ContainElement(subprojects[0]))
			Expect(m.State()).To(Equal(states.Provisional))
		})
	})

	Describe("#AddSections", func() {
		It("adds section information to the metadata", func() {
			author := "Dawn Chen"
			title := "kubelet"
			owningSIG := "sig-node"
			subprojects := []string{"kubelet"}

			info := newMockRoutingInfoProvider()
			info.OwningSIGOutput.Ret0 <- owningSIG
			info.AffectedSubprojectsOutput.Ret0 <- subprojects
			info.SIGWideOutput.Ret0 <- true
			info.KubernetesWideOutput.Ret0 <- false
			info.ParticipatingSIGsOutput.Ret0 <- []string{}
			info.ContentDirOutput.Ret0 <- ""

			m, err := metadata.New([]string{author}, title, info)
			Expect(err).ToNot(HaveOccurred())

			ss := []sections.Info{&testSection{}}
			m.AddSections(ss)
			Expect(m.Sections()).To(HaveLen(1))
		})

		It("dedupes by section name", func() {
			author := "Dawn Chen"
			title := "kubelet"
			owningSIG := "sig-node"
			subprojects := []string{"kubelet"}

			info := newMockRoutingInfoProvider()
			info.OwningSIGOutput.Ret0 <- owningSIG
			info.AffectedSubprojectsOutput.Ret0 <- subprojects
			info.SIGWideOutput.Ret0 <- true
			info.KubernetesWideOutput.Ret0 <- false
			info.ParticipatingSIGsOutput.Ret0 <- []string{}
			info.ContentDirOutput.Ret0 <- ""

			m, err := metadata.New([]string{author}, title, info)
			Expect(err).ToNot(HaveOccurred())

			ss := []sections.Info{&testSection{}}
			m.AddSections(ss)
			m.AddSections(ss)
			Expect(m.Sections()).To(HaveLen(1))
		})
	})

	Describe("#Persist()", func() {
		It("writes a YAML representation to disk", func() {
			tmpDir, err := ioutil.TempDir("", "kep-content")
			Expect(err).ToNot(HaveOccurred())
			defer os.RemoveAll(tmpDir)

			author := "Dawn Chen"
			title := "kubelet"
			owningSIG := "sig-node"
			subprojects := []string{"kubelet"}
			now := time.Now()

			info := newMockRoutingInfoProvider()
			info.OwningSIGOutput.Ret0 <- owningSIG
			info.AffectedSubprojectsOutput.Ret0 <- subprojects
			info.SIGWideOutput.Ret0 <- true
			info.KubernetesWideOutput.Ret0 <- false
			info.ParticipatingSIGsOutput.Ret0 <- []string{}
			info.ContentDirOutput.Ret0 <- tmpDir

			m, err := metadata.New([]string{author}, title, info)
			Expect(err).ToNot(HaveOccurred())

			ss := []sections.Info{&testSection{}}
			m.AddSections(ss)

			err = m.Persist()
			Expect(err).ToNot(HaveOccurred())

			contentDir := m.ContentDir()
			Expect(contentDir).To(Equal(tmpDir))

			readMetadata, err := metadata.Open(tmpDir)
			Expect(err).ToNot(HaveOccurred())

			Expect(readMetadata.Authors()).To(ContainElement(author))
			Expect(readMetadata.Title()).To(Equal(title))
			Expect(readMetadata.OwningSIG()).To(Equal(owningSIG))
			Expect(readMetadata.AffectedSubprojects()).To(ContainElement(subprojects[0]))
			Expect(readMetadata.Sections()).To(HaveLen(1))

			lastUpdated := readMetadata.LastUpdated()
			lastUpdatedMinute := lastUpdated.Round(time.Minute)
			nowish := now.Round(time.Minute)

			Expect(lastUpdatedMinute.Equal(nowish)).To(BeTrue())
		})
	})
})

type testSection struct{}

func (s *testSection) Filename() string { return testSectionName }
func (s *testSection) Name() string     { return testSectionFilename }

const (
	testSectionName     = "Test Section"
	testSectionFilename = "test_section.md"
)
