package sigs

import (
	"path/filepath"
	"strings"

	"github.com/calebamiles/keps/pkg/sigs/internal/generated"
)

func All() []string {
	return generated.SIGList
}

func Exists(s string) bool {
	return generated.SIGSet[s]
}

func ExtractNameFromPath(p string) string {
	sigName, _ := extractSIGandSubprojectNamesFromPath(p)

	return sigName
}

func ExtractSubprojectNameFromPath(p string) string {
	_, subprojectName := extractSIGandSubprojectNamesFromPath(p)

	return subprojectName
}

func extractSIGandSubprojectNamesFromPath(p string) (string, string) {
	var sigName string
	var subprojectName string

	sigs := map[string]bool{}
	subprojects := map[string]bool{}

	components := strings.Split(p, string(filepath.Separator))
	for i, possibleSIG := range components {
		var possibleSubproject string
		if i < len(components) - 1  {
			possibleSubproject = components[i + 1]
		}

		if generated.SIGSet[possibleSIG] {
			sigs[possibleSIG] = true
			sigName = possibleSIG
		}

		if generated.SIGSubprojectMapping[possibleSIG][possibleSubproject] {
			subprojects[possibleSubproject] = true
			subprojectName = possibleSubproject
		}
	}

	switch {
	case len(sigs) > 1:
		sigName = "" // couldn't determine SIG uniquely
	case len(subprojects) > 1:
		subprojectName = "" // couldn't determine subproject uniquely
	}

	return sigName, subprojectName
}
