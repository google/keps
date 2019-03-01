package hermetic

import (
	"path/filepath"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/calebamiles/keps/pkg/changes/git"
	"github.com/calebamiles/keps/pkg/changes/github"
	"github.com/calebamiles/keps/pkg/changes/routing"
)

// A Repo abstracts the basic Git and GitHub operations required to submit a GitHub pull request against
// an upstream repository
type Repo interface {
	CopyFrom(fromLocation string) error
	SubmitChanges() (string, error)
	ReplaceRoutingInfo(routingInfo routing.Info) error
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
	routingInfo routing.Info
	underlying  git.Repo
	createPR    github.PullRequestCreator

	locker sync.Locker
}

func (r *repository) CopyFrom(fromLocation string) error {
	r.locker.Lock()
	defer r.locker.Unlock()

	repoLocation := r.routingInfo.PathToRepo()
	toLocationAfterRoot := r.routingInfo.ChangesUnderPath()

	// we manually join these strings so that filepath.Clean doesn't strip off the trailing filepath.Separator we need to recognize directories
	return RecursiveCopy(fromLocation, strings.Join([]string{repoLocation, toLocationAfterRoot}, string(filepath.Separator)))
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

	// push origin
	sourceBranch := r.routingInfo.SourceBranch()
	targetBranch := r.routingInfo.TargetBranch()
	token := r.routingInfo.Token()

	err = r.underlying.PushOrigin(token, sourceBranch, targetBranch)
	if err != nil {
		return "", err
	}

	// TODO: eventually have ChangeReceipt() return a typed error not empty string
	receipt, err := r.routingInfo.ChangeReceipt()
	if err != nil {
		return "", err
	}

	if receipt != "" {
		// we're not quitting early because this condition is not expected with the hermetic flow
		log.Warnf("existing PR may have already been created: %s", receipt)
	}

	prTitle := github.PullRequestTitle(r.routingInfo.Title())
	prDescription := github.PullRequestDescription(r.routingInfo.FullDescription())

	return r.createPR(r.routingInfo, prTitle, prDescription)
}

func (r *repository) ReplaceRoutingInfo(routingInfo routing.Info) error {
	r.locker.Lock()
	defer r.locker.Unlock()

	r.routingInfo = routingInfo
	return nil
}
