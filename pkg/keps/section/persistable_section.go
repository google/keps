package section

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

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
