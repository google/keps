package settings

const (
	Dirname  = "kep"
	Filename = "settings.yaml"
)

type User struct {
	ContentRoot  string `yaml:"content_root"`
	GitHubHandle string `yaml:"github_handle"`
}
