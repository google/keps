package metadata

type PullRequest struct {
	URL     string `yaml:"url"`
	Status  string `yaml:"status"`
	Release string `yaml:"release"`
}

type Issue struct {
	URL     string `yaml:"url"`
	Status  string `yaml:"status"`
	Release string `yaml:"release"`
}
