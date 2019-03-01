package inplace_test

import (
	"github.com/calebamiles/keps/pkg/changes/github"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/calebamiles/keps/pkg/changes/auth/authfakes"
	"github.com/calebamiles/keps/pkg/changes/git/gitfakes"
	"github.com/calebamiles/keps/pkg/changes/routing/routingfakes"

	"github.com/calebamiles/keps/pkg/changes/inplace"
)

var _ = Describe("Submitting changes in place", func() {
	Describe("staging, committing, and creating a PR", func() {
		It("submits changes through a pull request", func() {
			calledWith := &fakePullRequestCreatorCallArgs{}
			retValues := &fakePullRequestCreatorRetVals{}

			fakeUrl := "this isn't a URL!"
			retValues.PullRequestURL = fakeUrl

			sourceOwner := "kep-tester"
			sourceRepository := "enhancements"
			sourceBranch := "master"

			targetOwner := "Planctae"
			targetRepository := "enhancements"
			targetBranch := "master"

			pathToChanges := "an existing path inside the repository"
			pathToRepo := "a path to an existing Git repo"

			commitMsg := "imporant changes to pull in"

			prTitle := "an unusually easy to find name"
			prDescription := "these are super important I swear"

			owningSIG := "sig-architecture"

			principalName := "OSS KEP Tester"
			principalEmail := "kubernetes-sig-architecture+keps@googlegroups.com"

			githubToken := "this isn't a valid GitHub auth token"

			fakeToken := &authfakes.FakeTokenProvider{}
			fakeToken.ValueReturns(githubToken, nil)

			fakeRoutingInfo := &routingfakes.FakeInfo{}
			fakeRoutingInfo.ChangesUnderPathReturns(pathToChanges)
			fakeRoutingInfo.FullDescriptionReturns(prDescription)
			fakeRoutingInfo.ShortSummaryReturns(commitMsg)
			fakeRoutingInfo.OwningSIGReturns(owningSIG)
			fakeRoutingInfo.PathToRepoReturns(pathToRepo)
			fakeRoutingInfo.PrincipalEmailReturns(principalEmail)
			fakeRoutingInfo.PrincipalNameReturns(principalName)
			fakeRoutingInfo.SourceBranchReturns(sourceBranch)
			fakeRoutingInfo.SourceRepositoryReturns(sourceRepository)
			fakeRoutingInfo.SourceRepositoryOwnerReturns(sourceOwner)
			fakeRoutingInfo.TargetBranchReturns(targetBranch)
			fakeRoutingInfo.TargetRepositoryReturns(targetRepository)
			fakeRoutingInfo.TargetRepositoryOwnerReturns(targetOwner)
			fakeRoutingInfo.TitleReturns(prTitle)
			fakeRoutingInfo.TokenReturns(fakeToken)

			fakeRepo := &gitfakes.FakeRepo{}
			fakeRepo.AddReturns(nil)        // no error
			fakeRepo.CommitReturns(nil)     // no error
			fakeRepo.PushOriginReturns(nil) // no error

			fakePrCreator := newFakePullRequestCreator(retValues, calledWith)

			repo, err := inplace.NewRepo(fakeRoutingInfo, fakeRepo, fakePrCreator)
			Expect(err).ToNot(HaveOccurred(), "expected no error when creating a new inplace.Repo with fake but valid routing info, git.Repo, and PR creator")

			prUrl, err := repo.SubmitChanges()
			Expect(err).ToNot(HaveOccurred(), "expected no error when submitting changes in place using fake routing.Info, git.Repo, and PR creator")

			Expect(prUrl).To(Equal(fakeUrl), "expected returned PR URL to match stub from fake PR creator")
			Expect(calledWith.Title).To(Equal(prTitle), "expected github.PullRequestCreator to be called with argument to inplace.Repo")
			Expect(calledWith.Description).To(Equal(prDescription), "expected github.PullRequestCreator to be called with argument to inplace.Repo")

			Expect(fakeRepo.AddCallCount()).To(Equal(1), "expected git.Repo#Add() to be called only once when submitting changes inplace")
			Expect(fakeRepo.AddArgsForCall(0)).To(Equal(pathToChanges), "expected git.Repo#Add() to be called with the path given to inplace.Repo#Add()")

			Expect(fakeRepo.CommitCallCount()).To(Equal(1), "expected git.Repo#Commit() to be called only once when submitting changes inplace")

			givenName, givenEmail, givenMessage := fakeRepo.CommitArgsForCall(0)
			Expect(givenEmail).To(Equal(principalEmail), "expected git.Repo#Commit() to be called with the committer email from routing.Info")
			Expect(givenName).To(Equal(principalName), "expected git.Repo#Commit() to be called with the committer name from routing.Info")
			Expect(givenMessage).To(Equal(commitMsg), "expected git.Repo#Commit() to be called with the commit message from routing.Info")

			Expect(fakeRepo.PushOriginCallCount()).To(Equal(1), "expected git.Repo#PushOrigin() to be called only once when submitting changes inplace")
			Expect(calledWith.Title).To(Equal(prTitle), "expected github.PullRequestCreator to be called with the title given to inplace.Repo#CreatePR()")
			Expect(calledWith.Description).To(Equal(prDescription), "expected github.PullRequestCreator to be called with the description given to inplace.Repo#CreatePR()")
		})

		Context("when an existing change receipt (Pull Request) is returned by the routing.Info", func() {
			It("returns the change receipt", func() {

			})
		})
	})

	Describe("opening an existing Git repository", func() {
		XIt("prepares the repository to submit changes to GitHub via Pull Request", func() {
			Fail("test not written")
		})
	})
})

type fakePullRequestCreatorCallArgs struct {
	Title       string
	Description string
}

type fakePullRequestCreatorRetVals struct {
	PullRequestURL string
	Error          error
}

func newFakePullRequestCreator(
	returnVals *fakePullRequestCreatorRetVals,
	calledWith *fakePullRequestCreatorCallArgs,
) github.PullRequestCreator {
	return func(_ github.PullRequestRoutingInfo, prTitle github.PullRequestTitle, prDescription github.PullRequestDescription) (string, error) {
		calledWith.Title = string(prTitle)
		calledWith.Description = string(prDescription)

		return returnVals.PullRequestURL, returnVals.Error
	}
}
