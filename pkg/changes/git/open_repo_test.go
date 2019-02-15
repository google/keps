package git_test

import (
	"io/ioutil"
	"os"

	libgit "gopkg.in/src-d/go-git.v4"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/calebamiles/keps/pkg/changes/git"
)

var _ = Describe("opening an existing Git repo on disk", func() {
	It("opens the Git repo", func() {
		By("expecting the repo to already exist on disk")

		tmpDir, err := ioutil.TempDir("", "keps-git-test")
		Expect(err).ToNot(HaveOccurred(), "expected no error when creating a temp directory for test")
		defer os.RemoveAll(tmpDir)

		_, err = libgit.PlainClone(tmpDir, false, &libgit.CloneOptions{
			URL: exampleRepoUrl,
		})

		Expect(err).ToNot(HaveOccurred(), "expected no error when cloning example repository from GitHub")

		By("opening the repository")

		_, err = git.Open(tmpDir)
		Expect(err).ToNot(HaveOccurred(), "expected no error when opening a previously cloned repository")
	})

	It("obtains a file lock on the repository", func() {
		By("expecting the repo to already exist on disk")

		tmpDir, err := ioutil.TempDir("", "keps-git-test")
		Expect(err).ToNot(HaveOccurred(), "expected no error when creating a temp directory for test")
		defer os.RemoveAll(tmpDir)

		_, err = libgit.PlainClone(tmpDir, false, &libgit.CloneOptions{
			URL: exampleRepoUrl,
		})

		Expect(err).ToNot(HaveOccurred(), "expected no error when cloning example repository from GitHub")

		_, err = git.Open(tmpDir)
		Expect(err).ToNot(HaveOccurred(), "expected no error when opening a Git repo for the first time using the KEP tool")

		By("obtaining an exclusive file lock")

		_, err = git.Open(tmpDir)
		Expect(err).To(MatchError("could not obtain exlusive file lock when opening repository"), "expected to receive error when attempting to open Git repo that has already been opened by the KEP tool")
	})
})
