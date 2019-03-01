package changes

type Submitter interface {
	SubmitterName() string
	SubmitChanges() (string, error)
}
