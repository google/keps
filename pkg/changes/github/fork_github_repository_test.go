package github_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/calebamiles/keps/pkg/changes/auth/authfakes"

	"github.com/calebamiles/keps/pkg/changes/github"
)

var _ = Describe("Forking a GitHub repository", func() {
	Describe("Fork()", func() {
		const (
			exampleRepoOwner = "Planctae"
			exampleRepoName  = "Hello-World"
		)

		It("forks the upstream repository", func() {
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

			_, err := github.Fork(token, exampleRepoOwner, exampleRepoName)
			Expect(err).ToNot(HaveOccurred(), "expected no error when forking an existing GitHub repository")

			expectedRepoApiUrl := github.RepoApiUrl(githubHandle, exampleRepoName)

			resp, err := http.Get(expectedRepoApiUrl)
			Expect(err).ToNot(HaveOccurred(), "expected no error when getting information about forked repository")
			defer resp.Body.Close()

			var getRepoResponse struct {
				Forked bool `json:"fork"`
			}

			Expect(resp.StatusCode).To(Equal(http.StatusOK), "expected `200 OK` when getting forked repo information")

			bodyBytes, err := ioutil.ReadAll(resp.Body)
			Expect(err).ToNot(HaveOccurred(), "expected no error when reading GET repo HTTP response")

			err = json.Unmarshal(bodyBytes, &getRepoResponse)
			Expect(err).ToNot(HaveOccurred(), "expected no error when unmarshalling HTTP response from JSON")

			Expect(getRepoResponse.Forked).To(BeTrue(), "expected forked repository to be listed as a fork by GitHub")

			err = deleteGithubRepo(token, githubHandle, exampleRepoName)
			Expect(err).ToNot(HaveOccurred(), "expected no error after deleting forked GitHub repo during cleanup")
		})

		Context("when the repository has already been forked", func() {
			It("returns no error", func() {
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

				_, err := github.Fork(token, exampleRepoOwner, exampleRepoName)
				Expect(err).ToNot(HaveOccurred(), "expected no error when forking an existing GitHub repository")

				_, err = github.Fork(token, exampleRepoOwner, exampleRepoName)
				Expect(err).ToNot(HaveOccurred(), "expected no error when forking a repository that has already been forked")

				err = deleteGithubRepo(token, githubHandle, exampleRepoName)
				Expect(err).ToNot(HaveOccurred(), "expected no error after deleting forked GitHub repo during cleanup")
			})
		})
	})
})
