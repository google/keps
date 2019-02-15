package git_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	libgit "gopkg.in/src-d/go-git.v4"
	libgitconfig "gopkg.in/src-d/go-git.v4/config"
	libgithttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"

	"github.com/calebamiles/keps/pkg/changes/auth"
	"github.com/calebamiles/keps/pkg/changes/auth/authfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/calebamiles/keps/pkg/changes/git"
)

var _ = Describe("working with a Git repository", func() {
	Describe("#SetOrigin()", func() {
		It("updates the remote `origin`", func() {
			githubToken := os.Getenv("KEP_TEST_GITHUB_TOKEN")
			if githubToken == "" {
				Skip("KEP_TEST_GITHUB_TOKEN unset and required for test")
			}

			tmpDir, err := ioutil.TempDir("", "kep-git-test")
			Expect(err).ToNot(HaveOccurred(), "expected no error when creating a temp directory for test")
			defer os.RemoveAll(tmpDir)

			exampleRepoLocation := filepath.Join(tmpDir, exampleRepoName)

			token := &authfakes.FakeTokenProvider{}
			token.ValueReturns(githubToken, nil)

			repo, err := git.Clone(token, exampleRepoUrl, exampleRepoLocation)
			Expect(err).ToNot(HaveOccurred(), "expected no error when cloning a repository with a valid token and nonexistent location")

			fakeGitUrl := "https://github.com/octocat/fakeRepo"

			By("updating the remote `origin`")

			err = repo.SetOrigin(fakeGitUrl)
			Expect(err).ToNot(HaveOccurred(), "expected no error when updating the remote `origin`")

			underlyingRepo, err := libgit.PlainOpen(exampleRepoLocation)
			Expect(err).ToNot(HaveOccurred(), "expected no error when opening a previously cloned repository")

			remote, err := underlyingRepo.Remote(git.OriginRemoteName)
			Expect(err).ToNot(HaveOccurred(), "expected no error when fetching remote `origin` after changing URL")

			remoteUrls := remote.Config().URLs

			Expect(remoteUrls).To(HaveLen(1), "expected remote `origin` to have one associated URL after changing URL")
			Expect(remoteUrls[0]).To(Equal(fakeGitUrl), "expected remote `origin` to have the updated URL after change")
		})

		Context("when `origin` is already set", func() {
			It("updates the remote `origin`", func() {
				githubToken := os.Getenv("KEP_TEST_GITHUB_TOKEN")
				if githubToken == "" {
					Skip("KEP_TEST_GITHUB_TOKEN unset and required for test")
				}

				tmpDir, err := ioutil.TempDir("", "kep-git-test")
				Expect(err).ToNot(HaveOccurred(), "expected no error when creating a temp directory for test")
				defer os.RemoveAll(tmpDir)

				exampleRepoLocation := filepath.Join(tmpDir, exampleRepoName)

				token := &authfakes.FakeTokenProvider{}
				token.ValueReturns(githubToken, nil)

				repo, err := git.Clone(token, exampleRepoUrl, exampleRepoLocation)
				Expect(err).ToNot(HaveOccurred(), "expected no error when cloning a repository with a valid token and nonexistent location")

				fakeGitUrlOne := "https://github.com/octocat/fakeRepo"
				err = repo.SetOrigin(fakeGitUrlOne)
				Expect(err).ToNot(HaveOccurred(), "expected no error when updating the remote `origin`")

				By("updating the remote `origin`")

				fakeGitUrlTwo := "https://github.com/octocat/anotherFakeRepo"
				err = repo.SetOrigin(fakeGitUrlTwo)
				Expect(err).ToNot(HaveOccurred(), "expected no error when updating the remote `origin`")

				underlyingRepo, err := libgit.PlainOpen(exampleRepoLocation)
				Expect(err).ToNot(HaveOccurred(), "expected no error when opening a previously cloned repository")

				remote, err := underlyingRepo.Remote(git.OriginRemoteName)
				Expect(err).ToNot(HaveOccurred(), "expected no error when fetching remote `origin` after changing URL")

				remoteUrls := remote.Config().URLs

				Expect(remoteUrls).To(HaveLen(1), "expected remote `origin` to have one associated URL after changing URL")
				Expect(remoteUrls[0]).To(Equal(fakeGitUrlTwo), "expected remote `origin` to have the updated URL after change")
			})
		})

	})

	Describe("#Checkout()", func() {
		It("checkes out the branch", func() {
			githubToken := os.Getenv("KEP_TEST_GITHUB_TOKEN")
			if githubToken == "" {
				Skip("KEP_TEST_GITHUB_TOKEN unset and required for test")
			}

			tmpDir, err := ioutil.TempDir("", "kep-git-test")
			Expect(err).ToNot(HaveOccurred(), "expected no error when creating a temp directory for test")
			defer os.RemoveAll(tmpDir)

			exampleRepoLocation := filepath.Join(tmpDir, exampleRepoName)

			token := &authfakes.FakeTokenProvider{}
			token.ValueReturns(githubToken, nil)

			repo, err := git.Clone(token, exampleRepoUrl, exampleRepoLocation)
			Expect(err).ToNot(HaveOccurred(), "expected no error when cloning a repository with a valid token and nonexistent location")

			err = repo.Checkout(git.DefaultBranchName)
			Expect(err).ToNot(HaveOccurred(), "expected no error when checking out the existing `master` branch")
		})

		Context("when the branch does not already exist", func() {
			It("creates and checks out the branch", func() {
				githubToken := os.Getenv("KEP_TEST_GITHUB_TOKEN")
				if githubToken == "" {
					Skip("KEP_TEST_GITHUB_TOKEN unset and required for test")
				}

				tmpDir, err := ioutil.TempDir("", "kep-git-test")
				Expect(err).ToNot(HaveOccurred(), "expected no error when creating a temp directory for test")
				defer os.RemoveAll(tmpDir)

				exampleRepoLocation := filepath.Join(tmpDir, exampleRepoName)

				token := &authfakes.FakeTokenProvider{}
				token.ValueReturns(githubToken, nil)

				repo, err := git.Clone(token, exampleRepoUrl, exampleRepoLocation)
				Expect(err).ToNot(HaveOccurred(), "expected no error when cloning a repository with a valid token and nonexistent location")

				newBranchName := "the-new-amazing-thing"

				err = repo.Checkout(newBranchName)
				Expect(err).ToNot(HaveOccurred(), "expected no error when creating and checking out a branch")

				underlyingRepo, err := libgit.PlainOpen(exampleRepoLocation)
				Expect(err).ToNot(HaveOccurred(), "expected no error when opening an existing Git repository")

				_, err = underlyingRepo.Branch(newBranchName)
				Expect(err).ToNot(HaveOccurred(), "expected no error when fetching newly created branch information")

				underlyingConfig, err := underlyingRepo.Config()

				branches := underlyingConfig.Branches
				Expect(branches).To(HaveLen(2), "expected `master` and `the-amazing-new-thing` to exist as branches")

				newBranch := branches[newBranchName]
				Expect(newBranch).ToNot(BeNil(), "expected to be able to retrieve info on branch `the-amazing-new-thing`")
			})
		})

	})

	Describe("#Add()", func() {
		Context("when the path is a directory", func() {
			It("adds all the files in the directory", func() {
				githubToken := os.Getenv("KEP_TEST_GITHUB_TOKEN")
				if githubToken == "" {
					Skip("KEP_TEST_GITHUB_TOKEN unset and required for test")
				}

				tmpDir, err := ioutil.TempDir("", "kep-git-test")
				Expect(err).ToNot(HaveOccurred(), "expected no error when creating a temp directory for test")
				defer os.RemoveAll(tmpDir)

				exampleRepoLocation := filepath.Join(tmpDir, exampleRepoName)

				token := &authfakes.FakeTokenProvider{}
				token.ValueReturns(githubToken, nil)

				repo, err := git.Clone(token, exampleRepoUrl, exampleRepoLocation)
				Expect(err).ToNot(HaveOccurred(), "expected no error when cloning a repository with a valid token and nonexistent location")

				testDirectoryName := "test_directory"
				err = os.MkdirAll(filepath.Join(exampleRepoLocation, testDirectoryName), os.ModePerm)
				Expect(err).ToNot(HaveOccurred(), "expected no error when creating a test directory within an existing Git repository")

				testFilename := "test_file.md"
				err = ioutil.WriteFile(filepath.Join(exampleRepoLocation, testDirectoryName, testFilename), []byte("some test content\n"), os.ModePerm)
				Expect(err).NotTo(HaveOccurred(), "expected no error when writing a test file within an existing Git repository")

				By("adding the given directory to the Git staging area")

				err = repo.Add(testDirectoryName)
				Expect(err).ToNot(HaveOccurred(), "expected no error when staging a directory containing a file for commit")

				underlyingRepo, err := libgit.PlainOpen(exampleRepoLocation)
				Expect(err).ToNot(HaveOccurred(), "expected no error when opening an existing Git repository")

				underlyingWorktree, err := underlyingRepo.Worktree()
				Expect(err).ToNot(HaveOccurred(), "expected no error when opening work tree of Git repository")

				underlyingStatus, err := underlyingWorktree.Status()
				Expect(err).ToNot(HaveOccurred(), "expected no error when determining status of Git repository")

				fileStatus := underlyingStatus.File(filepath.Join(testDirectoryName, testFilename))
				Expect(string(fileStatus.Staging)).To(Equal(string(libgit.Added)), "expected newly added file to have status `Added` in the staging area")
				Expect(string(fileStatus.Worktree)).To(Equal(string(libgit.Unmodified)), "expected newly added file to have status `Unmodified` in the working tree")
			})
		})

		Context("when the path is a regular file", func() {
			It("adds the file", func() {
				githubToken := os.Getenv("KEP_TEST_GITHUB_TOKEN")
				if githubToken == "" {
					Skip("KEP_TEST_GITHUB_TOKEN unset and required for test")
				}

				tmpDir, err := ioutil.TempDir("", "kep-git-test")
				Expect(err).ToNot(HaveOccurred(), "expected no error when creating a temp directory for test")
				defer os.RemoveAll(tmpDir)

				exampleRepoLocation := filepath.Join(tmpDir, exampleRepoName)

				token := &authfakes.FakeTokenProvider{}
				token.ValueReturns(githubToken, nil)

				repo, err := git.Clone(token, exampleRepoUrl, exampleRepoLocation)
				Expect(err).ToNot(HaveOccurred(), "expected no error when cloning a repository with a valid token and nonexistent location")

				testFilename := "test_file.md"
				err = ioutil.WriteFile(filepath.Join(exampleRepoLocation, testFilename), []byte("some test content\n"), os.ModePerm)
				Expect(err).NotTo(HaveOccurred(), "expected no error when writing a test file within an existing Git repository")

				By("adding the given file to the Git staging area")

				err = repo.Add(testFilename)
				Expect(err).ToNot(HaveOccurred(), "expected no error when staging a file for commit")

				underlyingRepo, err := libgit.PlainOpen(exampleRepoLocation)
				Expect(err).ToNot(HaveOccurred(), "expected no error when opening an existing Git repository")

				underlyingWorktree, err := underlyingRepo.Worktree()
				Expect(err).ToNot(HaveOccurred(), "expected no error when opening work tree of Git repository")

				underlyingStatus, err := underlyingWorktree.Status()
				Expect(err).ToNot(HaveOccurred(), "expected no error when determining status of Git repository")

				fileStatus := underlyingStatus.File(testFilename)
				Expect(string(fileStatus.Staging)).To(Equal(string(libgit.Added)), "expected newly added file to have status `Added` in the staging area")
				Expect(string(fileStatus.Worktree)).To(Equal(string(libgit.Unmodified)), "expected newly added file to have status `Unmodified` in the working tree")
			})
		})

		Context("when there are no changes to the path", func() {
			It("returns no error", func() {
				githubToken := os.Getenv("KEP_TEST_GITHUB_TOKEN")
				if githubToken == "" {
					Skip("KEP_TEST_GITHUB_TOKEN unset and required for test")
				}

				tmpDir, err := ioutil.TempDir("", "kep-git-test")
				Expect(err).ToNot(HaveOccurred(), "expected no error when creating a temp directory for test")
				defer os.RemoveAll(tmpDir)

				exampleRepoLocation := filepath.Join(tmpDir, exampleRepoName)

				token := &authfakes.FakeTokenProvider{}
				token.ValueReturns(githubToken, nil)

				repo, err := git.Clone(token, exampleRepoUrl, exampleRepoLocation)
				Expect(err).ToNot(HaveOccurred(), "expected no error when cloning a repository with a valid token and nonexistent location")

				existingFilename := "README"

				err = repo.Add(existingFilename)
				Expect(err).ToNot(HaveOccurred(), "expected no error when staging a file for commit")
			})
		})

	})

	Describe("#Commit()", func() {
		It("commits the staged changes", func() {
			githubToken := os.Getenv("KEP_TEST_GITHUB_TOKEN")
			if githubToken == "" {
				Skip("KEP_TEST_GITHUB_TOKEN unset and required for test")
			}

			tmpDir, err := ioutil.TempDir("", "kep-git-test")
			Expect(err).ToNot(HaveOccurred(), "expected no error when creating a temp directory for test")
			defer os.RemoveAll(tmpDir)

			exampleRepoLocation := filepath.Join(tmpDir, exampleRepoName)

			token := &authfakes.FakeTokenProvider{}
			token.ValueReturns(githubToken, nil)

			repo, err := git.Clone(token, exampleRepoUrl, exampleRepoLocation)
			Expect(err).ToNot(HaveOccurred(), "expected no error when cloning a repository with a valid token and nonexistent location")

			testFilename := "test_file.md"
			err = ioutil.WriteFile(filepath.Join(exampleRepoLocation, testFilename), []byte("some test content\n"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred(), "expected no error when writing a test file within an existing Git repository")

			err = repo.Add(testFilename)
			Expect(err).ToNot(HaveOccurred(), "expected no error when staging a file for commit")

			testCommitName := "OSS KEP Tool"
			testCommitEmail := "oss-kep-tool@noreply.com"
			testCommitMessage := "a great idea to share"

			err = repo.Commit(testCommitName, testCommitEmail, testCommitMessage)
			Expect(err).ToNot(HaveOccurred(), "expected no error when committing a newly created and added file to an existing Git repo")

			underlyingRepo, err := libgit.PlainOpen(exampleRepoLocation)
			Expect(err).ToNot(HaveOccurred(), "expected no error when opening an existing Git repository")

			repoHead, err := underlyingRepo.Head()
			Expect(err).ToNot(HaveOccurred(), "expected no error when determining HEAD of existing Git repository")

			commit, err := underlyingRepo.CommitObject(repoHead.Hash())
			Expect(err).ToNot(HaveOccurred(), "expected no error when retreiving commit at HEAD of existing Git repository")

			Expect(commit.Message).To(Equal(testCommitMessage), "expected commit at HEAD to have commit message matching last commit")

		})

		Context("when there are no changes", func() {
			It("returns no error and does not create a commit", func() {
				githubToken := os.Getenv("KEP_TEST_GITHUB_TOKEN")
				if githubToken == "" {
					Skip("KEP_TEST_GITHUB_TOKEN unset and required for test")
				}

				tmpDir, err := ioutil.TempDir("", "kep-git-test")
				Expect(err).ToNot(HaveOccurred(), "expected no error when creating a temp directory for test")
				defer os.RemoveAll(tmpDir)

				exampleRepoLocation := filepath.Join(tmpDir, exampleRepoName)

				token := &authfakes.FakeTokenProvider{}
				token.ValueReturns(githubToken, nil)

				repo, err := git.Clone(token, exampleRepoUrl, exampleRepoLocation)
				Expect(err).ToNot(HaveOccurred(), "expected no error when cloning a repository with a valid token and nonexistent location")

				testFilename := "test_file.md"
				err = ioutil.WriteFile(filepath.Join(exampleRepoLocation, testFilename), []byte("some test content\n"), os.ModePerm)
				Expect(err).NotTo(HaveOccurred(), "expected no error when writing a test file within an existing Git repository")

				err = repo.Add(testFilename)
				Expect(err).ToNot(HaveOccurred(), "expected no error when staging a file for commit")

				testCommitName := "OSS KEP Tool"
				testCommitEmail := "oss-kep-tool@noreply.com"
				testCommitMessageOne := "a great idea to share"
				testCommitMessageTwo := "a slightly less developed idea to share"

				err = repo.Commit(testCommitName, testCommitEmail, testCommitMessageOne)
				Expect(err).ToNot(HaveOccurred(), "expected no error when committing a newly created and added file to an existing Git repo")

				err = repo.Commit(testCommitName, testCommitEmail, testCommitMessageTwo)
				Expect(err).ToNot(HaveOccurred(), "expected no error when committing with no staged changes")

				underlyingRepo, err := libgit.PlainOpen(exampleRepoLocation)
				Expect(err).ToNot(HaveOccurred(), "expected no error when opening an existing Git repository")

				repoHead, err := underlyingRepo.Head()
				Expect(err).ToNot(HaveOccurred(), "expected no error when determining HEAD of existing Git repository")

				commit, err := underlyingRepo.CommitObject(repoHead.Hash())
				Expect(err).ToNot(HaveOccurred(), "expected no error when retreiving commit at HEAD of existing Git repository")

				Expect(commit.Message).To(Equal(testCommitMessageOne), "expected commit at HEAD to have commit message matching last commit with changes")
			})
		})

	})

	Describe("#PushOrigin()", func() {
		It("pushes the changes from the current branch to the remote `origin`", func() {
			By("cloning the Git repository, adding a file, and comitting")

			githubToken := os.Getenv("KEP_TEST_GITHUB_TOKEN")
			if githubToken == "" {
				Skip("KEP_TEST_GITHUB_TOKEN unset and required for test")
			}

			tmpDir, err := ioutil.TempDir("", "kep-git-test")
			Expect(err).ToNot(HaveOccurred(), "expected no error when creating a temp directory for test")
			defer os.RemoveAll(tmpDir)

			exampleRepoLocation := filepath.Join(tmpDir, exampleRepoName)

			token := &authfakes.FakeTokenProvider{}
			token.ValueReturns(githubToken, nil)

			repo, err := git.Clone(token, exampleRepoUrl, exampleRepoLocation)
			Expect(err).ToNot(HaveOccurred(), "expected no error when cloning a repository with a valid token and nonexistent location")

			branchName := fmt.Sprintf("a-great-idea-%s", uuid.New().String())

			err = repo.Checkout(branchName)
			Expect(err).ToNot(HaveOccurred(), "expected no error creating a new branch in an existing Git repository")

			testFilename := "test_file.md"
			err = ioutil.WriteFile(filepath.Join(exampleRepoLocation, testFilename), []byte("some test content\n"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred(), "expected no error when writing a test file within an existing Git repository")

			err = repo.Add(testFilename)
			Expect(err).ToNot(HaveOccurred(), "expected no error when staging a file for commit")

			testCommitName := "OSS KEP Tool"
			testCommitEmail := "oss-kep-tool@noreply.com"
			testCommitMessage := "a great idea to share"

			err = repo.Commit(testCommitName, testCommitEmail, testCommitMessage)
			Expect(err).ToNot(HaveOccurred(), "expected no error when committing a newly created and added file to an existing Git repo")

			By("pushing committed changes to `origin`")

			err = repo.PushOrigin(token, branchName, branchName)
			Expect(err).ToNot(HaveOccurred(), "expected no error when pushing committed changes to new remote branch")

			fetchUrl := fmt.Sprintf("%s/%s/%s", exampleRepoBlobUrl, branchName, testFilename)

			fetchResp, err := http.Get(fetchUrl)
			Expect(err).ToNot(HaveOccurred(), "expected no error when fetching newly pushed commit as a blob from GitHub")
			defer fetchResp.Body.Close()

			Expect(fetchResp.StatusCode).To(Equal(http.StatusOK), "expected to receieve `200 OK` when fetching blob from GitHub")

			underlyingRepo, err := libgit.PlainOpen(exampleRepoLocation)
			Expect(err).ToNot(HaveOccurred(), "expected no error when opening an existing Git repository")

			err = underlyingRepo.Push(&libgit.PushOptions{
				RefSpecs:   []libgitconfig.RefSpec{libgitconfig.RefSpec(fmt.Sprintf(":refs/heads/%s", branchName))},
				RemoteName: git.OriginRemoteName,
				Auth: &libgithttp.BasicAuth{
					Username: auth.ArbitraryUsername,
					Password: githubToken,
				},
			})

			Expect(err).ToNot(HaveOccurred(), "expected no error deleting a remote Git branch after test")

		})

		Context("when there are no changes to push", func() {
			It("returns no error", func() {
				githubToken := os.Getenv("KEP_TEST_GITHUB_TOKEN")
				if githubToken == "" {
					Skip("KEP_TEST_GITHUB_TOKEN unset and required for test")
				}

				tmpDir, err := ioutil.TempDir("", "kep-git-test")
				Expect(err).ToNot(HaveOccurred(), "expected no error when creating a temp directory for test")
				defer os.RemoveAll(tmpDir)

				exampleRepoLocation := filepath.Join(tmpDir, exampleRepoName)

				token := &authfakes.FakeTokenProvider{}
				token.ValueReturns(githubToken, nil)

				repo, err := git.Clone(token, exampleRepoUrl, exampleRepoLocation)
				Expect(err).ToNot(HaveOccurred(), "expected no error when cloning a repository with a valid token and nonexistent location")

				err = repo.PushOrigin(token, git.DefaultBranchName, git.DefaultBranchName)
				Expect(err).ToNot(HaveOccurred(), "expected no error when pushing no new changes to an existing remote branch")
			})
		})
	})

})
