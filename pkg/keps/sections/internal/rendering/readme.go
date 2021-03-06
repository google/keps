package rendering

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/calebamiles/keps/pkg/keps/sections/internal/unrendered"
)

const (
	ReadmeName     = "README"
	ReadmeFilename = "README.md"
)

func NewReadme(info InfoProvider) ([]byte, error) {
	sectionContent := &bytes.Buffer{}

	t, err := template.New(ReadmeName).Funcs(funcMap).Parse(unrendered.Readme)
	if err != nil {
		return nil, err
	}

	err = t.Execute(sectionContent, info)
	if err != nil {
		return nil, err
	}

	return sectionContent.Bytes(), nil
}

func joinComma(ss []string) string {
	return strings.Join(ss, ", ")
}

func removeReadme(ss []string) []string {
	secs := []string{}

	for i := range ss {
		if ss[i] == ReadmeFilename {
			continue
		}

		secs = append(secs, ss[i])
	}

	return secs
}

// TODO() remove this once generating SIG info from path stores display name there
func sigDisplayName(s string) string {
	return strings.Title(strings.TrimSpace(strings.Replace(strings.Replace(s, "-", " ", -1), "sig", "", 1)))
}

var funcMap = template.FuncMap{
	"joinComma":    joinComma,
	"displayName":  sigDisplayName,
	"sectionName":  NameForFilename,
	"removeReadme": removeReadme,
}
