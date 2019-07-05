package plugins

import (
	"io"
	"strings"
	"text/template"

	"github.com/theliebeskind/genfig/types"
)

type substitutorPlugin struct {
	s   types.SchemaMap
	tpl *template.Template
}

var (
	substitutor = substitutorPlugin{
		s: types.SchemaMap{},
		tpl: template.Must(template.
			New("substitutor").
			Funcs(template.FuncMap{
				// Convert an env var name to a substitution path
				"makeSubstPath": func(s string) string {
					return strings.ToLower(strings.Join(strings.Split(s, "_")[1:], "."))
				},
				// Convert an env var name to a Config path
				"makePath": func(s string) string {
					return strings.Join(strings.Split(s, "_")[1:], ".")
				},
			}).
			Parse(`import (
	"strings"
)

const (
	maxSubstitutionIteraions = 5
)

var (
	raw Config
)

// Substitute replaces all.
// The return value informs, whether all substitutions could be
// applied within {maxSubstitutionIteraions} or not
func (c *Config) Substitute() bool {
	c.ResetSubstitution()

	// backup the "raw" configuration
	raw = *c

	run := 0
	for {
		if run == maxSubstitutionIteraions {
			return false
		}
		if c.substitute() == 0 {
			return true
		}
	}
}

// ResetSubstitution resets the configuration to the state,
// before the substitution was applied
func (c *Config) ResetSubstitution() {
	c = &raw
}

// substitute tries to replace all substitutions in strings
func (c *Config) substitute() int {
	cnt := 0

	r := strings.NewReplacer({{range $_, $v := .}}{{if eq $v.Content "string"}}
		"${{"{"}}{{makeSubstPath $v.Path}}{{"}"}}", c.{{makePath $v.Path}},
	{{end}}{{end}})

	{{range $_, $v := .}}{{if eq $v.Content "string"}}
	if strings.Contains(c.{{makePath $v.Path}}, "${") {
		cnt += 1
		c.{{makePath $v.Path}} = r.Replace(c.{{makePath $v.Path}})
		if !strings.Contains(c.{{makePath $v.Path}}, "${") {
			cnt -= 1
		} 
	}
	{{end}}{{end}}

	return cnt
}
`))}
)

func init() {
	// "register" plugin
	Plugins["substitutor"] = &substitutor
}

// GetInitCall returns the availibility and the string of the
// function to be called on init
func (p *substitutorPlugin) GetInitCall() (string, bool) {
	return "Current.Substitute()", true
}

// SetSchemaMap sets the schema to be used when WriteTo is called
func (p *substitutorPlugin) SetSchemaMap(s types.SchemaMap) {
	p.s = s
}

// WriteTo performs the acutal writing to a buffer (or io.Writer).
// For this plugin, the template is simply "rendered" into the writer.
func (p *substitutorPlugin) WriteTo(w io.Writer) (l int64, err error) {
	err = p.tpl.Execute(w, p.s)
	return
}
