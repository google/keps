package settings

type Runtime interface {
	Title() string
	Authors() []string
	TargetDir() string
	ContentRoot() string
}

func NewRuntime(contentRoot string, targetDir string, authors []string, title string) Runtime {
	return &runtime{
		authors:     authors,
		title:       title,
		targetDir:   targetDir,
		contentRoot: contentRoot,
	}
}

type runtime struct {
	authors     []string
	title       string
	targetDir   string
	contentRoot string
}

func (r *runtime) Authors() []string   { return r.authors }
func (r *runtime) Title() string       { return r.title }
func (r *runtime) TargetDir() string   { return r.targetDir }
func (r *runtime) ContentRoot() string { return r.contentRoot }
