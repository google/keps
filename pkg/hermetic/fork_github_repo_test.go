package porcelain_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	git "gopkg.in/src-d/go-git.v4"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/calebamiles/keps/pkg/porcelain"
)

var _ = Describe("working with a Git repository hosted on GitHub", func() {
	Describe("Fork()", func() {
		It("creates a copy of an upstream GitHub repository in a user account", func() {
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

			tmpDir, err := ioutil.TempDir("", "keps-fork-test")
			Expect(err).ToNot(HaveOccurred())
			defer os.RemoveAll(tmpDir)

			toLocation := filepath.Join(tmpDir, "forked-repo")
			withBranchName := "keps-porcelain-fork-test"

			By("forking a remote Git repository locally")

			owner := "Charkha"
			repo := "Hello-World"

			forkedRepo, err := porcelain.Fork(githubHandle, tokenProvider, owner, repo, toLocation, withBranchName)
			Expect(err).ToNot(HaveOccurred(), "forking GitHub repository in test")

			defer forkedRepo.DeleteRemote()
			defer forkedRepo.DeleteLocal()

			expectedGitDir := filepath.Join(toLocation, ".git")
			Expect(expectedGitDir).To(BeADirectory(), "expected to find .git directory after fork")

			gitRepo, err := git.PlainOpen(toLocation)
			Expect(err).ToNot(HaveOccurred(), "expected to open a Git repository")

			By("setting the repoURL as the `upstream` remote")

			remotes, err := gitRepo.Remotes()
			Expect(err).ToNot(HaveOccurred(), "listing git remotes")
			Expect(remotes).To(HaveLen(2), "expected git repo to have two configured remotes")

			remoteNames := []string{remotes[0].Config().Name, remotes[1].Config().Name}
			Expect(remoteNames).To(ContainElement(porcelain.UpstreamRemoteName), "expected configured remotes to contain name `upstream`")
			Expect(remoteNames).To(ContainElement(porcelain.OriginRemoteName), "expected second configured remotes to contain name `origin`")

			By("checking out a new branch that tracks upstream")

			head, err := gitRepo.Head()
			Expect(err).ToNot(HaveOccurred(), "reading HEAD of Git repository")

			Expect(head.String()).ToNot(Equal(fmt.Sprintf("refs/heads/%s", withBranchName)))
		})

		Context("when the location to clone the repo exists", func() {
			It("returns an error", func() {
				tmpDir, err := ioutil.TempDir("", "keps-fork-test")
				Expect(err).ToNot(HaveOccurred())
				defer os.RemoveAll(tmpDir)

				tokenProvider := newMockTokenProvider()
				githubHandle := "doesnt-matter"
				withBranchName := "keps-porcelain-fork-test"

				owner := "Charkha"
				repo := "Hello-World"

				_, err = porcelain.Fork(githubHandle, tokenProvider, owner, repo, tmpDir, withBranchName)
				Expect(err.Error()).To(ContainSubstring("may exist already, refusing to overwrite"), "expected error to contain refusal to clone over existing location")

			})
		})
	})
})
