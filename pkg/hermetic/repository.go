package hermetic

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// A Repo abstracts the basic Git and GitHub operations required to submit a GitHub pull request against
// an upstream repository
//
// Besides two obvious cleanup operations (delete forked GitHub repository, delete local Git repository), A
// Repo implementation should expose just  enough functionality to add files to the repository and create
// a GitHub Pull Request containing only those files
type Repo interface {
	Add(fromLocation string, toLocationAfterRoot string) error
	CreatePR(description string, body string) (string, error) // push local changes, hit PR API
	DeleteLocal() error
	DeleteRemote() error
}

const (
	DefaultBranchName  = "master"
	UpstreamRemoteName = "upstream"
	OriginRemoteName   = "origin"
)

type repository struct {
	gitRepo           *git.Repository
	branchName        string
	deleteGithubRepo  deleteUserRepoFunc
	createPullRequest createPullRequestFunc
	pushLocalChanges  pushRepoChangesFunc
	token             tokenProvider
	locker            *sync.Mutex
	localPath         string
	remoteDeleted     bool
	localDeleted      bool
}

func (r *repository) Add(fromLocation string, toLocation string) error {
	r.locker.Lock()
	defer r.locker.Unlock()

	from, err := os.Open(fromLocation)
	if err != nil {
		return err
	}

	worktree, err := r.gitRepo.Worktree()
	if err != nil {
		return err
	}

	dst := filepath.Join(r.localPath, toLocation)

	to, err := os.Create(dst)
	if err != nil {
		return err
	}

	_, err = io.Copy(from, to)
	if err != nil {
		return err
	}

	err = to.Close()
	if err != nil {
		return err
	}

	err = from.Close()
	if err != nil {
		return err
	}

	_, err = worktree.Add(toLocation)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) CreatePR(prTitle string, prDescription string) (string, error) {
	r.locker.Lock()
	defer r.locker.Unlock()

	if r.localDeleted || r.remoteDeleted {
		log.Errorf("requested PR creation after deleting local or remote repository. Local repository deleted: %t. Remote repository deleted: %t", r.localDeleted, r.remoteDeleted)
		return "", fmt.Errorf("cannot create pull request if either local or remote Git repository has been deleted. Local deleted: %t. Remote deleted: %t", r.localDeleted, r.remoteDeleted)
	}

	// make sure we're on the right branch
	worktree, err := r.gitRepo.Worktree()
	if err != nil {
		return "", err
	}

	commitMsg := fmt.Sprintf("%s\n\n%s", prTitle, prDescription)
	_, err = worktree.Commit(commitMsg, &git.CommitOptions{
		// TODO read principal user info and set KEP tool as committer
		Author: &object.Signature{
			Name:  arbitraryBasicAuthUsername,
			Email: kepToolEmail,
			When:  time.Now(),
		},
	})

	if err != nil {
		log.Errorf("committing files before pushing: %s", err)
		return "", err
	}

	// push local changes
	err = r.pushLocalChanges()
	if err != nil {
		log.Errorf("pushing local changes: %s", err)
		return "", err
	}

	// use GitHub API to create PR
	prUrl, err := r.createPullRequest(prTitle, prDescription)
	if err != nil {
		log.Errorf("creating GitHub Pull Request: %s", err)
		return "", err
	}

	return prUrl, nil
}

func (r *repository) DeleteLocal() error {
	r.locker.Lock()
	defer r.locker.Unlock()

	if r.localDeleted {
		log.Info("local Git repository should already be deleted, skipping request")
		return nil
	}

	err := os.RemoveAll(r.localPath)
	if err != nil {
		log.Errorf("deleting local Git repository: %s", err)
		return err
	}

	r.localDeleted = true

	return nil
}

func (r *repository) DeleteRemote() error {
	r.locker.Lock()
	defer r.locker.Unlock()

	if r.remoteDeleted {
		log.Info("GitHub repository should already be deleted, skipping request")
		return nil
	}

	err := r.deleteGithubRepo()
	if err != nil {
		log.Errorf("deleting GitHub fork: %s", err)
		return err
	}

	r.remoteDeleted = true

	return nil
}

type tokenProvider interface {
	Value() (string, error)
}
