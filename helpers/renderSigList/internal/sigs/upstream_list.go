package sigs

type upstreamSIGList struct {
	SIGs []upstreamSIGEntry `yaml:"sigs"`
}

type upstreamSIGEntry struct {
	Name        string                    `yaml:"name"` // we actually want to look at what the SIG is called on disk
	Subprojects []upstreamSubprojectEntry `yaml:"subprojects"`
}

type upstreamSubprojectEntry struct {
	Name string `yaml:"name"`
}

const UpstreamSIGListURL = "https://raw.githubusercontent.com/kubernetes/community/master/sigs.yaml"
