package workflow

import (
	"github.com/calebamiles/keps/pkg/keps"
	"github.com/calebamiles/keps/pkg/keps/states"
	"github.com/calebamiles/keps/pkg/settings"
)

// Propose prepares a KEP for initial proposal to the Kubernetes Community.
// The KEP process has been designed as an iterative process; taking
// inspiration from https://blog.golang.org/toward-go2, Propose prepares
// the author to explain the importance of their change through a KEP
// Propose currently:
//  - sets KEP state to `draft`
// Errors returned by Propose are likely due to file i/o
// Eventually, Propose may also handle git and GitHub operations
func Propose(runtime settings.Runtime) error {
	p, err := keps.Path(runtime.ContentRoot(), runtime.TargetDir())
	if err != nil {
		return err
	}

	kep, err := keps.Open(p)
	if err != nil {
		return err
	}

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
