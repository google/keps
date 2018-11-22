package sections

import (
	"path/filepath"
	"strings"

	"github.com/calebamiles/keps/pkg/keps/sections/internal/rendering"
)

type section interface {
	Filename() string
	Name() string
	Content() []byte
	Persist() error
	Erase() error
}

// TODO add tests
func sectionNameForFilename(filename string) string {
	knownSectionName := toplevelFilenameSet[filename]
	if knownSectionName == "" {
		return guessNameFromFilename(filename)
	}

	return knownSectionName
}

var toplevelFilenameSet = map[string]string{
	rendering.DeveloperGuideFilename:     rendering.DeveloperGuideName,
	rendering.GraduationCriteriaFilename: rendering.GraduationCriteriaName,
	rendering.MotivationFilename:         rendering.MotivationName,
	rendering.OperatorGuideFilename:      rendering.OperatorGuideName,
	rendering.ReadmeFilename:             rendering.ReadmeName,
	rendering.SummaryFilename:            rendering.SummaryName,
	rendering.TeacherGuideFilename:       rendering.TeacherGuideName,
}

func guessNameFromFilename(filename string) string {
	probableFile := filepath.Base(filename)

	filenameNoExtension := strings.Replace(probableFile, ".md", "", 1)
	return strings.Title(strings.Replace(filenameNoExtension, "_", " ", -1))
}

const (
	DeveloperGuide     = rendering.DeveloperGuideName
	GraduationCriteria = rendering.GraduationCriteriaName
	Motivation         = rendering.MotivationName
	OperatorGuide      = rendering.OperatorGuideName
	Readme             = rendering.ReadmeName
	Summary            = rendering.SummaryName
	TeacherGuide       = rendering.TeacherGuideName
)
