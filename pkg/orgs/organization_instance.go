package orgs

import (
	"github.com/calebamiles/keps/pkg/settings"
)

type Instance interface {
	// where KEP content lives
	EnhancementsRepositoryOwner() string
	EnhancementsRepository() string
	EnhancementsRepositoryDefaultBranch() string

	// where release tracking happens
	EnhancementsTrackingRepositoryOwner() string
	EnhancementsTrackingRepository() string
	EnhancementsTrackingRepositoryDefaultBranch() string

	// where API review happens
	ApiReviewRepositoryOwner() string
	ApiReviewRepository() string
	ApiReviewDefaultBranch() string

	// security
	IsAuthorized(runtime settings.Runtime) (bool, error)
}

