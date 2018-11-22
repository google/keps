package metadata

import (
	"time"
)

type Old struct {
	KepNumber         string    `yaml:"kep_number,omitempty"`
	Title             string    `yaml:"title,omitempty"`
	Status            string    `yaml:"status,omitempty"`
	Authors           []string  `yaml:"authors"`
	OwningSIG         string    `yaml:"owning-sig,omitempty"`
	ParticipatingSIGs []string  `yaml:"participating-sigs,omitempty"`
	Reviewers         []string  `yaml:"reviewers,omitempty"`
	Approvers         []string  `yaml:"approvers,omitempty"`
	Editor            string    `yaml:"editor,omitempty"`
	CreationDate      time.Time `yaml:"creation-date,omitempty"`
	LastUpdated       time.Time `yaml:"last-updated,omitempty"`
	SeeAlso           []string  `yaml:"see-also,omitempty"`
	Replaces          []string  `yaml:"replaces,omitempty"`
	SupersededBy      []string  `yaml:"superseded-by,omitempty"`
}
