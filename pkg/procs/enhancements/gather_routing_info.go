package enhancements

import (
	"fmt"
	"strings"


	"github.com/calebamiles/keps/pkg/changes/routing"
	"github.com/calebamiles/keps/pkg/keps"
	"github.com/calebamiles/keps/pkg/keps/events"
	"github.com/calebamiles/keps/pkg/orgs"
	"github.com/calebamiles/keps/pkg/settings"
)

func GatherRoutingFrom(runtime settings.Runtime, kep keps.Instance, org orgs.Instance, event events.Lifecycle) (routing.Info, error) {
	sourceRepositoryOwner := runtime.PrincipalGithubHandle()
	sourceRepository := org.EnhancementsRepository()
	sourceBranch := org.EnhancementsRepositoryDefaultBranch()

	targetRepositoryOwner := org.EnhancementsRepositoryOwner()
	targetRepository := org.EnhancementsRepository()
	targetBranch := org.EnhancementsRepositoryDefaultBranch()

	committerEmail := runtime.PrincipalEmail()
	committerDisplayName :=  runtime.PrincipalDisplayName()

	title := strings.Title(fmt.Sprintf("propose %s for sponsorship by %s", kep.Title(), kep.OwningSIG()))
	fullDescription := strings.Title(fmt.Sprintf("propose %s for sponsorship by %s", kep.Title(), kep.OwningSIG()))
	shortSummary := fmt.Sprintf("kep-lifecycle: propose %s for sponsorship by %s", kep.Title(), kep.OwningSIG())

	owningSig := kep.OwningSIG()

	localRepositoryLocation := runtime.ContentRoot()
	changesUnderPath := kep.ContentDir()

	var changeReceipt = func() string {
		return kep.GetLifecyclePR(event)
	}

	info, err := routing.NewInfo(
		changeReceipt,
		routing.SourceOwner(sourceRepositoryOwner),
		routing.SourceRepository(sourceRepository),
		routing.SourceBranch(sourceBranch),
		routing.TargetOwner(targetRepositoryOwner),
		routing.TargetRepository(targetRepository),
		routing.TargetBranch(targetBranch),
		routing.Title(title),
		routing.FullDescription(fullDescription),
		routing.ShortSummary(shortSummary),
		routing.OwningSIG(owningSig),
		routing.PrincipalName(committerDisplayName),
		routing.PrincipalEmail(committerEmail),
		routing.PathToRepo(localRepositoryLocation),
		routing.PathToChanges(changesUnderPath),
	)

	if err != nil {
		return nil, err
	}

	return info, nil
}

