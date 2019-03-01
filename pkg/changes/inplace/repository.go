package inplace

import (
	"sync"

	"github.com/calebamiles/keps/pkg/changes/auth"
	"github.com/calebamiles/keps/pkg/changes/git"
	"github.com/calebamiles/keps/pkg/changes/github"
	"github.com/calebamiles/keps/pkg/changes/routing"
)

// An inplace.Repo abstracts the details of
//	- staging changes
//	- making a Git commit
//	- rebasing from upstream/targetBranch (stretch) (TODO)
//	- pushing to the remote `origin`
//	- if no ChangeReceipt() then create a pull request
//	- if ChangeReceipt() then just push origin (probably the KEP tool adding the lifecycle PR to metadata)
//	- return the PR
type Repo interface {
	SubmitChanges() (string, error)
	SubmitterName() string
}

func NewRepo(routingInfo routing.Info, underlying git.Repo, submitChanges github.PullRequestCreator) (Repo, error) {
	r := &repository{
		routingInfo: routingInfo,
		underlying:  underlying,
		createPR:    submitChanges,
		locker:      new(sync.Mutex),
	}

	return r, nil
}

type repository struct {
	token       auth.TokenProvider
	routingInfo routing.Info
	underlying  git.Repo
	createPR    github.PullRequestCreator

	locker sync.Locker
}

func (r *repository) SubmitChanges() (string, error) {
	r.locker.Lock()
	defer r.locker.Unlock()

	// add
	loc := r.routingInfo.ChangesUnderPath()

	err := r.underlying.Add(loc)
	if err != nil {
		return "", err
	}

	// commit
	name := r.routingInfo.PrincipalName()
	email := r.routingInfo.PrincipalEmail()
	message := r.routingInfo.ShortSummary()

	err = r.underlying.Commit(name, email, message)
	if err != nil {
		return "", err
	}

	// TODO it would be wonderful to rebase here

	// push origin
	sourceBranch := r.routingInfo.SourceBranch()
	targetBranch := r.routingInfo.TargetBranch()

	err = r.underlying.PushOrigin(r.token, sourceBranch, targetBranch)
	if err != nil {
		return "", err
	}

	// TODO have the empty change receipt return a typed error
	receipt := r.routingInfo.ChangeReceipt()
	switch {
	case receipt != "":
		// we've already created the PR so our push should have just updated it
		return receipt, nil
	default:
		prTitle := github.PullRequestTitle(r.routingInfo.Title())
		prDescription := github.PullRequestDescription(r.routingInfo.FullDescription())

		return r.createPR(r.token, r.routingInfo, prTitle, prDescription)
	}
}

func (r *repository) SubmitterName() string {
	r.locker.Lock()
	defer r.locker.Unlock()

	return r.routingInfo.PrincipalName()
}
