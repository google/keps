package rendering

import (
	"bytes"
	"text/template"

	"github.com/calebamiles/keps/pkg/keps/section/internal/unrendered"
)

const (
	OperatorGuideName     = "Operator Guide"
	OperatorGuideFilename = "guides/operator.md"
)

func NewOperatorGuide(info InfoProvider) ([]byte, error) {
	sectionContent := &bytes.Buffer{}

	t, err := template.New(OperatorGuideName).Parse(unrendered.OperatorGuide)
	if err != nil {
		return nil, err
	}

	err = t.Execute(sectionContent, info)
	if err != nil {
		return nil, err
	}

	return sectionContent.Bytes(), nil
}
