package changeset

type Description interface {
	Title() string
	FullDescription() string
	ShortSummary() string
	Receipt() string
}

type Title string
type FullDescription string
type ShortSummary string

func Describe(title Title, fullDescription FullDescription, shortSummary ShortSummary, receipt func() string) (Description, error) {

	d := &description{
		title:           string(title),
		fullDescription: string(fullDescription),
		shortSummary:    string(shortSummary),
		receipt:         receipt,
	}

	return d, nil
}

type description struct {
	title           string
	fullDescription string
	shortSummary    string
	receipt         func() string
}

func (d *description) Title() string           { return d.title }
func (d *description) FullDescription() string { return d.fullDescription }
func (d *description) ShortSummary() string    { return d.shortSummary }
func (d *description) Receipt() string         { return d.receipt() }
