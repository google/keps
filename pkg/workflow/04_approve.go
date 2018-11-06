package workflow

import (
	"github.com/calebamiles/keps/pkg/index"
	"github.com/calebamiles/keps/pkg/keps"
	"github.com/calebamiles/keps/pkg/keps/states"
	"github.com/calebamiles/keps/pkg/settings"
)

func Approve(runtime settings.Runtime) error {
	kep, err := keps.Open(runtime.TargetDir())
	if err != nil {
		return err
	}

	err = kep.SetState(states.Implementable)
	if err != nil {
		return err
	}

	err = kep.Persist()
	if err != nil {
		return err
	}

	// our KEP should now be added to the index
	idx, err := index.Rebuild(runtime)
	if err != nil {
		return err
	}

	err = idx.Persist()
	if err != nil {
		return err
	}

	// TODO add mechanics for creating PR

	return nil
}
