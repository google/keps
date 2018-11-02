package settings

type Runtime interface {
	Principal() string
	TargetDir() string
	ContentRoot() string
}

func NewRuntime(contentRoot string, targetDir string, principal string) Runtime {
	return &runtime{
		principal:   principal,
		targetDir:   targetDir,
		contentRoot: contentRoot,
	}
}

type runtime struct {
	principal   string
	targetDir   string
	contentRoot string
}

func (r *runtime) Principal() string   { return r.principal }
func (r *runtime) TargetDir() string   { return r.targetDir }
func (r *runtime) ContentRoot() string { return r.contentRoot }
