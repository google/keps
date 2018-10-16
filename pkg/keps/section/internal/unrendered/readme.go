package unrendered

const Readme = `
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
