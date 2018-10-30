package sections

type commonSectionInfo struct {
	filename   string
	name       string
	content    []byte
	contentDir string
}

func (i *commonSectionInfo) Filename() string { return i.filename }
func (i *commonSectionInfo) Name() string     { return i.name }
func (i *commonSectionInfo) Content() []byte  { return i.content }
