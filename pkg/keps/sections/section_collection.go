package sections

import (
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/hashicorp/go-multierror"
)

type Collection interface {
	// write
	Persist() error
	Erase() error

	// read
	Sections() []string
	ContentDir() string
}

func OpenCollection(locations locationProvider) (Collection, error) {
	var errs *multierror.Error

	contentDir := locations.ContentDir()
	secs := []section{}
	for _, sectionFilename := range locations.Sections() {
		contentBytes, readErr := ioutil.ReadFile(filepath.Join(contentDir, sectionFilename))
		if readErr != nil {
			errs = multierror.Append(errs, readErr)
			continue // try and read the rest of the sections
		}

		s := &readOnlySection{
			commonSectionInfo: &commonSectionInfo{
				filename:   sectionFilename,
				content:    contentBytes,
				name:       sectionNameForFilename(sectionFilename), // this will be a best effort guess for any non top level sections
				contentDir: contentDir,
			},
		}

		secs = append(secs, s)
	}

	return NewCollection(contentDir, secs...), errs.ErrorOrNil()
}

// NewCollection creates a Collection from a variable number of sections along with the content directory where the sections reside
// NewCollection is a [variadic function](https://gobyexample.com/variadic-functions) in order to facilitate testing since we can't pass a slice of a non exported interface into this function from test
func NewCollection(contentDir string, sections ...section) Collection {
	return &collection{
		contentDir: contentDir,
		sections:   sections,
	}
}

type locationProvider interface {
	// Sections are expected to be section file names
	ContentDir() string
	Sections() []string
}

type collection struct {
	sections   []section
	contentDir string
	locker     sync.RWMutex
}

func (c *collection) Persist() error {
	c.locker.Lock()
	defer c.locker.Unlock()

	var errs *multierror.Error

	errs = multierror.Append(errs, c.persist())
	if errs.ErrorOrNil() != nil {
		errs = multierror.Append(errs, c.erase())
	}

	return errs.ErrorOrNil()
}

func (c *collection) Erase() error {
	c.locker.Lock()
	defer c.locker.Unlock()

	return c.erase()
}

func (c *collection) Sections() []string {
	c.locker.RLock()
	defer c.locker.RUnlock()

	sectionFilenames := []string{}
	for i := range c.sections {
		sectionFilenames = append(sectionFilenames, c.sections[i].Filename())
	}

	return sectionFilenames
}

func (c *collection) ContentDir() string {
	c.locker.RLock()
	defer c.locker.RUnlock()

	return c.contentDir
}

func (c *collection) persist() error {
	var errs *multierror.Error

	for _, s := range c.sections {
		errs = multierror.Append(errs, s.Persist())
	}

	return errs.ErrorOrNil()
}

func (c *collection) erase() error {
	var errs *multierror.Error

	for _, s := range c.sections {
		errs = multierror.Append(errs, s.Erase())
	}

	return errs.ErrorOrNil()
}
