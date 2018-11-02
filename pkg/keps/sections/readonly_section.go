package sections

type readOnlySection struct {
	*commonSectionInfo
}

// TODO add info level log that persist/erase called
func (s *readOnlySection) Persist() error { return nil }
func (s *readOnlySection) Erase() error   { return nil }
