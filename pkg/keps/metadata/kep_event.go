package metadata

import (
	"time"

	"github.com/calebamiles/keps/pkg/keps/events"
)

type kepEvent struct {
	PrincipalField string           `yaml:"principal,omitempty"`
	DateField      time.Time        `yaml:"date,omitempty"`
	TypeField      events.Lifecycle `yaml:"event_type,omitempty"`
	NotesField     string           `yaml:"notes,omitempty"`
}

func (e *kepEvent) Type() events.Lifecycle { return e.TypeField }
func (e *kepEvent) Principal() string      { return e.PrincipalField }
func (e *kepEvent) At() time.Time          { return e.DateField }
func (e *kepEvent) Notes() string          { return e.NotesField }

type byOldestFirst []*kepEvent

func (evts byOldestFirst) Len() int      { return len(evts) }
func (evts byOldestFirst) Swap(i, j int) { evts[i], evts[j] = evts[j], evts[i] }
func (evts byOldestFirst) Less(i, j int) bool {
	return evts[i].DateField.Before(evts[j].DateField)
}
