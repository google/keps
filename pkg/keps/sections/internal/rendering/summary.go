package rendering

import (
	"bytes"
	"text/template"

	"github.com/calebamiles/keps/pkg/keps/sections/internal/unrendered"
)

const (
	SummaryName     = "Summary"
	SummaryFilename = "summary.md"
)

func NewSummary(info InfoProvider) ([]byte, error) {
	sectionContent := &bytes.Buffer{}

	t, err := template.New(SummaryName).Parse(unrendered.Summary)
	if err != nil {
		return nil, err
	}

	err = t.Execute(sectionContent, info)
	if err != nil {
		return nil, err
	}

	return sectionContent.Bytes(), nil
}
