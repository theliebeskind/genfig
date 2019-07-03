package writers

import (
	"io"
	"text/template"
)

var (
	initTpl = template.Must(template.New("init").Parse(`import (
		"os"
	)
	
	// Current is the current config, selected by the curren env and
	// updated by the availalbe env vars
	var Current *Config
	
	// This init tries to retreive the current environemnt via the
	// common env var 'ENV' and applies activated plugins
	func init() {
		Current, _ = Get(os.Getenv("ENV"))
	
		// apply activated plugins
		Current.UpdateFromEnv()
	}
	
`))
)

//WriteInit writes
func WriteInit(w io.Writer, pluginCalls map[string]string) error {
	return initTpl.Execute(w, struct {
		PluginCalls map[string]string
	}{PluginCalls: pluginCalls})
}
