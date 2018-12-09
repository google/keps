package porcelain

import (
	"fmt"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"

	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

const (
	arbitraryBasicAuthUsername = "ossKEPtool" // Username can be anything except empty
)

const (
	// TODO extract this to a higher level caller
	UpstreamEnhancementsTrackingRepoApiUrl = "https://api.github.com/repos/Charkha/enhancements-tracking/forks"
)

func Fork(githubHandle string, token tokenProvider, upstreamApiUrl string, upstreamGitUrl string, toLocation string, withBranchName string) (Repo, error) {
	if _, err := os.Stat(toLocation); !os.IsNotExist(err) {
		log.Errorf("location: %s may exist already, refusing to overwrite", toLocation)
		return nil, err
	}

	// call fork API
	forkUpstream, err := newCreateForkFunc(packageHttpClient, token, upstreamApiUrl)
	if err != nil {
		log.Errorf("creating fork function: %s", err)
		return nil, err
	}

	forkedRepoApiUrl, forkedRepoGitUrl, err := forkUpstream()
	if err != nil {
		log.Errorf("forking upstream %s: %s", upstreamApiUrl, err)
		return nil, err
	}

	authToken, err := token.Value()
	if err != nil {
		return nil, err
	}

	gitRepo, err := git.PlainClone(toLocation, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: arbitraryBasicAuthUsername,
			Password: authToken,
		},
		URL: upstreamGitUrl,
	})

	if err != nil {
		log.Errorf("cloning `upstream` remote: %s", err)
		return nil, err
	}

	err = gitRepo.DeleteRemote(OriginRemoteName)
	if err != nil {
		log.Errorf("deleting remote `origin`: %s", err)
		return nil, err
	}

	_, err = gitRepo.CreateRemote(&config.RemoteConfig{Name: UpstreamRemoteName, URLs: []string{upstreamGitUrl}})
	if err != nil {
		log.Errorf("creating `upstream` remote: %s", err)
		return nil, err
	}

	// create new branch
	err = gitRepo.CreateBranch(&config.Branch{Name: withBranchName, Remote: UpstreamRemoteName})
	if err != nil {
		log.Errorf("creating Git branch: %s", err)
		return nil, err
	}

	// construct source location
	sourceLocation := fmt.Sprintf("%s:master", githubHandle)

	// create pull request callback
	createPullRequest, err := newCreatePRFunc(packageHttpClient, token, forkedRepoApiUrl, sourceLocation)
	if err != nil {
		log.Errorf("creating pull request creator callback: %s", err)
		return nil, err
	}

	// set origin to forked repo
	_, err = gitRepo.CreateRemote(&config.RemoteConfig{Name: OriginRemoteName, URLs: []string{forkedRepoGitUrl}})
	if err != nil {
		log.Errorf("setting origin URL to forked location: %s", err)
		return nil, err
	}

	// create delete repo callback
	deleteGithubRepo, err := newDeleteGithubUserRepoFunc(packageHttpClient, token, forkedRepoApiUrl)
	if err != nil {
		log.Errorf("creating delete forked repo callback: %s", err)
		return nil, err
	}

	// create push local changes callback
	pushLocalChanges, err := newPushRepoChangesFunc(gitRepo, token)
	if err != nil {
		log.Errorf("creating push local changes callback: %s", err)
		return nil, err
	}

	r := &repo{
		gitRepo:           gitRepo,
		createPullRequest: createPullRequest,
		deleteGithubRepo:  deleteGithubRepo,
		pushLocalChanges:  pushLocalChanges,
		localPath:         toLocation,
		token:             token,
		locker:            &sync.Mutex{},
	}

	return r, nil
}
