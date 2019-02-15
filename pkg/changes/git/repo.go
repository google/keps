package git

import (
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	libgit "gopkg.in/src-d/go-git.v4"
	libgitconfig "gopkg.in/src-d/go-git.v4/config"
	libgitplumbing "gopkg.in/src-d/go-git.v4/plumbing"
	libgitobject "gopkg.in/src-d/go-git.v4/plumbing/object"
	libgithttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"

	"github.com/calebamiles/keps/pkg/changes/auth"
)

type Repo interface {
	SetOrigin(originLocation string) error

	Checkout(branchName string) error // needed by the hermetic.Fork workflow
	Add(path string) error
	Commit(name string, email string, message string) error
	PushOrigin(token auth.TokenProvider, localBranch string, remoteBranch string) error
}

type repository struct {
	underlying *libgit.Repository
	localPath  string

	locker sync.Locker
}

func (r *repository) SetOrigin(originLocation string) error {
	r.locker.Lock()
	defer r.locker.Unlock()

	_, err := r.underlying.CreateRemote(&libgitconfig.RemoteConfig{Name: OriginRemoteName, URLs: []string{originLocation}})
	switch err {
	case nil:
		return nil

	case libgit.ErrRemoteExists:
		delErr := r.underlying.DeleteRemote(OriginRemoteName)
		if delErr != nil {
			return delErr
		}

		_, createErr := r.underlying.CreateRemote(&libgitconfig.RemoteConfig{Name: OriginRemoteName, URLs: []string{originLocation}})
		if createErr != nil {
			return createErr
		}

		return nil

	default:
		log.Errorf("creating `origin` remote with URL: %s due to error: %s", originLocation, err)
		return err
	}
}

func (r *repository) Checkout(branchName string) error {
	r.locker.Lock()
	defer r.locker.Unlock()

	worktree, err := r.underlying.Worktree()
	if err != nil {
		log.Errorf("opening work tree of repository: %s", err)
		return err
	}

	fullBranchName := fmt.Sprintf("refs/heads/%s", branchName)
	err = worktree.Checkout(&libgit.CheckoutOptions{
		Branch: libgitplumbing.ReferenceName(fullBranchName),
	})

	switch err {
	case libgitplumbing.ErrReferenceNotFound:
		// TODO we should fetch upstream and create new branches from that in the future
		headHash, err := r.underlying.ResolveRevision(libgitplumbing.Revision(DefaultBranchName))
		if err != nil {
			log.Errorf("resolving revision at HEAD: %s", err)
			return err
		}

		err = worktree.Checkout(&libgit.CheckoutOptions{
			Branch: libgitplumbing.ReferenceName(fullBranchName),
			Create: true,
			Hash:   *headHash,
		})

		if err != nil {
			log.Errorf("checking out newly created branch: %s with error: %s", branchName, err)
			return err
		}

		// it really feels wrong to have to do this but the tests don't seem to be convinced that we can create a branch without doing this as well the the checkout.
		// TODO fix this
		err = r.underlying.CreateBranch(&libgitconfig.Branch{Name: branchName, Remote: OriginRemoteName, Merge: libgitplumbing.ReferenceName(fullBranchName)})
		if err != nil {
			log.Errorf("creating branch: %s with error: %s", branchName, err)
			return err
		}

		return nil

	case nil:
		return nil

	default:
		log.Errorf("checking out or creating branch: %s with error: %s", branchName, err)
		return err
	}
}

func (r *repository) Add(p string) error {
	r.locker.Lock()
	defer r.locker.Unlock()

	worktree, err := r.underlying.Worktree()
	if err != nil {
		// TODO add log
		return err
	}

	_, err = worktree.Add(p)
	if err != nil {
		return err
	}

	return nil

}

func (r *repository) Commit(name string, email string, message string) error {
	r.locker.Lock()
	defer r.locker.Unlock()

	worktree, err := r.underlying.Worktree()
	if err != nil {
		return err
	}

	status, err := worktree.Status()
	if err != nil {
		log.Errorf("determining work tree status: %s", err)
		return err
	}

	if status.IsClean() {
		log.Info("returning early from repo.Commit because status is clean")
		return nil
	}

	_, err = worktree.Commit(message, &libgit.CommitOptions{
		Author: &libgitobject.Signature{
			Name:  name,
			Email: email,
			When:  time.Now(),
		},

		// TODO add test for this
		Committer: &libgitobject.Signature{
			Name:  "OSS KEP Tool",
			Email: "kubernetes-sig-architecture@googlegroups.com",
			When:  time.Now(),
		},
	})

	if err != nil {
		log.Errorf("committing files: %s", err)
		return err
	}

	return nil
}

func (r *repository) PushOrigin(token auth.TokenProvider, localBranch string, remoteBranch string) error {
	return r.Push(token, OriginRemoteName, localBranch, remoteBranch)
}

func (r *repository) Push(token auth.TokenProvider, remoteName string, localBranch string, remoteBranch string) error {
	r.locker.Lock()
	defer r.locker.Unlock()

	authToken, err := token.Value()
	if err != nil {
		return err
	}

	worktree, err := r.underlying.Worktree()
	if err != nil {
		return err
	}

	repoStatus, err := worktree.Status()
	if err != nil {
		return err
	}

	uncommittedFiles := !repoStatus.IsClean()
	if uncommittedFiles {
		return fmt.Errorf("uncommitted files exist. Please commit or remove them before pushing changes: %s", repoStatus)
	}

	// push local changes
	err = r.underlying.Push(&libgit.PushOptions{
		RefSpecs:   []libgitconfig.RefSpec{libgitconfig.RefSpec(fmt.Sprintf("+refs/heads/%s:refs/heads/%s", localBranch, remoteBranch))},
		RemoteName: remoteName,
		Auth: &libgithttp.BasicAuth{
			Username: auth.ArbitraryUsername,
			Password: authToken,
		},
	})

	switch err {
	case nil:
		return nil
	case libgit.NoErrAlreadyUpToDate:
		return nil
	default:
		return err
	}
}
