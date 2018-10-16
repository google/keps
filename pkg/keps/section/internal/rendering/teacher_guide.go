package rendering

import (
	"bytes"
	"text/template"

	"github.com/calebamiles/keps/pkg/keps/section/internal/unrendered"
)

const (
	TeacherGuideName     = "Teacher Guide"
	TeacherGuideFilename = "guides/teacher.md"
)

func NewTeacherGuide(info InfoProvider) ([]byte, error) {
	sectionContent := &bytes.Buffer{}

	t, err := template.New(TeacherGuideName).Parse(unrendered.TeacherGuide)
	if err != nil {
		return nil, err
	}

	err = t.Execute(sectionContent, info)
	if err != nil {
		return nil, err
	}

	return sectionContent.Bytes(), nil
}
