package metadata

import (
	"time"
)

type New struct {
	AuthorsField           []string  `yaml:"authors,omitempty"`
	TitleField             string    `yaml:"title,omitempty"`
	ShortIDField           *int      `yaml:"kep_number,omitempty"`
	ReviewersField         []string  `yaml:"reviewers,omitempty"`
	ApproversField         []string  `yaml:"approvers,omitempty"`
	EditorsField           []string  `yaml:"editors,omitempty"`
	StateField             string    `yaml:"state,omitempty"`
	ReplacesField          []string  `yaml:"replaces,omitempty"`
	SupersededByField      []string  `yaml:"superseded_by,omitempty"`
	DevelopmentThemesField []string  `yaml:"development_themes,omitempty"`
	LastUpdatedField       time.Time `yaml:"last_updated,omitempty"`
	CreatedField           time.Time `yaml:"created,omitempty"`
	UniqueIDField          string    `yaml:"uuid,omitempty"`
	SectionsField          []string  `yaml:"sections,omitempty"`

	OwningSIGField           string   `yaml:"owning_sig,omitempty"`
	AffectedSubprojectsField []string `yaml:"affected_subprojects,omitempty"`
	ParticipatingSIGsField   []string `yaml:"participating_sigs,omitempty"`
	KubernetesWideField      bool     `yaml:"kubernetes_wide,omitempty"`
	SIGWideField             bool     `yaml:"sig_wide,omitempty"`
}

func (n *New) Authors() []string      { return n.AuthorsField }
func (n *New) Title() string          { return n.TitleField }
func (n *New) State() string          { return n.StateField }
func (n *New) LastUpdated() time.Time { return n.LastUpdatedField }
func (n *New) Created() time.Time     { return n.CreatedField }
func (n *New) OwningSIG() string      { return n.OwningSIGField }
