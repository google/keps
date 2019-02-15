package git_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	libgit "gopkg.in/src-d/go-git.v4"

	"github.com/calebamiles/keps/pkg/changes/auth/authfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/calebamiles/keps/pkg/changes/git"
)

var _ = Describe("cloning a Git repo to disk", func() {
	Describe("Clone()", func() {
		It("clones a Git repository to disk", func() {
			githubToken := os.Getenv("KEP_TEST_GITHUB_TOKEN")
			if githubToken == "" {
				Skip("KEP_TEST_GITHUB_TOKEN unset and required for test")
			}

			tmpDir, err := ioutil.TempDir("", "keps-clone-test")
			Expect(err).ToNot(HaveOccurred(), "expected no error to occur creating a temp directory")
			defer os.RemoveAll(tmpDir)

			exampleRepoLocation := filepath.Join(tmpDir, exampleRepoName)

			token := &authfakes.FakeTokenProvider{}
			token.ValueReturns(githubToken, nil)

			_, err = git.Clone(token, exampleRepoUrl, exampleRepoLocation)
			Expect(err).ToNot(HaveOccurred(), "expected no error when cloning an example Git repo using a valid GitHub token")

			Expect(filepath.Join(exampleRepoLocation, gitRepoConfigDir)).To(BeADirectory(), "expected to find `.git` directory inside of a cloned Git directory")

			underlyingRepo, err := libgit.PlainOpen(exampleRepoLocation)
			Expect(err).ToNot(HaveOccurred(), "expected no error when opening a previously cloned repository")

			_, err = underlyingRepo.Head()
			Expect(err).ToNot(HaveOccurred(), "expected no error when retrieving `HEAD` of a previously cloned repository")
		})

		It("obtains an exclusive file lock for the repository", func() {
			githubToken := os.Getenv("KEP_TEST_GITHUB_TOKEN")
			if githubToken == "" {
				Skip("KEP_TEST_GITHUB_TOKEN unset and required for test")
			}

			tmpDir, err := ioutil.TempDir("", "keps-clone-test")
			Expect(err).ToNot(HaveOccurred(), "expected no error to occur creating a temp directory")
			defer os.RemoveAll(tmpDir)

			exampleRepoLocation := filepath.Join(tmpDir, exampleRepoName)

			token := &authfakes.FakeTokenProvider{}
			token.ValueReturns(githubToken, nil)

			_, err = git.Clone(token, exampleRepoUrl, exampleRepoLocation)
			Expect(err).ToNot(HaveOccurred(), "expected no error when cloning an example Git repo using a valid GitHub token")

			_, err = git.Open(exampleRepoLocation)
			Expect(err).To(MatchError("could not obtain exlusive file lock when opening repository"), "expected to receieve file lock error when attempting to open an already open repository")
		})

		Context("when the provided location already exists", func() {
			It("returns an error", func() {
				githubToken := os.Getenv("KEP_TEST_GITHUB_TOKEN")
				if githubToken == "" {
					Skip("KEP_TEST_GITHUB_TOKEN unset and required for test")
				}

				tmpDir, err := ioutil.TempDir("", "keps-clone-test")
				Expect(err).ToNot(HaveOccurred(), "expected no error to occur creating a temp directory")
				defer os.RemoveAll(tmpDir)

				token := &authfakes.FakeTokenProvider{}
				token.ValueReturns(githubToken, nil)

				_, err = git.Clone(token, exampleRepoUrl, tmpDir)
				Expect(err.Error()).To(ContainSubstring("may exist already, refusing to overwrite"), "expected error when attempting to clone a Git repository to an existing directory")
			})
		})
	})
})
