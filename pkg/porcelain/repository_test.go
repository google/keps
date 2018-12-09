package porcelain_test

import (
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

			// call #1: repo clone
			tokenProvider.ValueOutput.Ret0 <- githubToken
			tokenProvider.ValueOutput.Ret1 <- nil

			// call #1: delete repo callback
			tokenProvider.ValueOutput.Ret0 <- githubToken
			tokenProvider.ValueOutput.Ret1 <- nil

			tmpDir, err := ioutil.TempDir("", "keps-fork-test")
			Expect(err).ToNot(HaveOccurred())
			defer os.RemoveAll(tmpDir)

			toLocation := filepath.Join(tmpDir, "forked-repo")
			withBranchName := "keps-porcelain-fork-test"

			By("forking a remote Git repository locally")

			repoGitUrl := "https://github.com/octocat/Hello-World"
			repoApiUrl := "https://api.github.com/repos/octocat/Hello-World/forks"

			forkedRepo, err := porcelain.Fork(githubHandle, tokenProvider, repoApiUrl, repoGitUrl, toLocation, withBranchName)
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

		})

		Context("when the location to clone the repo exists", func() {
			XIt("returns an error", func() {
				//toLocation, err := ioutil.TempDir("", "keps-clone-test")
				//toLocatExpect(err).ToNot(HaveOccurred(), "creating temp dir before clone test")
				//toLocatdefer os.RemoveAll(toLocation)

				//toLocat_, err = porcelain.Clone(tokenProvider, repoUrl, porcelain.DefaultBranch, toLocation)
				//toLocatExpect(err).ToNot(HaveOccurred(), "cloning GitHub repository for test")
			})
		})
	})

	Describe("operations on Git repositories", func() {
		Describe("#AddFiles()", func() {
			XIt("adds the files to the repository and commits", func() {
				/*
					githubToken := os.Getenv("KEP_TEST_GITHUB_TOKEN")
					Expect(githubToken).ToNot(BeEmpty(), "KEP_TEST_GITHUB_TOKEN unset and required for test")

					tokenProvider := newMockTokenProvider()

					// call #1: repo clone
					tokenProvider.ValueOutput.Ret0 <- githubToken
					tokenProvider.ValueOutput.Ret1 <- nil

					toLocation := filepath.Join(os.TempDir(), "keps-add-files-test")
					defer os.RemoveAll(toLocation)

					repoUrl := "https://github.com/octocat/Hello-World"
					repo, err := porcelain.Clone(tokenProvider, repoUrl, porcelain.DefaultBranchName, toLocation)
					Expect(err).ToNot(HaveOccurred(), "cloning GitHub repository for test")

					tmpDir, err := ioutil.TempDir(toLocation, "example-add-dir")
					Expect(err).ToNot(HaveOccurred(), "creating sub directory in Git repo")

					exampleFilename := "example.md"
					err = ioutil.WriteFile(filepath.Join(tmpDir, exampleFilename), []byte("example content"), os.ModePerm)
					Expect(err).ToNot(HaveOccurred(), "writing a temp file for a test git commit")

					pathToDir, err := filepath.Rel(toLocation, tmpDir)
					Expect(err).ToNot(HaveOccurred(), "determining relative path to created subdirectory from test repository root")

					commitMsg := "example commit message"
					err = repo.AddPaths(commitMsg, []string{pathToDir})
					Expect(err).ToNot(HaveOccurred(), "adding a test file to a Git repository")

					gitRepo, err := git.PlainOpen(toLocation)
					Expect(err).ToNot(HaveOccurred(), "expected to find Git repository")

					By("adding the files to the Git repository and making a commit")
					head, err := gitRepo.Head()
					Expect(err).ToNot(HaveOccurred(), "fetching HEAD of git repo")

					lastCommit, err := gitRepo.CommitObject(head.Hash())
					Expect(err).ToNot(HaveOccurred(), "reading commit at HEAD")
					println(lastCommit)

					_, err = lastCommit.File(filepath.Join(pathToDir, exampleFilename))
					Expect(err).ToNot(HaveOccurred(), "reading expected committed file info from test commit")
				*/
			})
		})

		Describe("#CreatePR()", func() {
			XIt("pushes any local changes to the origin and creates a PR", func() {
				/*
					githubToken := os.Getenv("KEP_TEST_GITHUB_TOKEN")
					Expect(githubToken).ToNot(BeEmpty(), "KEP_TEST_GITHUB_TOKEN unset and required for test")

					tokenProvider := newMockTokenProvider()

					// call #1: repo clone
					tokenProvider.ValueOutput.Ret0 <- githubToken
					tokenProvider.ValueOutput.Ret1 <- nil

					toLocation := filepath.Join(os.TempDir(), "keps-add-files-test")
					defer os.RemoveAll(toLocation)

					repoUrl := "https://github.com/octocat/Hello-World"
					repo, err := porcelain.Clone(tokenProvider, repoUrl, porcelain.DefaultBranchName, toLocation)
					Expect(err).ToNot(HaveOccurred(), "cloning GitHub repository for test")

					tmpDir, err := ioutil.TempDir(toLocation, "example-add-dir")
					Expect(err).ToNot(HaveOccurred(), "creating sub directory in Git repo")

					exampleFilename := "example.md"
					err = ioutil.WriteFile(filepath.Join(tmpDir, exampleFilename), []byte("example content"), os.ModePerm)
					Expect(err).ToNot(HaveOccurred(), "writing a temp file for a test git commit")

					pathToDir, err := filepath.Rel(toLocation, tmpDir)
					Expect(err).ToNot(HaveOccurred(), "determining relative path to created subdirectory from test repository root")

					commitMsg := "example commit message"
					err = repo.AddPaths(commitMsg, []string{pathToDir})
					Expect(err).ToNot(HaveOccurred(), "adding a test file to a Git repository")

					By("pushing local changes and creating a new PR")

					// creating a PR should return the PR URL
					// curling the URL should return some kind of OK
				*/
			})
		})

	})
})
