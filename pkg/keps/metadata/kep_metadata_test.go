package metadata_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"
	"os"
	"time"

	"github.com/calebamiles/keps/pkg/keps/metadata"
	"github.com/calebamiles/keps/pkg/keps/states"
)

var _ = Describe("KEP Metdata", func() {
	Describe("New()", func() {
		It("creates metadata", func() {
			author := "dchen1107"
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
			Expect(m.State()).To(Equal(states.Draft))
		})
	})

	Describe("#AddSections()", func() {
		It("adds section information to the metadata and dedupes", func() {
			author := "dchen1107"
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

			m.AddSectionLocations([]string{"test_section.md", "test_section.md"})
			Expect(m.SectionLocations()).To(HaveLen(1))
		})
	})

	Describe("#AddApprovers()", func() {
		It("adds approvers and dedupes", func() {
			author := "dchen1107"
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

			m.AddApprovers([]string{"bgrant0607", "bgrant0607"})
			Expect(m.Approvers()).To(HaveLen(1))
		})
	})

	Describe("#AddReviewers()", func() {
		It("adds reviewers and dedupes", func() {
			author := "dchen1107"
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

			m.AddReviewers([]string{"smarterclayton", "smarterclayton"})
			Expect(m.Reviewers()).To(HaveLen(1))
		})
	})

	Describe("#AddEvent()", func() {
		It("adds a KEP lifecycle event to metadata", func() {
			Fail("test not written")
		})
	})

	Describe("#AssociatePR()", func() {
		It("adds an associated PR to metadata", func() {
			Fail("test not written")
		})
	})

	Describe("#Events()", func() {
		It("returns KEP lifecycle event names in the order they occurred", func() {
			Fail("test not written")
		})
	})

	Describe("#SetState()", func() {
		It("sets the state on the metadata", func() {
			author := "dchen1107"
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

			m.SetState(states.Provisional)
			Expect(m.State()).To(Equal(states.Provisional))
		})
	})

	Describe("#Persist()", func() {
		It("writes a YAML representation to disk", func() {
			tmpDir, err := ioutil.TempDir("", "kep-content")
			Expect(err).ToNot(HaveOccurred())
			defer os.RemoveAll(tmpDir)

			author := "dchen1107"
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

			m.AddSectionLocations([]string{"test_section.md"})

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
			Expect(readMetadata.SectionLocations()).To(HaveLen(1))

			lastUpdated := readMetadata.LastUpdated()
			lastUpdatedMinute := lastUpdated.Round(time.Minute)
			nowish := now.Round(time.Minute)

			Expect(lastUpdatedMinute.Equal(nowish)).To(BeTrue())
		})
	})
})
