package render

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/calebamiles/keps/helpers/convert/internal/metadata"
)

func NewReadme(meta *metadata.New, sectionLocations map[string]string) ([]byte, error) {
	info := &renderingInfo{
		New: meta, // TODO think about renaming metadata.New to metadata.CurrentSpec or similar
	}

	for name, location := range sectionLocations {
		info.Sections = append(info.Sections, &sectionInfo{
			Name:     sectionName(name),
			Filename: location,
		})
	}

	t, err := template.New("README").Funcs(funcMap).Parse(readmeTemplate)
	if err != nil {
		return nil, err
	}

	rendered := &bytes.Buffer{}
	err = t.Execute(rendered, info)
	if err != nil {
		return nil, err
	}

	return rendered.Bytes(), nil
}

func sectionName(s string) string {
	return strings.Title(strings.TrimPrefix(s, "## "))
}

func joinComma(ss []string) string {
	return strings.Join(ss, ", ")
}

// TODO() remove this once generating SIG info from path stores display name there
func sigDisplayName(s string) string {
	return strings.Title(strings.TrimSpace(strings.Replace(strings.Replace(s, "-", " ", -1), "sig", "", 1)))
}

var funcMap = template.FuncMap{
	"joinComma":   joinComma,
	"displayName": sigDisplayName,
}

type sectionInfo struct {
	Name     string
	Filename string
}

type renderingInfo struct {
	*metadata.New
	Sections []*sectionInfo
}

const readmeTemplate = `
# {{.Title}}

- **Authors: {{joinComma .Authors}}**
- **Sponsoring SIG: [{{displayName .OwningSIG}}](https://github.com/kubernetes/community/tree/master/{{.OwningSIG}}/README.md)**
- **Status: {{.State}}**
- **Last Updated: {{.LastUpdated}}**

## Table of Contents
{{- with .Sections}}
{{range .}}
1. [{{.Name}}]({{.Filename -}})
{{end -}}
{{end -}}
`
