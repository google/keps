package section

import (
	"time"

	"github.com/calebamiles/keps/pkg/keps/section/internal/rendering"
	"github.com/calebamiles/keps/pkg/keps/states"
)

type renderingInfoProvider interface {
	Title() string
	Authors() []string
	OwningSIG() string
	State() states.Name
	ContentDir() string
	LastUpdated() time.Time
}

type renderingInfoAndSectionProvider struct {
	ss []rendering.SectionProvider
	renderingInfoProvider
}

func (i *renderingInfoAndSectionProvider) Sections() []rendering.SectionProvider {
	return i.ss
}
