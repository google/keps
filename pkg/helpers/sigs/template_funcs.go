package sigs

import (
	"strings"
)

func canonicalSIGName(raw string) string {
	return "sig-" + strings.Replace(strings.ToLower(raw), " ", "-", -1)
}

func canonicalSubprojectName(raw string) string {
	return strings.Replace(strings.ToLower(raw), " ", "-", -1)
}
