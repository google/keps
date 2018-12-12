package hermetic_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/calebamiles/keps/pkg/hermetic"
)

var _ = Describe("submitting changes to an upstream GitHub repository", func() {
	Describe("#CreatePR()", func() {
		It("pushes any local changes to the origin and creates a PR", func() {
			By("performing a lot of setup")

			githubToken := os.Getenv("KEP_TEST_GITHUB_TOKEN")
			Expect(githubToken).ToNot(BeEmpty(), "KEP_TEST_GITHUB_TOKEN unset and required for test")

			githubHandle := os.Getenv("KEP_TEST_GITHUB_HANDLE")
			Expect(githubHandle).ToNot(BeEmpty(), "KEP_TEST_GITHUB_HANDLE unset and required for test")

			tokenProvider := newMockTokenProvider()

			// call #1: repo fork
			tokenProvider.ValueOutput.Ret0 <- githubToken
			tokenProvider.ValueOutput.Ret1 <- nil

			// call #2: repo clone
			tokenProvider.ValueOutput.Ret0 <- githubToken
			tokenProvider.ValueOutput.Ret1 <- nil

			// call #3: delete repo callback
			tokenProvider.ValueOutput.Ret0 <- githubToken
			tokenProvider.ValueOutput.Ret1 <- nil

			// call #4: push repo callback
			tokenProvider.ValueOutput.Ret0 <- githubToken
			tokenProvider.ValueOutput.Ret1 <- nil

			// call #5: create pull request callback
			tokenProvider.ValueOutput.Ret0 <- githubToken
			tokenProvider.ValueOutput.Ret1 <- nil

			tmpDir, err := ioutil.TempDir("", "keps-fork-test")
			Expect(err).ToNot(HaveOccurred())
			//defer os.RemoveAll(tmpDir)

			toLocation := filepath.Join(tmpDir, "forked-repo")
			withBranchName := "keps-hermetic-fork-test"

			owner := "Charkha"
			repo := "Hello-World"

			forkedRepo, err := hermetic.Fork(githubHandle, tokenProvider, owner, repo, toLocation, withBranchName)
			Expect(err).ToNot(HaveOccurred(), "forking GitHub repository in test")

			defer forkedRepo.DeleteRemote()
			defer forkedRepo.DeleteLocal()

			exampleDir, err := ioutil.TempDir("", "example-add-dir")
			Expect(err).ToNot(HaveOccurred(), "creating directory for example pull request file")

			exampleFilename := "example.md"
			exampleLocation := filepath.Join(exampleDir, exampleFilename)

			err = ioutil.WriteFile(exampleLocation, []byte("example content"), os.ModePerm)
			Expect(err).ToNot(HaveOccurred(), "writing a temp file for a test git commit")
			defer os.RemoveAll(exampleDir)

			By("adding a file to the repo")

			err = forkedRepo.Add(exampleLocation, exampleFilename) // add to the root of the repo
			Expect(err).ToNot(HaveOccurred(), "adding a test file to a Git repository")

			By("creating a GitHub Pull Request")

			prTitle := "example pull request"
			prDescription := "these are important changes"

			prUrl, err := forkedRepo.CreatePR(prTitle, prDescription)
			Expect(err).ToNot(HaveOccurred(), "creating pull request for test")

			resp, err := http.Get(prUrl)
			Expect(err).ToNot(HaveOccurred(), "GET-ting pull request URL")
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusOK), "expected status code when GET-ting pull request URL to be 200 OK")

			prApiUrl := strings.Replace(prUrl, "https://github", "https://api.github", 1)
			prApiUrl = strings.Replace(prUrl, "github.com", "github.com/repos", 1)

			var payloadBody struct {
				State string `json:"state"`
			}

			payloadBody.State = "closed"
			payloadBytes, err := json.Marshal(payloadBody)
			Expect(err).ToNot(HaveOccurred(), "marshalling close PR request payload")

			payload := ioutil.NopCloser(bytes.NewBuffer(payloadBytes))

			// create HTTP request
			req, err := http.NewRequest(http.MethodPatch, prApiUrl, payload)
			Expect(err).ToNot(HaveOccurred(), "creating close PR HTTP request")

			// add auth header
			req.Header.Add("Authorization", fmt.Sprintf("token %s", githubToken))

			// set context
			ctx, cancel := context.WithTimeout(context.Background(), 7000*time.Millisecond)
			createPullRequest := req.WithContext(ctx)
			defer cancel()

			// Do request
			c := &http.Client{}

			resp, err = c.Do(createPullRequest)
			Expect(err).ToNot(HaveOccurred(), "closing pull request")

			defer resp.Body.Close()
		})
	})

})
