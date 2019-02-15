package github_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/calebamiles/keps/pkg/changes/git"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/calebamiles/keps/pkg/changes/auth/authfakes"
	"github.com/calebamiles/keps/pkg/changes/github/githubfakes"

	"github.com/calebamiles/keps/pkg/changes/github"
)

var _ = Describe("Creating a GitHub Pull Request", func() {
	Describe("CreatePullRequest()", func() {
		It("creates a GitHub pull request from an existing forked repository", func() {
			By("doing a lot of tedious setup")

			githubToken := os.Getenv("KEP_TEST_GITHUB_TOKEN")
			if githubToken == "" {
				Skip("KEP_TEST_GITHUB_TOKEN unset and required for test")
			}

			githubHandle := os.Getenv("KEP_TEST_GITHUB_HANDLE")
			if githubToken == "" {
				Skip("KEP_TEST_GITHUB_HANDLE unset and required for test")
			}

			committerName := "OSS KEP Tool Test"
			committerEmail := "kubernetes-sig-architecture@googlegroups.com"
			commitMessage := "adds an important change"

			pullRequestTitle := github.PullRequestTitle("A Really Great Change")
			pullRequestDescription := github.PullRequestDescription("Please Pull This In!")

			sourceBranch := "master"
			sourceRepository := "Hello-World"
			sourceRepositoryOwner := githubHandle
			targetBranch := "master"
			targetRepository := "Hello-World"
			targetRepositoryOwner := "Planctae"

			tmpDir, err := ioutil.TempDir("", "github-create-pr-test")
			Expect(err).ToNot(HaveOccurred(), "expected no error when creating a temporary directory")

			localRepoLocation := filepath.Join(tmpDir, "example_repo")

			token := &authfakes.FakeTokenProvider{}
			token.ValueReturns(githubToken, nil)

			routingInfo := &githubfakes.FakePullRequestRoutingInfo{}
			routingInfo.SourceBranchReturns(sourceBranch)
			routingInfo.SourceRepositoryReturns(sourceRepository)
			routingInfo.SourceRepositoryOwnerReturns(sourceRepositoryOwner)
			routingInfo.TargetBranchReturns(targetBranch)
			routingInfo.TargetRepositoryReturns(targetRepository)
			routingInfo.TargetRepositoryOwnerReturns(targetRepositoryOwner)
			routingInfo.TokenReturns(token)

			repoUrl, err := github.Fork(token, targetRepositoryOwner, targetRepository)
			Expect(err).ToNot(HaveOccurred(), "expected no error when forking what should be an existing repository to the user account")
			defer deleteGithubRepo(token, sourceRepositoryOwner, sourceRepository)

			repo, err := git.Clone(token, repoUrl, localRepoLocation)
			Expect(err).ToNot(HaveOccurred(), "expected no error when cloning a previously forked repository")

			testFileContent := "some example content"
			testFilename := "test_content.md"

			err = ioutil.WriteFile(filepath.Join(localRepoLocation, testFilename), []byte(testFileContent), os.ModePerm)
			Expect(err).ToNot(HaveOccurred(), "expected no error when writing a test file to include in pr")

			err = repo.Add(testFilename)
			Expect(err).ToNot(HaveOccurred(), "expected no error when adding a newly created file to the repository")

			err = repo.Commit(committerName, committerEmail, commitMessage)
			Expect(err).ToNot(HaveOccurred(), "expected no error when committing a file that was newly added to the repository")

			err = repo.PushOrigin(token, sourceBranch, targetBranch)
			Expect(err).ToNot(HaveOccurred(), "expected no error when pushing newly committed changes to origin")

			By("creating a Pull Request using the GitHub API")

			prUrl, err := github.CreatePullRequest(routingInfo, pullRequestTitle, pullRequestDescription)
			Expect(err).ToNot(HaveOccurred(), "expected no error when creating a pull request from recently pushed changes")
			defer closeGithubPullRequest(token, prUrl)

			apiUrl := strings.Replace(prUrl, "github.com", "api.github.com/repos", 1)
			apiUrl = strings.Replace(apiUrl, "pull", "pulls", 1)

			resp, err := http.Get(apiUrl)
			Expect(err).ToNot(HaveOccurred(), "expected no error when getting PR details from the GitHub API")
			defer resp.Body.Close()

			respBytes, err := ioutil.ReadAll(resp.Body)
			Expect(err).ToNot(HaveOccurred(), "expected no error when reading http response from get PR details request")

			response := &getPrResponse{}
			err = json.Unmarshal(respBytes, response)
			Expect(err).ToNot(HaveOccurred(), "expected no error when unmarshalling JSON payload from git PR details response into struct")

			Expect(string(response.Title)).To(Equal(string(pullRequestTitle)), "expected upstream pull request title to match given title")
			Expect(string(response.Description)).To(Equal(string(pullRequestDescription)), "expected upstream pull request description to match given description")
			Expect(response.Head.Label).To(Equal(fmt.Sprintf("%s:%s", sourceRepositoryOwner, sourceBranch)), "expected pull request head to be of the form <githubOwner:sourceBranch>")
			Expect(response.Base.Label).To(Equal(fmt.Sprintf("%s:%s", targetRepositoryOwner, targetBranch)), "expected pull request base to be of the form <upstreamOwner:targetBranch>")
		})
	})
})

type getPrResponse struct {
	Title       string        `json:"title"`
	Description string        `json:"body"`
	Base        baseReference `json:"base"`
	Head        headReference `json:"head"`
}

type baseReference struct {
	Ref   string `json:"ref"`
	Label string `json:"label"`
}

type headReference struct {
	Ref   string `json:"ref"`
	Label string `json:"label"`
}
