package hermetic_test

import (
	"github.com/calebamiles/keps/pkg/changes/github"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/calebamiles/keps/pkg/changes/auth/authfakes"
	"github.com/calebamiles/keps/pkg/changes/git/gitfakes"
	"github.com/calebamiles/keps/pkg/changes/routing/routingfakes"

	"github.com/calebamiles/keps/pkg/changes/hermetic"
)

var _ = Describe("Subitting changes hermetically", func() {
	Describe("staging, comitting, and creating a PR", func() {
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

			repo, err := hermetic.NewRepo(fakeRoutingInfo, fakeRepo, fakePrCreator)
			Expect(err).ToNot(HaveOccurred(), "expected no error when creating a new hermetic.Repo with fake but valid routing info, git.Repo, and PR creator")

			prUrl, err := repo.SubmitChanges()
			Expect(err).ToNot(HaveOccurred(), "expected no error when submitting changes hermetically using fake routing.Info, git.Repo, and PR creator")

			Expect(prUrl).To(Equal(fakeUrl), "expected returned PR URL to match stub from fake PR creator")
			Expect(calledWith.Title).To(Equal(prTitle), "expected github.PullRequestCreator to be called with argument to hermetic.Repo")
			Expect(calledWith.Description).To(Equal(prDescription), "expected github.PullRequestCreator to be called with argument to hermetic.Repo")

			Expect(fakeRepo.AddCallCount()).To(Equal(1), "expected git.Repo#Add() to be called only once when copying files to the hermetic.Repo")
			Expect(fakeRepo.AddArgsForCall(0)).To(Equal(pathToChanges), "expected git.Repo#Add() to be called with the path given to hermetic.Repo#Add()")

			Expect(fakeRepo.CommitCallCount()).To(Equal(1), "expected git.Repo#Commit() to be called only once when submitting changes hermetically")

			givenName, givenEmail, givenMessage := fakeRepo.CommitArgsForCall(0)
			Expect(givenEmail).To(Equal(principalEmail), "expected git.Repo#Commit() to be called with the committer email from routing.Info")
			Expect(givenName).To(Equal(principalName), "expected git.Repo#Commit() to be called with the committer name from routing.Info")
			Expect(givenMessage).To(Equal(commitMsg), "expected git.Repo#Commit() to be called with the commit message from routing.Info")

			Expect(fakeRepo.PushOriginCallCount()).To(Equal(1), "expected git.Repo#PushOrigin() to be called only once when submitting changes hermetically")

			Expect(calledWith.Title).To(Equal(prTitle), "expected github.PullRequestCreator to be called with the title given to hermetic.Repo#CreatePR()")
			Expect(calledWith.Description).To(Equal(prDescription), "expected github.PullRequestCreator to be called with the description given to hermetic.Repo#CreatePR()")

		})
	})

	Describe("setting up a prestine, or hermetic, git repository for submitting changes", func() {
		Describe("Fork()", func() {
			XIt("prepares a prestine copy of an existing Git repository hosted on GitHub", func() {
				Fail("test not written")
			})

		})
	})

	/*
		Describe("Forking, copying content, and creating a GitHub Pull Request", func() {
			It("allows changes to be made to a prestine copy of an upstream repository", func() {

				Skip("this test should be rewritten to target the `sandbox` (user account)")

				githubToken := os.Getenv("KEP_TEST_GITHUB_TOKEN")
				if githubToken == "" {
					Skip("KEP_TEST_GITHUB_TOKEN unset and required for test")
				}

				githubHandle := os.Getenv("KEP_TEST_GITHUB_HANDLE")
				if githubToken == "" {
					Skip("KEP_TEST_GITHUB_HANDLE unset and required for test")
				}

				token := &authfakes.FakeTokenProvider{}
				token.ValueReturns(githubToken, nil)

				tmpDir, err := ioutil.TempDir("", "hermetic-fork-repo")
				Expect(err).ToNot(HaveOccurred(), "expected no error when creating a temporary directory for test")

				repoOwnerName := "Planctae"
				repoName := "Hello-World"
				withBranchName := "a-great-idea"

				toLocation := filepath.Join(tmpDir, "example_repo")

				// fork repo to location
				repo, err := hermetic.Fork(
					token,
					repoOwnerName,
					repoName,
					githubHandle,
					toLocation,
					withBranchName,
				)

				Expect(err).ToNot(HaveOccurred(), "expected no error when forking an existing GitHub repo to a user account")

				tmpFile, err := ioutil.TempFile(tmpDir, "submit-hermetic-changes-test")
				Expect(err).ToNot(HaveOccurred(), "expected no error when creating a temporary file to commit")

				exampleContent := "some example content"
				_, err = tmpFile.WriteString(exampleContent)
				Expect(err).ToNot(HaveOccurred(), "expected no error writing test content to a temporary file to commit")

				err = tmpFile.Close()
				Expect(err).ToNot(HaveOccurred(), "expected no error closing temp file with test content")

				testFileLocation := tmpFile.Name()

				// copy example file to repo
				targetLocation := "release-x.y/"
				err = repo.Copy(testFileLocation, targetLocation)
				Expect(err).ToNot(HaveOccurred(), "expected no error when copying an existing file to a hermetically forked repository")

				// commit
				committerEmail := "kubernetes-sig-architecture@googlegroups.com"
				committerName := "OSS KEP Tool Tester"
				commitMsg := "propose Test File to release-x.y"
				err = repo.Commit(committerName, committerEmail, commitMsg)
				Expect(err).ToNot(HaveOccurred(), "expected no error when committing a newly copied file to a hermetically forked repository")

				// create pr
				testPRBody := "this is a really important change"
				testPRTitle := "please pull this in"
				prUrl, err := repo.CreatePR(testPRTitle, testPRBody)
				Expect(err).ToNot(HaveOccurred(), "expected no error when creating a pull request against the upstream repository from a hermetically forked repository")

				prFilesUrl := fmt.Sprintf("%s/files", prUrl)

				resp, err := http.Get(prFilesUrl)
				Expect(err).ToNot(HaveOccurred(), "expected no error when getting list of pull request files")
				defer resp.Body.Close()

				bodyBytes, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred(), "expected no error reading the HTTP response from GitHub")

				listFilesResponse := &listPRFilesResponse{}
				err = json.Unmarshal(bodyBytes, &listFilesResponse)
				Expect(err).ToNot(HaveOccurred(), "expected no error unmarshalling the HTTP response from JSON")

				Expect(listFilesResponse.Files).To(HaveLen(1), "expected pull request to contain a single file")

				fileEntry := listFilesResponse.Files[0]
				Expect(fileEntry.Filename).To(Equal(tmpFile.Name()), "expected pull request to contain file with name of test file")
				Expect(fileEntry.Sha).To(BeEmpty(), "expected pull request to contain file with SHA of test file")

				// assert file content matches
				// close PR
				// delete forked repo
			})
		})
	*/

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

/*
type listPRFilesResponse struct {
	Files []*listPRFilesEntry
}

type listPRFilesEntry struct {
	Sha      string `json:"sha"`
	Filename string `json:"filename"`
}
*/
