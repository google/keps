package porcelain

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	log "github.com/sirupsen/logrus"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// TODO add branch/upstreamURL override helpers
type Repo interface {
	AddPaths(commitMsg string, paths []string) error
	CreatePR(description string, body string) (string, error) // push local changes, hit PR API
	DeleteLocal() error
	DeleteRemote() error
}

const (
	DefaultBranchName  = "master"
	UpstreamRemoteName = "upstream"
	OriginRemoteName   = "origin"
)

type repo struct {
	gitRepo           *git.Repository
	upstreamApiUrl    string
	deleteGithubRepo  deleteUserRepoFunc
	createPullRequest createPullRequestFunc
	pushLocalChanges  pushRepoChangesFunc
	token             tokenProvider
	locker            *sync.Mutex
	localPath         string
	remoteDeleted     bool
	localDeleted      bool
}

func (r *repo) AddPaths(commitMsg string, paths []string) error {
	r.locker.Lock()
	defer r.locker.Unlock()

	w, err := r.gitRepo.Worktree()
	if err != nil {
		return err
	}

	var errs *multierror.Error
	for _, p := range paths {
		_, addErr := w.Add(p)
		errs = multierror.Append(errs, addErr)
	}

	if errs.ErrorOrNil() != nil {
		return errs
	}

	_, err = w.Commit(commitMsg, &git.CommitOptions{
		// TODO read principal user info and set KEP tool as committer
		Author: &object.Signature{
			Name:  arbitraryBasicAuthUsername,
			Email: kepToolEmail,
			When:  time.Now(),
		},
	})

	if err != nil {
		return err
	}

	// TODO create cute repo return to chain
	// package.Fork("enhancements-tracking").AddPaths(<my-local-kep-content>).CreatePR("KEP: <KEP_TITLE>", <KEP_SUMMARY>)

	return nil
}

func (r *repo) CreatePR(prTitle string, prDescription string) (string, error) {
	r.locker.Lock()
	defer r.locker.Unlock()

	panic("not implemented")

	if r.localDeleted || r.remoteDeleted {
		log.Errorf("requested PR creation after deleting local or remote repository. Local repository deleted: %t. Remote repository deleted: %t", r.localDeleted, r.remoteDeleted)
		return "", fmt.Errorf("cannot create pull request if either local or remote Git repository has been deleted. Local deleted: %t. Remote deleted: %t", r.localDeleted, r.remoteDeleted)
	}

	// push local changes
	err := r.pushLocalChanges()
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

func (r *repo) DeleteLocal() error {
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

func (r *repo) DeleteRemote() error {
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
