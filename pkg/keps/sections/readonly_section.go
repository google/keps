package sections

type readOnlySection struct {
	*commonSectionInfo
}

func (s *readOnlySection) Persist() error { return nil }
func (s *readOnlySection) Erase() error   { return nil }
