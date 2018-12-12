package hermetic_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	git "gopkg.in/src-d/go-git.v4"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/calebamiles/keps/pkg/hermetic"
)

var _ = Describe("working with a Git repository", func() {
	Describe("#AddFiles()", func() {
		It("adds the files to the repository and commits", func() {
			By("performing a bunch of set up")
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
			withBranchName := "keps-hermetic-fork-test"

			owner := "Charkha"
			repo := "Hello-World"

			forkedRepo, err := hermetic.Fork(githubHandle, tokenProvider, owner, repo, toLocation, withBranchName)
			Expect(err).ToNot(HaveOccurred(), "forking GitHub repository in test")

			defer forkedRepo.DeleteRemote()
			defer forkedRepo.DeleteLocal()

			exampleDir, err := ioutil.TempDir("", "example-add-dir")
			Expect(err).ToNot(HaveOccurred(), "creating example directory")
			defer os.RemoveAll(exampleDir)

			exampleFilename := "example.md"
			exampleLocation := filepath.Join(exampleDir, exampleFilename)

			err = ioutil.WriteFile(exampleLocation, []byte("example content"), os.ModePerm)
			Expect(err).ToNot(HaveOccurred(), "writing a temp file for a test git commit")

			By("adding the files to the Git repository and making a commit")

			err = forkedRepo.Add(exampleLocation, exampleFilename)
			Expect(err).ToNot(HaveOccurred(), "adding a test file to a Git repository")

			gitRepo, err := git.PlainOpen(toLocation)
			Expect(err).ToNot(HaveOccurred(), "expected to find Git repository")

			worktree, err := gitRepo.Worktree()
			Expect(err).ToNot(HaveOccurred(), "opening Git worktree")

			statusOf, err := worktree.Status()
			Expect(err).ToNot(HaveOccurred(), "getting worktree status")

			fileStatus := statusOf.File(exampleFilename)
			Expect(fileStatus.Staging).To(Equal(git.Added))
		})
	})
})
