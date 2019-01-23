package rendering

import (
	"path/filepath"
	"strings"
)

// TODO move after all tests green
func NameForFilename(filename string) string {
	knownSectionName := toplevelFilenameSet[filename]
	if knownSectionName == "" {
		return guessNameFromFilename(filename)
	}

	return knownSectionName
}

func guessNameFromFilename(filename string) string {
	probableFile := filepath.Base(filename)

	filenameNoExtension := strings.Replace(probableFile, ".md", "", 1)
	return strings.Title(strings.Replace(filenameNoExtension, "_", " ", -1))
}

var toplevelFilenameSet = map[string]string{
	DeveloperGuideFilename:     DeveloperGuideName,
	GraduationCriteriaFilename: GraduationCriteriaName,
	MotivationFilename:         MotivationName,
	OperatorGuideFilename:      OperatorGuideName,
	ReadmeFilename:             ReadmeName,
	SummaryFilename:            SummaryName,
	TeacherGuideFilename:       TeacherGuideName,
}
