package hermetic

import (
	"fmt"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/calebamiles/keps/pkg/changes/git"
	"github.com/calebamiles/keps/pkg/changes/github"
	"github.com/calebamiles/keps/pkg/changes/routing"
)

const (
	arbitraryBasicAuthUsername = "ossKEPtool" // Username can be anything except empty
)

// Fork is responsible for "forking" an upstream Git repository that is hosted on GitHub into the account of the authenticated GitHub user
// with the goal of proposing an isolated change to upstream. Initial upstream targets are imagined to be an "enhancements tracking"
// repository, used for managing release over release changes to Kubernetes, as well as an "API review tracking" repository used by SIG
// Architecture to manage API changes with a Kubernetes wide scope
//
// Fork:
//   - Issues a fork request to the GitHub API
//   - Clones the default branch from upstream
//   - Creates a new branch to add changes
//   - Sets the Git remote name "origin" to the forked repository in the account of `githubHandle`
func Fork(routingInfo routing.Info) (Repo, error) {
	toLocation := routingInfo.PathToRepo()
	if _, err := os.Stat(toLocation); !os.IsNotExist(err) {
		log.Errorf("location: %s may exist already, refusing to overwrite", toLocation)
		return nil, fmt.Errorf("location: %s may exist already, refusing to overwrite", toLocation)
	}

	token := routingInfo.Token()
	owner := routingInfo.SourceRepositoryOwner()
	repo := routingInfo.SourceRepository()
	githubHandle := routingInfo.TargetRepositoryOwner()
	withBranchName := routingInfo.PathToRepo()

	// TODO think about dropping the URL return value
	_, err := github.Fork(token, owner, repo)
	if err != nil {
		log.Errorf("forking upstream %s: %s", github.ForkUrl(owner, repo), err)
		return nil, err
	}

	gitRepo, err := git.Clone(token, github.GitUrl(owner, repo), toLocation)
	if err != nil {
		log.Errorf("cloning `upstream` remote: %s", err)
		return nil, err
	}

	err = gitRepo.SetOrigin(github.GitUrl(githubHandle, repo))
	if err != nil {
		log.Errorf("setting origin url: %s, with error: %s", github.GitUrl(githubHandle, repo), err)
		return nil, err
	}

	err = gitRepo.Checkout(withBranchName)
	if err != nil {
		log.Errorf("checking out new branch: %s, with error: %s", withBranchName, err)
		return nil, err
	}

	return NewRepo(routingInfo, gitRepo, github.CreatePullRequest)
}
