package check_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"

	"github.com/calebamiles/keps/pkg/keps/check"
	"github.com/calebamiles/keps/pkg/keps/metadata/metadatafakes"
)

var _ = Describe("Checking Metadata", func() {
	Describe("Checking that sections exist", func() {
		It("ensures that each section exists with content", func() {
			tmpDir, err := ioutil.TempDir("", "kep-checks")
			Expect(err).ToNot(HaveOccurred())
			defer os.RemoveAll(tmpDir)

			meta := &metadatafakes.FakeKEP{}
			meta.ContentDirReturns(tmpDir)

			By("returning no error for a valid section")
			testSectionFilename := "test_section.md"
			testSectionContent := []byte("some test section content")

			err = ioutil.WriteFile(filepath.Join(tmpDir, testSectionFilename), testSectionContent, os.ModePerm)
			Expect(err).ToNot(HaveOccurred())

			meta.SectionsReturns([]string{testSectionFilename})

			err = check.ThatAllSectionsExistWithContent(meta)
			Expect(err).ToNot(HaveOccurred())

			By("returning an error if there is no content")
			testSectionWithNoContentFilename := "no_content.md"
			f, err := os.Create(filepath.Join(tmpDir, testSectionWithNoContentFilename))
			Expect(err).ToNot(HaveOccurred())
			err = f.Close()
			Expect(err).ToNot(HaveOccurred())

			meta.SectionsReturns([]string{testSectionFilename, testSectionWithNoContentFilename})

			err = check.ThatAllSectionsExistWithContent(meta)
			merr, ok := err.(*multierror.Error)
			Expect(ok).To(BeTrue())

			Expect(merr.Errors).To(HaveLen(1))
			Expect(merr.Errors[0].Error()).To(ContainSubstring("Section contains no content"))

			By("returning an error if there is no section on disk")
			testMissingSectionFilename := "no_section.md"
			meta.SectionsReturns([]string{testSectionFilename, testSectionWithNoContentFilename, testMissingSectionFilename})

			err = check.ThatAllSectionsExistWithContent(meta)
			merr, ok = err.(*multierror.Error)
			Expect(ok).To(BeTrue())
			Expect(merr.Errors).To(HaveLen(2))
			Expect(merr.Errors[1].Error()).To(ContainSubstring("Section does not exist on disk"))
		})
	})

	Describe("Checking that SIGs exist", func() {
		It("ensures that each SIG exists", func() {
			meta := &metadatafakes.FakeKEP{}

			By("returning no errors if all SIGs listed exist")
			meta.OwningSIGReturns("sig-architecture")

			err := check.ThatAllSIGsExist(meta)
			Expect(err).ToNot(HaveOccurred())

			By("returning an error if the owning SIG does not exist")
			meta.OwningSIGReturns("sig-not-real-at-all")

			err = check.ThatAllSIGsExist(meta)
			merr, ok := err.(*multierror.Error)
			Expect(ok).To(BeTrue())
			Expect(merr.Errors).To(HaveLen(1))
			Expect(merr.Errors[0].Error()).To(ContainSubstring("invalid SIG: sig-not-real-at-all"))

			By("returning an error if a participating SIG does not exist")
			meta.ParticipatingSIGsReturns([]string{"sig-architecture", "sig-not-real-sorry"})

			err = check.ThatAllSIGsExist(meta)
			merr, ok = err.(*multierror.Error)
			Expect(ok).To(BeTrue())
			Expect(merr.Errors).To(HaveLen(2))
			Expect(merr.Errors[0].Error()).To(ContainSubstring("invalid SIG: sig-not-real-at-all"))
			Expect(merr.Errors[1].Error()).To(ContainSubstring("invalid SIG: sig-not-real-sorry"))
		})
	})

	Describe("Checking that the title is set", func() {
		It("ensures that the title is non empty", func() {
			meta := &metadatafakes.FakeKEP{}

			By("returing an error if the title is unset")
			err := check.ThatTitleIsSet(meta)

			Expect(err.Error()).To(ContainSubstring("no title set"))

			By("returning no error if the title is set")
			meta.TitleReturns("a great and large idea")

			err = check.ThatTitleIsSet(meta)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("Checking that the authors exist", func() {
		It("ensures that there is at least one non empty author", func() {
			meta := &metadatafakes.FakeKEP{}

			By("returing an error if there are no authors")
			err := check.ThatAuthorsExist(meta)
			Expect(err.Error()).To(ContainSubstring("no authors listed"))

			By("returning an error if there is an empty author")
			meta.AuthorsReturns([]string{"smartCoder", "", "anotherSmartCoder"})
			err = check.ThatAuthorsExist(meta)
			Expect(err.Error()).To(ContainSubstring("empty string given for author"))

			By("not returning an error if evertything is fine")
			meta.AuthorsReturns([]string{"smartCoder"})
			err = check.ThatAuthorsExist(meta)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("Checking that a KEP exists upstream", func() {
		XIt("ensures that the KEP exists upstream", func() {
			Fail("test not written until migration")
		})
	})

	Describe("Checking that a KEP has been accepted upstream", func() {
		XIt("checks that the KEP exists upstream with provisional status or greater", func() {
			Fail("test not written until migration")
		})
	})

	Describe("Checking that subprojects exist", func() {
		It("ensures that subproject information is valid", func() {
			meta := &metadatafakes.FakeKEP{}

			By("returning no error if no subproject information is given")
			meta.AffectedSubprojectsReturns([]string{})
			err := check.ThatSubprojectsExist(meta)
			Expect(err).ToNot(HaveOccurred())

			By("returning no error if the subproject is valid")
			meta.AffectedSubprojectsReturns([]string{"kubelet"})
			err = check.ThatSubprojectsExist(meta)
			Expect(err).ToNot(HaveOccurred())

			By("returning an error if the subproject is invalid")
			meta.AffectedSubprojectsReturns([]string{"nobody-would-fund-this"})
			err = check.ThatSubprojectsExist(meta)
			Expect(err.Error()).To(ContainSubstring("invalid subproject"))
		})
	})

	Describe("Checking that the owning SIG is set", func() {
		It("ensures that an owning SIG has been set", func() {
			meta := &metadatafakes.FakeKEP{}

			By("returning an error if no owning SIG has been set")
			meta.OwningSIGReturns("")
			err := check.ThatHasOwningSIG(meta)
			Expect(err.Error()).To(ContainSubstring("Empty SIG information"))

			By("returning an error if the owning SIG does not exist")
			meta.OwningSIGReturns("nope-not-real")
			err = check.ThatHasOwningSIG(meta)
			Expect(err.Error()).To(ContainSubstring("Invalid owning SIG"))

			By("returning no error if the owning SIG is set and valid")
			meta.OwningSIGReturns("sig-node")
			err = check.ThatHasOwningSIG(meta)
			Expect(err).ToNot(HaveOccurred())
		})

	})

	Describe("Checking that the author is not an approver", func() {
		It("ensures that the author is not in the set of approvers", func() {
			meta := &metadatafakes.FakeKEP{}

			By("returning an error if the author is in the set of approvers")
			meta.AuthorsReturns([]string{"aSmartCoder"})
			meta.ApproversReturns([]string{"aFairReviewer", "aToughReviewer", "aSmartCoder"})

			err := check.ThatAuthorIsNotApprover(meta)
			Expect(err.Error()).To(ContainSubstring("is listed as both an author and approver"))

			By("not returning an error when the author is not in the set of approvers")
			meta.AuthorsReturns([]string{"aSmartCoder"})
			meta.ApproversReturns([]string{"aFairReviewer", "aToughReviewer"})

			err = check.ThatAuthorIsNotApprover(meta)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("Checking that the author is not a reviewer", func() {
		It("ensures that the reviewer is not in the set of reviewers", func() {
			meta := &metadatafakes.FakeKEP{}

			By("returning an error if the author is in the set of reviewers")
			meta.AuthorsReturns([]string{"aSmartCoder"})
			meta.ReviewersReturns([]string{"aFairReviewer", "aToughReviewer", "aSmartCoder"})

			err := check.ThatAuthorIsNotReviewer(meta)
			Expect(err.Error()).To(ContainSubstring("is listed as both an author and reviewer"))

			By("not returning an error when the author is not in the set of reviewers")
			meta.AuthorsReturns([]string{"aSmartCoder"})
			meta.ReviewersReturns([]string{"aFairReviewer", "aToughReviewer"})

			err = check.ThatAuthorIsNotReviewer(meta)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("Checking that there are editors", func() {
		It("ensures that there are editors", func() {
			meta := &metadatafakes.FakeKEP{}

			By("returning an error if there are no editors")
			meta.EditorsReturns([]string{})
			err := check.ThatThereAreEditors(meta)
			Expect(err.Error()).To(ContainSubstring("no editors"))

			By("returning an error if an editor is invalid")

			By("returning no error if there are editors")
			meta.EditorsReturns([]string{})
			err = check.ThatThereAreEditors(meta)
			Expect(err.Error()).To(ContainSubstring("no editors"))
		})
	})

	Describe("Checking that there are reviewers", func() {
		It("ensures that there are reviewers", func() {
			meta := &metadatafakes.FakeKEP{}

			By("returning an error if there are no reviewers")
			meta.ReviewersReturns([]string{})
			err := check.ThatThereAreReviewers(meta)
			Expect(err.Error()).To(ContainSubstring("no reviewers"))

			By("returning an error if a reviewer is invalid")
			meta.ReviewersReturns([]string{""})
			err = check.ThatThereAreReviewers(meta)
			Expect(err.Error()).To(ContainSubstring("invalid reviewer"))

			By("returning no error if there are reviewers")
			meta.ReviewersReturns([]string{"aFairReviewer"})
			err = check.ThatThereAreReviewers(meta)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("Checking that there are approvers", func() {
		It("ensures that there are approvers", func() {
			meta := &metadatafakes.FakeKEP{}

			By("returning an error if there are no approvers")
			meta.ApproversReturns([]string{})
			err := check.ThatThereAreApprovers(meta)
			Expect(err.Error()).To(ContainSubstring("no approvers"))

			By("returning an error if an approver is invalid")
			meta.ApproversReturns([]string{""})
			err = check.ThatThereAreApprovers(meta)
			Expect(err.Error()).To(ContainSubstring("invalid approver"))

			By("returning no error if there are approvers")
			meta.ApproversReturns([]string{"aGoodApprover"})
			err = check.ThatThereAreApprovers(meta)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("Checking that there is an owning SIG, approvers, and reviewers", func() {
		It("checks that the chain of ownership has been established", func() {
			meta := &metadatafakes.FakeKEP{}

			err := check.ThatThereAreOwners(meta)
			merr, ok := err.(*multierror.Error)
			Expect(ok).To(BeTrue())
			Expect(merr.Errors).To(HaveLen(3))

			By("checking the authors, approvers, reviewers, and owning SIG are valid")
			meta.AuthorsReturns([]string{"aGoodCoder"})
			meta.ApproversReturns([]string{"aBusyApprover"})
			meta.ReviewersReturns([]string{"aFairReviewer"})
			meta.OwningSIGReturns("sig-architecture")

			err = check.ThatThereAreOwners(meta)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("Checking that the KEP has a UUID", func() {
		It("checks that UniqueID is set", func() {
			meta := &metadatafakes.FakeKEP{}
			meta.UniqueIDReturns("")
			err := check.ThatKEPHasUUID(meta)
			Expect(err.Error()).To(ContainSubstring("empty string given as UUID"))

			uuidString := uuid.New().String()
			meta.UniqueIDReturns(uuidString)
			err = check.ThatKEPHasUUID(meta)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("Checking that timestamps have been set correctly", func() {
		It("checks that created time has been set", func() {
			meta := &metadatafakes.FakeKEP{}
			meta.CreatedReturns(time.Time{})

			err := check.ThatCreatedTimeExists(meta)
			Expect(err.Error()).To(ContainSubstring("created at time is invalid: not before empty time"))

			meta.CreatedReturns(time.Now())
			err = check.ThatCreatedTimeExists(meta)
			Expect(err).ToNot(HaveOccurred())
		})

		It("checks that last updated is after created time", func() {
			meta := &metadatafakes.FakeKEP{}

			now := time.Now()
			before := now.Add(-time.Minute)
			after := now.Add(time.Minute)

			meta.CreatedReturns(after)
			meta.LastUpdatedReturns(before)

			err := check.ThatLastUpdatedAfterCreated(meta)
			Expect(err.Error()).To(ContainSubstring("Created at time is after last updated"))

			meta.CreatedReturns(before)
			meta.LastUpdatedReturns(after)
			err = check.ThatLastUpdatedAfterCreated(meta)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("Checking that the KEP has development themes", func() {
		It("checks that a development theme has been set", func() {
			meta := &metadatafakes.FakeKEP{}
			meta.DevelopmentThemesReturns([]string{})
			err := check.ThatKEPHasDevelopmentThemes(meta)
			Expect(err.Error()).To(ContainSubstring("no development themes set"))

			meta.DevelopmentThemesReturns([]string{""})
			err = check.ThatKEPHasDevelopmentThemes(meta)
			Expect(err.Error()).To(ContainSubstring("Invalid development theme"))

			meta.DevelopmentThemesReturns([]string{"stability"})
			err = check.ThatKEPHasDevelopmentThemes(meta)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("Checking that the KEP has a stability development theme", func() {
		It("checks that a stability development theme has been set", func() {
			meta := &metadatafakes.FakeKEP{}
			meta.DevelopmentThemesReturns([]string{"stability"})
			err := check.ThatKEPHasDevelopmentThemes(meta)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
