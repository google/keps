package keps_test

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/calebamiles/keps/pkg/keps"
	"github.com/calebamiles/keps/pkg/keps/states"

	"github.com/calebamiles/keps/pkg/keps/metadata/metadatafakes"
	"github.com/calebamiles/keps/pkg/keps/sections/sectionsfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("A KEP", func() {
	Describe("Check()", func() {
		It("ensures that basic KEP invariants are satisified", func() {
			now := time.Now()
			before := now.Add(-time.Hour)
			title := "Deprecate The Kubernetes Enhancement Proposal Process"
			owningSIG := "sig-architecture"
			authors := []string{"jbeda", "calebamiles"}
			state := states.Rejected

			fakeContent := &sectionsfakes.FakeCollection{}
			fakeMetadata := &metadatafakes.FakeKEP{}
			fakeMetadata.AuthorsReturns(authors)
			fakeMetadata.CreatedReturns(before)
			fakeMetadata.LastUpdatedReturns(now)
			fakeMetadata.TitleReturns(title)
			fakeMetadata.OwningSIGReturns(owningSIG)
			fakeMetadata.StateReturns(state)

			k, err := keps.New(fakeMetadata, fakeContent)
			Expect(err).ToNot(HaveOccurred())

			By("returning an error when an invariant is not satisfied")
			err = k.Check()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("empty string given as UUID"))
		})

		It("ensures that the KEP is valid for its current status", func() {
			now := time.Now()
			before := now.Add(-time.Hour)
			uniqueID := uuid.New().String()
			title := "Deprecate The Kubernetes Enhancement Proposal Process"
			owningSIG := "sig-architecture"
			authors := []string{"jbeda", "calebamiles"}
			state := states.Provisional

			fakeContent := &sectionsfakes.FakeCollection{}
			fakeMetadata := &metadatafakes.FakeKEP{}
			fakeMetadata.AuthorsReturns(authors)
			fakeMetadata.CreatedReturns(before)
			fakeMetadata.LastUpdatedReturns(now)
			fakeMetadata.TitleReturns(title)
			fakeMetadata.OwningSIGReturns(owningSIG)
			fakeMetadata.StateReturns(state)
			fakeMetadata.UniqueIDReturns(uniqueID)

			k, err := keps.New(fakeMetadata, fakeContent)
			Expect(err).ToNot(HaveOccurred())

			By("returning an error if a required condition for KEP at a given state is not met")
			err = k.Check()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("missing Motivation"))
			Expect(err.Error()).To(ContainSubstring("missing Summary"))
			Expect(err.Error()).To(ContainSubstring("has 0 short ID which should be unset for provisional KEPs"))
		})
	})

	Describe("#AddCheck()", func() {
		XIt("adds checks that are called by Check()", func() {
			Fail("test not written")
		})
	})

	Describe("#SetState()", func() {
		XIt("attempts to set the state on the KEP", func() {
			Fail("test not written")
		})
	})

	Describe("Reading metadata from an instance", func() {
		It("exposes a limited set of metadata to provide human consumable index information", func() {
			now := time.Now()
			before := now.Add(-time.Hour)
			uniqueID := uuid.New().String()
			title := "The Kubernetes Enhancement Proposal Process"
			owningSIG := "sig-architecture"
			authors := []string{"jbeda", "calebamiles"}
			contentDir := "content/kubernetes-wide/kubernetes-enhancement-proposal-proccess"

			fakeContent := &sectionsfakes.FakeCollection{}
			fakeMetadata := &metadatafakes.FakeKEP{}
			fakeMetadata.AuthorsReturns(authors)
			fakeMetadata.ContentDirReturns(contentDir)
			fakeMetadata.CreatedReturns(before)
			fakeMetadata.LastUpdatedReturns(now)
			fakeMetadata.TitleReturns(title)
			fakeMetadata.OwningSIGReturns(owningSIG)
			fakeMetadata.UniqueIDReturns(uniqueID)

			By("exposing a limited set of fields of KEP metadata")
			k, err := keps.New(fakeMetadata, fakeContent)
			Expect(err).ToNot(HaveOccurred())

			Expect(k.UniqueID()).To(Equal(uniqueID))
			Expect(k.Title()).To(Equal(title))
			Expect(k.OwningSIG()).To(Equal(owningSIG))
			Expect(k.Authors()).To(Equal(authors))
			Expect(k.Created()).To(Equal(before))
			Expect(k.LastUpdated()).To(Equal(now))
			Expect(k.ContentDir()).To(Equal(contentDir))
		})
	})

	Describe("Persist()", func() {
		It("persists the metadata and section content", func() {
			fakeContent := &sectionsfakes.FakeCollection{}
			fakeMetadata := &metadatafakes.FakeKEP{}

			fakeContent.PersistReturns(nil)
			fakeMetadata.PersistReturns(nil)

			kep, err := keps.New(fakeMetadata, fakeContent)
			Expect(err).ToNot(HaveOccurred())

			err = kep.Persist()
			Expect(err).ToNot(HaveOccurred())
			Expect(fakeContent.PersistCallCount()).To(Equal(1))
			Expect(fakeMetadata.PersistCallCount()).To(Equal(1))
		})

		It("attempts to roll back changes if an error occured", func() {
			fakeContent := &sectionsfakes.FakeCollection{}
			fakeMetadata := &metadatafakes.FakeKEP{}

			fakeContent.PersistReturns(errors.New("fake error"))
			fakeMetadata.PersistReturns(nil)

			kep, err := keps.New(fakeMetadata, fakeContent)
			Expect(err).ToNot(HaveOccurred())

			err = kep.Persist()
			Expect(err.Error()).To(ContainSubstring("fake error"))
			Expect(fakeContent.EraseCallCount()).To(Equal(1))
		})
	})
})

const fakeSectionFilename = "fake_section.md"
