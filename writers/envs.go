package writers

import (
	"io"
	"text/template"
)

var (
	envsTpl = template.Must(template.New("envs").Parse(`// Envs holds the envirnoment-specific configurations so that
// they can easily be accessed by e.g. Envs.Default
var Envs = struct{ 
{{range $_, $k := .Envs}}	{{$k}} Config
{{end}}}{}

var envMap = map[string]*Config{
{{range $k, $v  := .Envs}}	"{{$k}}": &Envs.{{$v}},
{{end}}}

// Get returns the config matching 'env' if found, otherwie the default config.
// The bool return value indicates, if a match was found (true) or the default config
// is returned
func Get(env string) (*Config, bool) {
	if c, ok := envMap[env]; ok {
		return c, true
	}
	return &Envs.Default, false
}
`))
)

//WriteEnvs writes
func WriteEnvs(w io.Writer, envs map[string]string) error {
	return envsTpl.Execute(w, struct {
		Envs map[string]string
	}{Envs: envs})
}
