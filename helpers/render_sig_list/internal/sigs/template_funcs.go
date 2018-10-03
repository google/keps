package sigs

import (
	"strings"
	"text/template"
)

func canonicalName(raw string) string {
	return strings.Replace(strings.ToLower(raw), " ", "-", -1)
}

func canonicalSIGName(raw string) string {
	return "sig-" + canonicalName(raw)
}

var upstreamListTemplateFuncs = template.FuncMap{
		"canonicalName": canonicalName,
		"canonicalSIGName": canonicalSIGName,
}
