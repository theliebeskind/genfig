package plugins

import (
	"io"
	"strings"
	"text/template"

	"github.com/thlcodes/genfig/models"
)

type writeToEnvPlugin struct {
	s   models.SchemaMap
	tpl *template.Template
}

var (
	writeToEnv = writeToEnvPlugin{
		s: models.SchemaMap{},
		tpl: template.Must(template.
			New("writeToEnv").
			Funcs(template.FuncMap{
				"upper":     strings.ToUpper,
				"hasPrefix": strings.HasPrefix,
				// Remove root (usually "Config_") from env var name
				"cleanPrefixEnv": func(s string) string {
					return strings.Join(strings.Split(s, "_")[1:], "_")
				},
				// Converte an env var name to a Config path
				"makePath": func(s string) string {
					return strings.Join(strings.Split(s, "_")[1:], ".")
				},
			}).
			Parse(`import (
	"fmt"
	"os"
	"encoding/json"
)

var (
	_ = os.Setenv
	_ = fmt.Sprintf
	_ = json.Marshal
)

func (c *Config) WriteToEnv() {
	var buf []byte
	_ = buf
{{range $_, $v := .}}{{if not $v.IsStruct}}
	{{if hasPrefix $v.Content "[]"}}
	buf, _ = json.Marshal(c.{{makePath $v.Path}})
	_ = os.Setenv("{{cleanPrefixEnv (upper $v.Path)}}", string(buf))
	{{else}}
	_ = os.Setenv("{{cleanPrefixEnv (upper $v.Path)}}", fmt.Sprintf("%v", c.{{makePath $v.Path}}))
	{{end}}
{{end}}{{end}}
}
`))}
)

func init() {
	// "register" plugin
	Plugins["write_to_env"] = &writeToEnv
}

// GetInitCall returns the availibility and the string of the
// function to be called on init
func (p *writeToEnvPlugin) GetInitCall() (string, bool) {
	return "", false
}

// SetSchemaMap sets the schema to be used when WriteTo is called
func (p *writeToEnvPlugin) SetSchemaMap(s models.SchemaMap) {
	p.s = s
}

// WriteTo performs the acutal writing to a buffer (or io.Writer).
// For this plugin, the template is simply "rendered" into the writer.
func (p *writeToEnvPlugin) WriteTo(w io.Writer) (l int64, err error) {
	err = p.tpl.Execute(w, p.s)
	return
}
