package render

import (
	"bytes"
	"text/template"
)

func DeveloperGuide(sections map[string][]byte) ([]byte, error) {
	t, err := template.New("developer guide").Parse(developerGuideTemplate)
	if err != nil {
		return nil, err
	}

	rendered := &bytes.Buffer{}
	err = t.Execute(rendered, sections)
	if err != nil {
		return nil, err
	}

	return rendered.Bytes(), nil
}

const developerGuideTemplate = `
{{ $proposalHeading := "## Proposal" -}}
{{ $drawbacksHeading := "## Drawbacks" -}}
{{ $alternativesHeading := "## Alternatives" -}}

{{- $proposalHeading }}

{{ index . $proposalHeading | printf "%s" }}

{{ $drawbacksHeading }}

{{ index . $drawbacksHeading | printf "%s" }}

{{ $alternativesHeading | printf "%s" }}

{{ index . $alternativesHeading | printf "%s" }}
`
