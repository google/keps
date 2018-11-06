package workflow

import (
	"github.com/calebamiles/keps/pkg/keps"
	"github.com/calebamiles/keps/pkg/keps/states"
	"github.com/calebamiles/keps/pkg/settings"
)

func Accept(runtime settings.Runtime) error {
	kep, err := keps.Open(runtime.TargetDir())
	if err != nil {
		return err
	}

	kep.AddApprovers(runtime.Principal())
	kep.AddReviewers(runtime.Principal())

	err = kep.SetState(states.Provisional)
	if err != nil {
		return err
	}

	err = kep.Persist()
	if err != nil {
		return err
	}

	// TODO add mechanics for creating PR

	return nil
}
