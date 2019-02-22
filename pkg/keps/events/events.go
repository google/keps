package events

import (
	"errors"
	"time"
)

// implemented by the metadata.kepEvent struct
type Occurrence interface {
	Type() Lifecycle
	Principal() string
	At() time.Time
	Notes() string
}

type Recorder func(event Lifecycle, actor string, notes string) error

// TODO move these to a centralized error package
var (
	ErrNoOccurrence = errors.New("no such event occured")
	ErrNoSuchType   = errors.New("no such event type exists")
)

