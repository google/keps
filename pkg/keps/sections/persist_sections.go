package sections

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
)

// Persist attempts persist all given sections; returning all errors
func Persist(ss []Entry) error {
	var errs *multierror.Error

	for _, section := range ss {
		errs = multierror.Append(errs, section.Persist())
	}

	return errs.ErrorOrNil()
}

type persistableSection struct {
	*commonSectionInfo
}

func (s *persistableSection) Persist() error {
	loc := filepath.Join(s.contentDir, s.filename)
	return ioutil.WriteFile(loc, s.content, os.ModePerm)
}

func (s *persistableSection) Erase() error {
	loc := filepath.Join(s.contentDir, s.filename)
	return os.Remove(loc)
}
