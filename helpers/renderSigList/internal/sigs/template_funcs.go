package sigs

import (
	"path/filepath"
	"strings"
	"text/template"
)

func canonicalName(raw string) string {
	return strings.Replace(strings.ToLower(raw), " ", "-", -1)
}

func canonicalSIGName(raw string) string {
	return "sig-" + canonicalName(raw)
}

func joinPath(elem ...string) string {
	return filepath.Join(elem...)
}

var upstreamListTemplateFuncs = template.FuncMap{
	"canonicalName":    canonicalName,
	"canonicalSIGName": canonicalSIGName,
	"joinPath":         joinPath,
}
