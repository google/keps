package workflow

import (
	"github.com/calebamiles/keps/pkg/keps"
	"github.com/calebamiles/keps/pkg/keps/states"
	"github.com/calebamiles/keps/pkg/settings"
)

// Accept signifies that the rationale for a KEP has been accepted by a sponsoring SIG
// and that work is underway to design and document an approach to realize the value
// described by the introduction of a KEP.
// Currently Accept:
//  - adds the principal as both an approver and reviewer
//  - sets the KEP state to `provisional`
//  - persists the KEP to disk
func Accept(runtime settings.Runtime) error {
	p, err := keps.Path(runtime.ContentRoot(), runtime.TargetDir())
	if err != nil {
		return err
	}

	kep, err := keps.Open(p)
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
