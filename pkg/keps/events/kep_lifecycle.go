package events

import (
	"time"
)

type Lifecycle string

const (
	Proposal        Lifecycle = "proposal"
	ApprovalRequest Lifecycle = "approval request"
)

type Sortable interface {
	Date() time.Time
}

type ByOldestFirst []Sortable

func (evts ByOldestFirst) Len() int      { return len(evts) }
func (evts ByOldestFirst) Swap(i, j int) { evts[i], evts[j] = evts[j], evts[i] }
func (evts ByOldestFirst) Less(i, j int) bool {
	return evts[i].Date().Before(evts[j].Date())
}
