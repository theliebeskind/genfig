package genfig

import (
	"io"
	"sort"
	"text/template"
)

var (
	envsTpl = template.Must(template.New("envs").Parse(`var Envs = struct{ 
{{range $env := .Envs}}	{{$env}} Config
{{end}}}{}
`))
)

func writeEnvs(w io.Writer, envs []string) error {
	sort.Strings(envs)
	return envsTpl.Execute(w, struct {
		Envs []string
	}{Envs: envs})
}
