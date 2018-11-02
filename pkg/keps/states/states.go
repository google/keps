package states

type Name string

const (
	Draft         Name = "draft"
	Provisional   Name = "provisional"
	Implementable Name = "implementable"
	Implemented   Name = "implemented"
	Deferred      Name = "deferred"
	Rejected      Name = "rejected"
	Withdrawn     Name = "withdrawn"
	Replaced      Name = "replaced"
)
