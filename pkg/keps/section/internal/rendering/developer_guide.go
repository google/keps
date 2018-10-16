package rendering

import (
	"bytes"
	"text/template"

	"github.com/calebamiles/keps/pkg/keps/section/internal/unrendered"
)

const (
	DeveloperGuideName     = "Developer Guide"
	DeveloperGuideFilename = "guides/developer.md"
)

func NewDeveloperGuide(info InfoProvider) ([]byte, error) {
	sectionContent := &bytes.Buffer{}

	t, err := template.New(DeveloperGuideName).Parse(unrendered.DeveloperGuide)
	if err != nil {
		return nil, err
	}

	err = t.Execute(sectionContent, info)
	if err != nil {
		return nil, err
	}

	return sectionContent.Bytes(), nil
}
