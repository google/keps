package sections

import (
	"github.com/calebamiles/keps/pkg/keps/sections/internal/rendering"
)

func IsAutogenerated(name string) bool {
	switch name {
	case Readme:
		return true
	// TODO insert hugo autogenerated frontmatter
	default:
		return false
	}
}

func AutoGeneratedFrom(info renderingInfoProvider) ([]Entry, error) {
	readme, err := newReadme(info)
	if err != nil {
		return nil, err
	}

	entries := []Entry{
		readme,
		// TODO insert hugo autogenerated sections
	}

	return entries, nil
}

func newReadme(info renderingInfoProvider) (Entry, error) {
	readmeBytes, err := rendering.NewReadme(info)
	if err != nil {
		return nil, err
	}

	sec := &persistableSection{
		commonSectionInfo: &commonSectionInfo{
			filename:   rendering.ReadmeFilename,
			name:       rendering.ReadmeName,
			contentDir: info.ContentDir(),
			content:    readmeBytes,
		},
	}

	return sec, nil
}