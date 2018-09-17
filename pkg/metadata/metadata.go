package metadata

type KEP struct {
	Title string `yaml:"title"`
	Authors []string `yaml:"authors"`
	Number *int `yaml:"number"`
	OwningSIG string `yaml:"owning_sig"`
	ParticipatingSIGs []string `yaml:"participating_sigs"`
	Reviewers []string `yaml:"reviewers"`
	Approvers []string `yaml:"approvers"`
	Editors []string `yaml:"editors"`
	Status string `yaml:"status"`
	Replaces string `yaml:"replaces"`
}

type Section struct {
	Name string `yaml:"section_name"`
	File string `yaml:"file"`
}

type PullRequest struct {
	URL string `yaml:"url"`
	Status string `yaml:"status"`
	Release string `yaml:"release"`
}

type Issue struct {
	URL string `yaml:"url"`
	Status string `yaml:"status"`
	Release string `yaml:"release"`
}
