package unrendered

const Readme = `
# {{.Title}}

- **Authors: {{joinComma .Authors}}**
- **Sponsoring SIG: [{{displayName .OwningSIG}}](https://github.com/kubernetes/community/tree/master/{{.OwningSIG}}/README.md)**
- **Status: {{.State}}**
- **Last Updated: {{.LastUpdated}}**

## Table of Contents
{{- with removeReadme .SectionLocations}}
{{range .}}
1. [{{sectionName .}}]({{. -}})
{{end -}}
{{end -}}
`
