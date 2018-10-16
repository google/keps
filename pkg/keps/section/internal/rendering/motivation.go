package rendering

import (
	"bytes"
	"text/template"

	"github.com/calebamiles/keps/pkg/keps/section/internal/unrendered"
)

const (
	MotivationName     = "Motivation"
	MotivationFilename = "motivation.md"
)

func NewMotivation(info InfoProvider) ([]byte, error) {
	sectionContent := &bytes.Buffer{}

	t, err := template.New(MotivationName).Parse(unrendered.Motivation)
	if err != nil {
		return nil, err
	}

	err = t.Execute(sectionContent, info)
	if err != nil {
		return nil, err
	}

	return sectionContent.Bytes(), nil
}
