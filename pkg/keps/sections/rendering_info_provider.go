package sections

import (
	"time"

	"github.com/calebamiles/keps/pkg/keps/sections/internal/rendering"
	"github.com/calebamiles/keps/pkg/keps/states"
)

//TODO clean up this interface (e.g. whether it should be exported or not)
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
