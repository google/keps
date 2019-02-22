package enhancements

import (
	log "github.com/sirupsen/logrus"

	"github.com/calebamiles/keps/pkg/changes/auth"
	"github.com/calebamiles/keps/pkg/changes/inplace"
	"github.com/calebamiles/keps/pkg/changes/routing"
	"github.com/calebamiles/keps/pkg/keps"
	"github.com/calebamiles/keps/pkg/keps/events"
	"github.com/calebamiles/keps/pkg/orgs"
	"github.com/calebamiles/keps/pkg/settings"
)

func Propose(runtime settings.Runtime, org orgs.Instance, kep keps.Instance) (string, error) {
	routingInfo, err := GatherRoutingFrom(runtime, kep, org, events.Proposal)
	if err != nil {
		log.Errorf("gathering information required to propose KEP: %s", err)
		return "", err
	}

	token, err := auth.NewProvideTokenFromPath(runtime.TokenPath())
	if err != nil {
		return "", err
	}

	return ProposeFrom(token, kep.AddLifecyclePR, routingInfo)
}

func ProposeFrom(token auth.TokenProvider, record events.Recorder, routingInfo routing.Info) (string, error) {
	repo, err := inplace.Open(routingInfo)
	if err != nil {
		return "", err
	}

	// create PR
	prUrl, err := repo.SubmitChanges()
	if err != nil {
		return "", err
	}

	err = record(events.Proposal, routingInfo.SourceRepositoryOwner(), prUrl)
	if err != nil {
		return "", err
	}

	// add commit to existing PR
	prUrl, err = repo.SubmitChanges()
	if err != nil {
		return "", err
	}

	return prUrl, nil
}

