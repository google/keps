package orgs

import (
	"github.com/calebamiles/keps/pkg/settings"
)

type sandbox struct {
	enhancementsRepositoryOwner         string
	enhancementsRepository              string
	enhancementsRepositoryDefaultBranch string

	enhancementsTrackingRepositoryOwner         string
	enhancementsTrackingRepository              string
	enhancementsTrackingRepositoryDefaultBranch string

	apiReviewRepositoryOwner string
	apiReviewRepository      string
	apiReviewDefaultBranch   string
}

func (s *sandbox) EnhancementsRepositoryOwner() string { return s.enhancementsRepositoryOwner }
func (s *sandbox) EnhancementsRepository() string      { return s.enhancementsRepository }
func (s *sandbox) EnhancementsRepositoryDefaultBranch() string {
	return s.enhancementsRepositoryDefaultBranch
}

func (s *sandbox) EnhancementsTrackingRepositoryOwner() string {
	return s.enhancementsTrackingRepositoryOwner
}
func (s *sandbox) EnhancementsTrackingRepository() string { return s.enhancementsTrackingRepository }
func (s *sandbox) EnhancementsTrackingRepositoryDefaultBranch() string {
	return s.enhancementsTrackingRepositoryDefaultBranch
}

func (s *sandbox) ApiReviewRepositoryOwner() string { return s.apiReviewRepositoryOwner }
func (s *sandbox) ApiReviewRepository() string      { return s.apiReviewRepository }
func (s *sandbox) ApiReviewDefaultBranch() string   { return s.apiReviewDefaultBranch }

func (s *sandbox) IsAuthorized(_ settings.Runtime) (bool, error) {

	// the user is always allowed to use the sandbox

	return true, nil
}
