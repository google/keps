package routing

import (
	"github.com/calebamiles/keps/pkg/changes/changeset"
	"github.com/calebamiles/keps/pkg/keps"
	"github.com/calebamiles/keps/pkg/orgs"
	"github.com/calebamiles/keps/pkg/settings"
)

func NewInfoFrom(runtime settings.Runtime, kep keps.Instance, org orgs.Instance, changeSet changeset.Description) (Info, error) {
	sourceRepositoryOwner := runtime.PrincipalGithubHandle()
	sourceRepository := org.EnhancementsRepository()
	sourceBranch := org.EnhancementsRepositoryDefaultBranch()

	targetRepositoryOwner := org.EnhancementsRepositoryOwner()
	targetRepository := org.EnhancementsRepository()
	targetBranch := org.EnhancementsRepositoryDefaultBranch()

	committerEmail := runtime.PrincipalEmail()
	committerDisplayName := runtime.PrincipalDisplayName()

	title := changeSet.Title()
	fullDescription := changeSet.FullDescription()
	shortSummary := changeSet.ShortSummary()

	owningSig := kep.OwningSIG()

	localRepositoryLocation := runtime.ContentRoot()
	changesUnderPath := kep.ContentDir()

	info, err := NewInfo(
		changeSet.Receipt,
		SourceOwner(sourceRepositoryOwner),
		SourceRepository(sourceRepository),
		SourceBranch(sourceBranch),
		TargetOwner(targetRepositoryOwner),
		TargetRepository(targetRepository),
		TargetBranch(targetBranch),
		Title(title),
		FullDescription(fullDescription),
		ShortSummary(shortSummary),
		OwningSIG(owningSig),
		PrincipalName(committerDisplayName),
		PrincipalEmail(committerEmail),
		PathToRepo(localRepositoryLocation),
		PathToChanges(changesUnderPath),
	)

	if err != nil {
		return nil, err
	}

	return info, nil
}
