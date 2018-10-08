package metadata_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"
	"os"
	"time"

	"github.com/calebamiles/keps/pkg/keps/metadata"
	"github.com/calebamiles/keps/pkg/keps/states"
	"github.com/calebamiles/keps/pkg/sigs/sigsfakes"
)

var _ = Describe("KEP Metdata", func() {
	Describe("New()", func() {
		It("creates metadata", func() {
			author := "Dawn Chen"
			title := "kubelet"
			owningSIG := "sig-node"
			subproject := "kubelet"

			info := new(sigsfakes.FakeRoutingInfo)
			info.OwningSIGReturns(owningSIG)
			info.AffectedSubprojectsReturns([]string{subproject})
			info.SIGWideReturns(true)

			m, err := metadata.New([]string{author}, title, info)
			Expect(err).ToNot(HaveOccurred())

			Expect(m.Authors()).To(ContainElement(author))
			Expect(m.Title()).To(Equal(title))
			Expect(m.OwningSIG()).To(Equal(owningSIG))
			Expect(m.AffectedSubprojects()).To(ContainElement(subproject))
			Expect(m.State()).To(Equal(states.Provisional))
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
			subproject := "kubelet"
			now := time.Now()

			info := new(sigsfakes.FakeRoutingInfo)
			info.OwningSIGReturns(owningSIG)
			info.AffectedSubprojectsReturns([]string{subproject})
			info.SIGWideReturns(true)
			info.ContentDirReturns(tmpDir)

			m, err := metadata.New([]string{author}, title, info)
			Expect(err).ToNot(HaveOccurred())

			err = m.Persist()
			Expect(err).ToNot(HaveOccurred())

			contentDir := m.ContentDir()
			Expect(contentDir).To(Equal(tmpDir))

			readMetadata, err := metadata.Open(tmpDir)
			Expect(err).ToNot(HaveOccurred())

			Expect(readMetadata.Authors()).To(ContainElement(author))
			Expect(readMetadata.Title()).To(Equal(title))
			Expect(readMetadata.OwningSIG()).To(Equal(owningSIG))
			Expect(readMetadata.AffectedSubprojects()).To(ContainElement(subproject))

			lastUpdated := readMetadata.LastUpdated()
			lastUpdatedMinute := lastUpdated.Round(time.Minute)
			nowish := now.Round(time.Minute)

			Expect(lastUpdatedMinute.Equal(nowish)).To(BeTrue())
		})
	})
})
