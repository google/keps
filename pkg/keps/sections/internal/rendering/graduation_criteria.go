package rendering

import (
	"bytes"
	"text/template"

	"github.com/calebamiles/keps/pkg/keps/sections/internal/unrendered"
)

const (
	GraduationCriteriaName     = "Graduation Criteria"
	GraduationCriteriaFilename = "graduation_criteria.md"
)

func NewGraduationCriteria(info InfoProvider) ([]byte, error) {
	sectionContent := &bytes.Buffer{}

	t, err := template.New(GraduationCriteriaName).Parse(unrendered.GraduationCriteria)
	if err != nil {
		return nil, err
	}

	err = t.Execute(sectionContent, info)
	if err != nil {
		return nil, err
	}

	return sectionContent.Bytes(), nil
}
