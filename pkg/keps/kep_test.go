package keps_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"

	"github.com/calebamiles/keps/pkg/keps/metadata"
	"github.com/calebamiles/keps/pkg/keps/metadata/metadatafakes"
	"github.com/calebamiles/keps/pkg/keps/sections"
	"github.com/calebamiles/keps/pkg/keps/sections/sectionsfakes"
	"github.com/calebamiles/keps/pkg/keps/states"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/calebamiles/keps/pkg/keps"
)

var _ = Describe("A KEP", func() {
	Describe("New()", func() {
		It("ensures that the KEP is valid for its current status", func() {
			now := time.Now()
			before := now.Add(-time.Hour)
			uniqueID := uuid.New().String()
			title := "Deprecate The Kubernetes Enhancement Proposal Process"
			owningSIG := "sig-architecture"
			authors := []string{"jbeda", "calebamiles"}
			state := states.Provisional

			emptyContent := []sections.Entry{}

			fakeMetadata := &metadatafakes.FakeKEP{}
			fakeMetadata.AuthorsReturns(authors)
			fakeMetadata.CreatedReturns(before)
			fakeMetadata.LastUpdatedReturns(now)
			fakeMetadata.TitleReturns(title)
			fakeMetadata.OwningSIGReturns(owningSIG)
			fakeMetadata.StateReturns(state)
			fakeMetadata.UniqueIDReturns(uniqueID)

			By("returning an error if a required condition for KEP at a given state is not met")

			_, err := keps.New(fakeMetadata, emptyContent)
			Expect(err).To(HaveOccurred(), "creating a KEP should return an error if the metadata is invalid for the given state")
			Expect(err.Error()).To(ContainSubstring("missing Motivation"), "a valid KEP with `provisional` state should include a motivation")
			Expect(err.Error()).To(ContainSubstring("missing Summary"), "a valid KEP with `provisional` state should include a summary")
			Expect(err.Error()).To(ContainSubstring("has 0 short ID which should be unset for provisional KEPs"), "a valid KEP with `provisional` state should not have claimed a short ID (int)")
		})

		It("adds any section locations to the metadata", func() {
			now := time.Now()
			before := now.Add(-time.Hour)
			uniqueID := uuid.New().String()
			title := "The Kubernetes Enhancement Proposal Process"
			owningSIG := "sig-architecture"
			authors := []string{"jbeda", "calebamiles"}
			contentDir := "content/kubernetes-wide/kubernetes-enhancement-proposal-proccess"

			sectionOneName := "Section One"
			sectionOneFilename := "section_one.md"

			sectionOne := &sectionsfakes.FakeEntry{}
			sectionOne.NameReturns(sectionOneName)
			sectionOne.FilenameReturns(sectionOneFilename)

			sectionTwoName := "Section Two"
			sectionTwoFilename := "section_two.md"

			sectionTwo := &sectionsfakes.FakeEntry{}
			sectionTwo.NameReturns(sectionTwoName)
			sectionTwo.FilenameReturns(sectionTwoFilename)

			fakeSections := []sections.Entry{sectionOne, sectionTwo}

			fakeMetadata := &metadatafakes.FakeKEP{}
			fakeMetadata.AuthorsReturns(authors)
			fakeMetadata.ContentDirReturns(contentDir)
			fakeMetadata.CreatedReturns(before)
			fakeMetadata.LastUpdatedReturns(now)
			fakeMetadata.TitleReturns(title)
			fakeMetadata.OwningSIGReturns(owningSIG)
			fakeMetadata.UniqueIDReturns(uniqueID)
			fakeMetadata.StateReturns(states.Draft) // use `draft` to avoid failing checks for `provisional` state

			k, err := keps.New(fakeMetadata, fakeSections)
			Expect(err).ToNot(HaveOccurred(), "expected no error when creating a new KEP with existing sections")

			includedSections := k.Sections()
			Expect(includedSections).To(HaveLen(2), "expected two sections passed to kep.New() to be returned")
		})
	})

	Describe("#AddCheck()", func() {
		It("adds checks that are called by Check()", func() {
			now := time.Now()
			before := now.Add(-time.Hour)
			uniqueID := uuid.New().String()
			title := "The Kubernetes Enhancement Proposal Process"
			owningSIG := "sig-architecture"
			authors := []string{"jbeda", "calebamiles"}
			contentDir := "content/kubernetes-wide/kubernetes-enhancement-proposal-proccess"

			fakeMetadata := &metadatafakes.FakeKEP{}
			fakeMetadata.AuthorsReturns(authors)
			fakeMetadata.ContentDirReturns(contentDir)
			fakeMetadata.CreatedReturns(before)
			fakeMetadata.LastUpdatedReturns(now)
			fakeMetadata.TitleReturns(title)
			fakeMetadata.OwningSIGReturns(owningSIG)
			fakeMetadata.UniqueIDReturns(uniqueID)
			fakeMetadata.StateReturns(states.Draft) // use `draft` to avoid failing checks for `provisional` state

			emptySections := []sections.Entry{}

			k, err := keps.New(fakeMetadata, emptySections)
			Expect(err).ToNot(HaveOccurred(), "expected no error when creating a new KEP with valid metadata and no existing sections")

			expectedCheckError := errors.New("an expected error occurred")
			var fakeCheckOne = func(_ metadata.KEP) error {
				return expectedCheckError
			}

			var fakeCheckTwo = func(_ metadata.KEP) error {
				return expectedCheckError
			}

			k.AddChecks(fakeCheckOne, fakeCheckTwo)
			err = k.Check()

			Expect(err).To(HaveOccurred(), "expected calling Check() on a KEP with failing checks to return error")

			merr, ok := err.(*multierror.Error)
			Expect(ok).To(BeTrue(), "typcasting an error returned by Check() as a hashicorp/go-multierror should be ok")

			Expect(merr.Errors).To(HaveLen(2), "expected two errors to be contained in the error returned from Check()")
			Expect(merr.Errors[0]).To(MatchError(expectedCheckError), "expected error returned from Check() to contain the injected error")
			Expect(merr.Errors[1]).To(MatchError(expectedCheckError), "expected error returned from Check() to contain the injected error")
		})
	})

	Describe("#SetState()", func() {
		It("attempts to set the state on the KEP", func() {
			now := time.Now()
			before := now.Add(-time.Hour)
			uniqueID := uuid.New().String()
			title := "The Kubernetes Enhancement Proposal Process"
			owningSIG := "sig-architecture"
			authors := []string{"jbeda", "calebamiles"}
			contentDir := "content/kubernetes-wide/kubernetes-enhancement-proposal-proccess"

			fakeMetadata := &metadatafakes.FakeKEP{}
			fakeMetadata.AuthorsReturns(authors)
			fakeMetadata.ContentDirReturns(contentDir)
			fakeMetadata.CreatedReturns(before)
			fakeMetadata.LastUpdatedReturns(now)
			fakeMetadata.TitleReturns(title)
			fakeMetadata.OwningSIGReturns(owningSIG)
			fakeMetadata.UniqueIDReturns(uniqueID)
			fakeMetadata.StateReturns(states.Draft) // use `draft` to avoid failing checks for `provisional` state

			k, err := keps.New(fakeMetadata, []sections.Entry{})
			Expect(err).ToNot(HaveOccurred(), "expected no error when creating a KEP with valid metadata and no sections")

			By("adding any missing sections and running consistency checks for the desired state")

			err = k.SetState(states.Provisional)
			Expect(err).ToNot(HaveOccurred(), "expected no error when setting a new KEP to `provisional` state")

			includedSections := k.Sections()
			Expect(includedSections).To(ConsistOf(sections.Summary, sections.Motivation), "expected KEP to have two sections: Summary, and Motivation")
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

			fakeMetadata := &metadatafakes.FakeKEP{}
			fakeMetadata.AuthorsReturns(authors)
			fakeMetadata.ContentDirReturns(contentDir)
			fakeMetadata.CreatedReturns(before)
			fakeMetadata.LastUpdatedReturns(now)
			fakeMetadata.TitleReturns(title)
			fakeMetadata.OwningSIGReturns(owningSIG)
			fakeMetadata.UniqueIDReturns(uniqueID)
			fakeMetadata.StateReturns(states.Draft) // use `draft` to avoid failing checks for `provisional` state

			fakeContent := []sections.Entry{}

			By("exposing a limited set of fields of KEP metadata")

			k, err := keps.New(fakeMetadata, fakeContent)
			Expect(err).ToNot(HaveOccurred(), "creating a KEP with keps.New() should not return an error with valid imput")

			Expect(k.UniqueID()).To(Equal(uniqueID), "a KEP should expose the UUID of its metadata")
			Expect(k.Title()).To(Equal(title), "a KEP should expose the title of its metadata")
			Expect(k.OwningSIG()).To(Equal(owningSIG), "a KEP should expose the owning SIG of its metadata")
			Expect(k.Authors()).To(Equal(authors), "a KEP should expose the authors of its metadata")
			Expect(k.Created()).To(Equal(before), "a KEP should expose the created date of its metadata")
			Expect(k.LastUpdated()).To(Equal(now), "a KEP should expose the last updated date of its metadata")
			Expect(k.ContentDir()).To(Equal(contentDir), "a KEP should expose the contentDir of its metadata")
		})
	})

	Describe("Persist()", func() {
		It("persists the metadata and section content", func() {
			tmpDir, err := ioutil.TempDir("", "kep-persist-test")
			Expect(err).ToNot(HaveOccurred(), "creating a temp directory should not return an error")
			defer os.RemoveAll(tmpDir)

			now := time.Now()
			before := now.Add(-time.Hour)
			uniqueID := uuid.New().String()
			title := "The Kubernetes Enhancement Proposal Process"
			owningSIG := "sig-architecture"
			authors := []string{"jbeda", "calebamiles"}

			fakeMetadata := &metadatafakes.FakeKEP{}
			fakeMetadata.AuthorsReturns(authors)
			fakeMetadata.ContentDirReturns(tmpDir)
			fakeMetadata.CreatedReturns(before)
			fakeMetadata.LastUpdatedReturns(now)
			fakeMetadata.TitleReturns(title)
			fakeMetadata.OwningSIGReturns(owningSIG)
			fakeMetadata.UniqueIDReturns(uniqueID)
			fakeMetadata.StateReturns(states.Draft) // use `draft` to avoid failing checks for `provisional` state

			fakeMetadata.PersistReturns(nil)

			sectionOneName := "Section One"
			sectionOneFilename := "section_one.md"

			sectionOne := &sectionsfakes.FakeEntry{}
			sectionOne.NameReturns(sectionOneName)
			sectionOne.FilenameReturns(sectionOneFilename)
			sectionOne.PersistReturns(nil)

			sectionTwoName := "Section Two"
			sectionTwoFilename := "section_two.md"

			sectionTwo := &sectionsfakes.FakeEntry{}
			sectionTwo.NameReturns(sectionTwoName)
			sectionTwo.FilenameReturns(sectionTwoFilename)
			sectionTwo.PersistReturns(nil)

			fakeSections := []sections.Entry{sectionOne, sectionTwo}

			kep, err := keps.New(fakeMetadata, fakeSections)
			Expect(err).ToNot(HaveOccurred(), "creating a new KEP with valid metadata and no section content should not return error")

			By("persisting KEP metadata")

			err = kep.Persist()
			Expect(err).ToNot(HaveOccurred(), "expected no error when persisting a valid KEP")

			Expect(fakeMetadata.PersistCallCount()).To(Equal(1), "expected Persist() on KEP metadata to be called once when Persist() is called on the parent KEP")

			By("persisting section content")

			Expect(sectionOne.PersistCallCount()).To(Equal(1), "expected Persist() to be called on Section One during parent KEP Persist()")
			Expect(sectionTwo.PersistCallCount()).To(Equal(1), "expected Persist() to be called on Section Two during parent KEP Persist()")

			err = kep.Persist()
			Expect(err).ToNot(HaveOccurred(), "expected no error when persisting a valid KEP")

			expectedReadmePath := filepath.Join(tmpDir, "README.md")
			Expect(expectedReadmePath).To(BeARegularFile(), "expected README.md to be autogenerated during Persist()")
		})

		XIt("attempts to roll back changes if an error occured", func() {
			Skip("not implemented")
		})
	})
})
