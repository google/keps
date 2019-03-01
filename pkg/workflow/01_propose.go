package workflow

import (
	"github.com/calebamiles/keps/pkg/keps"
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
func Propose(runtime settings.Runtime, org orgs.Instance, kep keps.Instance) (string, error) {
	repo, err := enhancements.OpenRepo(runtime, org, kep)
	if err != nil {
		return "", err
	}

	prUrl, err := enhancements.Propose(repo, kep)
	if err != nil {
		return "", err
	}

	return prUrl, nil
}
