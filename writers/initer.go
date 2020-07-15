package writers

import (
	"io"
	"sort"
	"text/template"
)

var (
	initTpl = template.Must(template.New("init").Parse(`import (
	"os"
	"fmt"
)

var (
	_ = os.Getenv
	_ = fmt.Printf
)

// Current is the current config, selected by the curren env and
// updated by the availalbe env vars
var Current *Config

// This init tries to retrieve the current environment via the
// common env var 'ENV' and applies activated plugins
func init() {
	Current, _ = Get(os.Getenv("ENV"))
	{{range $_, $v := .PluginCalls}}
	{{$v}}
	{{end}}
}
	
`))
)

//WriteInit writes
func WriteInit(w io.Writer, pluginCalls map[string]string) error {
	keys := []string{}
	for k := range pluginCalls {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	calls := []string{}
	for _, k := range keys {
		calls = append(calls, pluginCalls[k])
	}
	return initTpl.Execute(w, struct {
		PluginCalls []string
	}{PluginCalls: calls})
}
