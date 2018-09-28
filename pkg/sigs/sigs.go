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
	var possibleSubproject string

	sigs := map[string]bool{}
	subprojects := map[string]bool{}

	components := strings.Split(p, string(filepath.Separator))
	for i, possibleSIG := range components {
		isSIG := generated.SIGSet[possibleSIG]
		if !isSIG {
			continue // check the next component for a SIG quickly
		}

		sigs[possibleSIG] = true
		sigName = possibleSIG

		if i > len(components) - 2 {
			break // no need to look for a nested subproject past here
		}

		possibleSubproject = components[i + 1]
		sigSubprojects := generated.SIGSubprojectMapping[sigName]
		isSubproject := sigSubprojects[possibleSubproject]

		if isSubproject {
			subprojects[possibleSubproject] = true
			subprojectName = possibleSubproject
		}
	}

	switch {
	case len(sigs) > 1:
		sigName = "" // couldn't determine SIG uniquely
		subprojectName = "" // couldn't determine SIG/subproject uniquely
	case len(subprojects) > 1:
		// consider subproject information to be optional for now
		subprojectName = "" // couldn't determine subproject uniquely
	}

	return sigName, subprojectName
}
