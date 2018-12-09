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

var _ = Describe("working with a Git repository", func() {
	Describe("#AddFiles()", func() {
		FIt("adds the files to the repository and commits", func() {
			githubToken := os.Getenv("KEP_TEST_GITHUB_TOKEN")
			Expect(githubToken).ToNot(BeEmpty(), "KEP_TEST_GITHUB_TOKEN unset and required for test")

			githubHandle := os.Getenv("KEP_TEST_GITHUB_HANDLE")
			Expect(githubHandle).ToNot(BeEmpty(), "KEP_TEST_GITHUB_HANDLE unset and required for test")

			tokenProvider := newMockTokenProvider()

			// call #1: repo clone
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

			repoGitUrl := "https://github.com/Charkha/Hello-World"
			repoApiUrl := "https://api.github.com/repos/Charkha/Hello-World/forks"

			forkedRepo, err := porcelain.Fork(githubHandle, tokenProvider, repoApiUrl, repoGitUrl, toLocation, withBranchName)
			Expect(err).ToNot(HaveOccurred(), "forking GitHub repository in test")

			defer forkedRepo.DeleteRemote()
			defer forkedRepo.DeleteLocal()

			exampleDir, err := ioutil.TempDir(toLocation, "example-add-dir")
			Expect(err).ToNot(HaveOccurred(), "creating sub directory in Git repo")

			exampleFilename := "example.md"
			err = ioutil.WriteFile(filepath.Join(exampleDir, exampleFilename), []byte("example content"), os.ModePerm)
			Expect(err).ToNot(HaveOccurred(), "writing a temp file for a test git commit")

			pathToDir, err := filepath.Rel(toLocation, exampleDir)
			Expect(err).ToNot(HaveOccurred(), "determining relative path to created subdirectory from test repository root")

			commitMsg := "example commit message"
			err = forkedRepo.AddPaths(commitMsg, []string{pathToDir})
			Expect(err).ToNot(HaveOccurred(), "adding a test file to a Git repository")

			gitRepo, err := git.PlainOpen(toLocation)
			Expect(err).ToNot(HaveOccurred(), "expected to find Git repository")

			By("adding the files to the Git repository and making a commit")
			head, err := gitRepo.Head()
			Expect(err).ToNot(HaveOccurred(), "fetching HEAD of git repo")

			lastCommit, err := gitRepo.CommitObject(head.Hash())
			Expect(err).ToNot(HaveOccurred(), "reading commit at HEAD")

			_, err = lastCommit.File(filepath.Join(pathToDir, exampleFilename))
			Expect(err).ToNot(HaveOccurred(), "reading expected committed file info from test commit")
		})
	})
})
