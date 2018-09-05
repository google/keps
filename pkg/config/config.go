package config

import (
	"github.com/calebamiles/keps/pkg/state"
	"github.com/calebamiles/keps/pkg/sigs"
)

type Config interface {
	RootDir() string
	OwningSIG() []sigs.SIG
	ParticipatingSIGs() []sigs.SIG
	State() state.State
}

