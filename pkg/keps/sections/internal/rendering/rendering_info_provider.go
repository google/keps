package rendering

import (
	"time"

	"github.com/calebamiles/keps/pkg/keps/states"
)

// TODO delete this interface if unused
type SectionProvider interface {
	Name() string
	Filename() string
	Content() []byte
}

type InfoProvider interface {
	Title() string
	Authors() []string
	OwningSIG() string
	State() states.Name
	ContentDir() string
	LastUpdated() time.Time
	SectionLocations() []string
}
