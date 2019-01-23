package workflow

import (
	"github.com/calebamiles/keps/pkg/keps"
	"github.com/calebamiles/keps/pkg/keps/states"
	"github.com/calebamiles/keps/pkg/settings"
)

// Approve allows an approver to signal that a KEP is
// approved for implementation. The KEP is checked for
// consistency before the state is updated, which can
// return an error
func Approve(runtime settings.Runtime) error {
	p, err := keps.Path(runtime.ContentRoot(), runtime.TargetDir())
	if err != nil {
		return err
	}

	kep, err := keps.Open(p)
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

	// TODO add mechanics for creating PR

	return nil
}
