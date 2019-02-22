package inplace

import (
	"github.com/calebamiles/keps/pkg/changes/git"
	"github.com/calebamiles/keps/pkg/changes/github"
	"github.com/calebamiles/keps/pkg/changes/routing"
)

func Open(routingInfo routing.Info) (Repo, error) {
	panic("not tested")

	pathToRepo := routingInfo.PathToRepo()
	gitRepo, err := git.Open(pathToRepo)
	if err != nil {
		return nil, err
	}

	return NewRepo(routingInfo, gitRepo, github.CreatePullRequest)
}
