package release

import (
	"fmt"

	"github.com/calebamiles/keps/pkg/changes/hermetic"
	"github.com/calebamiles/keps/pkg/keps"
	"github.com/calebamiles/keps/pkg/orgs/planctae"
	"github.com/calebamiles/keps/pkg/settings"
)

func CreateKepTrackingPr(runtime settings.Runtime, targetRelease string) (string, error) {
	p, err := keps.Path(runtime.ContentRoot(), runtime.TargetDir())
	if err != nil {
		return "", err
	}

	kep, err := keps.Open(p)
	if err != nil {
		return "", err
	}

	// add check that KEP is releasable
	kep.AddChecks(check.ThatIsReleasable)
	err = kep.Check()
	if err != nil {
		return "", err
	}

	githubHandle := runtime.Principal()
	token := runtime.TokenProvider()

	canOperate, err := planctae.IsAuthorized(githubHandle, token)
	if err != nil {
		return "", err
	}

	if !canOperate {
		return "", fmt.Errorf("the organization %s, believes that user %s is not authorized to create an enhancement tracking PR", planctae.Organization, githubHandle)
	}

	repo, err := hermetic.Fork(
		githubHandle,
		token,
		planctae.Organization,
		TrackingRepo,
		cloneLocation,
		branchName(kep.Title(), targetRelease),
	)

	if err != nil {
		return "", err
	}

	defer repo.DeleteLocal()

	summaryLocation := filepath.Join(runtime.TargetDir(), sections.SummaryFilename)
	releaseTrackingDir := trackingDir(targetRelease)

	err = repo.Add(summaryLocation, releaseTrackingDir, commitMessage(kep.Title(), targetRelease)
	if err != nil {
		return "", err
	}

	// create PR
	prUrl, err := repo.CreatePR(prMessage(kep.Title(), targetRelease))
	if err != nil {
		return "", err
	}

	// associate PR with KEP (to write) keps.LinkPR()
	err = kep.LinkPR(prUrl)
	if err != nil {
		return "", err
	}

	// target release in KEP (to write) keps.TargetRelease()
	err = kep.TargetRelease(targetRelease)
	if err != nil {
		return "", err
	}

	return prUrl, nil
}
