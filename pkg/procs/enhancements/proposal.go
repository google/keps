package enhancements

import (
	log "github.com/sirupsen/logrus"

	"github.com/calebamiles/keps/pkg/changes"
	"github.com/calebamiles/keps/pkg/keps"
	"github.com/calebamiles/keps/pkg/keps/events"
	"github.com/calebamiles/keps/pkg/keps/states"
)

func Propose(repo changes.Submitter, kep keps.Instance) (string, error) {
	err := kep.SetState(states.Provisional)
	if err != nil {
		log.Errorf("setting kep to state %s: %s", states.Provisional, err)
		return "", nil
	}

	err = kep.Persist()
	if err != nil {
		log.Errorf("persisting KEP locally: %s", err)
		return "", nil
	}

	prUrl, err := repo.SubmitChanges()
	if err != nil {
		log.Errorf("creating pull request: %s", err)
		return "", nil
	}

	err = kep.AddLifecyclePR(events.Proposal, repo.SubmitterName(), prUrl)
	if err != nil {
		log.Errorf("adding record of proposal PR to KEP: %s", err)
		return "", nil
	}

	_, err = repo.SubmitChanges()
	if err != nil {
		log.Errorf("updating proposal PR with self reference in KEP metadata: %s", err)
		return "", nil
	}

	return prUrl, nil
}
