package sigs

import (
	"io"
	"text/template"
)

func RenderList(l *upstreamSIGList, f io.Writer) error {
	t := template.Must(template.New("upstream SIG list").Funcs(upstreamListTemplateFuncs).Parse(upstreamListTemplate))

	return t.Execute(f, l)
}
