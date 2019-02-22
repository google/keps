package workflow

import (
	"github.com/calebamiles/keps/pkg/keps"
	"github.com/calebamiles/keps/pkg/keps/states"
	"github.com/calebamiles/keps/pkg/orgs"
	"github.com/calebamiles/keps/pkg/procs/enhancements"
	"github.com/calebamiles/keps/pkg/settings"
)

// Propose prepares a KEP for initial proposal to the Kubernetes Community.
// The KEP process has been designed as an iterative process; taking
// inspiration from https://blog.golang.org/toward-go2, Propose prepares
// the author to explain the importance of their change through a KEP
// Propose currently:
//  - sets KEP state to `provisional`
//  - creates a pull request against the enhancements repository
func Propose(runtime settings.Runtime) (string, error) {
	p, err := keps.Path(runtime.ContentRoot(), runtime.TargetDir())
	if err != nil {
		return "", err
	}

	kep, err := keps.Open(p)
	if err != nil {
		return "", err
	}

	err = kep.SetState(states.Provisional)
	if err != nil {
		return "", err
	}

	err = kep.Persist()
	if err != nil {
		return "", err
	}

	var org orgs.Instance

	// this isn't my favorite place for this since it'd be duped for every workflow command, TODO move after walking through the workflow with folks
	switch runtime.IsSandboxed() {
	case true:
		org, err := orgs.NewSandbox(runtime, kep)
		if err != nil {
			return "", err
		}

	default:
		// will be replaced with orgs.Kubernetes when upstreamed
		org = planctae.Organization()
	}

	prUrl, err := enhancements.Propose(runtime, org, kep)
	if err != nil {
		return "", err
	}

	return prUrl, nil
}

