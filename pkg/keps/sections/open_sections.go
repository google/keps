package sections

import (
	"io/ioutil"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
)

// metadata.KEP should satisfy this interface
type sectionLocationsProvider interface {
	SectionLocations() []string
	ContentDir() string
}

func Open(prv sectionLocationsProvider) ([]Entry, error) {
	locs := prv.SectionLocations()
	contentDir := prv.ContentDir()
	entries := []Entry{}

	var errs *multierror.Error
	for _, loc := range locs {
		readBytes, err := ioutil.ReadFile(filepath.Join(contentDir, loc))
		if err != nil {
			errs = multierror.Append(errs, err)
			// TODO add log
			continue
		}

		entry := &readOnlySection{
			commonSectionInfo: &commonSectionInfo{
				filename:   loc,
				name:       NameForFilename(loc),
				content:    readBytes,
				contentDir: contentDir,
			},
		}

		entries = append(entries, entry)
	}

	if errs.ErrorOrNil() != nil {
		return nil, errs
	}

	return entries, nil
}

type readOnlySection struct {
	*commonSectionInfo
}

// TODO add info level log that persist/erase called
func (s *readOnlySection) Persist() error { return nil }
func (s *readOnlySection) Erase() error   { return nil }
