package states

type State string

const (
	Provisional   State = "provisional"
	Implementable State = "implementable"
	Implemented   State = "implemented"
	Deferred      State = "deferred"
	Rejected      State = "rejected"
	Withdrawn     State = "withdrawn"
	Replaced      State = "replaced"
)
