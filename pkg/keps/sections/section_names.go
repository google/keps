package sections

import (
	"path/filepath"
	"strings"

	"github.com/calebamiles/keps/pkg/keps/sections/internal/rendering"
)

const (
	Summary            = rendering.SummaryName
	Motivation         = rendering.MotivationName
	DeveloperGuide     = rendering.DeveloperGuideName
	OperatorGuide      = rendering.OperatorGuideName
	TeacherGuide       = rendering.TeacherGuideName
	GraduationCriteria = rendering.GraduationCriteriaName
	Readme             = rendering.ReadmeName
)

func Filename(name string) string {
	switch name {
	case Summary:
		return rendering.SummaryFilename
	case Motivation:
		return rendering.MotivationFilename
	case TeacherGuide:
		return rendering.TeacherGuideFilename
	case DeveloperGuide:
		return rendering.DeveloperGuideFilename
	case OperatorGuide:
		return rendering.OperatorGuideFilename
	case GraduationCriteria:
		return rendering.GraduationCriteriaFilename
	case Readme:
		return rendering.ReadmeFilename
	default:
		return filepath.Clean(strings.Replace(strings.ToLower(name), " ", "_", -1) + ".md")
	}
}

func NameForFilename(filename string) string {
	// TODO move out of rendering/readme.go
	return rendering.NameForFilename(filename)
}
