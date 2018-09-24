package metadata

type KEP struct {
	Authors   []string `yaml:"authors"`
	Title     string   `yaml"title"`
	KEPNumber int      `yaml:"key_number"`
	Reviewers []string `yaml:"reviewers"`
	Approvers []string `yaml:"approvers"`
	Editors   []string `yaml:"editors"`
	Status    string   `yaml:"status"`
	Replaces  string   `yaml:"replaces"`
	Routing
}
